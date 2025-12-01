CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(225),
    email VARCHAR(100),
    password VARCHAR(8),
    created_at DATE,
    updated_at DATE
);