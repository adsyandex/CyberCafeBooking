package models

// Booking модель данных для таблицы бронирований
type Booking struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Contact   string `json:"contact"`
	Computer  string `json:"computer"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}
