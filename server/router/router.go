package router

import (
	"log"
	"server/db"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Router struct {
	Router *gin.Engine
	pool   *pgxpool.Pool
}

func NewRouter() *Router {
	s := &Router{
		Router: gin.Default(),
	}

	s.connectToDatabase()
	s.registerRoutes()
	return s
}

func (r *Router) registerRoutes() {
	r.registerAdminRoutes()
	r.registerPublicRoutes()
	r.registerPrivateRoutes()
}

func (r *Router) connectToDatabase() {
	pool, err := db.Connect()
	if err != nil {
		log.Printf("Failed to connect to database: %v\n", err)
		log.Println("\033[31;1;4mDATABASE IS NOT ACTIVE, THINGS WILL BREAK\033[0m")
	}
	r.pool = pool
}

func (r *Router) PGClose() {
	if r.pool != nil {
		r.pool.Close()
	}
}
