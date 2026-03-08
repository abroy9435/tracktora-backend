package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB is the global database connection pool
var DB *pgxpool.Pool

// ConnectDB initializes the database connection
func ConnectDB() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// Create a connection pool
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Ping the database to verify the connection
	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Database is not responding: %v\n", err)
	}

	DB = pool
	fmt.Println("✅ Successfully connected to PostgreSQL!")
}
