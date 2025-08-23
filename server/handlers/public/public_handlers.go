package public

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	EIGHT_HOUR  = 3600 * 8
	MILE_APPROX = 69
)

type Handler struct {
	pool *pgxpool.Pool
}

func NewHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{
		pool: pool,
	}
}

func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "public pong",
	})
}

type Crime struct {
	Address      string  `json:"address"`
	Neighborhood string  `json:"neighborhood"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	CrimeType    string  `json:"crime_type"`
	Date         string  `json:"date"`
	Time         string  `json:"time"`
}

type CrimeDump struct {
	Case          string  `json:"case"`
	CrimeCategory string  `json:"crimeCategory"`
	Neighborhood  string  `json:"neighborhood"`
	Street        string  `json:"street"`
	City          string  `json:"city"`
	Zip           string  `json:"zip"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Date          string  `json:"date"`
	Time          string  `json:"time"`
	Source        string  `json:"source"`
}

type CrimeWithDistance struct {
	Case          string  `json:"case"`
	CrimeCategory string  `json:"crimeCategory"`
	Neighborhood  string  `json:"neighborhood"`
	Street        string  `json:"street"`
	City          string  `json:"city"`
	Zip           string  `json:"zip"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Date          string  `json:"date"`
	Time          string  `json:"time"`
	Source        string  `json:"source"`
	Distance      float64 `json:"distance"`
}

type CrimeResponse struct {
	Count  int     `json:"count"`
	Crimes []Crime `json:"crimes"`
}

type CrimeDumpResponse struct {
	Count  int         `json:"count"`
	Crimes []CrimeDump `json:"crimes"`
}

type OrderedPair struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

type CrimeStats struct {
	TotalCrimes   int           `json:"total_crimes"`
	CrimesByType  []OrderedPair `json:"crimes_by_type"`
	CrimesByDate  []OrderedPair `json:"crimes_by_date"`
	CrimesByHour  []OrderedPair `json:"crimes_by_hour"`
	MostDangerous []string      `json:"most_dangerous_areas"`
	SafestAreas   []string      `json:"safest_areas"`
}

type HeatMapPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Intensity int     `json:"intensity"`
	Radius    float64 `json:"radius"`
}

type CrimesByArea struct {
	Area   string  `json:"area"`
	Count  int     `json:"count"`
	Crimes []Crime `json:"crimes"`
}

func getCurrentYear() string {
	return strconv.Itoa(time.Now().Year())
}

func validateYears(year string) []string {
	if year == "" {
		return []string{getCurrentYear()}
	}

	years := strings.Split(year, ",")
	validYears := []string{}

	for i := range years {
		log.Printf("current year %s\n", years[i])
		if len(years[i]) == 4 {
			if yearInt, err := strconv.Atoi(years[i]); err == nil && yearInt >= 2018 && yearInt <= time.Now().Year() {
				validYears = append(validYears, years[i])
			}
		}
	}
	return validYears
}

func validateYear(year string) string {
	if year == "" {
		return getCurrentYear()
	}

	if len(year) == 4 {
		if yearInt, err := strconv.Atoi(year); err == nil && yearInt >= 2018 && yearInt <= time.Now().Year() {
			return year
		}
	}

	return getCurrentYear()
}

func (h *Handler) GetCrimes(c *gin.Context) {
	year := validateYears(c.Query("year"))
	crimeType := c.Query("type")
	date := c.Query("date")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	city := c.Query("city")
	neighborhood := c.Query("neighborhood")
	limitStr := c.Query("limit")

	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
	}

	var crimes []Crime

	crimes, err := h.getCrimesWithFilters(year, crimeType, date, startDate, endDate, city, neighborhood, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve crimes",
		})
		return
	}

	response := CrimeResponse{
		Crimes: crimes,
		Count:  len(crimes),
	}

	c.JSON(http.StatusOK, response)
}

// Gets crimes with many details
func (h *Handler) GetDetailedCrime(c *gin.Context) {
	years := validateYears(c.Query("year"))
	crimeTypes := strings.Split(c.Query("crimeType"), ",")
	cities := strings.Split(c.Query("city"), ",")
	neighborhoods := strings.Split(c.Query("neighborhood"), ",")
	sources := strings.Split(c.Query("source"), ",")

	date := c.Query("date")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	caseNumber := c.Query("caseNumber")
	street := c.Query("street")
	zipCode := c.Query("zipCode")

	limitStr := c.Query("limit")
	limit := 1000
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	if len(crimeTypes) == 1 && crimeTypes[0] == "" {
		crimeTypes = []string{}
	}
	if len(cities) == 1 && cities[0] == "" {
		cities = []string{}
	}
	if len(neighborhoods) == 1 && neighborhoods[0] == "" {
		neighborhoods = []string{}
	}
	if len(sources) == 1 && sources[0] == "" {
		sources = []string{}
	}

	if len(years) == 0 {
		years = []string{"2025"}
	}

	crimes, err := h.getDetailedCrimesWithAdvancedFilters(
		years, crimeTypes, cities, neighborhoods, sources,
		date, startDate, endDate, caseNumber, street, zipCode, limit,
	)
	if err != nil {
		http.Error(c.Writer, "Failed to fetch crimes", http.StatusInternalServerError)
		return
	}

	response := CrimeDumpResponse{Crimes: crimes, Count: len(crimes)}
	c.JSON(http.StatusOK, response)
}

func (h *Handler) getCrimesWithFilters(year []string, crimeType, date, startDate, endDate, city, neighborhood string, limit int) ([]Crime, error) {
	var query string
	if len(year) == 1 {
		query = `
		SELECT 
			COALESCE(a.street_address, 'Unknown Address') as address,
			COALESCE(n.neighborhood_name, 'Unknown neighborhood') as neighborhood,
			ci.incident_date::text,
			COALESCE(ci.incident_time, '00:00') as incident_time,
			COALESCE(l.latitude, 47.2529) as latitude,
			COALESCE(l.longitude, -122.4443) as longitude,
			COALESCE(cc.category_name, 'Other') as crime_type
		FROM crime_incidents_` + year[0] + ` ci
		JOIN addresses a ON ci.address_id = a.address_id
		JOIN neighborhoods n ON a.neighborhood_id = n.neighborhood_id
		JOIN cities c ON a.city_id = c.city_id
		LEFT JOIN locations l ON a.location_id = l.location_id
		LEFT JOIN crime_categories cc ON ci.crime_category_id = cc.crime_category_id
		WHERE 1=1
		`
	} else {
		sb := strings.Builder{}
		p1 := `
		SELECT 
			COALESCE(a.street_address, 'Unknown Address') as address,
			COALESCE(n.neighborhood_name, 'Unknown neighborhood') as neighborhood,
			ci.incident_date::text,
			COALESCE(ci.incident_time, '00:00') as incident_time,
			COALESCE(l.latitude, 47.2529) as latitude,
			COALESCE(l.longitude, -122.4443) as longitude,
			COALESCE(cc.category_name, 'Other') as crime_type
		FROM (
		`
		_, err := sb.WriteString(p1)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		var i int
		for i = 0; i < len(year)-1; i++ {
			_, err = sb.WriteString("SELECT incident_date, incident_time, address_id, crime_category_id FROM crime_incidents_" + year[i] + "\nUNION ALL\n")
			if err != nil {
				log.Println(err)
				return nil, err
			}
		}
		_, err = sb.WriteString("SELECT incident_date, incident_time, address_id, crime_category_id FROM crime_incidents_" + year[i])
		if err != nil {
			log.Println(err)
			return nil, err
		}
		p2 := `
		) as ci
		JOIN addresses a ON ci.address_id = a.address_id
		JOIN neighborhoods n ON a.neighborhood_id = n.neighborhood_id
		JOIN cities c ON a.city_id = c.city_id
		LEFT JOIN locations l ON a.location_id = l.location_id
		LEFT JOIN crime_categories cc ON ci.crime_category_id = cc.crime_category_id
		WHERE 1=1
		`
		_, err = sb.WriteString(p2)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		query = sb.String()
	}

	var args []any
	argIndex := 1

	if crimeType != "" {
		query += fmt.Sprintf(" AND LOWER(cc.category_name) = LOWER($%d)", argIndex)
		args = append(args, crimeType)
		argIndex++
	}

	if date != "" {
		query += fmt.Sprintf(" AND ci.incident_date = $%d", argIndex)
		args = append(args, date)
		argIndex++
	}

	if startDate != "" {
		query += fmt.Sprintf(" AND ci.incident_date >= $%d", argIndex)
		args = append(args, startDate)
		argIndex++
	}

	if endDate != "" {
		query += fmt.Sprintf(" AND ci.incident_date <= $%d", argIndex)
		args = append(args, endDate)
		argIndex++
	}

	if neighborhood != "" {
		query += fmt.Sprintf(" AND n.neighborhood_name LIKE $%d", argIndex)
		args = append(args, "%"+neighborhood+"%")
		argIndex++
	}

	if city != "" {
		query += fmt.Sprintf(" AND c.city_name LIKE $%d", argIndex)
		args = append(args, "%"+city+"%")
		argIndex++
	}

	query += " ORDER BY ci.incident_date DESC, a.street_address DESC"
	query += fmt.Sprintf(" LIMIT $%d", argIndex)
	args = append(args, limit)

	rows, err := h.pool.Query(context.Background(), query, args...)
	if err != nil {
		log.Printf("Error querying crimes: %v", err)
		return h.getAllMockCrimes(), nil
	}
	defer rows.Close()

	var crimes []Crime
	for rows.Next() {
		var crime Crime
		err := rows.Scan(
			&crime.Address,
			&crime.Neighborhood,
			&crime.Date,
			&crime.Time,
			&crime.Latitude,
			&crime.Longitude,
			&crime.CrimeType,
		)
		if err != nil {
			log.Printf("Error scanning crime row: %v", err)
			continue
		}
		crimes = append(crimes, crime)
	}

	if len(crimes) == 0 {
		return h.getAllMockCrimes(), nil
	}

	return crimes, nil
}

func (h *Handler) getDetailedCrimesWithAdvancedFilters(
	years []string,
	crimeTypes, cities, neighborhoods, sources []string,
	date, startDate, endDate, caseNumber, street, zipCode string,
	limit int,
) ([]CrimeDump, error) {
	var query string

	if len(years) == 1 {
		query = `
		SELECT
			ci.case_num,
			cc.category_name,
			COALESCE(n.neighborhood_name, 'Unknown neighborhood') as neighborhood,
			COALESCE(a.street_address, 'Unknown Address') as address,
			c.city_name,
			a.postal_code,
			l.latitude as latitude,
			l.longitude as longitude,
			ci.incident_date::text,
			ci.incident_time::text,
			s.source_name
		FROM crime_incidents_` + years[0] + ` ci
		JOIN data_sources s on ci.source_id = s.source_id
		JOIN addresses a ON ci.address_id = a.address_id
		JOIN neighborhoods n ON a.neighborhood_id = n.neighborhood_id
		JOIN cities c ON a.city_id = c.city_id
		LEFT JOIN locations l ON a.location_id = l.location_id
		LEFT JOIN crime_categories cc ON ci.crime_category_id = cc.crime_category_id
		WHERE 1=1
		`
	} else {
		sb := strings.Builder{}
		p1 := `
		SELECT
			ci.case_num,
			cc.category_name,
			COALESCE(n.neighborhood_name, 'Unknown neighborhood') as neighborhood,
			COALESCE(a.street_address, 'Unknown Address') as address,
			c.city_name,
			a.postal_code,
			l.latitude as latitude,
    		l.longitude as longitude,
			ci.incident_date::text,
			s.source_name
		FROM (
		`
		_, err := sb.WriteString(p1)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		var i int
		for i = 0; i < len(years)-1; i++ {
			_, err = sb.WriteString("SELECT incident_date, incident_time, address_id, crime_category_id FROM crime_incidents_" + years[i] + "\nUNION ALL\n")
			if err != nil {
				log.Println(err)
				return nil, err
			}
		}
		_, err = sb.WriteString("SELECT incident_date, incident_time, address_id, crime_category_id FROM crime_incidents_" + years[i])
		if err != nil {
			log.Println(err)
			return nil, err
		}
		p2 := `
		) as ci
		JOIN data_sources s on ci.source_id = s.source_id
		JOIN addresses a ON ci.address_id = a.address_id
		JOIN neighborhoods n ON a.neighborhood_id = n.neighborhood_id
		JOIN cities c ON a.city_id = c.city_id
		LEFT JOIN locations l ON a.location_id = l.location_id
		LEFT JOIN crime_categories cc ON ci.crime_category_id = cc.crime_category_id
		WHERE 1=1
		`
		_, err = sb.WriteString(p2)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		query = sb.String()
	}

	var args []any
	argIndex := 1

	if len(crimeTypes) > 0 && crimeTypes[0] != "" {
		placeholders := make([]string, len(crimeTypes))
		for i, crimeType := range crimeTypes {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, strings.ToLower(crimeType))
			argIndex++
		}
		query += fmt.Sprintf(" AND LOWER(cc.category_name) ILIKE ANY(ARRAY[%s])", strings.Join(placeholders, ","))
	}

	if len(cities) > 0 && cities[0] != "" {
		placeholders := make([]string, len(cities))
		for i, city := range cities {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, "%"+city+"%")
			argIndex++
		}
		query += fmt.Sprintf(" AND c.city_name ILIKE ANY(ARRAY[%s])", strings.Join(placeholders, ","))
	}

	if len(neighborhoods) > 0 && neighborhoods[0] != "" {
		placeholders := make([]string, len(neighborhoods))
		for i, neighborhood := range neighborhoods {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, "%"+neighborhood+"%")
			argIndex++
		}
		query += fmt.Sprintf(" AND n.neighborhood_name ILIKE ANY(ARRAY[%s])", strings.Join(placeholders, ","))
	}

	if len(sources) > 0 && sources[0] != "" {
		placeholders := make([]string, len(sources))
		for i, source := range sources {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, source)
			argIndex++
		}
		query += fmt.Sprintf(" AND s.source_name IN (%s)", strings.Join(placeholders, ","))
	}

	if date != "" {
		query += fmt.Sprintf(" AND ci.incident_date = $%d", argIndex)
		args = append(args, date)
		argIndex++
	}
	if startDate != "" {
		query += fmt.Sprintf(" AND ci.incident_date >= $%d", argIndex)
		args = append(args, startDate)
		argIndex++
	}
	if endDate != "" {
		query += fmt.Sprintf(" AND ci.incident_date <= $%d", argIndex)
		args = append(args, endDate)
		argIndex++
	}

	if caseNumber != "" {
		query += fmt.Sprintf(" AND CAST(ci.case_num AS TEXT) ILIKE $%d", argIndex)
		args = append(args, "%"+caseNumber+"%")
		argIndex++
	}
	if street != "" {
		query += fmt.Sprintf(" AND a.street_address ILIKE $%d", argIndex)
		args = append(args, "%"+street+"%")
		argIndex++
	}
	if zipCode != "" {
		query += fmt.Sprintf(" AND a.postal_code = $%d", argIndex)
		args = append(args, zipCode)
		argIndex++
	}

	query += " ORDER BY ci.incident_date DESC, a.street_address DESC"
	query += fmt.Sprintf(" LIMIT $%d", argIndex)
	args = append(args, limit)

	rows, err := h.pool.Query(context.Background(), query, args...)
	if err != nil {
		log.Printf("Error querying crimes: %v", err)
		return nil, err
	}
	defer rows.Close()

	var crimes []CrimeDump
	for rows.Next() {
		var crime CrimeDump
		err := rows.Scan(
			&crime.Case,
			&crime.CrimeCategory,
			&crime.Neighborhood,
			&crime.Street,
			&crime.City,
			&crime.Zip,
			&crime.Latitude,
			&crime.Longitude,
			&crime.Date,
			&crime.Time,
			&crime.Source,
		)
		if err != nil {
			log.Printf("Error scanning crime row: %v", err)
			continue
		}
		crimes = append(crimes, crime)
	}
	return crimes, nil
}

// Haversine formula to calculate distance between two points in MILES
func haversineDistanceMiles(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 3959

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

// Geographic radius filtering - crimes within X miles of a point
func (h *Handler) GetCrimesInRadius(c *gin.Context) {
	year := validateYears(c.Query("year"))
	latStr := c.Query("lat")
	lonStr := c.Query("lng")
	radiusStr := c.Query("radius")
	crimeType := c.Query("type")
	limitStr := c.Query("limit")

	if latStr == "" || lonStr == "" || radiusStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "lat, lng, and radius parameters are required",
			"example": "/api/public/crimes/radius?lat=47.2529&lng=-122.4443&radius=1.5&year=2025",
		})
		return
	}

	centerLat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
		return
	}

	centerLon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
		return
	}

	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil || radius <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid radius"})
		return
	}

	limit := 500
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 2000 {
			limit = l
		}
	}

	crimes, err := h.getCrimesInRadius(year, centerLat, centerLon, radius, crimeType, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve crimes in radius",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"center": gin.H{
			"latitude":  centerLat,
			"longitude": centerLon,
		},
		"radius_miles": radius,
		"crimes":       crimes,
		"count":        len(crimes),
		"year":         year,
	})
}

func (h *Handler) getCrimesInRadius(year []string, centerLat, centerLon, radius float64, crimeType string, limit int) ([]CrimeWithDistance, error) {
	var query string
	if len(year) == 1 {
		query = `
		SELECT
			ci.case_num,
			cc.category_name,
			COALESCE(n.neighborhood_name, 'Unknown neighborhood') as neighborhood,
			COALESCE(a.street_address, 'Unknown Address') as address,
			c.city_name,
			a.postal_code,
			l.latitude as latitude,
			l.longitude as longitude,
			ci.incident_date::text,
			ci.incident_time::text
		FROM crime_incidents_` + year[0] + ` ci
		JOIN addresses a ON ci.address_id = a.address_id
		JOIN neighborhoods n ON a.neighborhood_id = n.neighborhood_id
		JOIN cities c ON a.city_id = c.city_id
		LEFT JOIN locations l ON a.location_id = l.location_id
		LEFT JOIN crime_categories cc ON ci.crime_category_id = cc.crime_category_id
		WHERE l.latitude IS NOT NULL AND l.longitude IS NOT NULL
	`
	} else {
		sb := strings.Builder{}
		p1 := `
		SELECT 
			ci.case_num,
			cc.category_name,
			COALESCE(n.neighborhood_name, 'Unknown neighborhood') as neighborhood,
			COALESCE(a.street_address, 'Unknown Address') as address,
			c.city_name,
			a.postal_code,
			l.latitude as latitude,
			l.longitude as longitude,
			ci.incident_date::text,
			ci.incident_time::text
		FROM (
		`
		_, err := sb.WriteString(p1)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		for i := range year {
			_, err = sb.WriteString(`crime_incidents_` + year[i] + `UNION ALL `)
			if err != nil {
				log.Println(err)
				return nil, err
			}
		}
		p2 := `
		) as ci
		JOIN addresses a ON ci.address_id = a.address_id
		JOIN neighborhoods n ON a.neighborhood_id = n.neighborhood_id
		JOIN cities c ON a.city_id = c.city_id
		LEFT JOIN locations l ON a.location_id = l.location_id
		LEFT JOIN crime_categories cc ON ci.crime_category_id = cc.crime_category_id
		WHERE l.latitude IS NOT NULL AND l.longitude IS NOT NULL
		`
		_, err = sb.WriteString(p2)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		query = sb.String()
	}

	var args []any
	argIndex := 1

	if crimeType != "" {
		query += fmt.Sprintf(" AND LOWER(cc.category_name) = LOWER($%d)", argIndex)
		args = append(args, crimeType)
		argIndex++
	}

	query += " ORDER BY ci.incident_date DESC"
	query += fmt.Sprintf(" LIMIT $%d", argIndex)
	args = append(args, limit*2)

	rows, err := h.pool.Query(context.Background(), query, args...)
	if err != nil {
		log.Printf("Error querying crimes for radius: %v", err)
		return []CrimeWithDistance{}, err
	}
	defer rows.Close()

	var crimesInRadius []CrimeWithDistance
	for rows.Next() {
		var crime CrimeWithDistance
		err := rows.Scan(
			&crime.Case,
			&crime.CrimeCategory,
			&crime.Neighborhood,
			&crime.Street,
			&crime.City,
			&crime.Zip,
			&crime.Latitude,
			&crime.Longitude,
			&crime.Date,
			&crime.Time,
		)
		if err != nil {
			log.Printf("Error scanning crime row: %v", err)
			continue
		}

		distance := haversineDistanceMiles(centerLat, centerLon, crime.Latitude, crime.Longitude)
		if distance <= radius {
			crime.Distance = distance
			crimesInRadius = append(crimesInRadius, crime)
			if len(crimesInRadius) >= limit {
				break
			}
		}
	}

	return crimesInRadius, nil
}

// Crime statistics endpoint
func (h *Handler) GetCrimeStats(c *gin.Context) {
	year := validateYear(c.Query("year"))

	stats, err := h.calculateCrimeStats(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to calculate crime statistics",
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *Handler) calculateCrimeStats(year string) (CrimeStats, error) {
	query := `
		WITH crime_type_stats AS (
			SELECT 
				'type' as stat_type,
				COALESCE(cc.category_name, 'Other') as key,
				COUNT(*) as count,
				ROW_NUMBER() OVER (ORDER BY COUNT(*) DESC) as rank
			FROM crime_incidents_` + year + ` ci
			LEFT JOIN crime_categories cc ON ci.crime_category_id = cc.crime_category_id
			GROUP BY cc.category_name
		),
		date_stats AS (
			SELECT 
				'date' as stat_type,
				ci.incident_date::text as key,
				COUNT(*) as count,
				ROW_NUMBER() OVER (ORDER BY COUNT(*) DESC) as rank
			FROM crime_incidents_` + year + ` ci
			GROUP BY ci.incident_date
		),
		hour_stats AS (
			SELECT 
				'hour' as stat_type,
				LPAD(EXTRACT(HOUR FROM COALESCE(ci.incident_time, '00:00:00')::time)::text, 2, '0') as key,
				COUNT(*) as count,
				EXTRACT(HOUR FROM COALESCE(ci.incident_time, '00:00:00')::time) as sort_order
			FROM crime_incidents_` + year + ` ci
			GROUP BY EXTRACT(HOUR FROM COALESCE(ci.incident_time, '00:00:00')::time)
		),
		area_stats AS (
			SELECT 
				'area' as stat_type,
				COALESCE(n.neighborhood_name, 'Unknown Area') as key,
				COUNT(*) as count,
				ROW_NUMBER() OVER (ORDER BY COUNT(*) DESC) as rank
			FROM crime_incidents_` + year + ` ci
			JOIN addresses a ON ci.address_id = a.address_id
			LEFT JOIN neighborhoods n ON a.neighborhood_id = n.neighborhood_id
			GROUP BY n.neighborhood_name
			HAVING COUNT(*) > 0
		),
		total_crimes AS (
			SELECT COUNT(*) as total FROM crime_incidents_` + year + `
		)
		SELECT stat_type, key, count, 
			   COALESCE(rank, sort_order) as order_val,
			   (SELECT total FROM total_crimes) as total_count
		FROM (
			SELECT stat_type, key, count, rank, NULL::numeric as sort_order FROM crime_type_stats
			UNION ALL
			SELECT stat_type, key, count, rank, NULL::numeric as sort_order FROM date_stats WHERE rank <= 30
			UNION ALL
			SELECT stat_type, key, count, NULL::bigint as rank, sort_order FROM hour_stats
			UNION ALL
			SELECT stat_type, key, count, rank, NULL::numeric as sort_order FROM area_stats
		) combined
		ORDER BY stat_type, order_val
	`

	var crimesByType []OrderedPair
	var crimesByDate []OrderedPair
	var crimesByHour []OrderedPair
	var mostDangerous, safestAreas []string
	var totalCrimes int
	var areas []struct {
		name  string
		count int
	}

	rows, err := h.pool.Query(context.Background(), query)
	if err != nil {
		return CrimeStats{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var statType, key string
		var count int
		var orderVal *float64
		var totalCount int

		if err := rows.Scan(&statType, &key, &count, &orderVal, &totalCount); err != nil {
			continue
		}

		totalCrimes = totalCount

		switch statType {
		case "type":
			crimesByType = append(crimesByType, OrderedPair{Key: key, Value: count})
		case "date":
			crimesByDate = append(crimesByDate, OrderedPair{Key: key, Value: count})
		case "hour":
			crimesByHour = append(crimesByHour, OrderedPair{Key: key, Value: count})
		case "area":
			areas = append(areas, struct {
				name  string
				count int
			}{key, count})
		}
	}

	if len(areas) > 0 {
		maxAreas := min(len(areas), 5)
		for i := range maxAreas {
			mostDangerous = append(mostDangerous, areas[i].name)
		}
		if len(areas) > 5 {
			start := len(areas) - 5
			for i := start; i < len(areas); i++ {
				safestAreas = append([]string{areas[i].name}, safestAreas...)
			}
		}
	}

	if len(mostDangerous) == 0 {
		mostDangerous = []string{"No data available"}
	}
	if len(safestAreas) == 0 {
		safestAreas = []string{"No data available"}
	}

	stats := CrimeStats{
		TotalCrimes:   totalCrimes,
		CrimesByType:  crimesByType,
		CrimesByDate:  crimesByDate,
		CrimesByHour:  crimesByHour,
		MostDangerous: mostDangerous,
		SafestAreas:   safestAreas,
	}

	return stats, nil
}

// Heat map data endpoint - returns points with intensity
func (h *Handler) GetHeatMapData(c *gin.Context) {
	year := validateYear(c.Query("year"))
	crimeType := c.Query("type")
	gridSizeStr := c.Query("grid_size")

	gridSize := 0.005
	if gridSizeStr != "" {
		if gs, err := strconv.ParseFloat(gridSizeStr, 64); err == nil && gs > 0 && gs <= 0.1 {
			gridSize = gs
		}
	}

	heatPoints, err := h.generateHeatMapData(year, crimeType, gridSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate heat map data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"heat_points":     heatPoints,
		"grid_size_deg":   gridSize,
		"grid_size_miles": gridSize * MILE_APPROX,
		"total_points":    len(heatPoints),
		"year":            year,
		"crime_type":      crimeType,
	})
}

func (h *Handler) generateHeatMapData(year, crimeType string, gridSize float64) ([]HeatMapPoint, error) {
	query := `
		SELECT 
			COALESCE(l.latitude, 47.2529) as latitude,
			COALESCE(l.longitude, -122.4443) as longitude
		FROM crime_incidents_` + year + ` ci
		JOIN addresses a ON ci.address_id = a.address_id
		LEFT JOIN locations l ON a.location_id = l.location_id
		LEFT JOIN crime_categories cc ON ci.crime_category_id = cc.crime_category_id
		WHERE l.latitude IS NOT NULL AND l.longitude IS NOT NULL
	`

	var args []any
	if crimeType != "" {
		query += " AND LOWER(cc.category_name) = LOWER($1)"
		args = append(args, crimeType)
	}

	rows, err := h.pool.Query(context.Background(), query, args...)
	if err != nil {
		log.Printf("Error querying crimes for heat map: %v", err)
		return []HeatMapPoint{}, err
	}
	defer rows.Close()

	intensityMap := make(map[string]int)

	for rows.Next() {
		var lat, lon float64
		if err := rows.Scan(&lat, &lon); err != nil {
			continue
		}

		gridLat := math.Round(lat/gridSize) * gridSize
		gridLon := math.Round(lon/gridSize) * gridSize
		key := fmt.Sprintf("%.6f,%.6f", gridLat, gridLon)
		intensityMap[key]++
	}

	var heatPoints []HeatMapPoint
	for key, intensity := range intensityMap {
		coords := strings.Split(key, ",")
		lat, _ := strconv.ParseFloat(coords[0], 64)
		lon, _ := strconv.ParseFloat(coords[1], 64)

		heatPoints = append(heatPoints, HeatMapPoint{
			Latitude:  lat,
			Longitude: lon,
			Intensity: intensity,
			Radius:    math.Min(float64(intensity)*50+100, 500),
		})
	}

	return heatPoints, nil
}

// Crime trends endpoint - crimes over time
func (h *Handler) GetCrimeTrends(c *gin.Context) {
	year := validateYear(c.Query("year"))
	crimeType := c.Query("type")
	period := c.DefaultQuery("period", "daily") // daily, weekly, monthly

	trends, err := h.calculateCrimeTrends(year, crimeType, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to calculate crime trends",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"trends":     trends,
		"year":       year,
		"crime_type": crimeType,
		"period":     period,
	})
}

func (h *Handler) calculateCrimeTrends(year, crimeType, period string) (map[string]map[string]int, error) {
	var dateFormat string
	switch period {
	case "weekly":
		dateFormat = "YYYY-\"Week\"-WW"
	case "monthly":
		dateFormat = "YYYY-MM"
	default:
		dateFormat = "YYYY-MM-DD"
	}

	query := `
		SELECT 
			TO_CHAR(ci.incident_date, '` + dateFormat + `') as time_period,
			COALESCE(cc.category_name, 'Other') as crime_type,
			COUNT(*) as count
		FROM crime_incidents_` + year + ` ci
		LEFT JOIN crime_categories cc ON ci.crime_category_id = cc.crime_category_id
		WHERE 1=1
	`

	var args []interface{}
	if crimeType != "" {
		query += " AND LOWER(cc.category_name) = LOWER($1)"
		args = append(args, crimeType)
	}

	query += " GROUP BY time_period, cc.category_name ORDER BY time_period, cc.category_name"

	rows, err := h.pool.Query(context.Background(), query, args...)
	if err != nil {
		log.Printf("Error querying crime trends: %v", err)
		return make(map[string]map[string]int), err
	}
	defer rows.Close()

	trends := make(map[string]map[string]int)

	for rows.Next() {
		var timePeriod, crimeTypeResult string
		var count int

		if err := rows.Scan(&timePeriod, &crimeTypeResult, &count); err != nil {
			continue
		}

		if trends[timePeriod] == nil {
			trends[timePeriod] = make(map[string]int)
		}
		trends[timePeriod][crimeTypeResult] = count
		trends[timePeriod]["total"] += count
	}

	return trends, nil
}

func (h *Handler) GetDangerousAreas(c *gin.Context) {
	year := validateYear(c.Query("year"))
	includeDetails := c.DefaultQuery("include_details", "false") == "true"
	limitStr := c.Query("limit")

	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	areaStats, err := h.getAreaStatistics(year, includeDetails, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve area statistics",
		})
		return
	}

	sort.Slice(areaStats, func(i, j int) bool {
		return areaStats[i].Count > areaStats[j].Count
	})

	mostDangerous := "No data available"
	safest := "No data available"
	if len(areaStats) > 0 {
		mostDangerous = areaStats[0].Area
		for i := len(areaStats) - 1; i >= 0; i-- {
			if areaStats[i].Count > 0 {
				safest = areaStats[i].Area
				break
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"areas":           areaStats,
		"most_dangerous":  mostDangerous,
		"safest":          safest,
		"total_areas":     len(areaStats),
		"year":            year,
		"include_details": includeDetails,
	})
}

func (h *Handler) getAreaStatistics(year string, includeDetails bool, limit int) ([]CrimesByArea, error) {
	statsQuery := `
		SELECT 
			COALESCE(n.neighborhood_name, 'Unknown Area') as area_name,
			COUNT(*) as crime_count
		FROM crime_incidents_` + year + ` ci
		JOIN addresses a ON ci.address_id = a.address_id
		LEFT JOIN neighborhoods n ON a.neighborhood_id = n.neighborhood_id
		GROUP BY n.neighborhood_name
		HAVING COUNT(*) > 0
		ORDER BY crime_count DESC
		LIMIT $1
	`

	rows, err := h.pool.Query(context.Background(), statsQuery, limit)
	if err != nil {
		log.Printf("Error querying area statistics: %v", err)
		return []CrimesByArea{}, err
	}
	defer rows.Close()

	var areaStats []CrimesByArea

	for rows.Next() {
		var area CrimesByArea

		err := rows.Scan(&area.Area, &area.Count)
		if err != nil {
			log.Printf("Error scanning area stats: %v", err)
			continue
		}

		if includeDetails {
			crimes, err := h.getCrimesForNeighborhood(area.Area, year, 10)
			if err != nil {
				log.Printf("Error getting crimes for area %s: %v", area.Area, err)
				continue
			}
			area.Crimes = crimes
		}

		areaStats = append(areaStats, area)
	}

	return areaStats, nil
}

// Gets crimes for certain neighborhood
func (h *Handler) getCrimesForNeighborhood(neighborhood, year string, limit int) ([]Crime, error) {
	query := `
		SELECT 
			ci.incident_date::text, 
			COALESCE(ci.incident_time, '00:00') as incident_time,
			COALESCE(l.latitude, 47.2529) as latitude,
			COALESCE(l.longitude, -122.4443) as longitude,
			COALESCE(a.street_address, 'Unknown Address') as address,
			COALESCE(cc.category_name, 'Other') as crime_type
		FROM crime_incidents_` + year + ` ci
		JOIN addresses a ON ci.address_id = a.address_id
		LEFT JOIN neighborhoods n ON a.neighborhood_id = n.neighborhood_id
		LEFT JOIN locations l ON a.location_id = l.location_id
		LEFT JOIN crime_categories cc ON ci.crime_category_id = cc.crime_category_id
		WHERE COALESCE(n.neighborhood_name, 'Unknown Area') = $1
		ORDER BY ci.incident_date DESC
		LIMIT $2
	`

	rows, err := h.pool.Query(context.Background(), query, neighborhood, limit)
	if err != nil {
		log.Printf("Error querying crimes for area %s: %v", neighborhood, err)
		return []Crime{}, err
	}
	defer rows.Close()

	var crimes []Crime
	for rows.Next() {
		var crime Crime
		err := rows.Scan(
			&crime.Date,
			&crime.Time,
			&crime.Latitude,
			&crime.Longitude,
			&crime.Address,
			&crime.CrimeType,
		)
		if err != nil {
			log.Printf("Error scanning crime: %v", err)
			continue
		}
		crimes = append(crimes, crime)
	}

	return crimes, nil
}

// Backup mock data function
func (h *Handler) getAllMockCrimes() []Crime {
	return []Crime{
		{Address: "123 Main St, Tacoma, WA", Latitude: 47.2529, Longitude: -122.4443, CrimeType: "Theft", Date: "2024-08-10", Time: "14:30"},
		{Address: "456 Pacific Ave, Tacoma, WA", Latitude: 47.2563, Longitude: -122.4590, CrimeType: "Burglary", Date: "2024-08-11", Time: "22:15"},
		{Address: "789 Commerce St, Tacoma, WA", Latitude: 47.2528, Longitude: -122.4594, CrimeType: "Assault", Date: "2024-08-12", Time: "18:45"},
		{Address: "321 6th Ave, Tacoma, WA", Latitude: 47.2545, Longitude: -122.4580, CrimeType: "Vandalism", Date: "2024-08-09", Time: "03:20"},
		{Address: "654 Market St, Tacoma, WA", Latitude: 47.2510, Longitude: -122.4520, CrimeType: "Drug Offense", Date: "2024-08-08", Time: "16:10"},
	}
}
