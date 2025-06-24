# Movie Shows Discovery Backend

A comprehensive backend API for movie discovery, built with Go. This application provides real-time movie search, trending content, watchlist management, and personalized recommendations.

## 🚀 Features

### Core Features
- **🔍 Search**: Real-time search results from TMDB API
- **📄 Details Page**: Comprehensive movie information from TMDB and OMDB
- **⭐ Trending**: Fetch trending content from TMDB
- **📚 Genre/Category Browsing**: Genre-based content retrieval
- **📊 Ratings Integration**: Combined ratings from TMDB, OMDB, Rotten Tomatoes, IMDB, and Metacritic
- **🧠 Recommendation Engine**: Personalized recommendations based on user watchlist

### Technical Features
- **✅ API Error Handling**: Comprehensive error handling with timeouts and retry logic
- **✅ Pagination Support**: Full pagination support for all endpoints
- **✅ Response Caching**: In-memory caching for improved performance
- **✅ Rate Limiting**: Built-in rate limiting for API protection
- **✅ Secure Configuration**: Environment-based configuration management
- **✅ Data Validation**: Input validation and response sanitization
- **✅ Fallback Handling**: Graceful fallbacks for missing data


## 🛠️ Installation

### Prerequisites
- Go 1.24.4 or higher
- TMDB API key
- OMDB API key

### Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd Movie_Shows_Discovery/backend
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables**
   Create a `.env` file in the backend directory:
   ```env
   # API Keys
   TMDB_API_KEY=your_tmdb_api_key_here
   OMDB_API_KEY=your_omdb_api_key_here
   
   # Server Configuration
   PORT=8080
   HOST=localhost
   
   # Cache Configuration
   CACHE_TTL=3600
   SEARCH_CACHE_TTL=1800
   TRENDING_CACHE_TTL=3600
   
   # Rate Limiting
   TMDB_RATE_LIMIT=40
   OMDB_RATE_LIMIT=1000
   ```

4. **Get API Keys**
   - **TMDB**: Sign up at [themoviedb.org](https://www.themoviedb.org/settings/api)
   - **OMDB**: Sign up at [omdbapi.com](http://www.omdbapi.com/apikey.aspx)

5. **Run the application**
   ```bash
   go run cmd/server/main.go
   ```

The server will start on `http://localhost:8080`

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication
Most endpoints require a user ID in the `X-User-ID` header.

### Key Endpoints

#### Movies
- `GET /movies/search` - Search movies (supports `q`, `page`, `per_page`)
- `GET /movies/{id}` - Get movie details
- `GET /movies/{id}/similar` - Get similar movies
- `GET /movies/genres` - Get all genres
- `GET /movies/genres/{genreId}` - Get movies by genre

#### Trending
- `GET /trending` - Get trending movies
- `GET /trending/by-genre` - Get trending by genre
- `GET /trending/stats` - Get trending statistics
- `GET /trending/genres` - Get trending genres

#### Watchlist
- `POST /watchlist` - Create watchlist
- `GET /watchlist` - Get user's watchlist
- `POST /watchlist/items` - Add movie to watchlist
- `PUT /watchlist/items` - Update watchlist item
- `DELETE /watchlist/items` - Remove from watchlist
- `GET /watchlist/stats` - Get watchlist statistics
- `GET /watchlist/recommendations` - Get recommendations

### Example Requests

#### Search Movies
```bash
curl "http://localhost:8080/api/v1/movies/search?q=inception&page=1&per_page=10"
```
- Supports `q` (query), `page`, and `per_page` (default: 10) parameters.
- Results are paginated on the backend, so you can use `page` and `per_page` to navigate (e.g., for Back/Next buttons in the UI).

#### Get Movie Details
```bash
curl "http://localhost:8080/api/v1/movies/27205"
```

#### Get Trending Movies
```bash
curl "http://localhost:8080/api/v1/trending?timeframe=week&page=1"
```

#### Add to Watchlist
```bash
curl -X POST "http://localhost:8080/api/v1/watchlist/items" \
  -H "X-User-ID: user123" \
  -H "Content-Type: application/json" \
  -d '{"movie_id": 27205, "status": "to_watch", "rating": 0, "notes": "Want to watch this"}'
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `TMDB_API_KEY` | TMDB API key | Required |
| `OMDB_API_KEY` | OMDB API key | Required |
| `PORT` | Server port | 8080 |
| `HOST` | Server host | localhost |
| `CACHE_TTL` | Cache TTL in seconds | 3600 |
| `SEARCH_CACHE_TTL` | Search cache TTL | 1800 |
| `TRENDING_CACHE_TTL` | Trending cache TTL | 3600 |
| `TMDB_RATE_LIMIT` | TMDB requests per second | 40 |
| `OMDB_RATE_LIMIT` | OMDB requests per second | 1000 |

## 🚀 Performance Features

### Caching
- **Search Results**: 30 minutes
- **Trending Content**: 1 hour
- **Movie Details**: 1 hour
- **Genres**: 24 hours

### Rate Limiting
- **API Endpoints**: 100 requests per minute per IP
- **TMDB API**: Configurable rate limiting
- **OMDB API**: Configurable rate limiting

### Error Handling
- **Retry Logic**: Exponential backoff for failed requests
- **Timeout Handling**: Configurable timeouts for all API calls
- **Graceful Degradation**: Fallbacks for missing data

## 🧪 Testing

### Run Tests
```bash
go test ./...
```

### Test Coverage
```bash
go test -cover ./...
```

## 📊 Monitoring

### Health Check
```bash
curl http://localhost:8080/health
```

### Logging
The application provides comprehensive logging:
- Request/response logging
- Error logging with context
- Performance metrics
- API response logging

## 🔒 Security

### Rate Limiting
- Built-in rate limiting to prevent abuse
- Configurable limits per endpoint

### Input Validation
- All inputs are validated and sanitized
- SQL injection protection
- XSS protection

### Error Handling
- No sensitive information in error responses
- Proper HTTP status codes
- Structured error responses

## 🚀 Deployment

### Docker
```dockerfile
FROM golang:1.24.4-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

### Environment Setup
Ensure all required environment variables are set in your deployment environment.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📝 License

This project is licensed under the MIT License.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Check the API documentation
- Review the configuration guide

## 🔄 Updates

### Version History
- **v1.0.0**: Initial release with core features
- **v1.1.0**: Added recommendation engine
- **v1.2.0**: Enhanced caching and performance

