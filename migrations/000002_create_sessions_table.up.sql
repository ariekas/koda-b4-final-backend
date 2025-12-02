CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    userId INT NOT NULL REFERENCES users(id),
    refreshToken TEXT,
    created_at DATE,
    expires_at DATE, 
    updated_at DATE
)