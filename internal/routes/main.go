package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func MainApiRoutes(r *gin.Engine, pool *pgxpool.Pool){
	
	r.Static("/uploads", "./uploads")

	api := r.Group("/api/v1")
    {
        auth := api.Group("/auth")
        {
            AuthRoutes(auth, pool)
        }

        links := api.Group("/links")
        {
            LinkRoutes(links, pool)
        }
        dashboard := api.Group("/dashboard")
        {
            DashboardRouer(dashboard, pool)
        }
        user := api.Group("/users")
        {
            UserRouter(user, pool)
        }
    }
	
	
}