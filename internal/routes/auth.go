package routes

import (
	"shortlink/internal/handler"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AuthRoutes(r *gin.RouterGroup, pool *pgxpool.Pool){
	authController := handler.AuthController{Pool: pool}

	auth := r.Group("")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.POST("/refresh", authController.Refresh)
	}
}