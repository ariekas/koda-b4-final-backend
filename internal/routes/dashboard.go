package routes

import (
	"shortlink/internal/handler"
	"shortlink/internal/middelware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func DashboardRouer(r *gin.RouterGroup, pool *pgxpool.Pool) {
	dashboardController := handler.DashboardController{Pool: pool}

	dashboard := r.Group("/stats", middelware.VerifToken())
	{
		dashboard.GET("", dashboardController.GetDashboardStats)
	}
}
