package main

import (
	"data_parser/db"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	if _, err := os.ReadFile(".env"); err != nil {
		log.Fatalln("Need .env in data_parser folder\nwith: POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_HOST, POSTGRES_PORT, POSTGRES_DB!!")
	}
	p, err := db.Connect()
	if err != nil {
		log.Fatalln("error when connecting to db")
	}
	log.Println("Connected to db!")
	defer p.Close()

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		log.Fatalln("provide a csv!")
	}
	readCSV(p, argsWithoutProg[0])
	os.Exit(0)
}

func readCSV(p *pgxpool.Pool, path string) {
	dat, err := os.Open(path)
	if err != nil {
		log.Fatalf("cant read file %s\nerr: %s\n", path, err)
	}
	defer dat.Close()

	reader := csv.NewReader(dat)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(record)
	}
	// log.Println(records)
}

func insert(p *pgxpool.Pool, record []string) {
	// sql := `insert into `
}
