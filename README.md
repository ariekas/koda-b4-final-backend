# Backend - Shortlink Service

A URL shortening service built with Go, PostgreSQL, and Redis for caching.

## Tech Stack

- **Go**: 1.24.10
- **Database**: PostgreSQL
- **Cache**: Redis
- **Framework**: Gin

## Project Structure

```
shortlink/
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database connection & setup
│   ├── handler/         # HTTP request handlers
│   ├── middleware/      # HTTP middlewares
│   ├── models/          # Data models
│   ├── repository/      # Database operations
│   └── routes/          # Route definitions
├── migrations/          # Database migrations
├── uploads/
│   └── profiles/        # User profile uploads
├── .env                 # Environment variables
├── .gitignore
├── Dockerfile
├── go.mod
├── go.sum
└── main.go             # Application entry point
```

## Prerequisites

- Go 1.24.10 or higher
- PostgreSQL 12 or higher
- Redis 6 or higher

## Environment Variables

Create a `.env` file in the root directory:

```env
DATABASE_URL= # DATABASE URL
JWT_TOKEN= #JWT TOKEN

REDIS_URL= $ REDIS URL

ORIGIN_URL= # URL FRONT END
```

## Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd shortlink
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up PostgreSQL database**
   ```bash
   # Create database
   createdb shortlink_db
   ```

4. **Set up Redis**
   ```bash
   # Make sure Redis is running
   redis-cli ping
   # Should return: PONG
   ```

## Database Migrations

Run migrations to set up the database schema:

```bash
# Run all migrations
go run main.go migrate up

# Rollback last migration
go run main.go migrate down

# Check migration status
go run main.go migrate status
```

Or manually run migration files from the `migrations/` directory.

## Running the Application

### Development Mode

```bash
go run main.go
```

### Production Mode

```bash
# Build the application
go build -o shortlink main.go

# Run the binary
./shortlink
```

### Using Docker

```bash
# Build Docker image
docker build -t shortlink-backend .

# Run with Docker Compose
docker-compose up -d
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/logout` - User logout
- `POST /api/v1/auth/refres` - Refresh Token

### Dashboard
- `GET /api/v1/dashboard/stats` - Get dashboard data (requires auth)

### Home/Shortlinks
- `GET /api/v1/links` - Get all shortlinks (requires auth)
- `POST /api/v1/links` - Create new shortlink (requires auth)
- `GET /api//v1/links/:slug` - Get shortlink by slug (requires auth)
- `PATCH /api/v1/links/:slug` - Update shortlink (requires auth)
- `DELETE /api/v1/links/:slug` - Delete shortlink (requires auth)
- `GET /:slug` - Redirect to original URL (requires auth)

### User Profile
- `GET /api/v1/users/profile` - Get user profile (requires auth)
- `POST /api/v1/users/pic` - Upload profile picture (requires auth)

## Testing Endpoints

### Using cURL

**Register User:**
```bash
curl -X POST http://localhost:8082/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "username": "John Doe"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

**Create Shortlink:**
```bash
curl -X POST http://localhost:8080/api/v1/links \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  -d '{
    "originalUrl": "https://example.com/very-long-url",
  }'
```

**Access Shortlink:**
```bash
curl -L http://localhost:8080/mylink
```

## Redis Flushing Mechanism

The application uses Redis for caching shortlink data to improve performance.

### Manual Flush

**Flush all cache:**
```bash
redis-cli FLUSHALL
```

**Flush specific database:**
```bash
redis-cli -n 0 FLUSHDB
```

**Flush specific key:**
```bash
redis-cli DEL shortlink:<code>
```

### Automatic Cache Invalidation

The application automatically invalidates cache in these scenarios:
- When a shortlink is updated
- When a shortlink is deleted
- After a configurable TTL (Time To Live) expires

### Cache Keys Pattern

```
user:<id>:profile          # cached user profile
user:<id>:stats            # user statistics

link:<code>:destination    # original URL
link:<code>:clicks         # total clicks

analytics:user:<id>        # click analytics by user
analytics:link:<code>      # click analytics per link

device:<identity>:android  # device tracking
```
## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
