package middleware

import (
	"Auth/utils"
	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// read cookie (JWT from user)
		tokenString, err := ctx.Cookie("access_token")
		if err != nil {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "login required",
				},
			)
			ctx.Abort()
			return
		}

		// verify token
		token, err := utils.ParseAccessToken(
			tokenString,
		)
		// fmt.Println("TOKEN VALID:", token.Valid)
		// fmt.Println("PARSE ERROR:", err)
		if err != nil || !token.Valid {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "invalid token",
				},
			)
			ctx.Abort()
			return
		}
		// extract claims
		claims, ok := token.Claims.(jwt.MapClaims) // extract the payload(email,sub,type,iat,exp)
		// fmt.Println("CLAIMS:", claims)
		if !ok {
			ctx.JSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "invalid claims",
				},
			)
			ctx.Abort() //
			return
		}

		// save email in request context
		ctx.Set(
			"email", claims["sub"],
		)
		//we can also extract type,iat,exp
		// tokenType:=claims["type"].(string)
		// iat:=claims["iat"].(float64)
		// exp:=claims["exp"].(float64)
		// fmt.Println(tokenType,iat,exp)

		ctx.Next()
	}
}
