package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func Connect() (*pgxpool.Pool, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s\n", err)
	}

	port := os.Getenv("POSTGRES_PORT")
	if len(port) == 0 {
		port = "5432"
	}
	uri := os.Getenv("POSTGRES_URL")
	if len(uri) == 0 {
		log.Println("Using parts to construct db uri")

		user := os.Getenv("POSTGRES_USER")
		if len(user) == 0 {
			user = "postgres"
		}
		password := os.Getenv("POSTGRES_PASSWORD")
		if len(password) == 0 {
			password = "postgres"
		}
		host := os.Getenv("POSTGRES_HOST")
		if len(host) == 0 {
			host = "postgres"
		}
		db := os.Getenv("POSTGRES_DB")
		if len(db) == 0 {
			db = "postgres"
		}
		uri = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, password, host, port, db)
	} else {
		log.Println("Using POSTGRES_URL")
	}

	log.Println("uri: ", uri)

	config, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	log.Printf("\033[32mConnected to postgres on port %s\033[0m\n", port)
	return pool, nil
}
