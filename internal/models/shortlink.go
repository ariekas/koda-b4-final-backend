package models

import "time"

type ShortLink struct {
	Id          int       `json:"id"`
	UserId      int       `json:"userId"`
	OriginalUrl string    `json:"originalUrl"`
	ShortUrl    string    `json:"shortUrl"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	TotalClicks int64     `json:"totalClicks,omitempty"` 
}
