package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uint   `json:"uid,omitempty"`
	Email  string `json:"email,omitempty"`
	AppID  uint   `json:"app_id,omitempty"`
}

func NewToken(claims *Claims, secret string, duration time.Duration) (string, error) {
	const fn = "lib.jwt.NewToken"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":    "local-user-service",
		"uid":    claims.UserID,
		"email":  claims.Email,
		"app_id": claims.AppID,
		"exp":    time.Now().Add(duration).Unix(),
		"iat":    time.Now().Unix(),
	})

	signedString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	return signedString, nil
}

func ValidateToken(bearerToken, secret string) error {
	const fn = "lib.jwt.ValidateToken"

	splittedToken := strings.Split(bearerToken, " ")

	if len(splittedToken) != 2 {
		return fmt.Errorf("%s: %w", fn, errors.New("invalid bearer token"))
	}

	var claims Claims

	_, err := jwt.ParseWithClaims(splittedToken[1], &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return err
	}

	return nil
}
