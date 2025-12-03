package handler

import (
	"net/http"
	"shortlink/internal/models"
	"shortlink/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardController struct {
	Pool *pgxpool.Pool
}

func (dc *DashboardController) GetDashboardStats(c *gin.Context) {
	userIDInterface, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Unauthorized - user not authenticated",
		})
		return
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Invalid user ID format",
		})
		return
	}

	totalLinks, err := repository.GetTotalLinks(userID, dc.Pool)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch total links",
			"details": err.Error(),
		})
		return
	}

	totalVisits, err := repository.GetTotalVisits(userID, dc.Pool)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch total visits",
			"details": err.Error(),
		})
		return
	}

	avgClickRate, err := repository.GetAvgClickRate(userID, dc.Pool)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch average click rate",
			"details": err.Error(),
		})
		return
	}

	visitsGrowth, err := repository.GetVisitsGrowth(userID, dc.Pool)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch visits growth",
			"details": err.Error(),
		})
		return
	}

	last7Days, err := repository.GetLast7DaysVisits(userID, dc.Pool)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch last 7 days visits",
			"details": err.Error(),
		})
		return
	}

	response := models.Analytic{
		TotalLinks:   totalLinks,
		TotalVisits:  totalVisits,
		AvgClickRate: avgClickRate,
		VisitsGrowth: visitsGrowth,
		Last7Days:    last7Days,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}