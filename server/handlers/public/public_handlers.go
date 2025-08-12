package public

import (
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
	// Mock crime data for testing
	allCrimes := h.getAllMockCrimes()  //  helper function

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

		// Filter by exact date
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

// 1. Geographic radius filtering - crimes within X km of a point
func (h *Handler) GetCrimesInRadius(c *gin.Context) {
	// Get all crimes (same mock data as before)
	allCrimes := h.getAllMockCrimes()

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

// 2. Crime statistics endpoint
func (h *Handler) GetCrimeStats(c *gin.Context) {
	allCrimes := h.getAllMockCrimes()
	
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

// 3. Heat map data endpoint - returns points with intensity
func (h *Handler) GetHeatMapData(c *gin.Context) {
	allCrimes := h.getAllMockCrimes()
	
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

// 4. Crime trends endpoint - crimes over time
func (h *Handler) GetCrimeTrends(c *gin.Context) {
	allCrimes := h.getAllMockCrimes()
	
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
			"start": "2024-08-08",
			"end": "2024-08-12",
		},
	})
}

// 5. Dangerous areas endpoint
func (h *Handler) GetDangerousAreas(c *gin.Context) {
	allCrimes := h.getAllMockCrimes()
	
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

// Helper function to get all mock crimes (extract the data to reuse)
func (h *Handler) getAllMockCrimes() []Crime {
	return []Crime{
		{IncidentID: 1, Address: "123 Main St, Tacoma, WA", Latitude: 47.2529, Longitude: -122.4443, CrimeType: "Theft", Date: "2024-08-10", Time: "14:30"},
		{IncidentID: 2, Address: "456 Pacific Ave, Tacoma, WA", Latitude: 47.2563, Longitude: -122.4590, CrimeType: "Burglary", Date: "2024-08-11", Time: "22:15"},
		{IncidentID: 3, Address: "789 Commerce St, Tacoma, WA", Latitude: 47.2528, Longitude: -122.4594, CrimeType: "Assault", Date: "2024-08-12", Time: "18:45"},
		{IncidentID: 4, Address: "321 6th Ave, Tacoma, WA", Latitude: 47.2545, Longitude: -122.4580, CrimeType: "Vandalism", Date: "2024-08-09", Time: "03:20"},
		{IncidentID: 5, Address: "654 Market St, Tacoma, WA", Latitude: 47.2510, Longitude: -122.4520, CrimeType: "Drug Offense", Date: "2024-08-08", Time: "16:10"},
		{IncidentID: 6, Address: "987 Broadway, Tacoma, WA", Latitude: 47.2520, Longitude: -122.4470, CrimeType: "Theft", Date: "2024-08-12", Time: "09:15"},
		{IncidentID: 7, Address: "159 Union Ave, Tacoma, WA", Latitude: 47.2540, Longitude: -122.4600, CrimeType: "Burglary", Date: "2024-08-10", Time: "23:45"},
		{IncidentID: 8, Address: "753 6th Ave, Tacoma, WA", Latitude: 47.2555, Longitude: -122.4575, CrimeType: "Assault", Date: "2024-08-11", Time: "20:30"},
		{IncidentID: 9, Address: "147 Pine St, Tacoma, WA", Latitude: 47.2535, Longitude: -122.4495, CrimeType: "Theft", Date: "2024-08-09", Time: "12:45"},
		{IncidentID: 10, Address: "258 MLK Jr Way, Tacoma, WA", Latitude: 47.2548, Longitude: -122.4512, CrimeType: "Drug Offense", Date: "2024-08-10", Time: "01:15"},
	}
}