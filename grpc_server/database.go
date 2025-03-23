package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Подключение и настройка PostgreSQL через GORM

var DB *gorm.DB // Глобальное соединение с базой

// Структура модели User соответствует таблице в БД
type User struct {
	ID    string `gorm:"primaryKey"` // UUID как строка
	Name  string
	Email string `gorm:"unique"` // Email должен быть уникальным
	Age   int32
}

// Инициализация подключения к базе данных
func InitDB() {
	// Получаем значения из переменных окружения или ставим дефолт
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "root"),
		getEnv("DB_NAME", "go_test_db"),
		getEnv("DB_PORT", "5432"),
	)

	// Подключаемся к PostgreSQL через GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	// Автоматически создаём таблицу, если не существует
	db.AutoMigrate(&User{})
	DB = db
}

// Получение переменной окружения с fallback-значением
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
