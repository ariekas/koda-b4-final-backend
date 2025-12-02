CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    userId INT NOT NULL REFERENCES users(id),
    refreshToken TEXT,
    revoked BOOLEAN,
    created_at DATE,
    expires_at DATE, 
    updated_at DATE
)