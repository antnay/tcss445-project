package state

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *State) registerPublicRoutes() {
	s.Router.GET("/public/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "public ping",
		})
	})

}

func (s *State) registerPrivateRoutes() {
	{
		admin := s.Router.Group("/admin")
		admin.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "admin ping",
			})

		})

	}

	s.Router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "private ping",
		})
	})

}
