package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

type ConfigStruct struct {
	Host     string
	Database string
	Username string
	Password string
	Port     int
}

var Config *ConfigStruct

type NoteType struct {
	Id        int        `db:"id"`
	Title     string     `db:"title"`
	Body      string     `db:"body"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

func LoadConfig() *ConfigStruct {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env")
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("Ошибка преобразования DB_PORT из .env: %v\n", err)
	}

	return &ConfigStruct{
		Host:     os.Getenv("DB_HOST"),
		Database: os.Getenv("DB_DATABASE"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Port:     port,
	}
}
