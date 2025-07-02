package utils

import (
	"crypto/ed25519"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// var (
// HMAC
// ACCESS_KEY  = os.Getenv("JWT_ACCESS_SECRET")
// REFRESH_KEY = os.Getenv("JWT_REFRESH_SECRET")
// )

type TokenFactory struct {
	access_key_priv  ed25519.PrivateKey
	access_key_pub   ed25519.PublicKey
	refresh_key_priv ed25519.PrivateKey
	refresh_key_pub  ed25519.PublicKey
}

func (t *TokenFactory) NewTokenFactory() *TokenFactory {
	accessPrivKeyBytes := os.Getenv("JWT_ACCESS_DSA_KEY_PRIV")
	accessPrivKey := ed25519.PrivateKey(accessPrivKeyBytes)

	accessPubKeyBytes := os.Getenv("JWT_ACCESS_DSA_KEY_PUB")
	accessPubKey := ed25519.PublicKey(accessPubKeyBytes)

	refreshPrivKeyBytes := os.Getenv("JWT_REFRESH_DSA_KEY_PRIV")
	refreshPrivKey := ed25519.PrivateKey(refreshPrivKeyBytes)

	refreshPubKeyBytes := os.Getenv("JWT_REFRESH_DSA_KEY_PUB")
	refreshPubKey := ed25519.PublicKey(refreshPubKeyBytes)

	return &TokenFactory{
		access_key_priv:  accessPrivKey,
		access_key_pub:   accessPubKey,
		refresh_key_priv: refreshPrivKey,
		refresh_key_pub:  refreshPubKey,
	}

}

func (t *TokenFactory) CreateAccessToken(username string, role string, curTime time.Time) (string, error) {
	exp := curTime.Add(time.Hour * 8)
	log.Println("iat", curTime)
	log.Println("exp", exp)

	accessToken := jwt.NewWithClaims(jwt.SigningMethodEdDSA,
		jwt.MapClaims{
			"iss":      "server",
			"username": username,
			"role":     role,
			"iat":      curTime.Unix(),
			"exp":      exp.Unix(),
		})

	return accessToken.SignedString(t.access_key_priv)

	// // debug
	// signedToken, err := accessToken.SignedString(t.access_key_priv)
	// if err != nil {
	// 	log.Println("Error signing token:", err)
	// }

	// pubKeyBytes, err := x509.MarshalPKIXPublicKey(t.access_key_pub)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// pubKeyPEM := pem.EncodeToMemory(&pem.Block{
	// 	Type:  "PUBLIC KEY",
	// 	Bytes: pubKeyBytes,
	// })

	// fmt.Println("Public key (PEM):")
	// fmt.Println(string(pubKeyPEM))

	// parsedToken, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
	// 	if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
	// 		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	// 	}
	// 	return t.access_key_pub, nil
	// })
	// if err != nil {
	// 	log.Println("Error parsing/verifying token:", err)
	// }

	// if parsedToken.Valid {
	// 	log.Println("Token is valid")
	// } else {
	// 	log.Println("Token is invalid")

	// }

	// return signedToken, err
	// // end debug
}

func (t *TokenFactory) CreateRefreshToken(username string, role string, curTime time.Time) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodEdDSA,
		jwt.MapClaims{
			"iss":      "server",
			"username": username,
			"role":     role,
			"iat":      curTime.Unix(),
			"exp":      curTime.Add(168 * time.Hour).Unix(),
		})

	return accessToken.SignedString(t.refresh_key_priv)
}
