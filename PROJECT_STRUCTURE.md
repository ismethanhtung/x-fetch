# 📁 Project Structure - Cấu trúc Dự án

Tài liệu này mô tả chi tiết cấu trúc và chức năng của từng file/folder trong project.

## 📂 Cấu trúc Tổng Quát

```
x-twitter-backend/
├── config/                 # Configuration management
│   └── config.go          # Load và quản lý environment variables
│
├── handlers/              # HTTP request handlers
│   ├── tweets_handler.go  # Handlers cho tweets và user endpoints
│   └── middleware.go      # Logging, CORS, Recovery middlewares
│
├── models/                # Data structures
│   └── tweet.go          # Models cho Tweet, User, Response types
│
├── services/              # Business logic layer
│   └── twitter_service.go # Service để tương tác với Twitter API
│
├── main.go               # Entry point của application
├── go.mod                # Go module dependencies
├── go.sum                # Dependencies checksums
├── Makefile              # Build và development commands
├── Dockerfile            # Docker image configuration
├── docker-compose.yml    # Docker Compose setup
├── .air.toml             # Air hot reload configuration
├── .gitignore            # Git ignore rules
├── ENV_EXAMPLE           # Environment variables template
│
└── Documentation/
    ├── README.md         # Main documentation (đầy đủ nhất)
    ├── QUICKSTART_VI.md  # Quick start guide (5 phút)
    ├── TUTORIAL_VI.md    # Hướng dẫn chi tiết từng bước
    ├── EXAMPLES.md       # Code examples (nhiều ngôn ngữ)
    └── PROJECT_STRUCTURE.md  # File này

```

## 📄 Chi Tiết Từng File

### Root Level Files

#### `main.go`
**Mục đích**: Entry point của application, khởi tạo và chạy server.

**Chức năng chính**:
- Setup logging configuration
- Load configuration từ environment
- Khởi tạo Twitter service
- Setup HTTP router với các routes
- Start HTTP server
- Handle graceful shutdown

**Key functions**:
- `main()` - Entry point
- `setupRouter()` - Configure routes
- `setupLogging()` - Setup log format và level
- `gracefulShutdown()` - Handle SIGINT/SIGTERM
- `handleAPIDocs()` - API documentation endpoint

#### `go.mod`
**Mục đích**: Define Go module và dependencies.

**Dependencies**:
- `github.com/gorilla/mux` - HTTP router
- `github.com/michimani/gotwi` - Twitter API v2 client
- `github.com/joho/godotenv` - Environment variables loader
- `github.com/sirupsen/logrus` - Structured logging

#### `go.sum`
**Mục đích**: Checksums của dependencies để đảm bảo integrity.

#### `Makefile`
**Mục đích**: Automation commands cho development và deployment.

**Available commands**:
```bash
make help           # Hiển thị tất cả commands
make install        # Cài đặt dependencies
make build          # Build binary
make run            # Chạy application
make dev            # Hot reload với air
make test           # Chạy tests
make test-coverage  # Test với coverage report
make clean          # Xóa build artifacts
make lint           # Chạy linter
make fmt            # Format code
make vet            # Chạy go vet
make docker-build   # Build Docker image
make docker-run     # Run Docker container
```

#### `Dockerfile`
**Mục đích**: Multi-stage Docker build cho production deployment.

**Stages**:
1. **Builder stage**: Build Go binary
2. **Runtime stage**: Minimal Alpine image với binary

**Features**:
- Small image size (< 20MB)
- Built-in health check
- Non-root user
- Ca-certificates included

#### `docker-compose.yml`
**Mục đích**: Orchestrate Docker containers.

**Configuration**:
- Auto-restart policy
- Environment variables từ .env
- Health checks
- Network isolation

#### `.air.toml`
**Mục đích**: Configuration cho air hot reload tool.

**Features**:
- Auto rebuild khi code thay đổi
- Exclude test files và tmp directories
- Build error logging

#### `.gitignore`
**Mục đích**: Exclude files khỏi git tracking.

**Ignores**:
- Binaries (*.exe, *.dll, *.so)
- Dependencies (vendor/)
- Environment files (.env)
- IDE files (.vscode/, .idea/)
- OS files (.DS_Store)
- Logs (*.log, logs/)

#### `ENV_EXAMPLE`
**Mục đích**: Template cho environment variables.

**Variables**:
- `TWITTER_BEARER_TOKEN` - Twitter API Bearer Token (required)
- `SERVER_PORT` - Server port (default: 8080)
- `SERVER_HOST` - Server host (default: localhost)
- `APP_ENV` - Environment (development/production)
- `LOG_LEVEL` - Log level (debug/info/warn/error)
- `MAX_TWEETS_PER_REQUEST` - Max tweets limit
- `DEFAULT_TWEETS_COUNT` - Default tweets count

