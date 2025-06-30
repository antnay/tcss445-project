package utils

import (
	"crypto/ed25519"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ACCESS_KEY  = ed25519.NewKeyFromSeed([]byte(os.Getenv("JWT_ACCESS_SECRET")))
	REFRESH_KEY = ed25519.NewKeyFromSeed([]byte(os.Getenv("JWT_REFRESH_SECRET")))
)

func CreateAccessToken(username string, role string, curTime time.Time) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodEdDSA,
		jwt.MapClaims{
			"iss":        "server",
			"username":   username,
			"role":       role,
			"iat":        curTime.Unix(),
			"exp":        curTime.Add(60 * time.Hour).Unix(),
		})

	return accessToken.SignedString(ACCESS_KEY)
}

func CreateRefreshToken(username string, role string, curTime time.Time) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodEdDSA,
		jwt.MapClaims{
			"iss":        "server",
			"username":   username,
			"role":       role,
			"iat":        curTime.Unix(),
			"exp":        curTime.Add(168 * time.Hour).Unix(),
		})

	return accessToken.SignedString(REFRESH_KEY)
}
