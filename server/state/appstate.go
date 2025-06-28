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
	pool, err := db.PGConnect()
	if err != nil {
		log.Printf("Failed to connect to database: %v\n", err)
		log.Println("DATABASE IS NOT ACTIVE, THINGS WILL BREAK")
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
