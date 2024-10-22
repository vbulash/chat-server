package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Env Структура для распаковки .env
type Env struct {
	Host     string
	Database string
	Username string
	Password string
	Port     int
}

// Config Экземпляр конфигурации .env
var Config *Env

// ChatType Тип записи chats
type ChatType struct {
	ID        int        `db:"id"`
	Title     string     `db:"title"`
	Body      string     `db:"body"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

// LoadConfig Загрузка конфигурации из .env
func LoadConfig() *Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env")
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("Ошибка преобразования DB_PORT из .env: %v\n", err)
	}

	return &Env{
		Host:     os.Getenv("DB_HOST"),
		Database: os.Getenv("DB_DATABASE"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Port:     port,
	}
}
