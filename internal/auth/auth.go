package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func ValidateJWT(tokenString string) (string, error) {

	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")

	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(jwtSecret), nil },
	)
	if err != nil {
		return "", err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	return userIDString, nil
}

func CreateToken(ExpiresInSeconds *int, id int) (string, error) {

	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")

	timeToAdd := 24 * time.Hour
	if ExpiresInSeconds != nil {
		timeToAdd = time.Duration(*ExpiresInSeconds) * time.Second
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(timeToAdd)),
		Subject:   fmt.Sprint(id),
	})

	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
