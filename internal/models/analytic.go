package models

type Analytic struct {
	TotalLinks   int                `json:"total_links"`
	TotalVisits  int                `json:"total_visits"`
	AvgClickRate float64            `json:"avg_click_rate"`
	VisitsGrowth float64            `json:"visits_growth"`
	Last7Days    []DayVisit  `json:"last_7_days_chart"`
}

type DayVisit struct {
	Date  string `json:"date"`
	Count int    `json:"visitCount"`
}