package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/LilLebowski/loyalty-system/cmd/gophermart/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const CookieName = "UserID"

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func Authorization(config *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/api/user/register" || ctx.Request.URL.Path == "/api/user/login" {
			return
		}
		userID, err := getUserIDFromCookie(ctx, config)
		if err != nil {
			code := http.StatusUnauthorized
			contentType := ctx.Request.Header.Get("Content-Type")
			if contentType == "application/json" {
				ctx.Header("Content-Type", "application/json")
				ctx.JSON(code, gin.H{
					"message": fmt.Sprintf("Unauthorized %s", err),
					"code":    code,
				})
			} else {
				ctx.String(code, fmt.Sprintf("Unauthorized %s", err))
			}
			ctx.Abort()
			return
		}
		ctx.Set(CookieName, userID)
	}
}

func getUserIDFromCookie(ctx *gin.Context, config *config.Config) (string, error) {
	token, err := ctx.Cookie(CookieName)
	if err != nil {
		return "", err
	}
	userID, err := getUserID(token, config.SecretKey)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func BuildJWTString(config *config.Config, userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.TokenExpire)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func getUserID(tokenString string, secretKey string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})
	if err != nil {
		return "", fmt.Errorf("token is not valid")
	}

	if !token.Valid {
		return "", fmt.Errorf("token is not valid")
	}

	return claims.UserID, nil
}
