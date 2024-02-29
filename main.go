package main

import (
	userhandler "webserver/handler/userhandler"
	migration "webserver/migrations"

	"github.com/gin-gonic/gin"
)

func main() {
	migration.RunMigrations()
	router := gin.Default()
	router.POST("/user", userhandler.RegisterUser)
	router.GET("/user", userhandler.GetUser)
	router.DELETE("/user", userhandler.DeleteUser)

	router.Run("0.0.0.0:8085")
}
