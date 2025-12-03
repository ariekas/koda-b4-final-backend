package routes

import (
	"shortlink/internal/handler"
	"shortlink/internal/middelware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UserRouter(r *gin.RouterGroup, pool *pgxpool.Pool) {
	userController := handler.UserController{Pool: pool}

	dashboard := r.Group("", middelware.VerifToken())
	{
		dashboard.GET("/profile", userController.GetProfile)
		dashboard.GET("/:id", userController.GetUserById)
		dashboard.POST("/pic", userController.UploadProfilePicture)
	}
}
