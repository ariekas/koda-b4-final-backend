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
