package main

import (
	"shortlink/internal/models"

	"github.com/gin-gonic/gin"
)

func main(){
	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(201, models.Response{
			Success : true,
			Message : "back end runing",
		})
	})

	router.Run(":3231")
}