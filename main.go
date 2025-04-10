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
	redisClient *redis.Client
	ctx         = context.Background()
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or loading failed, using existing env vars")
	}

	// PostgreSQL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()

	// Redis
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379" // fallback if not set
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
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

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
