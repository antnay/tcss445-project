package router

import (
	adminHandlers "server/handlers/admin"
	privateHandlers "server/handlers/private"
	publicHandlers "server/handlers/public"
)

func (r *Router) registerAdminRoutes() {
	adminHandler := adminHandlers.NewHandler(r.pool)
	api := r.Router.Group("/api")

	// /admin/* group
	admin := api.Group("/admin")

	admin.GET("/ping", adminHandler.AdminPing)  // api/admin/ping group
	admin.POST("/pong", adminHandler.AdminPong) // api/admin/pong group
}

func (r *Router) registerPrivateRoutes() {
	privateHandler := privateHandlers.NewHandler(r.pool)
	api := r.Router.Group("/api")
	// api/ping
	api.GET("/ping", privateHandler.Ping)
	{
		// group := api.Group("/user")
	}

	// path parameter example
	api.GET("/hello/:name/:number", privateHandler.PathParamEx)
	// query parameters
	api.GET("/hello", privateHandler.QueryParamEx)

}

func (r *Router) registerPublicRoutes() {
	publicHandler := publicHandlers.NewHandler(r.pool)
	api := r.Router.Group("/api/public")

	api.GET("/ping", publicHandler.Ping) // api/public/ping
	api.POST("/register", publicHandler.Register)
	api.POST("/login", publicHandler.Login)
}
