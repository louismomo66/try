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

type User struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	// Auto-migrate schema
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("failed to migrate:", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from Go + GORM + PostgreSQL")
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		var users []User
		result := db.Find(&users)
		if result.Error != nil {
			http.Error(w, "Error fetching users", http.StatusInternalServerError)
			return
		}

		for _, user := range users {
			fmt.Fprintf(w, "ID: %d, Name: %s\n", user.ID, user.Name)
		}
	})

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
