package middelware

import "github.com/gin-gonic/gin"

func AllowPreflight(ctx *gin.Context){
	if ctx.Request.Method == "OPTIONS" {
		ctx.AbortWithStatus(200)
	} else {
		ctx.Next()
	}
}