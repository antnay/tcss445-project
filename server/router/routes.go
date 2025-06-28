package router

import (
	adminHandlers "server/handlers/admin"
	privateHandlers "server/handlers/private"
	publicHandlers "server/handlers/public"
)

func (r *Router) registerAdminRoutes() {
	adminHandler := adminHandlers.NewHandler(r.pool)

	// /admin/* group
	admin := r.Router.Group("/admin")

	admin.GET("/ping", adminHandler.AdminPing)  // /admin/ping group
	admin.POST("/pong", adminHandler.AdminPong) // /admin/pong group
}

func (r *Router) registerPrivateRoutes() {
	privateHandler := privateHandlers.NewHandler(r.pool)
	// /ping
	r.Router.GET("/ping", privateHandler.Ping)
	{
		group := r.Router.Group("/group")
		group.GET("/one", privateHandler.Ping)
	}

	// path parameter example
	r.Router.GET("/hello/:name/:number", privateHandler.PathParamEx)
	// query parameters
	r.Router.GET("/hello", privateHandler.QueryParamEx)

}
func (r *Router) registerPublicRoutes() {
	publicHandler := publicHandlers.NewHandler(r.pool)
	// /public/* group
	public := r.Router.Group("/public")

	public.GET("/ping", publicHandler.Ping) // /public/ping
}
