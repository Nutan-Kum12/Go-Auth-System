package main

import (
	"net/http"

	"Auth/config"
	"Auth/controller"
	"Auth/repository"
	"Auth/routes"
	"Auth/services"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	//repository
	repo := repository.MongoRepository{}
	//service
	userService := services.UserService{
		Repo: repo,
	}
	//controller
	UserController := controller.UserController{
		Service: userService,
	}
	router := gin.Default()
	//Gin router/server object banata haii
	// Is router ka kaam:
	// requests receive karna
	// endpoints manage karna
	// responses bhejna
	router.GET("/", func(c *gin.Context) {
		//Agar GET request "/" endpoint pe aaye
		//to ye function chalao

		// Ye:c *gin.Context
		// request + response context
		// hota hai.
		// Isse:
		// request read
		// JSON send
		// headers read
		// params access
		// karte hain.
		c.JSON(http.StatusOK, gin.H{
			"message": "Server Running",
		})
	})
	// routes
	routes.UserRoutes(
		router,
		UserController,
	)

	// run server
	router.Run(":8080")
}
