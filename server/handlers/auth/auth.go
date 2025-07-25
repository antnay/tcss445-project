package auth

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

const (
	EIGHT_HOUR = 3600 * 8
)

type Handler struct {
	pool         *pgxpool.Pool
	tokenFactory *utils.TokenFactory
}

type AuthError struct {
	Code    int
	Message string
	Err     error
}

type RegisterRequest struct {
	Email    string `form:"email" binding:"required,email"`
	Username string `form:"username" binding:"required,alphanum,min=3,max=20"`
	Password string `form:"password" binding:"required,min=8"`
}

func NewHandler(pool *pgxpool.Pool, tokenFactory *utils.TokenFactory) *Handler {
	return &Handler{
		pool:         pool,
		tokenFactory: tokenFactory,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Printf("Failed to bind form: %s", err)
		c.Error(utils.NewValidationError(err.Error()))
		return
	}

	appE := validateRegistrationForm(&req)
	if appE != nil {
		c.Error(appE)
		return
	}

	appE = h.checkUserExists(c, &req)
	if appE != nil {
		c.Error(appE)
		return
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("Failed to hash password: %s", err)
		c.Error(utils.NewInternalError())
		return
	}

	appE = h.createUser(&req, passwordHash)
	if appE != nil {
		c.Error(appE)
		return
	}

	access_token, appE := h.getAccessToken(req.Username)
	if appE != nil {
		c.Error(appE)
		return
	}

	// appE = h.setAccessCookie(c, req.Username)
	// if appE != nil {
	// 	c.Error(appE)
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"message":      fmt.Sprintf("registered %s", req.Username),
		"access_token": access_token,
	})
}

func (h *Handler) Login(c *gin.Context) {
}

func validateRegistrationForm(req *RegisterRequest) *utils.AppError {
	if len(req.Email) == 0 {
		log.Println("Missing email")
		return utils.NewValidationError("missing email")
	} else if len(req.Username) == 0 {
		log.Println("Missing username")
		return utils.NewValidationError("missing username")
	} else if len(req.Password) == 0 {
		log.Println("Missing password")
		return utils.NewValidationError("missing password")
	}
	return nil
}

func (h *Handler) checkUserExists(c *gin.Context, req *RegisterRequest) *utils.AppError {
	var emailExists, usernameExists bool
	err := h.pool.QueryRow(context.Background(),
		`SELECT 
        EXISTS(SELECT 1 FROM users WHERE email = $1) as email_exists,
        EXISTS(SELECT 1 FROM users WHERE username = $2) as username_exists`,
		req.Email, req.Username).Scan(&emailExists, &usernameExists)

	if err != nil {
		log.Printf("failed to check existence: %s\n", err)
		return utils.NewInternalError()
	} else if emailExists {
		return utils.NewConflictError("email already exists")
	} else if usernameExists {
		return utils.NewConflictError("username already exists")
	}
	return nil
}

func (h *Handler) createUser(req *RegisterRequest, passwordHash string) *utils.AppError {
	tx, err := h.pool.Begin(context.Background())
	if err != nil {
		log.Printf("Unable to start transaction: %s", err)
		return utils.NewInternalError()
	}
	defer tx.Rollback(context.Background())

	var userId int
	err = tx.QueryRow(context.Background(),
		`INSERT INTO users (username, email)
		VALUES ($1, $2)
		RETURNING id
		`,
		req.Username, req.Email).Scan(&userId)

	if err != nil {
		log.Printf("Insert user failed: %s", err)
		return utils.NewInternalError()
	}

	_, err = tx.Exec(context.Background(),
		`INSERT INTO password_login (user_id, password_hash)
		VALUES ($1, $2)
		`,
		userId, passwordHash)
	if err != nil {
		log.Printf("Insert password_login failed: %s", err)
		return utils.NewInternalError()
	}
	// err = tx.Commit(context.Background())
	if err != nil {
		log.Printf("Database commit failed: %s", err)
		return utils.NewInternalError()
	}
	return nil
}

func (h *Handler) setAccessCookie(c *gin.Context, username string) *utils.AppError {
	curTime := time.Now()
	accessToken, err := h.tokenFactory.CreateAccessToken(username, "user", curTime)
	if err != nil {
		log.Printf("Failed sign access token: %s", err)
		return utils.NewInternalError()
	}

	c.SetSameSite(http.SameSiteDefaultMode)
	c.SetCookie("access_token", accessToken, EIGHT_HOUR, "/", "localhost", true, true)
	return nil
}

func (h *Handler) getAccessToken(username string) (string, *utils.AppError) {
	curTime := time.Now()
	accessToken, err := h.tokenFactory.CreateAccessToken(username, "user", curTime)
	if err != nil {
		log.Printf("Failed sign access token: %s", err)
		return accessToken, utils.NewInternalError()
	}
	return accessToken, nil
}
