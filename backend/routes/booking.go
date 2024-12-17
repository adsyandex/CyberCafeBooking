package routes

import (
	"backend/database"
	"backend/models"
	"encoding/json"
	"net/http"
)

// BookingHandler обрабатывает запросы
func BookingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		// Обработка POST-запроса
		var booking models.Booking
		json.NewDecoder(r.Body).Decode(&booking)

		query := `INSERT INTO bookings (name, contact, computer, start_time, end_time)
				  VALUES (?, ?, ?, ?, ?)`
		_, err := database.DB.Exec(query, booking.Name, booking.Contact, booking.Computer, booking.StartTime, booking.EndTime)
		if err != nil {
			http.Error(w, "Ошибка сохранения бронирования: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	case http.MethodGet:
		// Обработка GET-запроса
		rows, err := database.DB.Query("SELECT id, name, contact, computer, start_time, end_time FROM bookings")
		if err != nil {
			http.Error(w, "Ошибка получения данных: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var bookings []models.Booking
		for rows.Next() {
			var booking models.Booking
			if err := rows.Scan(&booking.ID, &booking.Name, &booking.Contact, &booking.Computer, &booking.StartTime, &booking.EndTime); err != nil {
				http.Error(w, "Ошибка обработки данных: "+err.Error(), http.StatusInternalServerError)
				return
			}
			bookings = append(bookings, booking)
		}
		json.NewEncoder(w).Encode(bookings)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}