---

### `/config` Directory

#### `config/config.go`
**Mục đích**: Configuration management và environment loading.

**Main struct**:
```go
type Config struct {
    TwitterBearerToken  string
    ServerPort          string
    ServerHost          string
    AppEnv              string
    LogLevel            string
    MaxTweetsPerRequest int
    DefaultTweetsCount  int
}
```

**Functions**:
- `LoadConfig()` - Load và validate configuration
- `getEnv()` - Get environment variable với default value
- `getEnvAsInt()` - Get integer environment variable
- `GetAddress()` - Return server address (host:port)

**Validation**:
- Validate required fields (TWITTER_BEARER_TOKEN)
- Set default values cho optional fields
- Parse và validate integer values

---

### `/models` Directory

#### `models/tweet.go`
**Mục đích**: Define data structures cho API.

**Main structs**:

1. **Tweet** - Đại diện một tweet
   - ID, Text, AuthorID, CreatedAt
   - Metrics (likes, retweets, replies, quotes)
   - Entities (hashtags, mentions, URLs)
   - ReferencedTweets

2. **User** - Thông tin user
   - ID, Username, Name, Description
   - ProfileImageURL, Verified
   - UserMetrics (followers, following, tweets)

3. **TweetsResponse** - API response structure
   - Tweets array
   - User info
   - Meta (count, pagination)

4. **ErrorResponse** - Error response structure
   - Error code
   - Message
   - HTTP status code

**Design patterns**:
- JSON tags cho serialization
- Omitempty cho optional fields
- Consistent naming conventions

---

### `/services` Directory

#### `services/twitter_service.go`
**Mục đích**: Business logic layer, tương tác với Twitter API.

**Main struct**:
```go
type TwitterService struct {
    client *gotwi.Client
    config *config.Config
}
```

**Public methods**:

1. **NewTwitterService()** - Constructor
   - Initialize Twitter client
   - Setup authentication
   - Return service instance

2. **GetUserByUsername()** - Lấy thông tin user
   - Input: username
   - Output: *models.User
   - Fetch user data với all fields

3. **GetUserTweets()** - Lấy tweets của user
   - Input: username, maxResults
   - Output: *models.TweetsResponse
   - Fetch user info + tweets
   - Return complete response

4. **GetTweetsByUserID()** - Lấy tweets theo user ID
   - Input: userID, maxResults
   - Output: []models.Tweet
   - Direct fetch by ID

**Private methods**:

1. **convertToUser()** - Convert Twitter API user data sang models.User
2. **convertToTweet()** - Convert Twitter API tweet data sang models.Tweet

**Features**:
- Comprehensive error handling
- Field validation
- Rate limit awareness
- Detailed logging
- Null safety checks

---

### `/handlers` Directory

#### `handlers/tweets_handler.go`
**Mục đích**: HTTP request handlers cho API endpoints.

**Main struct**:
```go
type TweetsHandler struct {
    twitterService *services.TwitterService
}
```

**HTTP Handlers**:

1. **GetUserTweets()** - `GET /api/tweets/user/{username}`
   - Query params: count (optional)
   - Response: TweetsResponse
   - Error handling với proper status codes

2. **GetUserInfo()** - `GET /api/user/{username}`
   - Path param: username
   - Response: User
   - 404 nếu user không tồn tại

3. **HealthCheck()** - `GET /health`
   - No params
   - Response: Status object
   - Always returns 200 OK

**Helper methods**:
- `respondWithJSON()` - Send JSON response
- `respondWithError()` - Send error response

**Features**:
- Request validation
- Logging với context
- Consistent error responses
- Content-Type headers

#### `handlers/middleware.go`
**Mục đích**: HTTP middlewares cho cross-cutting concerns.

**Middlewares**:

1. **LoggingMiddleware** - Log tất cả HTTP requests
   - Log method, path, status, duration
   - Log IP address và User-Agent
   - Structured logging với logrus

2. **CORSMiddleware** - Handle CORS
   - Allow all origins (configurable)
   - Allow common methods
   - Handle preflight requests

3. **RecoveryMiddleware** - Recover từ panics
   - Catch panics
   - Log error
   - Return 500 response

**responseWriter wrapper**:
- Capture status code từ handlers
- Used by LoggingMiddleware

---

## 🔄 Request Flow

