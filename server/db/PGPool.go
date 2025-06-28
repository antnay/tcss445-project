package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func PGConnect() (*pgxpool.Pool, error) {
	port := os.Getenv("POSTGRES_PORT")
	if len(port) == 0 {
		port = "5434"
	}
	uri := os.Getenv("POSTGRES_URI")
	if len(uri) == 0 {

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
			host = "pg-db"
		}
		db := os.Getenv("POSTGRES_DB")
		if len(db) == 0 {
			db = "typesheet"
		}
		uri = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", user, password, host, port, db)
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

	log.Printf("Connected to postgres on port %s", port)
	return pool, nil
}
