package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func MainApiRoutes(r *gin.Engine, pool *pgxpool.Pool){
	
		auth := r.Group("/auth")
		{
			AuthRoutes(auth, pool)
		}
	
	
}