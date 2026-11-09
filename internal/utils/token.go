package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func GenerateToken(userID int64) (string, error) {
	if userID < 1 {
		return "", fmt.Errorf("user ID must be positive")
	}
	secret, err := jwtSecret()
	if err != nil {
		return "", err
	}
	tokenLifespan, err := strconv.Atoi(os.Getenv("TOKEN_LIFESPAN"))
	if err != nil {
		return "", fmt.Errorf("TOKEN_LIFESPAN must be a positive number of hours: %w", err)
	}
	if tokenLifespan < 1 {
		return "", fmt.Errorf("TOKEN_LIFESPAN must be a positive number of hours")
	}
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func ExtractToken(ctx *gin.Context) string {
	if token := strings.TrimSpace(ctx.Query("token")); token != "" {
		return token
	}
	fields := strings.Fields(ctx.GetHeader("Authorization"))
	if len(fields) == 2 && strings.EqualFold(fields[0], "Bearer") {
		return fields[1]
	}
	return ""
}

func TokenValid(ctx *gin.Context) error {
	_, err := parseTokenClaims(ctx)
	return err
}

func ExtractTokenUserID(ctx *gin.Context) (int64, error) {
	claims, err := parseTokenClaims(ctx)
	if err != nil {
		return 0, err
	}
	userID, err := strconv.ParseInt(fmt.Sprintf("%.0f", claims["user_id"]), 10, 64)
	if err != nil || userID < 1 {
		return 0, fmt.Errorf("invalid user_id claim")
	}
	return userID, nil
}

func ExtractTokenID(ctx *gin.Context) (uint, error) {
	userID, err := ExtractTokenUserID(ctx)
	if err != nil {
		return 0, err
	}
	if uint64(userID) > uint64(^uint(0)) {
		return 0, fmt.Errorf("user_id claim exceeds supported range")
	}
	return uint(userID), nil
}

func parseTokenClaims(ctx *gin.Context) (jwt.MapClaims, error) {
	tokenString := ExtractToken(ctx)
	if tokenString == "" {
		return nil, fmt.Errorf("missing or malformed Authorization: Bearer token")
	}
	secret, err := jwtSecret()
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid JWT claims")
	}
	return claims, nil
}

func jwtSecret() ([]byte, error) {
	secret := strings.TrimSpace(os.Getenv("SECRET_KEY"))
	if secret == "" {
		return nil, fmt.Errorf("SECRET_KEY is required for JWT authentication")
	}
	return []byte(secret), nil
}
