CREATE TABLE short_links (
    id SERIAL PRIMARY KEY,
    userId INT REFERENCES users(id),
    originalUrl TEXT,
    shortUrl TEXT,
    created_at DATE,
    updated_at DATE
)