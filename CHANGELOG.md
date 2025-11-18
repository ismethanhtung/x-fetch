# Changelog

Tất cả các thay đổi đáng chú ý của project này sẽ được ghi lại trong file này.

Format dựa trên [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
và project này tuân theo [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-01-15

### 🎉 Initial Release

#### Added
- ✅ Server backend hoàn chỉnh bằng Golang
- ✅ Integration với Twitter API v2
- ✅ Endpoint lấy thông tin user: `GET /api/user/{username}`
- ✅ Endpoint lấy tweets: `GET /api/tweets/user/{username}`
- ✅ Health check endpoint: `GET /health`
- ✅ API documentation endpoint: `GET /api/docs`
- ✅ Comprehensive logging với logrus
- ✅ CORS middleware
- ✅ Recovery middleware (panic handling)
- ✅ Graceful shutdown
- ✅ Environment-based configuration
- ✅ Docker support với Dockerfile
- ✅ Docker Compose configuration
- ✅ Makefile với automation commands
- ✅ Hot reload support với Air
- ✅ Structured project với clean architecture

#### Documentation
- ✅ README.md - Full documentation
- ✅ QUICKSTART_VI.md - Quick start guide
- ✅ TUTORIAL_VI.md - Detailed tutorial
- ✅ EXAMPLES.md - Code examples (cURL, JS, Python, Go, PHP, Java, Shell)
- ✅ PROJECT_STRUCTURE.md - Architecture documentation
- ✅ CHANGELOG.md - Version history

#### Features
- 🔐 Secure Bearer Token authentication
- 📊 Rich tweet data với metrics
- 👤 Detailed user information
- 🏷️ Hashtags, mentions, URLs extraction
- ⚡ High performance với Go
- 🛡️ Comprehensive error handling
- 📝 Structured logging
- 🌐 CORS enabled
- 🐳 Docker ready
- 🔄 Auto-reload trong development

### Technical Details

#### Dependencies
- `github.com/gorilla/mux v1.8.1` - HTTP router
- `github.com/michimani/gotwi v0.14.0` - Twitter API v2 client
- `github.com/joho/godotenv v1.5.1` - Environment loader
- `github.com/sirupsen/logrus v1.9.3` - Structured logging

#### API Endpoints
- `GET /health` - Health check
- `GET /api/user/{username}` - Get user info
- `GET /api/tweets/user/{username}?count=N` - Get user tweets
- `GET /api/docs` - API documentation

#### Configuration
- Environment-based configuration
- Validation of required fields
- Default values cho optional fields
- Support for .env file

#### Architecture
- Layered architecture (handlers → services → API)
- Dependency injection
- Single responsibility principle
- Clean code practices

---

## [Unreleased]

### Added
- Endpoint lấy danh sách tài khoản đang theo dõi: `GET /api/user/{username}/following`
- Service/handler logic mới + cập nhật documentation (README, Quickstart, Tutorial, Examples, Scripts)

### Planned Features
- [ ] Pagination support cho tweets
- [ ] Search tweets endpoint
- [ ] Twitter Streaming API support
- [ ] Redis caching layer
- [ ] Database persistence
- [ ] Rate limiting middleware
- [ ] Authentication/API keys
- [ ] Prometheus metrics
- [ ] Unit tests
- [ ] Integration tests
- [ ] CI/CD pipeline
- [ ] Kubernetes deployment

### Potential Improvements
- [ ] GraphQL API option
- [ ] WebSocket support cho real-time updates
- [ ] Tweet sentiment analysis
- [ ] Multiple account monitoring
- [ ] Scheduled tweet fetching
- [ ] Webhook notifications
- [ ] Admin dashboard
- [ ] Analytics và reporting

---

## Version History Format

```
## [Version] - YYYY-MM-DD

### Added
- New features

### Changed
- Changes in existing functionality

### Deprecated
- Soon-to-be removed features

### Removed
- Removed features

### Fixed
- Bug fixes

### Security
- Security fixes
```

---

**Legend:**
- ✅ Completed
- 🚧 In Progress
- 📝 Planned
- ❌ Cancelled

---

Maintained by: AI Assistant
Last Updated: 2024-01-15

