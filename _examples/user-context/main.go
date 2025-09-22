package main

import (
	"context"
	"fmt"
	"time"

	"github.com/SoeltanIT/agg-common-be/contek"
	"github.com/SoeltanIT/agg-common-be/types"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	token := createJWTToken()
	ctx := context.WithValue(context.Background(), "user", token)

	// Get user claims from context
	claims := contek.GetUserContext(ctx)
	if claims != nil {
		fmt.Println("=== User Claims ===")
		fmt.Printf("User ID: %s\n", claims.ID)
		fmt.Printf("Email: %s\n", claims.Email)
		fmt.Printf("Namespace: %s\n", claims.Namespace)
	}

	// Get raw token string
	rawToken := contek.GetUserRawToken(ctx)
	fmt.Println("\n=== Raw Token ===")
	fmt.Println(rawToken)
}

func createJWTToken() *jwt.Token {
	// Create a new token object with custom claims
	claims := &types.JWTClaims{
		ID:        "user-123",
		Email:     "user@example.com",
		Namespace: "test-namespace",
		Type:      "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// In a real application, you would sign the token with a secret key
	// tokenString, err := token.SignedString([]byte("your-secret-key"))

	// For demonstration, we'll just set the raw token string
	token.Raw = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InVzZXItMTIzIiwiZW1haWwiOiJ1c2VyQGV4YW1wbGUuY29tIiwibmFtZXNwYWNlIjoidGVzdC1uYW1lc3BhY2UiLCJ0eXBlIjoiYWRtaW4ifQ.EpM5kVzJesfnfNXCouZ7XhJjFmJQ6ZJ8X8XwQ6tqY0I"

	return token
}
