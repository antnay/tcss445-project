package router

import (
	adminHandlers "server/handlers/admin"
	//authHandlers "server/handlers/auth"
	privateHandlers "server/handlers/private"
	publicHandlers "server/handlers/public"
)

func (r *Router) registerAdminRoutes() {
	adminHandler := adminHandlers.NewHandler(r.pool)
	api := r.Router.Group("/api")
	admin := api.Group("/admin")

	admin.GET("/ping", adminHandler.AdminPing)
}

func (r *Router) registerPrivateRoutes() {
	privateHandler := privateHandlers.NewHandler(r.pool)
	api := r.Router.Group("/api")
	api.GET("/ping", privateHandler.Ping)
}

func (r *Router) registerPublicRoutes() {
	publicHandler := publicHandlers.NewHandler(r.pool)
	//authHandler := authHandlers.NewHandler(r.pool, r.tokenFactory)
	api := r.Router.Group("/api/public")

	api.GET("/ping", publicHandler.Ping) // api/public/ping
	// api.POST("/register", authHandler.Register)
	// api.POST("/login", authHandler.Login)
	api.GET("/crimes", publicHandler.GetCrimes)
	api.GET("/crimes/details", publicHandler.GetDetailedCrime)
	api.GET("/crimes/radius", publicHandler.GetCrimesInRadius) // Geographic filtering
	api.GET("/crimes/stats", publicHandler.GetCrimeStats)      // Statistics
	api.GET("/crimes/heatmap", publicHandler.GetHeatMapData)   // Heat map data
	api.GET("/crimes/trends", publicHandler.GetCrimeTrends)    // Time trends
	api.GET("/crimes/areas", publicHandler.GetDangerousAreas)
	api.GET("/options", publicHandler.GetCrimesFilterOptions)
}
