CREATE TABLE clicks (
    id SERIAL PRIMARY KEY,
    shortLinkId INT NOT NULL REFERENCES short_links(id) ON DELETE CASCADE,
    userId INT REFERENCES users(id),       
    ipAddress VARCHAR(50),
    referer TEXT,
    userAgent TEXT,
    country VARCHAR(100),
    city VARCHAR(100),
    deviceType VARCHAR(50),              
    browser VARCHAR(100),
    os VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
