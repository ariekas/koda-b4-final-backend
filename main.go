package main

import "github.com/gin-gonic/gin"

func main(){
	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(201, gin.H{
			"success" : true,
			"message": "back end runing",
		})
	})

	router.Run(":3231")
}