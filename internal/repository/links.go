package repository

import (
    "context"
    "errors"
    "fmt"
    "math/rand"
    "net/url"
    "shortlink/internal/models"
    "time"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

func FindShortLink(pool *pgxpool.Pool, slug string) (*models.ShortLink, error) {
    var link models.ShortLink

    err := pool.QueryRow(context.Background(),
        `SELECT id, userid, originalurl, shorturl, created_at, updated_at
         FROM short_links WHERE shorturl=$1`,
        slug,
    ).Scan(
        &link.Id,
        &link.UserId,
        &link.OriginalUrl,
        &link.ShortUrl,
        &link.CreatedAt,
        &link.UpdatedAt,
    )

    if errors.Is(err, pgx.ErrNoRows) {
        return nil, nil 
    }

    if err != nil {
        return nil, fmt.Errorf("query failed: %w", err)
    }

    return &link, nil
}

func CreateLink(pool *pgxpool.Pool, link models.ShortLink) error {
    now := time.Now()

    _, err := pool.Exec(context.Background(),
        `INSERT INTO short_links (userid, originalurl, shorturl, created_at, updated_at)
         VALUES ($1, $2, $3, $4, $5)`,
        link.UserId,
        link.OriginalUrl,
        link.ShortUrl,
        now,
        now,
    )

    if err != nil {
        return fmt.Errorf("insert failed: %w", err)
    }

    return nil
}

func ValidateURL(raw string) error {
    parsed, err := url.ParseRequestURI(raw)
    if err != nil {
        return fmt.Errorf("invalid URL format: %w", err)
    }

    if parsed.Scheme != "http" && parsed.Scheme != "https" {
        return fmt.Errorf("URL must start with http:// or https://")
    }

    return nil
}
func GenerateLink(n int) string {
    letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
    slug := make([]rune, n)

    for i := range slug {
        slug[i] = letters[rand.Intn(len(letters))]
    }

    return string(slug)
}

func CheckSlug(pool *pgxpool.Pool, slug string) (string, error) {
    exists, err := FindShortLink(pool, slug)
    if err != nil {
        return "", err
    }

    if exists != nil {
        return "", fmt.Errorf("slug already taken")
    }

    return slug, nil
}
func CreateShortLink(pool *pgxpool.Pool, userId int, originalUrl string) (*models.ShortLink, error) {
    if err := ValidateURL(originalUrl); err != nil {
        return nil, err
    }

    slug := GenerateLink(6)

    finalSlug, err := CheckSlug(pool, slug)
    if err != nil {
        return nil, err
    }

    link := models.ShortLink{
        UserId:      userId,
        OriginalUrl: originalUrl,
        ShortUrl:     finalSlug,
    }

    if err := CreateLink(pool, link); err != nil {
        return nil, err
    }

    return &link, nil
}

func ListLink(pool *pgxpool.Pool, userId int) ([]models.ShortLink, error) {
    rows, err := pool.Query(context.Background(), "SELECT id, userid, originalurl, shorturl, created_at, updated_at FROM short_links WHERE userid=$1", userId)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var links []models.ShortLink
    for rows.Next() {
        var l models.ShortLink
        rows.Scan(&l.Id, &l.UserId, &l.OriginalUrl, &l.ShortUrl, &l.CreatedAt, &l.UpdatedAt)
        links = append(links, l)
    }

    return links, nil
}

func DetailLink(pool *pgxpool.Pool, slug string, userId int) (*models.ShortLink, error) {
    var l models.ShortLink
    err := pool.QueryRow(context.Background(),
        "SELECT id, userid, originalurl, shorturl, created_at, updated_at FROM short_links WHERE shorturl=$1 AND userid=$2",
        slug, userId,
    ).Scan(&l.Id, &l.UserId, &l.OriginalUrl, &l.ShortUrl, &l.CreatedAt, &l.UpdatedAt)

    if err != nil {
        return nil, fmt.Errorf("link not found")
    }

    return &l, nil
}

func UpdateLink(pool *pgxpool.Pool, userId int, slug string, originalUrl string, customSlug *string) (*models.ShortLink, error) {
    link, err := DetailLink(pool, slug, userId)
    if err != nil {
        return nil, err
    }

    if originalUrl != "" {
        link.OriginalUrl = originalUrl
    }

    if customSlug != nil && *customSlug != "" {
        link.ShortUrl = *customSlug
    }

    _, err = pool.Exec(context.Background(),
        "UPDATE short_links SET originalurl=$1, shorturl=$2, updated_at=$3 WHERE id=$4",
        link.OriginalUrl, link.ShortUrl, time.Now(), link.Id,
    )

    if err != nil {
        return nil, err
    }

    return link, nil
}