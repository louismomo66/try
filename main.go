package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var (
	ctx         = context.Background()
	redisClient *redis.Client
)

func main() {
	// Load .env (optional on Render, but useful locally)
	_ = godotenv.Load()

	// PostgreSQL setup
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()

	// Redis setup (Render-friendly)
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL not set")
	}
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatal("Invalid REDIS_URL:", err)
	}
	redisClient = redis.NewClient(opts)
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	// Routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var now string
		err := db.QueryRow("SELECT NOW()").Scan(&now)
		if err != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Connected to Postgres! Time: %s", now)
	})

	http.HandleFunc("/visits", func(w http.ResponseWriter, r *http.Request) {
		count, err := redisClient.Incr(ctx, "visits").Result()
		if err != nil {
			http.Error(w, "Redis error", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Visit count: %d", count)
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
