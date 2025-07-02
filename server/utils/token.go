package utils

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
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

func NewTokenFactory() (*TokenFactory, error) {
	accessPrivKeyStr := os.Getenv("JWT_ACCESS_DSA_KEY_PRIV")
	if accessPrivKeyStr == "" {
		return nil, fmt.Errorf("JWT_ACCESS_DSA_KEY_PRIV not set")
	}
	accessPrivKey, err := parseEd25519PrivateKey(accessPrivKeyStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse access private key: %w", err)
	}

	accessPubKeyStr := os.Getenv("JWT_ACCESS_DSA_KEY_PUB")
	if accessPubKeyStr == "" {
		return nil, fmt.Errorf("JWT_ACCESS_DSA_KEY_PUB not set")
	}
	accessPubKey, err := parseEd25519PublicKey(accessPubKeyStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse access public key: %w", err)
	}

	refreshPrivKeyStr := os.Getenv("JWT_REFRESH_DSA_KEY_PRIV")
	if refreshPrivKeyStr == "" {
		return nil, fmt.Errorf("JWT_REFRESH_DSA_KEY_PRIV not set")
	}
	refreshPrivKey, err := parseEd25519PrivateKey(refreshPrivKeyStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse refresh private key: %w", err)
	}

	refreshPubKeyStr := os.Getenv("JWT_REFRESH_DSA_KEY_PUB")
	if refreshPubKeyStr == "" {
		return nil, fmt.Errorf("JWT_REFRESH_DSA_KEY_PUB not set")
	}
	refreshPubKey, err := parseEd25519PublicKey(refreshPubKeyStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse refresh public key: %w", err)
	}

	return &TokenFactory{
		access_key_priv:  accessPrivKey,
		access_key_pub:   accessPubKey,
		refresh_key_priv: refreshPrivKey,
		refresh_key_pub:  refreshPubKey,
	}, nil
}

func parseEd25519PrivateKey(keyStr string) (ed25519.PrivateKey, error) {
	block, _ := pem.Decode([]byte(keyStr))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PKCS8 private key: %w", err)
	}

	ed25519PrivKey, ok := parsedKey.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not Ed25519 private key")
	}

	if len(ed25519PrivKey) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("invalid Ed25519 private key size: got %d, want %d",
			len(ed25519PrivKey), ed25519.PrivateKeySize)
	}

	return ed25519PrivKey, nil
}

func parseEd25519PublicKey(keyStr string) (ed25519.PublicKey, error) {
	block, _ := pem.Decode([]byte(keyStr))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PKIX public key: %w", err)
	}

	ed25519PubKey, ok := parsedKey.(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("key is not Ed25519 public key")
	}

	if len(ed25519PubKey) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("invalid Ed25519 public key size: got %d, want %d",
			len(ed25519PubKey), ed25519.PublicKeySize)
	}

	return ed25519PubKey, nil
}

func (t *TokenFactory) CreateAccessToken(username string, role string, curTime time.Time) (string, error) {
	exp := curTime.Add(time.Hour * 8)

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
