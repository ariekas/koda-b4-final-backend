package handler

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"shortlink/internal/config"
	"shortlink/internal/middelware"
	"shortlink/internal/models"
	"shortlink/internal/repository"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthController struct {
	Pool *pgxpool.Pool
}

func (ac AuthController) Register(ctx *gin.Context) {
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
		UserId:       users.Id,
		RefreshToken: hash,
		CreatedAt:    now,
		ExpiresAt:    now.Add(24 * time.Hour),
		UpdatedAt:    now,
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
			"accessToken":  token,
			"refreshToken": refreshToken,
			"role":         users.Role,
		},
	})
}

func (ac AuthController) Refresh(ctx *gin.Context) {
	var refreshToken models.Session
	middelware.GenerateRefreshToken()
	err := ctx.BindJSON(&refreshToken)

	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "error failed type json",
		})
		return
	}

	if refreshToken.RefreshToken == "" {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Refresh token required",
		})
		return
	}

	hash := sha256.Sum256([]byte(refreshToken.RefreshToken))
	hashToken := base64.StdEncoding.EncodeToString(hash[:])

	session, err := repository.FindSesstionByToken(ac.Pool, hashToken)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	jwtToken := config.GetJwtToken()
	access, err := middelware.GenerateToken(jwtToken, "user", session.UserId)
	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "failed to generate token",
		})
		return
	}

	token, hashRef, err := middelware.GenerateRefreshToken()
	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "failed to generate refresh token",
		})
		return
	}

	err = repository.UpdateSesstion(ac.Pool, session.Id, hashRef)
	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "failed to update sesstion",
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "success to refresh token",
		Data: gin.H{
			"accessToken":  access,
			"refreshToken": token,
		},
	})
}
