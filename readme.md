# Comic Scraper API

A robust Golang-based web scraping service that extracts comic data from ReadComicOnline.li using go-rod, manages scraping jobs with Redis/AsynQ, stores data in PostgreSQL, and delivers results via webhooks.

## Features

- 🕷️ **Web Scraping**: Automated comic extraction from ReadComicOnline.li using go-rod
- 🗄️ **Database Storage**: PostgreSQL integration for persistent comic data storage
- ⚡ **Job Management**: Redis-backed job queue with AsynQ for scalable task processing
- 🔗 **Webhook Integration**: Automatic result delivery to specified endpoints
- 📊 **Real-time Monitoring**: Track scraping progress with detailed logging
- 🐳 **Docker Support**: Full containerization with Docker Compose
- 📚 **API Documentation**: Interactive Swagger UI available

## API Endpoints

### 1. Get Scraping Status
```
GET /api/v1/logger/:id
```

Retrieve the current status and progress of a scraping job by its log ID.

**Response:**
```json
{
  "ID": 0,
  "console": [
    "string"
  ],
  "hasInfo": true,
  "processedEpisodes": 0,
  "processedFiles": 0,
  "status": 0,
  "timeEstimated": 0,
  "totalEpisodes": 0,
  "totalFiles": 0,
  "webhookSend": true
}
```

### 2. Start Scraping Job
```
POST /api/v1/scrapper/request
```

Initiate a new comic scraping job with specified pages and webhook destination.

**Request:**
```json
{
  "authorization": "string",
  "pages": [
    "string"
  ],
  "webhookUrl": "string"
}
```

**Response:**
```json
[
  {
    "logId": 0
  }
]
```

## Quick Start

### Prerequisites

- Go 1.19+ (for local development)
- Docker and Docker Compose (for containerized deployment)
- Chrome browser (for local development only)

### Option 1: Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd comic-scraper
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Setup environment**
   ```bash
   cp .env.example .env
   # Edit .env file with your configuration
   ```

4. **Install Chrome browser**
   - Ensure Google Chrome is installed on your system
   - The application uses go-rod which requires Chrome for web scraping

5. **Run the application**
   ```bash
   go build ./cmd/app/main.go
   ./main
   ```

### Option 2: Docker Deployment (Recommended)

1. **Clone and start services**
   ```bash
   git clone <repository-url>
   cd comic-scraper
   docker compose up --build -d
   ```

2. **Verify deployment**
   ```bash
   docker compose ps
   ```

The application will be available at `http://localhost:8080` (or your configured port).

## Configuration

The application uses environment variables for configuration. Copy `.env.example` to `.env` and customize as needed:

```env
# Server Configuration
PORT=8080
GIN_MODE=release

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=comic_scraper
DB_SSLMODE=disable

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Scraping Configuration
CHROME_PATH=/usr/bin/google-chrome
SCRAPE_DELAY=1000
MAX_CONCURRENT_JOBS=5

# Webhook Configuration
WEBHOOK_TIMEOUT=30
WEBHOOK_RETRIES=3
```

## Docker Services

The Docker Compose setup includes:

- **App**: Main application container
- **PostgreSQL**: Database service
- **Redis**: Job queue and caching
- **Chrome**: Headless browser service for scraping

```yaml
# Key services in docker-compose.yml
services:
  app:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - postgres
      - redis
      - chrome
  
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: comic_scraper
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
  
  redis:
    image: redis:7-alpine
  
  chrome:
    image: chromedp/headless-shell:latest
```

## API Documentation

Interactive API documentation is available via Swagger UI:

```
http://localhost:8080/swagger/index.html
```

## Usage Examples

### Start a Scraping Job

```bash
curl -X POST http://localhost:8080/api/v1/scrapper/request \
  -H "Content-Type: application/json" \
  -d '{
    "authorization": "your-auth-token",
    "pages": [
      "https://readcomiconline.li/Comic/Batman",
      "https://readcomiconline.li/Comic/Superman"
    ],
    "webhookUrl": "https://your-webhook-endpoint.com/receive"
  }'
```

### Check Job Status

```bash
curl http://localhost:8080/api/v1/logger/123
```

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Client App    │───▶│   Comic API     │───▶│   PostgreSQL    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Webhook       │◀───│   AsynQ Worker  │───▶│     Redis       │
│   Destination   │    └─────────────────┘    └─────────────────┘
└─────────────────┘              │
                                 ▼
                       ┌─────────────────┐
                       │   go-rod        │
                       │   (Chrome)      │
                       └─────────────────┘
```

## Development

### Project Structure

```
.
├── cmd/
│   └── app/
│       └── main.go          # Application entry point
├── internal/
│   ├── api/                 # REST API handlers
│   ├── config/              # Configuration management
│   ├── database/            # Database models and migrations
│   ├── scraper/             # Web scraping logic
│   ├── queue/               # Job queue management
│   └── webhook/             # Webhook delivery
├── docs/                    # Swagger documentation
├── docker-compose.yml       # Docker services
├── Dockerfile              # Container build instructions
├── .env.example            # Environment template
├── go.mod                  # Go modules
└── README.md              # This file
```

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
CGO_ENABLED=0 GOOS=linux go build -o comic-scraper ./cmd/app/main.go
```

## Monitoring and Logging

The application provides comprehensive logging for monitoring scraping jobs:

- Real-time console output via WebSocket or polling
- Progress tracking (episodes processed, files downloaded)
- Error reporting and retry mechanisms
- Estimated completion times
- Webhook delivery status

## Troubleshooting

### Common Issues

1. **Chrome not found**: Ensure Chrome is installed or use Docker deployment
2. **Database connection failed**: Verify PostgreSQL is running and credentials are correct
3. **Redis connection failed**: Check Redis service status and configuration
4. **Scraping timeouts**: Adjust `SCRAPE_DELAY` and timeout settings

### Logs

Check application logs:
```bash
# Docker deployment
docker compose logs app

# Local deployment
./main 2>&1 | tee app.log
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For support and questions:
- Open an issue on GitHub
- Check the [Swagger documentation](http://localhost:8080/swagger/index.html)
- Review the application logs for detailed error information