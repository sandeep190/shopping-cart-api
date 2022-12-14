package middleware

import (
	"fmt"
	"net/http"
	"os"
	"shopping_cart/models"
	"strings"

	"shopping_cart/database"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("authUser")
		if exists && user.(models.User).ID != 0 {
			return
		} else {
			err, _ := c.Get("authErr")
			_ = c.AbortWithError(http.StatusUnauthorized, err.(error))
			return
		}
	}
}

func JwtTokenVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.Request.Header.Get("Authorization")
		if bearer != "" {
			jwtString := strings.Split(bearer, " ")
			if len(jwtString) == 2 {
				jwtEncoded := jwtString[1]
				token, err := jwt.Parse(jwtEncoded, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signin method %v", token.Header["alg"])
					}
					secret := []byte(os.Getenv("JWT_SECRET"))
					return secret, nil
				})

				if err != nil {
					println(err.Error())
					return
				}
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					userId := uint(claims["user_id"].(float64))
					fmt.Printf("authenticated user id is %d\n", userId)

					var user models.User
					if userId != 0 {
						database := database.GetConnection()
						database.First(&user, userId)
					}

					c.Set("authUser", user)
					c.Set("authUserId", user.ID)
				} else {
					fmt.Printf("request #%v", c.Request)
				}

			}
		}
	}
}
