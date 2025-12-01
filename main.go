package main

import (
	"shortlink/internal/database"
	"shortlink/internal/models"
	"shortlink/internal/routes"

	"github.com/gin-gonic/gin"
)

func main(){
	database := database.Database()
	
	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(201, models.Response{
			Success : true,
			Message : "back end runing",
		})
	})

	routes.MainApiRoutes(router, database)

	router.Run(":3231")
}