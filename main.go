package main

import (
	"shortlink/internal/config"
	"shortlink/internal/database"
	"shortlink/internal/handler"
	"shortlink/internal/middelware"
	"shortlink/internal/models"
	"shortlink/internal/routes"

	"github.com/gin-gonic/gin"
)

func main(){
	database := database.Database()
	config.InitRedis()
	
	router := gin.Default()

	router.MaxMultipartMemory = 8 << 20
	router.Use(middelware.Cors())
	router.Use(middelware.AllowPreflight)

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(201, models.Response{
			Success : true,
			Message : "back end runing",
		})
	})

	slc := handler.ShortLinkController{Pool: database}

	router.GET("/:slug", slc.Redirect)

	routes.MainApiRoutes(router, database)

	router.Run(":8082")
}