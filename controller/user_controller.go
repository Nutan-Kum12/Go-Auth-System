package controller

import (
	"Auth/model"
	"Auth/services"
	"Auth/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserController struct {
	Service services.UserService
}

func (c UserController) Register(context *gin.Context) {
	var user model.User
	err := context.ShouldBindJSON(&user) // convert request json to go struct
	if err != nil {
		context.JSON(
			http.StatusBadRequest, gin.H{
				"error": "Invalid JSON",
			})
		return
	}
	//call service
	err = c.Service.Register(user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//success response
	context.JSON(http.StatusOK, gin.H{
		"message": "User Registered Successfully",
	})
}

func (c UserController) Login(context *gin.Context) {
	var user model.User

	err := context.ShouldBindJSON(&user)// json->object(context->is used to read request body)
	if err != nil {
		context.JSON(
			http.StatusBadRequest, gin.H{
				"error": "Invalid JSON",
			})
		return
	}
	//call login service
	accessToken, refreshToken, err := c.Service.Login(
		user.Email,
		user.Password,
	)
	if err != nil {
		context.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}
	// set access token cookie (15 min)
	context.SetCookie(
		"access_token",  // cookie name
		accessToken,     // cookie value
		15*60,           // expiry (15 minutes)
		"/",             // path
		"localhost",
		false, // secure
		true,  // httpOnly
	)
	// set refresh token cookie (7 days)
	context.SetCookie(
		"refresh_token",  // cookie name
		refreshToken,     // cookie value
		7*24*3600,        // expiry (7 days)
		"/",              // path
		"localhost",
		false, // secure
		true,  // httpOnly
	)

	context.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
	})//send response
}

// RefreshToken validates the refresh token and issues a new access token
func (c UserController) RefreshToken(context *gin.Context) {
	// read refresh token from cookie
	refreshTokenString, err := context.Cookie("refresh_token")
	if err != nil {
		context.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "refresh token required",
			},
		)
		return
	}

	// verify refresh token
	token, err := utils.ParseRefreshToken(refreshTokenString)
	if err != nil || !token.Valid {
		context.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid or expired refresh token",
			},
		)
		return
	}

	// extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		context.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid claims",
			},
		)
		return
	}

	// verify token type
	tokenType, _ := claims["type"].(string)
	if tokenType != "refresh" {
		context.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid token type",
			},
		)
		return
	}

	// extract email
	email, _ := claims["sub"].(string)
	if email == "" {
		context.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid token payload",
			},
		)
		return
	}

	// generate new access token
	newAccessToken, err := utils.GenerateAccessToken(email)
	if err != nil {
		context.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "failed to generate access token",
			},
		)
		return
	}

	// set new access token cookie
	context.SetCookie(
		"access_token",
		newAccessToken,
		15*60,
		"/",
		"localhost",
		false,
		true,
	)

	context.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
	})
}

func (c UserController) Logout(context *gin.Context) {

	// clear access token cookie
	context.SetCookie(
		"access_token",
		"",
		-1,
		"/",
		"localhost",
		false,
		true,
	)
	// clear refresh token cookie
	context.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"localhost",
		false,
		true,
	)

	context.JSON(200, gin.H{
		"message": "Logout successful",
	})
}

func (c UserController) GetProfile(context *gin.Context) {
	//middleware se email lo
	email, exists := context.Get("email")
	if !exists {
		context.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "unauthorized",
			},
		)
		return
	}
	emailStr, ok := email.(string)
	//type assertion as ctx.Get return interface and service expect string
	if !ok {
		context.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": "invalid email",
			},
		)
		return
	}
	// service call
	user, err := c.Service.GetProfile(
		emailStr,
	)
	if err != nil {
		context.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}
	// response
	context.JSON(
		http.StatusOK,
		gin.H{
			"name":  user.Name,
			"email": user.Email,
		},
	)
}

