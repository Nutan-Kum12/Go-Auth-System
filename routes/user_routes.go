package routes

import (
	"Auth/controller"
	"Auth/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, controller controller.UserController) {

	router.POST(
		"user/register",
		controller.Register,
	)
	router.POST(
		"user/login",
		controller.Login,
	)
	router.POST(
		"/user/logout",
		middleware.AuthMiddleWare(),
		controller.Logout,
	)
	router.POST(
		"/user/refresh",
		controller.RefreshToken,
	)
	router.GET(
		"/user/profile",
		middleware.RateLimiter(),
		middleware.AuthMiddleWare(),
		controller.GetProfile,
	)
}
