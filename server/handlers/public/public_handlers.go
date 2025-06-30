package public

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"server/utils"
	"time"

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
		"message": "public pong",
	})
}

func (h *Handler) Register(c *gin.Context) {
	emailForm := c.PostForm("email")
	usernameForm := c.PostForm("username")
	passwordForm := c.PostForm("password")

	if len(emailForm) == 0 {
		log.Println("Missing email")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "missing email",
		})
		return
	}
	if len(usernameForm) == 0 {
		log.Println("Missing username")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "missing username",
		})
		return
	}
	if len(passwordForm) == 0 {
		log.Println("Missing password")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "missing password",
		})
		return
	}

	var emailExists, usernameExists bool
	err := h.pool.QueryRow(context.Background(),
		`SELECT 
        EXISTS(SELECT 1 FROM users WHERE email = $1) as email_exists,
        EXISTS(SELECT 1 FROM users WHERE username = $2) as username_exists`,
		emailForm, usernameForm).Scan(&emailExists, &usernameExists)

	if err != nil {
		log.Printf("failed to check existence: %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	if emailExists {
		c.JSON(http.StatusConflict, gin.H{
			"message": "email already exists",
		})
		return
	}
	if usernameExists {
		c.JSON(http.StatusConflict, gin.H{
			"message": "user already exists",
		})
		return
	}

	passwordHash, err := utils.HashPassword(passwordForm)
	if err != nil {
		log.Printf("Failed to hash password: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	tx, err := h.pool.Begin(context.Background())
	if err != nil {
		log.Printf("Unable to start transaction: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}
	defer tx.Rollback(context.Background())

	var userId int
	err = tx.QueryRow(context.Background(),
		`INSERT INTO users (username, email)
		VALUES ($1, $2)
		RETURNING id
		`,
		usernameForm, emailForm).Scan(&userId)

	if err != nil {
		log.Printf("Insert user failed: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	_, err = tx.Exec(context.Background(),
		`INSERT INTO password_login (user_id, password_hash)
		VALUES ($1, $2)
		`,
		userId, passwordHash)
	if err != nil {
		log.Printf("Insert password_login failed: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}
	// err = tx.Commit(context.Background())
	if err != nil {
		log.Printf("Database commit failed: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	curTime := time.Now()
	accessToken, err := utils.CreateAccessToken(usernameForm, "user", curTime)
	if err != nil {
		log.Printf("Failed sign access token: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}
	log.Println(accessToken)

	c.SetSameSite(http.SameSiteDefaultMode)
	c.SetCookie("access_token", accessToken, int(time.Hour*8), "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("registered %s", usernameForm),
	})

}

func (h *Handler) Login(c *gin.Context) {
}
