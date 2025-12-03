package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mssola/user_agent"
	"log"
	"net"
		"io"
		"encoding/json"
	"net/http"
	"shortlink/internal/config"
	"shortlink/internal/models"
	"shortlink/internal/repository"
	"strings"
	"time"
)

type ShortLinkController struct {
	Pool *pgxpool.Pool
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

func getClientIP(ctx *gin.Context) string {
	forwarded := ctx.GetHeader("X-Forwarded-For")
	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	realIP := ctx.GetHeader("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	ip, _, err := net.SplitHostPort(ctx.Request.RemoteAddr)
	if err != nil {
		return ctx.Request.RemoteAddr
	}
	return ip
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

func getLocationFromIP(ipAddress string) (country string, city string) {
	if ipAddress == "::1" || ipAddress == "127.0.0.1" || strings.HasPrefix(ipAddress, "192.168.") || strings.HasPrefix(ipAddress, "10.") {
		return "Local", "Local"
	}

	url := fmt.Sprintf("http://ip-api.com/json/%s", ipAddress)
	
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Failed to get geolocation: %v", err)
		return "", ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read geolocation response: %v", err)
		return "", ""
	}

	var geo models.GeoLocation
	if err := json.Unmarshal(body, &geo); err != nil {
		log.Printf("Failed to parse geolocation: %v", err)
		return "", ""
	}

	return geo.Country, geo.City
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
	keyID := fmt.Sprintf("link:%s:id", slug)
	keyClicks := fmt.Sprintf("link:%s:clicks", slug)

	var shortLinkID int
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

		redis.Set(context.Background(), keyDest, destination, 0)
		redis.Set(context.Background(), keyID, shortLinkID, 0)
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
			redis.Set(context.Background(), keyID, shortLinkID, 0)
		} else {
			fmt.Sscanf(idResult, "%d", &shortLinkID)
		}
	}

	redis.Incr(context.Background(), keyClicks)

	userAgent := ctx.Request.UserAgent()
	referer := ctx.Request.Referer()
	ipAddr := getClientIP(ctx)
	userID := getUserIDFromContext(ctx)


	go func() {
		ua := user_agent.New(userAgent)
		browserName, browserVersion := ua.Browser()
		
		country, city := getLocationFromIP(ipAddr)

		clickData := models.ClickData{
			ShortLinkID: shortLinkID,
			UserID:      userID,
			IPAddress:   ipAddr,
			Referer:     referer,
			UserAgent:   userAgent,
			Country:     country,
			City:        city,
			DeviceType:  getDeviceType(ua),
			Browser:     fmt.Sprintf("%s %s", browserName, browserVersion),
			OS:          ua.OS(),
			CreatedAt:   time.Now(),
		}

		if err := repository.InsertClick(sl.Pool, clickData); err != nil {
			log.Printf("Failed to record click for link %d: %v", shortLinkID, err)
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

	ctx.JSON(201, gin.H{"success": true, "message": "Short link deleted"})
}