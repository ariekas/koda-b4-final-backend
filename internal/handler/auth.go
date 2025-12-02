package handler

import (
	"fmt"
	"shortlink/internal/config"
	"shortlink/internal/middelware"
	"shortlink/internal/models"
	"shortlink/internal/repository"
	"time"

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
	now := time.Now()

	err := ctx.BindJSON(&input)
	if err != nil {
		fmt.Println("Error : Failed type much json")
	}

	jwtToken := config.GetJwtToken()

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

	refreshToken, hash, err := middelware.GenerateRefreshToken()
	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	session := models.Session{
		UserId: users.Id,
		RefreshToken: hash,
		CreatedAt: now,
		ExpiresAt: now.Add(24 * time.Hour),
		UpdatedAt: now,
	}

	err = repository.SaveSession(ac.Pool, session)
	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "failed to save sesstion",
		})
		return
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Login success",
		Data: gin.H{
			"accessToken": token,
			"refreshToken" : refreshToken,
			"role":  users.Role,
		},
	})
}