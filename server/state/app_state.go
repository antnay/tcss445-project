package state

import (
	"log"
	"server/db"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type State struct {
	Router *gin.Engine
	Pool   *pgxpool.Pool
}

func NewServer() *State {
	s := &State{
		Router: gin.Default(),
	}
	s.addPool()
	s.registerRoutes()
	return s
}

func (s *State) registerRoutes() {
	s.registerPublicRoutes()
	s.registerPrivateRoutes()
}

func (s *State) addPool() {
	pool, err := db.Connect()
	if err != nil {
		log.Printf("Failed to connect to database: %v\n", err)
		log.Println("\033[31m\033[\033[4mDATABASE IS NOT ACTIVE, THINGS WILL BREAK\033[0m")
	}
	s.Pool = pool
}

func (s *State) GetPgPool() *pgxpool.Pool {
	return s.Pool
}

func (s *State) PGClose() {
	if s.Pool != nil {
		s.Pool.Close()
	}
}
