package main

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func main() {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		fmt.Println("JWT_SECRET_KEY is not set")
		os.Exit(1)
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(10 * time.Minute).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("Error creating token:", err)
		os.Exit(1)
	}

	_, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("Error verifying token:", err)
		os.Exit(1)
	}

	fmt.Println("Successfully created and verified JWT token!")
}
