package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mssola/user_agent"
	"log"
	"shortlink/internal/config"
	"shortlink/internal/models"
	"shortlink/internal/repository"
	"time"
)

type ShortLinkController struct {
	Pool *pgxpool.Pool
}

func InvalidateUserDashboardCache(userID int) {
	rdb := config.RedisClient
	ctx := context.Background()
	
	statsKey := fmt.Sprintf("user:%d:stats", userID)
	last7DaysKey := fmt.Sprintf("analytics:%d:7d", userID)
	
	rdb.Del(ctx, statsKey, last7DaysKey)
	
	log.Printf("Invalidated dashboard cache for user %d\n", userID)
}

func getUserIDFromContext(ctx *gin.Context) *int {
	userID, exists := ctx.Get("userId")
	if !exists {
		return nil
	}
	fmt.Println("USER", userID)

	if id, ok := userID.(int); ok {
		return &id
	}

	return nil
}

func getDeviceType(ua *user_agent.UserAgent) string {
	if ua.Mobile() {
		return "mobile"
	}
	if ua.Bot() {
		return "bot"
	}
	return "desktop"
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

	InvalidateUserDashboardCache(userId)

	keyRedis := fmt.Sprintf("link:%s:destination", link.OriginalUrl)
	redis.Set(context.Background(), keyRedis, link.OriginalUrl, 0)

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
	redis := config.RedisClient

	links, err := repository.ListLink(sl.Pool, userId)
	if err != nil {
		ctx.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}

	enrichedLinks := make([]map[string]interface{}, 0)
	
	for _, link := range links {
		keyClicks := fmt.Sprintf("link:%s:clicks", link.ShortUrl)
		redisClicks, err := redis.Get(context.Background(), keyClicks).Int64()
		
		var totalClicks int64
		if err == nil {
			totalClicks = redisClicks
		} else {
			clickCount, err := repository.GetClickCount(sl.Pool, link.Id)
			if err == nil {
				totalClicks = clickCount
				redis.Set(context.Background(), keyClicks, totalClicks, 0)
			}
		}

		enrichedLink := map[string]interface{}{
			"id":          link.Id,
			"userId":      link.UserId,
			"originalUrl": link.OriginalUrl,
			"shortUrl":    link.ShortUrl,
			"createdAt":   link.CreatedAt,
			"updatedAt":   link.UpdatedAt,
			"totalClicks": totalClicks, 
		}
		enrichedLinks = append(enrichedLinks, enrichedLink)
	}

	ctx.JSON(200, gin.H{"success": true, "data": enrichedLinks})
}

func (sl ShortLinkController) DetailShortCode(ctx *gin.Context) {
	slug := ctx.Param("slug")
	userId := ctx.GetInt("userId")

	link, err := repository.DetailLink(sl.Pool, slug, userId)
	if err != nil {
		ctx.JSON(404, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"slug":         slug,
			"original_url": link.OriginalUrl,
		},
	})
}


func (sl ShortLinkController) Redirect(ctx *gin.Context) {
	slug := ctx.Param("slug")
	redis := config.RedisClient

	keyDest := fmt.Sprintf("link:%s:destination", slug)
	keyID := fmt.Sprintf("link:%s:id", slug)
	keyClicks := fmt.Sprintf("link:%s:clicks", slug)

	var shortLinkID int
	var linkOwnerID int
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
		shortLinkID = link.Id
		linkOwnerID = link.UserId 

		redis.Set(context.Background(), keyDest, destination, 0)
		redis.Set(context.Background(), keyID, shortLinkID, 0)
		redis.Set(context.Background(), fmt.Sprintf("link:%s:owner", slug), linkOwnerID, 0)
	} else {
		idResult, err := redis.Get(context.Background(), keyID).Result()
		if err != nil {
			link, err := repository.FindShortLink(sl.Pool, slug)
			if err != nil || link == nil {
				ctx.JSON(404, models.Response{
					Success: false,
					Message: "Short link not found",
				})
				return
			}
			shortLinkID = link.Id
			linkOwnerID = link.UserId
			redis.Set(context.Background(), keyID, shortLinkID, 0)
			redis.Set(context.Background(), fmt.Sprintf("link:%s:owner", slug), linkOwnerID, 0)
		} else {
			fmt.Sscanf(idResult, "%d", &shortLinkID)
			
			ownerResult, err := redis.Get(context.Background(), fmt.Sprintf("link:%s:owner", slug)).Result()
			if err != nil {
				link, _ := repository.FindShortLink(sl.Pool, slug)
				if link != nil {
					linkOwnerID = link.UserId
					redis.Set(context.Background(), fmt.Sprintf("link:%s:owner", slug), linkOwnerID, 0)
				}
			} else {
				fmt.Sscanf(ownerResult, "%d", &linkOwnerID)
			}
		}
	}

	redis.Incr(context.Background(), keyClicks)

	userAgent := ctx.Request.UserAgent()
	userID := getUserIDFromContext(ctx)

	go func() {
		ua := user_agent.New(userAgent)
		browserName, browserVersion := ua.Browser()
		

		clickData := models.ClickData{
			ShortLinkID: shortLinkID,
			UserID:      userID,
			IPAddress:   "",
			Referer:     "",
			UserAgent:   userAgent,
			Country:     "",
			City:        "",
			DeviceType:  getDeviceType(ua),
			Browser:     fmt.Sprintf("%s %s", browserName, browserVersion),
			OS:          ua.OS(),
			CreatedAt:   time.Now(),
		}

		if err := repository.InsertClick(sl.Pool, clickData); err != nil {
			log.Printf("Failed to record click for link %d: %v", shortLinkID, err)
		} else {
			if linkOwnerID > 0 {
				InvalidateUserDashboardCache(linkOwnerID)
			}
		}
	}()

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
	redis.Del(context.Background(), fmt.Sprintf("link:%s:id", slug))
	redis.Del(context.Background(), fmt.Sprintf("link:%s:clicks", slug))
	redis.Del(context.Background(), fmt.Sprintf("link:%s:owner", slug))
	
	InvalidateUserDashboardCache(userId)
	
	if body.CustomSlug != nil && *body.CustomSlug != slug {
		redis.Del(context.Background(), fmt.Sprintf("link:%s:destination", *body.CustomSlug))
		redis.Del(context.Background(), fmt.Sprintf("link:%s:id", *body.CustomSlug))
		redis.Del(context.Background(), fmt.Sprintf("link:%s:clicks", *body.CustomSlug))
		redis.Del(context.Background(), fmt.Sprintf("link:%s:owner", *body.CustomSlug))
	}

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
	redis.Del(context.Background(), fmt.Sprintf("link:%s:id", slug))
	redis.Del(context.Background(), fmt.Sprintf("link:%s:clicks", slug))
	redis.Del(context.Background(), fmt.Sprintf("link:%s:owner", slug))
	
	InvalidateUserDashboardCache(userId)

	ctx.JSON(201, gin.H{"success": true, "message": "Short link deleted"})
}