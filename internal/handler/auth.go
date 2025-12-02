package handler

import (
	"fmt"
	"shortlink/internal/config"
	"shortlink/internal/middelware"
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

func (ac AuthController) Login(ctx *gin.Context) {
	var input models.InputLogin


	err := ctx.BindJSON(&input)

	jwtToken := config.GetJwtToken()

	if err != nil {
		fmt.Println("Error : Failed type much json")
	}


	users, err := repository.FindUserEmail(ac.Pool, input.Email)
	if err != nil {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if !middelware.VerifPassword(users.Password, input.Password) {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Wrong password",
		})
		return
	}

	token, err := middelware.GenerateToken(jwtToken, users.Role, users.Id)
	if err != nil {
		fmt.Println("Error: Failed to generate token")
	} 

	
	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Login success",
		Data: gin.H{
			"token": token,
			"role":  users.Role,
		},
	})
}