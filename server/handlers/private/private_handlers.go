package private

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
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
		"message": "private pong",
	})
}

func (h *Handler) PathParamEx(c *gin.Context) {
	name := c.Param("name")
	num := c.Param("number")

	number, err := strconv.Atoi(num)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid number"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"number":  number,
		"message": fmt.Sprintf("hello %s", name),
	})
}

func (h *Handler) QueryParamEx(c *gin.Context) {
	name := c.Query("name")                 // Returns empty string if not found
	number := c.DefaultQuery("number", "5") // Returns default if not found
	// Query arrays: ?tags=go&tags=web&tags=api
	tags := c.QueryArray("tags")

	c.JSON(http.StatusOK, gin.H{
		"number":  number,
		"message": fmt.Sprintf("hello %s", name),
		"tags":    tags,
	})
}
