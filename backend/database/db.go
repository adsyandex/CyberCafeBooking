package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// Config структура для конфигурации
type Config struct {
	User     string `json:"db_user"`
	Password string `json:"db_password"`
	Name     string `json:"db_name"`
	Host     string `json:"db_host"`
	Port     string `json:"db_port"`
}

// InitDB подключается к MySQL базе данных
func InitDB() {
	// Читаем файл конфигурации
	file, err := os.Open("config/config.json")
	if err != nil {
		panic("Не удалось открыть файл конфигурации: " + err.Error())
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		panic("Не удалось декодировать конфигурацию: " + err.Error())
	}

	// Формируем строку подключения
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.User, config.Password, config.Host, config.Port, config.Name)

	// Подключаемся к базе данных
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		panic("Не удалось подключиться к базе данных: " + err.Error())
	}

	// Проверяем соединение
	if err := DB.Ping(); err != nil {
		panic("База данных недоступна: " + err.Error())
	}

	fmt.Println("Успешное подключение к базе данных")
}
