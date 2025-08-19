package main

import (
	"context"
	"data_parser/db"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"googlemaps.github.io/maps"
)

func main() {

	if _, err := os.ReadFile(".env"); err != nil {
		log.Fatalln("Need .env in data_parser folder\nwith: POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_HOST, POSTGRES_PORT, POSTGRES_DB!!")
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s\n", err)
	}
	p, err := db.Connect()
	if err != nil {
		log.Fatalln("error when connecting to db")
	}
	log.Println("Connected to db!")
	defer p.Close()

	log.Println(os.Getenv("MAPS_API"))
	c, err := maps.NewClient(maps.WithAPIKey(os.Getenv("MAPS_API")))
	if err != nil {
		log.Fatalf("fatal error: %s\n", err)
	}

	// addNeighborhood(p, c)

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		log.Fatalln("provide a csv!")
	}
	readCSV(p, c, argsWithoutProg[0])
	os.Exit(0)
}

func readCSV(p *pgxpool.Pool, c *maps.Client, path string) {
	dat, err := os.Open(path)
	if err != nil {
		log.Fatalf("cant read file %s\nerr: %s\n", path, err)
	}
	defer dat.Close()

	reader := csv.NewReader(dat)
	count := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(record)

		if count != 0 {
			// log.Println(len(record))
			insert(p, c, record)
			// for i := 0; i < 10; i++ {
			// 	_, err = reader.Read()
			// 	if err == io.EOF {
			// 		break
			// 	}
			// 	count++
			// }
		}
		log.Println(count)
		count++
	}
	// log.Println(records)
}

func insert(p *pgxpool.Pool, c *maps.Client, record []string) {
	caseno := record[0]
	lat := record[1]
	lon := record[2]
	category := record[3]
	address := record[4]
	date := record[5]
	time := record[6]
	neighborhood := record[8]

	// TODO:
	// if redacted continue
	// regex to filter out xxxx block of ...
	// reorganize categories

	req := maps.GeocodingRequest{Address: fmt.Sprintf("%s, tacoma, WA", address)}

	loc, err := c.Geocode(context.Background(), &req)
	if err != nil {
		log.Printf("fatal error: %s\n", err)
	}

	// log.Println(loc[0].AddressComponents[7])
	var zip string

	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()
	comp := loc[0].AddressComponents
	len := len(comp)
	goodZip := false
	goodNei := false
	for i := len - 1; i > 0; i-- {
		curr := comp[i]
		if curr.Types[0] == "postal_code" {
			// log.Println(curr.LongName)
			zip = curr.LongName
			goodZip = true
		}
		if curr.Types[0] == "neighborhood" {
			neighborhood = curr.LongName
			goodNei = true
		}
	}
	if !(goodZip && goodNei) {
		return
	}

	tx, err := p.Begin(context.Background())
	if err != nil {
		log.Fatalf("Could not start db transaction: %s", err)
	}
	sql := `SELECT *
		FROM add_crime_incident_partition(p_city_name := $1,
	                              p_state_name := $2,
	                              p_country_name := $3,
	                              p_source_name := $4,
	                              p_latitude := $5,
	                              p_longitude := $6,
	                              p_street_address := $7,
	                              p_postal_code := $8,
	                              p_crime_category_name :=$9,
								  p_neighborhood := $13,
	                              p_incident_date := $10::date,
	                              p_incident_time := $11::time,
	                              p_case_num := $12)`

	_, err = tx.Exec(context.Background(), sql,
		"Tacoma", "Washington", "United States", "City of Tacoma Reported Crime (Tacoma)",
		lat, lon, address, zip, category, date, time, caseno, neighborhood)

	if err != nil {
		log.Printf("errr rolling back! %s", err)
		tx.Rollback(context.Background())
	}
	tx.Commit(context.Background())

}

func addNeighborhood(p *pgxpool.Pool, c *maps.Client) {
	sql := `
	select street_address, postal_code,neighborhood_id 
	from addresses
	where neighborhood_id is null
	`

	q, err := p.Query(context.Background(), sql)
	if err != nil {
		log.Println("uh oh ", err)
	}
	for q.Next() {
		row, err := q.Values()
		if err != nil {
			log.Println("uh oh ", err)
		}

		log.Println(row...)
		// addr := row[len(row)-4]
		// log.Println(addr)

		req := maps.GeocodingRequest{Address: fmt.Sprintf("%s, tacoma, WA, %s", row[0], row[1])}

		loc, err := c.Geocode(context.Background(), &req)
		if err != nil {
			log.Printf("fatal error: %s\n", err)
		}
		log.Println(loc)

		defer func() {
			if err := recover(); err != nil {
				log.Println("panic occurred:", err)
			}
		}()
		comp := loc[0].AddressComponents
		len := len(comp)

		good := false
		var neighborhood string
		// log.Println(len)
		for i := range len {
			curr := comp[i]
			// log.Println(curr.Types)
			if curr.Types[0] == "neighborhood" {
				// neighborhood = curr.LongName
				log.Println(curr.LongName)
				good = true
			}
		}
		if !good {
			continue
		}

		continue

		// sql = `
		// select neighborhood_id from neighborhoods where neighborhood_name = $1
		// `
		var neighborhoodId int
		p.QueryRow(context.Background(), sql, neighborhood).Scan(&neighborhoodId)

		if neighborhoodId == 0 {
			sql = `
			WITH city_query AS (SELECT city_id
                    FROM cities
                    WHERE city_name = 'Tacoma')
			INSERT
			INTO neighborhoods (neighborhood_name, city_id)
			SELECT data.neighborhood_name,
		   		s.city_id
			FROM city_query s
		     	CROSS JOIN(VALUES ($1))
			AS data(neighborhood_name) 
			returning neighborhood_id`
			p.QueryRow(context.Background(), sql, neighborhood).Scan(&neighborhoodId)
		}

		sql = `
		UPDATE addresses
		SET neighborhood_id = $1
		WHERE street_address = $2 and postal_code = $3
		`
		var update []string
		p.QueryRow(context.Background(), sql, neighborhoodId, row[0], row[1]).Scan(update)
		log.Println()
	}
}
