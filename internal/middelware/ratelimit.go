package middelware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"shortlink/internal/config"
	"time"
)

func RateLimit(maxRequests int, duration time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		identity := ctx.PostForm("email")

		if identity == "" {
			identity = ctx.ClientIP()
		}

		endpoint := ctx.FullPath()
		key := fmt.Sprintf("ratelimit:%s:%s", identity, endpoint)

		count, err := config.RedisClient.Get(context.Background(), key).Int()
		if err != nil && err.Error() != "redis: nil" {
			ctx.JSON(500, gin.H{"message": "Internal server error"})
			ctx.Abort()
			return
		}

		if count >= maxRequests {
			ctx.JSON(429, gin.H{
				"success": false,
				"message": fmt.Sprintf("Rate limit exceeded, try again in %d seconds", int(duration.Seconds())),
			})
			ctx.Abort()
			return
		}

		tx := config.RedisClient.TxPipeline()
		tx.Incr(context.Background(), key)
		tx.Expire(context.Background(), key, duration)
		_, err = tx.Exec(context.Background())
		if err != nil {
			ctx.JSON(500, gin.H{"message": "Internal server error"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
