package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(role string) (string, error) {

	jwtSecret := os.Getenv("JWT_SECRET")
	jwtExpiry, err := strconv.Atoi(os.Getenv("JWT_EXPIRY_HOURS"))
	if err != nil || jwtExpiry <= 0 {
		return "", fmt.Errorf("Парсинг JWT_EXPIRY_HOURS: %w", err)
	}

	claims := jwt.MapClaims{
		"role": role,
		"exp":  time.Now().Add(time.Duration(jwtExpiry) * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Println("Ошибка создания токена:", err)
		return "", err
	}

	return signedToken, nil
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Неожиданный signing метод: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("Неправильный токен: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Невозможно распарсить claims")
	}

	return claims, nil
}
