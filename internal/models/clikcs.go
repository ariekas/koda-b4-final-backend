package models

import "time"


type ClickData struct {
	ID          int       `json:"id"`
	ShortLinkID int       `json:"shortLinkId"`
	UserID      *int      `json:"userId,omitempty"` 
	IPAddress   string    `json:"ipAddress"`
	Referer     string    `json:"referer"`
	UserAgent   string    `json:"userAgent"`
	Country     string    `json:"country"`
	City        string    `json:"city"`
	DeviceType  string    `json:"deviceType"`
	Browser     string    `json:"browser"`
	OS          string    `json:"os"`
	CreatedAt   time.Time `json:"createdAt"`
}

type GeoLocation struct {
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	Region      string `json:"region"`
	RegionName  string `json:"regionName"`
	City        string `json:"city"`
	Zip         string `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string `json:"timezone"`
	ISP         string `json:"isp"`
	Query       string `json:"query"`
}