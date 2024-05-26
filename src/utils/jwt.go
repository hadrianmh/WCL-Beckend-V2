package utils

import (
	"backend/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Tokens struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type SecretKey struct {
	JWTSecretKey    string
	JWTRefSecretKey string
}

// Global variables for secret keys
var (
	JWTSecretKey    []byte
	JWTRefSecretKey []byte
	JWTtokenEXP     int
	JWTreftokenEXP  int
)

func GetSecretKey() error {
	config, err := config.LoadConfig("./config.json")
	if err != nil {
		return err
	}

	JWTSecretKey = []byte(config.App.JWTSecretKey)
	JWTRefSecretKey = []byte(config.App.JWTRefSecretKey)
	JWTtokenEXP = config.App.JWTtokenexp
	JWTreftokenEXP = config.App.JWTreftokenexp

	return nil
}

// GenerateTokens generates a new access and refresh token for a given username
func GenerateTokens(uniqid int, username string, name string, role string) (Tokens, error) {
	if err := GetSecretKey(); err != nil {
		return Tokens{}, err
	}

	accessToken, err := createJWT(uniqid, username, name, role, JWTSecretKey, time.Minute*time.Duration(JWTtokenEXP))
	if err != nil {
		return Tokens{}, err
	}

	refreshToken, err := createJWT(uniqid, username, name, role, JWTRefSecretKey, time.Hour*24*time.Duration(JWTreftokenEXP))
	if err != nil {
		return Tokens{}, err
	}

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// createJWT creates a JWT token for a given username and expiration duration
func createJWT(uniqid int, username string, name string, role string, secretKey []byte, expiration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"uniqid":   uniqid,
		"username": username,
		"name":     name,
		"role":     role,
		"exp":      time.Now().Add(expiration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ValidateToken validates a given JWT token and returns the username if valid
func ValidateToken(tokenString string, request string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKey
		}

		if request == "refresh_token" {
			return JWTRefSecretKey, nil
		}

		return JWTSecretKey, nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", jwt.ErrInvalidKey
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", jwt.ErrInvalidKey
	}

	return username, nil
}

// AuthenticateJWT is a middleware that checks the validity of the JWT token
func AuthenticateJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":   http.StatusUnauthorized,
				"status": "error",
				"response": gin.H{
					"message": "access denied"}})
			ctx.Abort()
			return
		}

		username, err := ValidateToken(tokenString, "access_token")
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":   http.StatusUnauthorized,
				"status": "error",
				"response": gin.H{
					"message": "invalid token"}})
			ctx.Abort()
			return
		}

		ctx.Set("username", username)
		ctx.Next()
	}
}
