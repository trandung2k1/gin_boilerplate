package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func VerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var arrAuthorization = c.Request.Header["Authorization"]
		if len(arrAuthorization) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "failure",
				"message": "Authorization headers cannot be empty",
			})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		} else {
			var tokenString = c.Request.Header["Authorization"][0]
			if tokenString == "" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status":  "failure",
					"message": "Token string cannot be empty",
				})
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			} else {
				var token = strings.Split(tokenString, " ")[1]
				if token == "" {
					c.JSON(http.StatusUnauthorized, gin.H{
						"status":  "failure",
						"message": "Token cannot be empty",
					})
					c.AbortWithStatus(http.StatusUnauthorized)
					return
				} else {
					tk, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
						return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
					})
					if tk.Valid {
						claims := tk.Claims.(jwt.MapClaims)
						userId := fmt.Sprintf("%v", claims["userId"])
						c.Request.Header.Set("userId", userId)
						c.Next()
					} else if ve, ok := err.(*jwt.ValidationError); ok {
						if ve.Errors&jwt.ValidationErrorMalformed != 0 {
							c.JSON(http.StatusBadRequest, gin.H{
								"status":  "failure",
								"message": "Token is not valid",
							})
							c.AbortWithStatus(http.StatusUnauthorized)
							return
						} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
							c.JSON(http.StatusBadRequest, gin.H{
								"status":  "failure",
								"message": "Token is expired",
							})
							c.AbortWithStatus(http.StatusUnauthorized)
							return
						} else {
							c.JSON(http.StatusBadRequest, gin.H{
								"status":  "failure",
								"message": "Couldn't handle this token",
							})
							c.AbortWithStatus(http.StatusUnauthorized)
							return
						}
					} else {
						c.JSON(http.StatusBadRequest, gin.H{
							"status":  "failure",
							"message": "Couldn't handle this token",
						})
						c.AbortWithStatus(http.StatusUnauthorized)
						return
					}
				}
			}
		}
	}
}
