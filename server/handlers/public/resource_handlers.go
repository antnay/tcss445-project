package public

import (
	"context"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
)

// Returns dynamic filtering options
func (h *Handler) GetCrimesFilterOptions(c *gin.Context) {
	yearsQuery := `
		SELECT DISTINCT EXTRACT(YEAR FROM incident_date) as year 
		FROM (
			SELECT incident_date FROM crime_incidents_2020
			UNION ALL
			SELECT incident_date FROM crime_incidents_2021  
			UNION ALL
			SELECT incident_date FROM crime_incidents_2022
			UNION ALL
			SELECT incident_date FROM crime_incidents_2023
			UNION ALL
			SELECT incident_date FROM crime_incidents_2024
			UNION ALL
			SELECT incident_date FROM crime_incidents_2025
		) all_incidents
		WHERE incident_date IS NOT NULL
		ORDER BY year DESC
	`

	type FilterOptions struct {
		Years         []string `json:"years"`
		CrimeTypes    []string `json:"crimeTypes"`
		Cities        []string `json:"cities"`
		Neighborhoods []string `json:"neighborhoods"`
		Sources       []string `json:"sources"`
	}

	var wg sync.WaitGroup
	var options FilterOptions
	var mutex sync.Mutex

	wg.Add(1)
	go func() {
		defer wg.Done()
		rows, err := h.pool.Query(context.Background(), yearsQuery)
		if err == nil {
			defer rows.Close()
			var years []string
			for rows.Next() {
				var year int
				if err := rows.Scan(&year); err == nil {
					years = append(years, fmt.Sprintf("%d", year))
				}
			}
			mutex.Lock()
			options.Years = years
			mutex.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		rows, err := h.pool.Query(context.Background(), `
			SELECT DISTINCT category_name 
			FROM crime_categories 
			WHERE category_name IS NOT NULL 
			ORDER BY category_name ASC
		`)
		if err == nil {
			defer rows.Close()
			var crimeTypes []string
			for rows.Next() {
				var crimeType string
				if err := rows.Scan(&crimeType); err == nil {
					crimeTypes = append(crimeTypes, crimeType)
				}
			}
			mutex.Lock()
			options.CrimeTypes = crimeTypes
			mutex.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		rows, err := h.pool.Query(context.Background(), `
			SELECT DISTINCT c.city_name 
			FROM cities c
			INNER JOIN addresses a ON c.city_id = a.city_id
			WHERE c.city_name IS NOT NULL 
			ORDER BY c.city_name ASC
		`)
		if err == nil {
			defer rows.Close()
			var cities []string
			for rows.Next() {
				var city string
				if err := rows.Scan(&city); err == nil {
					cities = append(cities, city)
				}
			}
			mutex.Lock()
			options.Cities = cities
			mutex.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		rows, err := h.pool.Query(context.Background(), `
			SELECT DISTINCT n.neighborhood_name 
			FROM neighborhoods n
			INNER JOIN addresses a ON n.neighborhood_id = a.neighborhood_id
			WHERE n.neighborhood_name IS NOT NULL 
			AND n.neighborhood_name != 'Unknown neighborhood'
			ORDER BY n.neighborhood_name ASC
		`)
		if err == nil {
			defer rows.Close()
			var neighborhoods []string
			for rows.Next() {
				var neighborhood string
				if err := rows.Scan(&neighborhood); err == nil {
					neighborhoods = append(neighborhoods, neighborhood)
				}
			}
			mutex.Lock()
			options.Neighborhoods = neighborhoods
			mutex.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		rows, err := h.pool.Query(context.Background(), `
			SELECT DISTINCT source_name 
			FROM data_sources 
			WHERE source_name IS NOT NULL 
			ORDER BY source_name ASC
		`)
		if err == nil {
			defer rows.Close()
			var sources []string
			for rows.Next() {
				var source string
				if err := rows.Scan(&source); err == nil {
					sources = append(sources, source)
				}
			}
			mutex.Lock()
			options.Sources = sources
			mutex.Unlock()
		}
	}()

	wg.Wait()

	c.JSON(200, options)
}
