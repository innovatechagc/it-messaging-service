package main

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func main() {
	// Usar la misma clave secreta que en el .env
	secretKey := "dev-jwt-secret-key-change-in-production"
	issuer := "messaging-service"

	// Crear tokens para diferentes usuarios de prueba
	users := []struct {
		ID    string
		Email string
		Role  string
	}{
		{"user-1", "user1@example.com", "user"},
		{"user-2", "user2@example.com", "user"},
		{"admin-1", "admin@example.com", "admin"},
	}

	fmt.Println("=== JWT Tokens para Testing ===\n")

	for _, user := range users {
		token, err := generateToken(user.ID, user.Email, user.Role, secretKey, issuer)
		if err != nil {
			log.Printf("Error generando token para %s: %v", user.Email, err)
			continue
		}

		fmt.Printf("Usuario: %s (%s)\n", user.Email, user.Role)
		fmt.Printf("Token: %s\n\n", token)
	}

	fmt.Println("=== Instrucciones de uso ===")
	fmt.Println("1. Copia el token del usuario que quieras usar")
	fmt.Println("2. En Postman, agrega el header: Authorization: Bearer <token>")
	fmt.Println("3. Los tokens son válidos por 24 horas")
}

func generateToken(userID, email, role, secretKey, issuer string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}