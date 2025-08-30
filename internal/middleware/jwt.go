package middleware

import (
	"os"
	"strings"
	"time"
	"z-cube-backend/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte

func SetKey(key []byte) {
	if key != nil {
		jwtKey = key
	} else {
		jwtKey = []byte(os.Getenv("JWT_KEY"))
	}
}

func GenerateToken(userID int, username string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(ttl).Unix(),
	})
	return token.SignedString(jwtKey)
}

func jwtAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(403, response.BaseResponse{
				Data: "valid token",
			})
			return
		}

		list := strings.Split(authHeader, " ")
		if len(list) < 2 {
			ctx.AbortWithStatusJSON(403, response.BaseResponse{
				Data: "valid token",
			})
			return
		}

		token, err := jwt.Parse(list[1], func(t *jwt.Token) (any, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(403, response.BaseResponse{
				Data: "valid token",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			exp, ok := claims["exp"].(int64)
			if ok {
				if time.Now().Unix() > exp {
					ctx.AbortWithStatusJSON(403, response.BaseResponse{
						Data: "valid token",
					})
					return
				}
			}

			userId, ok := claims["user_id"]
			if ok {
				ctx.Set("userId", userId)
			}
			username, ok := claims["username"]
			if ok {
				ctx.Set("username", username)
			}
		}

		ctx.Next()
	}
}
