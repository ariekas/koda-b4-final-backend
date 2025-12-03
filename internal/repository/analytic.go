package repository

import (
	"context"
	"database/sql"
	"shortlink/internal/models"
	"time"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetTotalLinks(userID int, pool *pgxpool.Pool) (int, error) {
	var total int
	err := pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM short_links WHERE userid = $1", userID).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func GetTotalVisits(userID int, pool *pgxpool.Pool) (int, error) {
	var total int
	err := pool.QueryRow(context.Background(), `
		SELECT COUNT(c.id) 
		FROM clicks c
		INNER JOIN short_links sl ON c.shortlinkid = sl.id
		WHERE sl.userid = $1
	`, userID).Scan(&total)
	
	if err != nil {
		return 0, err
	}
	return total, nil
}

func GetAvgClickRate(userID int, pool *pgxpool.Pool) (float64, error) {
	var avgRate sql.NullFloat64
	err := pool.QueryRow(context.Background(), `
		SELECT AVG(click_count) as avg_rate
		FROM (
			SELECT sl.id, COUNT(c.id) as click_count
			FROM short_links sl
			LEFT JOIN clicks c ON sl.id = c.shortlinkid
			WHERE sl.userid = $1
			GROUP BY sl.id
		) as link_clicks
	`, userID).Scan(&avgRate)

	if err != nil {
		return 0, err
	}

	if !avgRate.Valid {
		return 0, nil
	}

	return avgRate.Float64, nil
}

func GetVisitsGrowth(userID int, pool *pgxpool.Pool) (float64, error) {
	var lastWeek, previousWeek int

	err := pool.QueryRow(context.Background(), `
		SELECT COUNT(c.id)
		FROM clicks c
		INNER JOIN short_links sl ON c.shortlinkid = sl.id
		WHERE sl.userid = $1 
		AND c.created_at >= NOW() - INTERVAL '7 days'
	`, userID).Scan(&lastWeek)

	if err != nil {
		return 0, err
	}

	err = pool.QueryRow(context.Background(), `
		SELECT COUNT(c.id)
		FROM clicks c
		INNER JOIN short_links sl ON c.shortlinkid = sl.id
		WHERE sl.userid = $1 
		AND c.created_at >= NOW() - INTERVAL '14 days'
		AND c.created_at < NOW() - INTERVAL '7 days'
	`, userID).Scan(&previousWeek)

	if err != nil {
		return 0, err
	}

	if previousWeek == 0 {
		if lastWeek > 0 {
			return 100.0, nil
		}
		return 0, nil
	}

	growth := float64(lastWeek-previousWeek) / float64(previousWeek) * 100
	return growth, nil
}
func GetLast7DaysVisits(userID int, pool *pgxpool.Pool) ([]models.DayVisit, error) {
	query := `
		SELECT 
			DATE(c.created_at) as visit_date,
			COUNT(c.id) as visit_count
		FROM clicks c
		INNER JOIN short_links sl ON c.shortlinkid = sl.id
		WHERE sl.userid = $1 
		AND c.created_at >= NOW() - INTERVAL '7 days'
		GROUP BY DATE(c.created_at)
		ORDER BY visit_date ASC
	`
	
	rows, err := pool.Query(context.Background(), query, userID)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	visits := make([]models.DayVisit, 0)
	for rows.Next() {
		var visit models.DayVisit
		var date time.Time
		err := rows.Scan(&date, &visit.Count)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		visit.Date = date.Format("2006-01-02")
		visits = append(visits, visit)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return visits, nil
}