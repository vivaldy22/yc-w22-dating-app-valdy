package jwt

import (
	"errors"
	"log"

	"github.com/golang-jwt/jwt/v5"

	ierror "yc-w22-dating-app-valdy/pkg/error"
)

// GenerateJWT generates JWT token with claims
func GenerateJWT(secret string, claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Printf("failed to sign access token: %s\n", err.Error())
		return "", ierror.ErrGeneral
	}

	return accessToken, nil
}

// ValidateJWT validates the JWT token and returns the claims if valid.
func ValidateJWT(secret, tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Unexpected signing method: %v", token.Header["alg"])
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Printf("invalid token")
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
