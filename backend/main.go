package main

import (
	"backend/database"
	"backend/routes"
	"log"
	"net/http"
)

func main() {
	// Подключение к базе данных
	database.InitDB()

	// Настройка маршрутов
	http.HandleFunc("/api/bookings", routes.BookingHandler)

	// Запуск сервера
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
