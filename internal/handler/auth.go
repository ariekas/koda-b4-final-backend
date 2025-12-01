package handler

import (
	"shortlink/internal/models"
	"shortlink/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthController struct{
	Pool *pgxpool.Pool
}

func (ac AuthController) Register(ctx *gin.Context){
	user, err := repository.CreateUser(ctx, ac.Pool)

	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success register",
		Data:    user,
	})
}