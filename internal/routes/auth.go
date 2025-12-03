package routes

import (
	"shortlink/internal/handler"
	"shortlink/internal/middelware"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AuthRoutes(r *gin.RouterGroup, pool *pgxpool.Pool){
	authController := handler.AuthController{Pool: pool}

	rateLimitLogin := middelware.RateLimit(5, 1*time.Minute)
	rateLimitRefresh := middelware.RateLimit(10, 1*time.Minute)
	auth := r.Group("")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", rateLimitLogin, authController.Login)
		auth.POST("/refresh", rateLimitRefresh,authController.Refresh)
		auth.POST("/logount", authController.Logout)
	}
}