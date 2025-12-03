package handler

import (
	"shortlink/internal/models"
	"time"

	"fmt"
	"path/filepath"
	"shortlink/internal/repository"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserController struct {
	Pool *pgxpool.Pool
}

func (uc UserController) GetProfile(ctx *gin.Context) {
	userIdInterface, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	userId := userIdInterface.(int)

	user, err := repository.GetUserWithPic(uc.Pool, userId)
	if err != nil {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Success get profile",
		Data:    user,
	})
}

func (uc UserController) GetUserById(ctx *gin.Context) {
	userIdStr := ctx.Param("id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	user, err := repository.GetUserWithPic(uc.Pool, userId)
	if err != nil {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Success get user",
		Data:    user,
	})
}

func (uc UserController) UploadProfilePicture(ctx *gin.Context) {
	userIdInterface, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	userId := userIdInterface.(int)

	file, err := ctx.FormFile("pic")
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "No file uploaded",
		})
		return
	}

	ext := filepath.Ext(file.Filename)
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true}
	if !allowedExts[ext] {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid file type. Only jpg, jpeg, png, gif allowed",
		})
		return
	}

	if file.Size > 5*1024*1024 {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "File size too large. Max 5MB",
		})
		return
	}

	filename := fmt.Sprintf("user_%d_%d%s", userId, time.Now().Unix(), ext)
	uploadPath := filepath.Join("uploads", "profiles", filename)

	err = ctx.SaveUploadedFile(file, uploadPath)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to save file",
		})
		return
	}

	err = repository.UpdateUserPic(uc.Pool, userId, filename)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Success upload profile picture",
		Data: gin.H{
			"pic": filename,
		},
	})
}
