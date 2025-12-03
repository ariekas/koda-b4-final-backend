CREATE TABLE pic_user (
    id SERIAL PRIMARY KEY,
    pic TEXT,
    userId INT REFERENCES users(id),
    created_at DATE,
    updated_at DATE
);