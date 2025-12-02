package routes

import (
	"shortlink/internal/handler"
	"shortlink/internal/middelware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func LinkRoutes(r *gin.RouterGroup, pool *pgxpool.Pool) {
	linkController := handler.ShortLinkController{Pool: pool}

	link := r.Group("", middelware.VerifToken())
	{
		link.POST("/", linkController.Create)
		link.GET("", linkController.GetAll)
		link.GET("/:slug", linkController.DetailShortCode) 
	}
}