```
Client Request
    ↓
RecoveryMiddleware (catch panics)
    ↓
LoggingMiddleware (log request)
    ↓
CORSMiddleware (add CORS headers)
    ↓
Router (gorilla/mux)
    ↓
Handler (tweets_handler.go)
    ↓
Service (twitter_service.go)
    ↓
Twitter API (via gotwi)
    ↓
Response (JSON)
```

## 🗂️ Data Flow

```
Twitter API Response
    ↓
gotwi Client (parse response)
    ↓
TwitterService (convert to models)
    ↓
Handler (add meta info)
    ↓
JSON Encoder
    ↓
HTTP Response
```

## 🏗️ Architecture Patterns

### 1. Layered Architecture
- **Handlers Layer**: HTTP concerns
- **Services Layer**: Business logic
- **Models Layer**: Data structures
- **Config Layer**: Configuration

### 2. Dependency Injection
- Services injected vào handlers
- Config injected vào services
- Easy testing và mocking

### 3. Single Responsibility
- Mỗi package có một responsibility rõ ràng
- Separation of concerns
- Easy maintenance

### 4. Error Handling
- Errors bubble up từ service → handler
- Proper HTTP status codes
- Structured error responses

## 📦 Package Dependencies

```
main
├── config
├── handlers
│   └── services
│       └── models
└── models
```

**Dependency rules**:
- `main` depends on tất cả
- `handlers` depends on `services` và `models`
- `services` depends on `models` và `config`
- `models` có no dependencies (pure data)
- `config` có no dependencies (except external)

## 🔐 Security Considerations

### Environment Variables
- Sensitive data trong .env (git ignored)
- Validation khi load
- No hardcoded secrets

### API Security
- Bearer Token authentication với Twitter
- CORS configured
- No data exposure trong errors

### Input Validation
- Validate user inputs
- Sanitize parameters
- Prevent injection attacks

## 🧪 Testing Strategy

### Unit Tests (TODO)
- Test services independently
- Mock Twitter API client
- Test error scenarios

### Integration Tests (TODO)
- Test full request flow
- Test với real API (optional)
- Test error handling

### E2E Tests (TODO)
- Test complete workflows
- Test với Docker container
- Test deployment

## 📊 Monitoring & Observability

### Logging
- Structured logs với logrus
- Request/response logging
- Error logging với context

### Health Checks
- `/health` endpoint
- Docker health check
- Ready for monitoring tools

### Metrics (TODO - Future)
- Request count
- Response time
- Error rate
- Twitter API usage

## 🚀 Deployment Options

### 1. Direct Binary
- Build với `go build`
- Run binary trực tiếp
- Simple, no overhead

### 2. Systemd Service
- Auto-restart
- Log management
- Production-ready

### 3. Docker Container
- Isolated environment
- Easy scaling
- Consistent deployment

### 4. Docker Compose
- Multi-container setup
- Easy configuration
- Development parity

### 5. Kubernetes (Future)
- Horizontal scaling
- Load balancing
- Auto-healing

## 📈 Scalability Considerations

### Current Bottlenecks
- Twitter API rate limits
- Single instance design

### Scaling Options
1. **Horizontal Scaling**
   - Deploy multiple instances
   - Add load balancer
   - Share cache layer

2. **Caching Layer**
   - Add Redis cache
   - Cache user info
   - Cache recent tweets

3. **Database**
   - Store tweets locally
   - Reduce API calls
   - Enable offline queries

4. **Message Queue**
   - Async tweet fetching
   - Handle spikes
   - Background processing

## 🔧 Configuration Management

### Development
- Use .env file
- Hot reload với air
- Debug logging

### Production
- Environment variables
- Secrets manager
- Info logging
- HTTPS enabled

## 📝 Code Style

### Conventions
- Go standard formatting
- Clear naming conventions
- Comprehensive comments
- Error handling everywhere

### Best Practices
- Small, focused functions
- DRY (Don't Repeat Yourself)
- SOLID principles
- Clean code principles

---

## 🎯 Quick Reference

### Add New Endpoint

1. Define model trong `models/tweet.go` (nếu cần)
2. Add service method trong `services/twitter_service.go`
3. Add handler trong `handlers/tweets_handler.go`
4. Register route trong `main.go` setupRouter()
5. Update documentation

### Add New Middleware

1. Create middleware function trong `handlers/middleware.go`
2. Apply middleware trong `main.go` setupRouter()

### Change Configuration

1. Add variable trong `config/config.go`
2. Add to `ENV_EXAMPLE`
3. Update README.md

---

**Đây là một project được thiết kế cẩn thận, dễ hiểu, dễ maintain và dễ scale!** 🚀

