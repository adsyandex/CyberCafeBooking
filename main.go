package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"context"
	"log"
	"os"
)

// Booking represents a single booking entry.
type Booking struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Contact   string    `json:"contact"`
	Computer  int       `json:"computer"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// BookingStore manages bookings in memory.
type BookingStore struct {
	mu       sync.Mutex
	bookings []Booking
	counter  int
	calendar *calendar.Service
}

func (s *BookingStore) AddBooking(b Booking) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check for conflicting bookings
	for _, existing := range s.bookings {
		if existing.Computer == b.Computer &&
			(b.StartTime.Before(existing.EndTime) && b.EndTime.After(existing.StartTime)) {
			return fmt.Errorf("conflicting booking for computer %d", b.Computer)
		}
	}

	// Add booking
	s.counter++
	b.ID = s.counter
	s.bookings = append(s.bookings, b)

	// Add event to Google Calendar
	if err := s.addEventToCalendar(b); err != nil {
		log.Println("Failed to add event to calendar:", err)
	}
	return nil
}

func (s *BookingStore) addEventToCalendar(b Booking) error {
	event := &calendar.Event{
		Summary:     fmt.Sprintf("Booking for %s", b.Name),
		Description: fmt.Sprintf("Computer %d booked by %s (%s)", b.Computer, b.Name, b.Contact),
		Start: &calendar.EventDateTime{
			DateTime: b.StartTime.Format(time.RFC3339),
			TimeZone: "UTC",
		},
		End: &calendar.EventDateTime{
			DateTime: b.EndTime.Format(time.RFC3339),
			TimeZone: "UTC",
		},
	}

	calendarID := "primary" // Replace with your calendar ID if needed
	_, err := s.calendar.Events.Insert(calendarID, event).Do()
	return err
}

func (s *BookingStore) ListBookings() []Booking {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]Booking{}, s.bookings...)
}

func setupGoogleCalendar() (*calendar.Service, error) {
	ctx := context.Background()

	// Load credentials.json
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return nil, fmt.Errorf("Unable to read client secret file: %v", err)
	}

	// Configure OAuth2
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)
	return calendar.NewService(ctx, option.WithHTTPClient(client))
}

func getClient(config *oauth2.Config) *http.Client {
	// Token file management for OAuth2
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", url)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	calendarService, err := setupGoogleCalendar()
	if err != nil {
		log.Fatalf("Unable to set up Google Calendar: %v", err)
	}

	store := &BookingStore{calendar: calendarService}

	http.HandleFunc("/bookings", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// List bookings
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(store.ListBookings())
			return
		} else if r.Method == http.MethodPost {
			// Add a new booking
			var b Booking
			if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
				http.Error(w, "Invalid input", http.StatusBadRequest)
				return
			}

			if b.StartTime.IsZero() || b.EndTime.IsZero() || b.Name == "" || b.Contact == "" || b.Computer <= 0 {
				http.Error(w, "Missing required fields", http.StatusBadRequest)
				return
			}

			if b.StartTime.After(b.EndTime) {
				http.Error(w, "Invalid time range", http.StatusBadRequest)
				return
			}

			if err := store.AddBooking(b); err != nil {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(b)
			return
		}

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Internet Cafe Booking System! Use /bookings to view or create bookings."))
	})

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

