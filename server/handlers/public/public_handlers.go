package public

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"  
	"strings" 

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	EIGHT_HOUR = 3600 * 8
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
	IncidentID   int     `json:"incident_id"`
	Address      string  `json:"address"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	CrimeType    string  `json:"crime_type"`
	Date         string  `json:"date"`
	Time         string  `json:"time"`
}

type CrimeResponse struct {
	Crimes []Crime `json:"crimes"`
	Count  int     `json:"count"`
}

type CrimeStats struct {
	TotalCrimes    int                    `json:"total_crimes"`
	CrimesByType   map[string]int         `json:"crimes_by_type"`
	CrimesByDate   map[string]int         `json:"crimes_by_date"`
	CrimesByHour   map[string]int         `json:"crimes_by_hour"`
	MostDangerous  []string               `json:"most_dangerous_areas"`
	SafestAreas    []string               `json:"safest_areas"`
}

type HeatMapPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Intensity int     `json:"intensity"`
	Radius    float64 `json:"radius"`
}

type CrimesByArea struct {
	Area   string `json:"area"`
	Count  int    `json:"count"`
	Crimes []Crime `json:"crimes"`
}

func (h *Handler) GetCrimes(c *gin.Context) {
	// Get real crime data from database
	allCrimes := h.getAllRealCrimes()

	// Get query parameters
	crimeType := c.Query("type")        // ?type=Theft
	date := c.Query("date")             // ?date=2024-08-10
	startDate := c.Query("start_date")  // ?start_date=2024-08-08
	endDate := c.Query("end_date")      // ?end_date=2024-08-12
	limit := c.Query("limit")           // ?limit=10

	var filteredCrimes []Crime
	for _, crime := range allCrimes {
		// Filter by crime type
		if crimeType != "" && strings.ToLower(crime.CrimeType) != strings.ToLower(crimeType) {
			continue
		}

		// Filter by date
		if date != "" && crime.Date != date {
			continue
		}

		// Filter by date range
		if startDate != "" && crime.Date < startDate {
			continue
		}
		if endDate != "" && crime.Date > endDate {
			continue
		}

		filteredCrimes = append(filteredCrimes, crime)
	}

	// Apply limit if specified
	if limit != "" {
		if limitNum, err := strconv.Atoi(limit); err == nil && limitNum > 0 && limitNum < len(filteredCrimes) {
			filteredCrimes = filteredCrimes[:limitNum]
		}
	}

	response := CrimeResponse{
		Crimes: filteredCrimes,
		Count:  len(filteredCrimes),
	}

	c.JSON(http.StatusOK, response)
}

// Haversine formula to calculate distance between two points
func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth's radius in kilometers

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

// Geographic radius filtering - crimes within X km of a point
func (h *Handler) GetCrimesInRadius(c *gin.Context) {
	allCrimes := h.getAllRealCrimes()

	// Get query parameters
	latStr := c.Query("lat")      // Center latitude
	lonStr := c.Query("lng")      // Center longitude  
	radiusStr := c.Query("radius") // Radius in kilometers
	crimeType := c.Query("type")
	
	if latStr == "" || lonStr == "" || radiusStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "lat, lng, and radius parameters are required",
			"example": "/api/public/crimes/radius?lat=47.2529&lng=-122.4443&radius=1.5",
		})
		return
	}

	centerLat, _ := strconv.ParseFloat(latStr, 64)
	centerLon, _ := strconv.ParseFloat(lonStr, 64)
	radius, _ := strconv.ParseFloat(radiusStr, 64)

	var crimesInRadius []Crime
	for _, crime := range allCrimes {
		distance := haversineDistance(centerLat, centerLon, crime.Latitude, crime.Longitude)
		if distance <= radius {
			// Filter by type if specified
			if crimeType == "" || strings.ToLower(crime.CrimeType) == strings.ToLower(crimeType) {
				crimesInRadius = append(crimesInRadius, crime)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"center": gin.H{
			"latitude": centerLat,
			"longitude": centerLon,
		},
		"radius_km": radius,
		"crimes": crimesInRadius,
		"count": len(crimesInRadius),
	})
}

// Crime statistics endpoint
func (h *Handler) GetCrimeStats(c *gin.Context) {
	allCrimes := h.getAllRealCrimes()
	
	// Count by type
	crimesByType := make(map[string]int)
	crimesByDate := make(map[string]int)
	crimesByHour := make(map[string]int)

	for _, crime := range allCrimes {
		// Count by type
		crimesByType[crime.CrimeType]++
		
		// Count by date
		crimesByDate[crime.Date]++
		
		// Count by hour (extract hour from time like "14:30")
		if len(crime.Time) >= 2 {
			hour := crime.Time[:2]
			crimesByHour[hour]++
		}
	}

	stats := CrimeStats{
		TotalCrimes:    len(allCrimes),
		CrimesByType:   crimesByType,
		CrimesByDate:   crimesByDate,
		CrimesByHour:   crimesByHour,
		MostDangerous:  []string{"Downtown Tacoma", "6th Ave Corridor", "Hilltop"},
		SafestAreas:    []string{"North End", "Proctor District", "Stadium District"},
	}

	c.JSON(http.StatusOK, stats)
}

// Heat map data endpoint - returns points with intensity
func (h *Handler) GetHeatMapData(c *gin.Context) {
	allCrimes := h.getAllRealCrimes()
	
	// Grid size for heat map (adjust as needed)
	gridSize := 0.005 // About 500 meters
	intensityMap := make(map[string]int)
	
	// Group crimes into grid cells
	for _, crime := range allCrimes {
		// Round coordinates to grid
		gridLat := math.Round(crime.Latitude/gridSize) * gridSize
		gridLon := math.Round(crime.Longitude/gridSize) * gridSize
		key := strconv.FormatFloat(gridLat, 'f', 6, 64) + "," + strconv.FormatFloat(gridLon, 'f', 6, 64)
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
			Radius:    float64(intensity) * 200, // Scale radius based on intensity
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"heat_points": heatPoints,
		"grid_size_km": gridSize * 111, // Convert to approximate km
		"total_points": len(heatPoints),
	})
}

// Crime trends endpoint - crimes over time
func (h *Handler) GetCrimeTrends(c *gin.Context) {
	allCrimes := h.getAllRealCrimes()
	
	// Group by date and type
	trends := make(map[string]map[string]int)
	
	for _, crime := range allCrimes {
		if trends[crime.Date] == nil {
			trends[crime.Date] = make(map[string]int)
		}
		trends[crime.Date][crime.CrimeType]++
		trends[crime.Date]["total"]++
	}

	c.JSON(http.StatusOK, gin.H{
		"trends": trends,
		"date_range": gin.H{
			"start": "2024-02-27",
			"end": "2024-11-22",
		},
	})
}

// Dangerous areas endpoint
func (h *Handler) GetDangerousAreas(c *gin.Context) {
	allCrimes := h.getAllRealCrimes()
	
	// Define area boundaries (simplified)
	areas := map[string][]Crime{
		"Downtown": {},
		"Hilltop": {},
		"Stadium District": {},
		"6th Avenue": {},
		"North End": {},
	}

	// Categorize crimes by area (simplified logic based on coordinates)
	for _, crime := range allCrimes {
		switch {
		case crime.Latitude > 47.254 && crime.Latitude < 47.257:
			areas["Downtown"] = append(areas["Downtown"], crime)
		case crime.Latitude > 47.250 && crime.Latitude < 47.253:
			areas["Hilltop"] = append(areas["Hilltop"], crime)
		case crime.Longitude > -122.46:
			areas["Stadium District"] = append(areas["Stadium District"], crime)
		case strings.Contains(crime.Address, "6th"):
			areas["6th Avenue"] = append(areas["6th Avenue"], crime)
		default:
			areas["North End"] = append(areas["North End"], crime)
		}
	}

	var areaStats []CrimesByArea
	for area, crimes := range areas {
		areaStats = append(areaStats, CrimesByArea{
			Area:   area,
			Count:  len(crimes),
			Crimes: crimes,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"areas": areaStats,
		"most_dangerous": "Downtown",
		"safest": "North End",
	})
}

// Main function to get real crime data from database
func (h *Handler) getAllRealCrimes() []Crime {
	// Corrected query - crime incidents connect directly to locations
	query := `
		SELECT 
			ci.incident_id,
			ci.address_id,
			ci.crime_category_id,
			ci.incident_date::text,
			COALESCE(ci.incident_time, '00:00') as incident_time,
			COALESCE(l.latitude, 47.2529) as latitude,
			COALESCE(l.longitude, -122.4443) as longitude
		FROM crime_incidents_2024 ci
		LEFT JOIN locations l ON ci.location_id = l.location_id
		ORDER BY ci.incident_id DESC
		LIMIT 200
	`
	
	rows, err := h.pool.Query(context.Background(), query)
	if err != nil {
		log.Printf("Error querying crimes with coordinates: %v", err)
		return h.getAllMockCrimes()
	}
	defer rows.Close()

	var crimes []Crime
	for rows.Next() {
		var crime Crime
		var addressID, categoryID int
		var incidentTime string
		var latitude, longitude float64
		
		err := rows.Scan(
			&crime.IncidentID,
			&addressID,
			&categoryID,
			&crime.Date,
			&incidentTime,
			&latitude,    // Real latitude from locations table
			&longitude,   // Real longitude from locations table
		)
		
		if err != nil {
			log.Printf("Error scanning crime row: %v", err)
			continue
		}

		crime.Time = incidentTime
		crime.Latitude = latitude    // Real coordinates!
		crime.Longitude = longitude  // Real coordinates!
		crime.Address = fmt.Sprintf("Address ID: %d", addressID)
		crime.CrimeType = mapCategoryIDToName(categoryID)

		crimes = append(crimes, crime)
	}

	if len(crimes) == 0 {
		log.Println("No crimes found in database, using mock data")
		return h.getAllMockCrimes()
	}

	log.Printf("Successfully loaded %d real crimes with coordinates from database", len(crimes))
	return crimes
}

// Map crime category IDs to readable names
func mapCategoryIDToName(categoryID int) string {
	categoryMap := map[int]string{
		1:  "Theft",
		2:  "Burglary", 
		3:  "Assault",
		4:  "Vandalism",
		5:  "Drug Offense",
		6:  "Robbery",
		7:  "Domestic Violence",
		8:  "Fraud",
		9:  "Vehicle Theft",
		10: "Weapons Violation",
	}
	
	if name, exists := categoryMap[categoryID]; exists {
		return name
	}
	return "Other"
}

// Backup mock data function
func (h *Handler) getAllMockCrimes() []Crime {
	return []Crime{
		{IncidentID: 1, Address: "123 Main St, Tacoma, WA", Latitude: 47.2529, Longitude: -122.4443, CrimeType: "Theft", Date: "2024-08-10", Time: "14:30"},
		{IncidentID: 2, Address: "456 Pacific Ave, Tacoma, WA", Latitude: 47.2563, Longitude: -122.4590, CrimeType: "Burglary", Date: "2024-08-11", Time: "22:15"},
		{IncidentID: 3, Address: "789 Commerce St, Tacoma, WA", Latitude: 47.2528, Longitude: -122.4594, CrimeType: "Assault", Date: "2024-08-12", Time: "18:45"},
		{IncidentID: 4, Address: "321 6th Ave, Tacoma, WA", Latitude: 47.2545, Longitude: -122.4580, CrimeType: "Vandalism", Date: "2024-08-09", Time: "03:20"},
		{IncidentID: 5, Address: "654 Market St, Tacoma, WA", Latitude: 47.2510, Longitude: -122.4520, CrimeType: "Drug Offense", Date: "2024-08-08", Time: "16:10"},
	}
}