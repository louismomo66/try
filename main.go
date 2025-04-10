package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or loading failed, using existing env vars")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var result struct {
			Now string
		}
		if err := db.Raw("SELECT NOW()").Scan(&result).Error; err != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Connected to Postgres! Time: %s", result.Now)
	})

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
