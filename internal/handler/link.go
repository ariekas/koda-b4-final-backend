package handler

import (
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

    link, err := repository.DetailLink(sl.Pool, slug, userId)
    if err != nil {
        ctx.JSON(404, gin.H{"success": false, "message": err.Error()})
        return
    }

    ctx.JSON(201, gin.H{"success": true, "data": link})
}

func (sl ShortLinkController) Redirect(ctx *gin.Context) {
    slug := ctx.Param("slug")

    link, err := repository.FindShortLink(sl.Pool, slug)
    if err != nil {
        ctx.JSON(500, gin.H{
            "success": false,
            "message": "Internal server error",
        })
        return
    }

    if link == nil {
        ctx.JSON(404, gin.H{
            "success": false,
            "message": "Short link not found",
        })
        return
    }

    ctx.Redirect(302, link.OriginalUrl)
}