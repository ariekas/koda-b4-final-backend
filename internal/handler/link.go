package handler

import (
	"context"
	"fmt"
	"shortlink/internal/config"
	"shortlink/internal/models"
	"shortlink/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShortLinkController struct {
	Pool *pgxpool.Pool
}

func (slc ShortLinkController) Create(ctx *gin.Context) {
	var input models.ShortLink

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "failed to parse JSON",
		})
		return
	}

	userId := ctx.GetInt("userId")

	link, err := repository.CreateShortLink(slc.Pool, userId, input.OriginalUrl)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	redis := config.RedisClient
	redis.Del(context.Background(), fmt.Sprintf("link:list:%d", userId))

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Short link created",
		Data: gin.H{
			"original_url": link.OriginalUrl,
			"short_url":    link.ShortUrl,
		},
	})
}

func (sl ShortLinkController) GetAll(ctx *gin.Context) {
	userId := ctx.GetInt("userId")

	links, err := repository.ListLink(sl.Pool, userId)
	if err != nil {
		ctx.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{"success": true, "data": links})
}

func (sl ShortLinkController) DetailShortCode(ctx *gin.Context) {
	slug := ctx.Param("slug")
	userId := ctx.GetInt("userId")
	redis := config.RedisClient

	keyRedis := fmt.Sprintf("link:%s:destination", slug)

	cacheVal, err := redis.Get(context.Background(), keyRedis).Result()
    fmt.Println(cacheVal)

	if err == nil {
		ctx.JSON(200, gin.H{
			"success": true,
			"source":  "cache",
			"data": gin.H{
				"slug":         slug,
				"original_url": cacheVal,
			},
		})
		return
	}

	link, err := repository.DetailLink(sl.Pool, slug, userId)
	if err != nil {
		ctx.JSON(404, gin.H{"success": false, "message": err.Error()})
		return
	}

	redis.Set(context.Background(), keyRedis, link.OriginalUrl, 0)

	ctx.JSON(201, gin.H{"success": true, "data": link})
}

func (sl ShortLinkController) Redirect(ctx *gin.Context) {
	slug := ctx.Param("slug")
	redis := config.RedisClient

	keyDest := fmt.Sprintf("link:%s:destination", slug)
	keyClicks := fmt.Sprintf("link:%s:clicks", slug)

	destination, err := redis.Get(context.Background(), keyDest).Result()

	if err != nil {
		link, err := repository.FindShortLink(sl.Pool, slug)
		if err != nil || link == nil {
			ctx.JSON(404, models.Response{
				Success: false,
				Message: "Short link not found",
			})
			return
		}
		destination = link.OriginalUrl
		redis.Set(context.Background(), keyDest, destination, 0)
	}
	redis.Incr(context.Background(), keyClicks)

	ctx.Redirect(302, destination)
}

func (sl ShortLinkController) Update(ctx *gin.Context) {
	slug := ctx.Param("slug")
	userId := ctx.GetInt("userId")

	var body struct {
		OriginalUrl string  `json:"originalUrl"`
		CustomSlug  *string `json:"customSlug"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(401, gin.H{"success": false, "message": "invalid request body"})
		return
	}

	link, err := repository.UpdateLink(sl.Pool, userId, slug, body.OriginalUrl, body.CustomSlug)
	if err != nil {
		ctx.JSON(401, gin.H{"success": false, "message": err.Error()})
		return
	}

	redis := config.RedisClient
	redis.Del(context.Background(), fmt.Sprintf("link:%s:destination", slug))
	redis.Del(context.Background(), fmt.Sprintf("link:%s:clicks", slug))
	ctx.JSON(201, gin.H{"success": true, "message": "Short link updated", "data": link})
}

func (sl ShortLinkController) Delete(ctx *gin.Context) {
	slug := ctx.Param("slug")
	userId := ctx.GetInt("userId")

	err := repository.DeleteLink(sl.Pool, userId, slug)
	if err != nil {
		ctx.JSON(401, gin.H{"success": false, "message": err.Error()})
		return
	}

	redis := config.RedisClient

	redis.Del(context.Background(), fmt.Sprintf("link:%s:destination", slug))
	redis.Del(context.Background(), fmt.Sprintf("link:%s:clicks", slug))

	ctx.JSON(201, gin.H{"success": true, "message": "Short link deleted"})
}
