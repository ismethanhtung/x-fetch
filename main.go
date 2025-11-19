package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"x-twitter-backend/config"
	"x-twitter-backend/handlers"
	"x-twitter-backend/services"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Setup logging
	setupLogging()

	log.Info("🚀 Khởi động X Twitter Backend Server...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.WithError(err).Fatal("❌ Không thể load configuration")
	}

	log.WithFields(log.Fields{
		"port":      cfg.ServerPort,
		"host":      cfg.ServerHost,
		"app_env":   cfg.AppEnv,
		"log_level": cfg.LogLevel,
	}).Info("✅ Configuration đã được load")

	// Initialize Twitter service
	twitterService, err := services.NewTwitterService(cfg)
	if err != nil {
		log.WithError(err).Fatal("❌ Không thể khởi tạo Twitter service")
	}

	// Initialize handlers
	tweetsHandler := handlers.NewTweetsHandler(twitterService)

	// Setup router
	router := setupRouter(tweetsHandler)

	// Create HTTP server
	server := &http.Server{
		Addr:         cfg.GetAddress(),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.WithField("address", cfg.GetAddress()).Info("🌐 Server đang lắng nghe...")
		log.Info("📝 API Documentation: http://" + cfg.GetAddress() + "/api/docs")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("❌ Lỗi khi start server")
		}
	}()

	// Graceful shutdown
	gracefulShutdown(server)
}

// setupRouter thiết lập tất cả các routes
func setupRouter(tweetsHandler *handlers.TweetsHandler) *mux.Router {
	router := mux.NewRouter()

	// Apply middlewares
	router.Use(handlers.RecoveryMiddleware)
	router.Use(handlers.LoggingMiddleware)
	router.Use(handlers.CORSMiddleware)

	// Health check
	router.HandleFunc("/health", tweetsHandler.HealthCheck).Methods("GET")

	// API routes
	api := router.PathPrefix("/api").Subrouter()

	// User routes
	api.HandleFunc("/user/{username}", tweetsHandler.GetUserInfo).Methods("GET")
	api.HandleFunc("/user/{username}/following", tweetsHandler.GetUserFollowing).Methods("GET")
	api.HandleFunc("/user/{username}/followers", tweetsHandler.GetUserFollowers).Methods("GET")
	api.HandleFunc("/user/{username}/liked", tweetsHandler.GetLikedTweets).Methods("GET")
	api.HandleFunc("/user/{username}/mentions", tweetsHandler.GetUserMentions).Methods("GET")
	api.HandleFunc("/user/{username}/timelines/reverse_chronological", tweetsHandler.GetUserTimelineReverseChronological).Methods("GET")
	api.HandleFunc("/user/{username}/tweets", tweetsHandler.GetUserTweets).Methods("GET")
	api.HandleFunc("/user/{username}/blocking", tweetsHandler.GetBlockingUsers).Methods("GET")
	api.HandleFunc("/user/{username}/muting", tweetsHandler.GetMutingUsers).Methods("GET")

	// Users routes
	api.HandleFunc("/users", tweetsHandler.ListUsers).Methods("GET")
	api.HandleFunc("/users/{user_id}", tweetsHandler.GetUserByID).Methods("GET")
	api.HandleFunc("/users/by/username/{username}", tweetsHandler.GetUserInfo).Methods("GET")
	api.HandleFunc("/users/me", tweetsHandler.GetMe).Methods("GET")
	api.HandleFunc("/users/search", tweetsHandler.SearchUsers).Methods("GET")
	api.HandleFunc("/users/reposts_of_me", tweetsHandler.GetRepostsOfMe).Methods("GET")

	// Tweets routes
	api.HandleFunc("/tweets", tweetsHandler.ListTweets).Methods("GET")
	api.HandleFunc("/tweets/user/{username}", tweetsHandler.GetUserTweets).Methods("GET")
	api.HandleFunc("/tweets/search", tweetsHandler.SearchTweets).Methods("GET")
	api.HandleFunc("/tweets/search/recent", tweetsHandler.SearchTweets).Methods("GET")
	api.HandleFunc("/tweets/{tweet_id}", tweetsHandler.GetTweetByID).Methods("GET")
	api.HandleFunc("/tweets/{tweet_id}/liking_users", tweetsHandler.GetLikingUsers).Methods("GET")
	api.HandleFunc("/tweets/{tweet_id}/quote_tweets", tweetsHandler.GetQuoteTweets).Methods("GET")
	api.HandleFunc("/tweets/{tweet_id}/retweeted_by", tweetsHandler.GetRetweetedBy).Methods("GET")
	api.HandleFunc("/tweets/{tweet_id}/hidden", tweetsHandler.HideTweet).Methods("PUT")
	api.HandleFunc("/tweets/counts/recent", tweetsHandler.GetTweetCounts).Methods("GET")

	// API documentation endpoint
	api.HandleFunc("/docs", handleAPIDocs).Methods("GET")
	
	// Test page route
	router.HandleFunc("/test", handleTestPage).Methods("GET")
	router.HandleFunc("/", handleTestPage).Methods("GET")

	log.Info("✅ Routes đã được thiết lập")
	return router
}

// setupLogging thiết lập logging configuration
func setupLogging() {
	// Set log format
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})

	// Set output
	log.SetOutput(os.Stdout)

	// Set log level
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

// gracefulShutdown xử lý graceful shutdown
func gracefulShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info("🛑 Đang shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.WithError(err).Error("❌ Lỗi khi shutdown server")
	}

	log.Info("✅ Server đã được shutdown thành công")
}

// handleAPIDocs trả về API documentation
func handleAPIDocs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	docs := `{
  "title": "X Twitter Backend API Documentation",
  "version": "2.0.0",
  "description": "API để lấy tweets và thông tin user từ X/Twitter - Phiên bản mở rộng",
  "endpoints": [
    {
      "path": "/health",
      "method": "GET",
      "description": "Health check endpoint",
      "example": "/health"
    },
    {
      "path": "/api/user/{username}",
      "method": "GET",
      "description": "Lấy thông tin user theo username",
      "parameters": {
        "username": "Username của tài khoản Twitter/X"
      },
      "example": "/api/user/elonmusk"
    },
    {
      "path": "/api/user/{username}/following",
      "method": "GET",
      "description": "Lấy danh sách tài khoản mà user đang theo dõi",
      "parameters": {
        "username": "Username của tài khoản Twitter/X",
        "count": "Số lượng accounts (default: 10, max: 1000)",
        "pagination_token": "Token phân trang (optional)"
      },
      "example": "/api/user/elonmusk/following?count=100"
    },
    {
      "path": "/api/user/{username}/followers",
      "method": "GET",
      "description": "Lấy danh sách người theo dõi (followers) của user",
      "parameters": {
        "username": "Username của tài khoản Twitter/X",
        "count": "Số lượng followers (default: 10, max: 1000)",
        "pagination_token": "Token phân trang (optional)"
      },
      "example": "/api/user/elonmusk/followers?count=50"
    },
    {
      "path": "/api/user/{username}/liked",
      "method": "GET",
      "description": "Lấy danh sách tweets mà user đã like",
      "parameters": {
        "username": "Username của tài khoản Twitter/X",
        "count": "Số lượng tweets (default: 10, max: 100)"
      },
      "example": "/api/user/elonmusk/liked?count=20"
    },
    {
      "path": "/api/user/{username}/mentions",
      "method": "GET",
      "description": "Lấy danh sách tweets có mention đến user",
      "parameters": {
        "username": "Username của tài khoản Twitter/X",
        "count": "Số lượng tweets (default: 10, max: 100)"
      },
      "example": "/api/user/elonmusk/mentions?count=20"
    },
    {
      "path": "/api/tweets/user/{username}",
      "method": "GET",
      "description": "Lấy tweets mới nhất của một user",
      "parameters": {
        "username": "Username của tài khoản Twitter/X",
        "count": "Số lượng tweets (default: 10, max: 100)"
      },
      "example": "/api/tweets/user/elonmusk?count=20"
    },
    {
      "path": "/api/tweets/search",
      "method": "GET",
      "description": "Tìm kiếm tweets theo từ khóa",
      "parameters": {
        "q": "Từ khóa tìm kiếm (bắt buộc)",
        "count": "Số lượng tweets (default: 10, max: 100)"
      },
      "example": "/api/tweets/search?q=golang&count=20"
    },
    {
      "path": "/api/tweets/{tweet_id}",
      "method": "GET",
      "description": "Lấy thông tin chi tiết của một tweet",
      "parameters": {
        "tweet_id": "ID của tweet (bắt buộc)"
      },
      "example": "/api/tweets/1234567890"
    },
    {
      "path": "/api/users/search",
      "method": "GET",
      "description": "Tìm kiếm users theo từ khóa",
      "parameters": {
        "q": "Từ khóa tìm kiếm (bắt buộc)",
        "count": "Số lượng users (default: 10, max: 100)"
      },
      "example": "/api/users/search?q=elon&count=10"
    }
  ],
  "authentication": "Yêu cầu TWITTER_BEARER_TOKEN trong environment variables",
  "notes": [
    "API tuân thủ rate limits của Twitter API",
    "Tất cả responses trả về dạng JSON",
    "Errors được trả về với format chuẩn: {error, message, code}",
    "Các API miễn phí và không bị giới hạn bởi Twitter API v2"
  ]
}`

	w.Write([]byte(docs))
}

// handleTestPage trả về trang test HTML
func handleTestPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	
	html := getTestPageHTML()
	w.Write([]byte(html))
}

// getTestPageHTML trả về HTML của trang test
func getTestPageHTML() string {
	return `<!DOCTYPE html>
<html lang="vi">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>X/Twitter API Testing Dashboard</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }

        .container {
            max-width: 1400px;
            margin: 0 auto;
        }

        .header {
            background: white;
            border-radius: 20px;
            padding: 30px;
            box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
            margin-bottom: 30px;
            text-align: center;
        }

        .header h1 {
            color: #1DA1F2;
            font-size: 2.5em;
            margin-bottom: 10px;
        }

        .header p {
            color: #666;
            font-size: 1.1em;
        }

        .api-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
            gap: 20px;
            margin-bottom: 30px;
        }

        .api-card {
            background: white;
            border-radius: 15px;
            padding: 25px;
            box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
            transition: transform 0.3s ease, box-shadow 0.3s ease;
        }

        .api-card:hover {
            transform: translateY(-5px);
            box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
        }

        .api-card h3 {
            color: #1DA1F2;
            margin-bottom: 15px;
            font-size: 1.3em;
            display: flex;
            align-items: center;
            gap: 10px;
        }

        .api-card .icon {
            font-size: 1.5em;
        }

        .api-card p {
            color: #666;
            margin-bottom: 15px;
            font-size: 0.9em;
        }

        .form-group {
            margin-bottom: 15px;
        }

        .form-group label {
            display: block;
            margin-bottom: 5px;
            color: #333;
            font-weight: 500;
            font-size: 0.9em;
        }

        .form-group input {
            width: 100%;
            padding: 10px 15px;
            border: 2px solid #e1e8ed;
            border-radius: 8px;
            font-size: 14px;
            transition: border-color 0.3s ease;
        }

        .form-group input:focus {
            outline: none;
            border-color: #1DA1F2;
        }

        .btn {
            width: 100%;
            padding: 12px;
            background: linear-gradient(135deg, #1DA1F2, #0d8bd9);
            color: white;
            border: none;
            border-radius: 8px;
            font-size: 1em;
            font-weight: 600;
            cursor: pointer;
            transition: transform 0.2s ease, box-shadow 0.2s ease;
        }

        .btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(29, 161, 242, 0.3);
        }

        .btn:active {
            transform: translateY(0);
        }

        .response-section {
            background: white;
            border-radius: 15px;
            padding: 25px;
            box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
            margin-top: 30px;
        }

        .response-section h2 {
            color: #1DA1F2;
            margin-bottom: 20px;
            display: flex;
            align-items: center;
            gap: 10px;
        }

        #response {
            background: #f5f8fa;
            border: 2px solid #e1e8ed;
            border-radius: 10px;
            padding: 20px;
            max-height: 600px;
            overflow-y: auto;
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
            font-size: 13px;
            line-height: 1.6;
            white-space: pre-wrap;
            word-wrap: break-word;
        }

        .loading {
            display: none;
            text-align: center;
            padding: 20px;
        }

        .loading.active {
            display: block;
        }

        .spinner {
            border: 4px solid #f3f3f3;
            border-top: 4px solid #1DA1F2;
            border-radius: 50%;
            width: 40px;
            height: 40px;
            animation: spin 1s linear infinite;
            margin: 0 auto 10px;
        }

        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }

        .error {
            color: #e74c3c;
            background: #ffe5e5;
            padding: 15px;
            border-radius: 8px;
            border-left: 4px solid #e74c3c;
        }

        .success {
            color: #27ae60;
            background: #e8f8f5;
            padding: 15px;
            border-radius: 8px;
            border-left: 4px solid #27ae60;
        }

        .badge {
            display: inline-block;
            padding: 3px 8px;
            background: #e8f5fe;
            color: #1DA1F2;
            border-radius: 5px;
            font-size: 0.75em;
            font-weight: 600;
            margin-left: 10px;
        }

        .endpoint {
            background: #f5f8fa;
            padding: 8px 12px;
            border-radius: 6px;
            font-family: monospace;
            font-size: 0.85em;
            margin-top: 10px;
            color: #555;
        }

        /* Scrollbar styling */
        #response::-webkit-scrollbar {
            width: 10px;
        }

        #response::-webkit-scrollbar-track {
            background: #f1f1f1;
            border-radius: 10px;
        }

        #response::-webkit-scrollbar-thumb {
            background: #1DA1F2;
            border-radius: 10px;
        }

        #response::-webkit-scrollbar-thumb:hover {
            background: #0d8bd9;
        }

        @media (max-width: 768px) {
            .api-grid {
                grid-template-columns: 1fr;
            }

            .header h1 {
                font-size: 1.8em;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🐦 X/Twitter API Testing Dashboard</h1>
            <p>Công cụ test API X/Twitter - Miễn phí & Không giới hạn</p>
            <p style="font-size: 0.9em; margin-top: 10px; color: #888;">📊 Tổng cộng: <strong>22 APIs</strong> sẵn sàng để test</p>
        </div>

        <div class="api-grid">
            <!-- User Info API -->
            <div class="api-card">
                <h3><span class="icon">👤</span> Thông tin User</h3>
                <p>Lấy thông tin chi tiết của một tài khoản Twitter/X</p>
                <div class="endpoint">GET /api/user/{username}</div>
                <div class="form-group">
                    <label>Username</label>
                    <input type="text" id="user-info-username" placeholder="elonmusk" value="elonmusk">
                </div>
                <button class="btn" onclick="getUserInfo()">Lấy thông tin</button>
            </div>

            <!-- User Tweets API -->
            <div class="api-card">
                <h3><span class="icon">📝</span> Tweets của User</h3>
                <p>Lấy danh sách tweets mới nhất của một user</p>
                <div class="endpoint">GET /api/tweets/user/{username}</div>
                <div class="form-group">
                    <label>Username</label>
                    <input type="text" id="user-tweets-username" placeholder="elonmusk" value="elonmusk">
                </div>
                <div class="form-group">
                    <label>Số lượng (max: 100)</label>
                    <input type="number" id="user-tweets-count" value="10" min="1" max="100">
                </div>
                <button class="btn" onclick="getUserTweets()">Lấy tweets</button>
            </div>

            <!-- Search Tweets API -->
            <div class="api-card">
                <h3><span class="icon">🔍</span> Tìm kiếm Tweets</h3>
                <p>Tìm kiếm tweets theo từ khóa</p>
                <div class="endpoint">GET /api/tweets/search</div>
                <div class="form-group">
                    <label>Từ khóa</label>
                    <input type="text" id="search-tweets-query" placeholder="golang" value="golang">
                </div>
                <div class="form-group">
                    <label>Số lượng (max: 100)</label>
                    <input type="number" id="search-tweets-count" value="10" min="1" max="100">
                </div>
                <button class="btn" onclick="searchTweets()">Tìm kiếm</button>
            </div>

            <!-- Tweet Detail API -->
            <div class="api-card">
                <h3><span class="icon">📄</span> Chi tiết Tweet</h3>
                <p>Lấy thông tin chi tiết của một tweet</p>
                <div class="endpoint">GET /api/tweets/{tweet_id}</div>
                <div class="form-group">
                    <label>Tweet ID</label>
                    <input type="text" id="tweet-id" placeholder="1234567890">
                </div>
                <button class="btn" onclick="getTweetDetail()">Xem chi tiết</button>
            </div>

            <!-- Following API -->
            <div class="api-card">
                <h3><span class="icon">➕</span> Danh sách Following</h3>
                <p>Lấy danh sách tài khoản mà user đang theo dõi</p>
                <div class="endpoint">GET /api/user/{username}/following</div>
                <div class="form-group">
                    <label>Username</label>
                    <input type="text" id="following-username" placeholder="elonmusk" value="elonmusk">
                </div>
                <div class="form-group">
                    <label>Số lượng (max: 1000)</label>
                    <input type="number" id="following-count" value="10" min="1" max="1000">
                </div>
                <button class="btn" onclick="getFollowing()">Xem Following</button>
            </div>

            <!-- Followers API -->
            <div class="api-card">
                <h3><span class="icon">👥</span> Danh sách Followers</h3>
                <p>Lấy danh sách người theo dõi user</p>
                <div class="endpoint">GET /api/user/{username}/followers</div>
                <div class="form-group">
                    <label>Username</label>
                    <input type="text" id="followers-username" placeholder="elonmusk" value="elonmusk">
                </div>
                <div class="form-group">
                    <label>Số lượng (max: 1000)</label>
                    <input type="number" id="followers-count" value="10" min="1" max="1000">
                </div>
                <button class="btn" onclick="getFollowers()">Xem Followers</button>
            </div>

            <!-- Liked Tweets API -->
            <div class="api-card">
                <h3><span class="icon">❤️</span> Tweets đã Like</h3>
                <p>Lấy danh sách tweets mà user đã like</p>
                <div class="endpoint">GET /api/user/{username}/liked</div>
                <div class="form-group">
                    <label>Username</label>
                    <input type="text" id="liked-username" placeholder="elonmusk" value="elonmusk">
                </div>
                <div class="form-group">
                    <label>Số lượng (max: 100)</label>
                    <input type="number" id="liked-count" value="10" min="1" max="100">
                </div>
                <button class="btn" onclick="getLikedTweets()">Xem Liked</button>
            </div>

            <!-- Mentions API -->
            <div class="api-card">
                <h3><span class="icon">@</span> Mentions</h3>
                <p>Lấy tweets có mention đến user</p>
                <div class="endpoint">GET /api/user/{username}/mentions</div>
                <div class="form-group">
                    <label>Username</label>
                    <input type="text" id="mentions-username" placeholder="elonmusk" value="elonmusk">
                </div>
                <div class="form-group">
                    <label>Số lượng (max: 100)</label>
                    <input type="number" id="mentions-count" value="10" min="1" max="100">
                </div>
                <button class="btn" onclick="getMentions()">Xem Mentions</button>
            </div>

            <!-- Search Users API -->
            <div class="api-card">
                <h3><span class="icon">🔎</span> Tìm kiếm Users</h3>
                <p>Tìm kiếm users theo từ khóa</p>
                <div class="endpoint">GET /api/users/search</div>
                <div class="form-group">
                    <label>Từ khóa</label>
                    <input type="text" id="search-users-query" placeholder="elon" value="elon">
                </div>
                <div class="form-group">
                    <label>Số lượng (max: 100)</label>
                    <input type="number" id="search-users-count" value="10" min="1" max="100">
                </div>
                <button class="btn" onclick="searchUsers()">Tìm Users</button>
            </div>

            <!-- List Tweets API -->
            <div class="api-card">
                <h3><span class="icon">📋</span> Danh sách Tweets</h3>
                <p>Lấy danh sách tweets theo IDs</p>
                <div class="endpoint">GET /api/tweets?ids=123,456</div>
                <div class="form-group">
                    <label>Tweet IDs (comma-separated)</label>
                    <input type="text" id="list-tweets-ids" placeholder="1234567890,9876543210">
                </div>
                <button class="btn" onclick="listTweets()">Lấy Tweets</button>
            </div>

            <!-- Liking Users API -->
            <div class="api-card" style="border-left: 4px solid #ff9800;">
                <h3><span class="icon">👍</span> Users đã Like <span style="font-size: 0.7em; color: #ff9800;">⚠️ Limited</span></h3>
                <p>Lấy danh sách users đã like một tweet <strong style="color: #ff9800;">(Có thể yêu cầu OAuth)</strong></p>
                <div class="endpoint">GET /api/tweets/{id}/liking_users</div>
                <div class="form-group">
                    <label>Tweet ID</label>
                    <input type="text" id="liking-users-tweet-id" placeholder="1234567890">
                </div>
                <div class="form-group">
                    <label>Số lượng (max: 100)</label>
                    <input type="number" id="liking-users-count" value="10" min="1" max="100">
                </div>
                <button class="btn" onclick="getLikingUsers()">Xem Liking Users</button>
            </div>

            <!-- Quote Tweets API -->
            <div class="api-card">
                <h3><span class="icon">💬</span> Quote Tweets</h3>
                <p>Lấy danh sách quote tweets của một tweet</p>
                <div class="endpoint">GET /api/tweets/{id}/quote_tweets</div>
                <div class="form-group">
                    <label>Tweet ID</label>
                    <input type="text" id="quote-tweets-tweet-id" placeholder="1234567890">
                </div>
                <div class="form-group">
                    <label>Số lượng (max: 100)</label>
                    <input type="number" id="quote-tweets-count" value="10" min="1" max="100">
                </div>
                <button class="btn" onclick="getQuoteTweets()">Xem Quote Tweets</button>
            </div>

            <!-- Retweeted By API -->
            <div class="api-card">
                <h3><span class="icon">🔄</span> Retweeted By</h3>
                <p>Lấy danh sách users đã retweet</p>
                <div class="endpoint">GET /api/tweets/{id}/retweeted_by</div>
                <div class="form-group">
                    <label>Tweet ID</label>
                    <input type="text" id="retweeted-by-tweet-id" placeholder="1234567890">
                </div>
                <div class="form-group">
                    <label>Số lượng (max: 100)</label>
                    <input type="number" id="retweeted-by-count" value="10" min="1" max="100">
                </div>
                <button class="btn" onclick="getRetweetedBy()">Xem Retweeted By</button>
            </div>

            <!-- Tweet Counts API -->
            <div class="api-card">
                <h3><span class="icon">📊</span> Tweet Counts</h3>
                <p>Lấy số lượng tweets theo query và time range</p>
                <div class="endpoint">GET /api/tweets/counts/recent</div>
                <div class="form-group">
                    <label>Query</label>
                    <input type="text" id="tweet-counts-query" placeholder="golang" value="golang">
                </div>
                <div class="form-group">
                    <label>Start Time (RFC3339)</label>
                    <input type="text" id="tweet-counts-start" placeholder="2024-01-01T00:00:00Z">
                </div>
                <div class="form-group">
                    <label>End Time (RFC3339)</label>
                    <input type="text" id="tweet-counts-end" placeholder="2024-01-02T00:00:00Z">
                </div>
                <button class="btn" onclick="getTweetCounts()">Lấy Counts</button>
            </div>

            <!-- List Users API -->
            <div class="api-card">
                <h3><span class="icon">👥</span> Danh sách Users</h3>
                <p>Lấy danh sách users theo IDs</p>
                <div class="endpoint">GET /api/users?ids=123,456</div>
                <div class="form-group">
                    <label>User IDs (comma-separated)</label>
                    <input type="text" id="list-users-ids" placeholder="44196397,12345678">
                </div>
                <button class="btn" onclick="listUsers()">Lấy Users</button>
            </div>

            <!-- User By ID API -->
            <div class="api-card">
                <h3><span class="icon">🆔</span> User theo ID</h3>
                <p>Lấy thông tin user theo ID</p>
                <div class="endpoint">GET /api/users/{user_id}</div>
                <div class="form-group">
                    <label>User ID</label>
                    <input type="text" id="user-by-id" placeholder="44196397">
                </div>
                <button class="btn" onclick="getUserByID()">Lấy User</button>
            </div>

            <!-- Get Me API -->
            <div class="api-card">
                <h3><span class="icon">👤</span> Authenticated User</h3>
                <p>Lấy thông tin authenticated user</p>
                <div class="endpoint">GET /api/users/me</div>
                <button class="btn" onclick="getMe()">Lấy Thông tin</button>
            </div>

            <!-- Blocking Users API -->
            <div class="api-card" style="border-left: 4px solid #ff9800;">
                <h3><span class="icon">🚫</span> Blocking Users <span style="font-size: 0.7em; color: #ff9800;">⚠️ OAuth Required</span></h3>
                <p>Lấy danh sách users bị block <strong style="color: #ff9800;">(Yêu cầu OAuth 1.0a)</strong></p>
                <div class="endpoint">GET /api/users/{username}/blocking</div>
                <div class="form-group">
                    <label>Username</label>
                    <input type="text" id="blocking-username" placeholder="elonmusk" value="elonmusk">
                </div>
                <div class="form-group">
                    <label>Số lượng (max: 1000)</label>
                    <input type="number" id="blocking-count" value="10" min="1" max="1000">
                </div>
                <button class="btn" onclick="getBlockingUsers()">Xem Blocking</button>
            </div>

            <!-- Muting Users API -->
            <div class="api-card" style="border-left: 4px solid #ff9800;">
                <h3><span class="icon">🔇</span> Muting Users <span style="font-size: 0.7em; color: #ff9800;">⚠️ OAuth Required</span></h3>
                <p>Lấy danh sách users bị mute <strong style="color: #ff9800;">(Yêu cầu OAuth 1.0a)</strong></p>
                <div class="endpoint">GET /api/users/{username}/muting</div>
                <div class="form-group">
                    <label>Username</label>
                    <input type="text" id="muting-username" placeholder="elonmusk" value="elonmusk">
                </div>
                <div class="form-group">
                    <label>Số lượng (max: 1000)</label>
                    <input type="number" id="muting-count" value="10" min="1" max="1000">
                </div>
                <button class="btn" onclick="getMutingUsers()">Xem Muting</button>
            </div>

            <!-- Timeline Reverse Chronological API -->
            <div class="api-card">
                <h3><span class="icon">⏰</span> Timeline Reverse</h3>
                <p>Lấy timeline reverse chronological</p>
                <div class="endpoint">GET /api/users/{username}/timelines/reverse_chronological</div>
                <div class="form-group">
                    <label>Username</label>
                    <input type="text" id="timeline-username" placeholder="elonmusk" value="elonmusk">
                </div>
                <div class="form-group">
                    <label>Số lượng (max: 100)</label>
                    <input type="number" id="timeline-count" value="10" min="1" max="100">
                </div>
                <button class="btn" onclick="getTimeline()">Xem Timeline</button>
            </div>

            <!-- Reposts Of Me API -->
            <div class="api-card" style="border-left: 4px solid #ff9800;">
                <h3><span class="icon">📤</span> Reposts Of Me <span style="font-size: 0.7em; color: #ff9800;">⚠️ OAuth Required</span></h3>
                <p>Lấy reposts của authenticated user <strong style="color: #ff9800;">(Yêu cầu OAuth 1.0a)</strong></p>
                <div class="endpoint">GET /api/users/reposts_of_me</div>
                <div class="form-group">
                    <label>Số lượng (max: 100)</label>
                    <input type="number" id="reposts-count" value="10" min="1" max="100">
                </div>
                <button class="btn" onclick="getRepostsOfMe()">Xem Reposts</button>
            </div>
        </div>

        <div class="response-section">
            <h2><span>📊</span> Response</h2>
            <div class="loading" id="loading">
                <div class="spinner"></div>
                <p>Đang xử lý...</p>
            </div>
            <div id="response">Chọn một API ở trên để bắt đầu test...</div>
        </div>
    </div>

    <script>
        const responseDiv = document.getElementById('response');
        const loadingDiv = document.getElementById('loading');

        async function makeRequest(url, description) {
            loadingDiv.classList.add('active');
            responseDiv.innerHTML = '';

            try {
                const response = await fetch(url);
                let data;
                
                // Thử parse JSON, nếu không được thì lấy text
                const contentType = response.headers.get('content-type');
                if (contentType && contentType.includes('application/json')) {
                    try {
                        data = await response.json();
                    } catch (e) {
                        const text = await response.text();
                        loadingDiv.classList.remove('active');
                        responseDiv.innerHTML = '<div class="error">❌ Lỗi parse JSON: ' + e.message + '</div>\n\n' + 
                            '<pre>' + text.substring(0, 500) + '</pre>';
                        return;
                    }
                } else {
                    const text = await response.text();
                    loadingDiv.classList.remove('active');
                    responseDiv.innerHTML = '<div class="error">❌ Response không phải JSON</div>\n\n' + 
                        '<pre>' + text.substring(0, 500) + '</pre>';
                    return;
                }

                loadingDiv.classList.remove('active');

                if (response.ok) {
                    responseDiv.innerHTML = '<div class="success">✅ ' + description + ' thành công!</div>\n\n' + 
                        JSON.stringify(data, null, 2);
                } else {
                    // Kiểm tra nếu là lỗi OAuth required
                    const errorMsg = data.message || 'Không xác định';
                    let errorClass = 'error';
                    let errorIcon = '❌';
                    
                    if (errorMsg.includes('OAuth') || errorMsg.includes('403') || errorMsg.includes('Forbidden')) {
                        errorClass = 'error';
                        errorIcon = '⚠️';
                    }
                    
                    responseDiv.innerHTML = '<div class="' + errorClass + '">' + errorIcon + ' ' + errorMsg + '</div>\n\n' + 
                        JSON.stringify(data, null, 2);
                }
            } catch (error) {
                loadingDiv.classList.remove('active');
                let errorMsg = error.message;
                if (errorMsg.includes('JSON') && errorMsg.includes('position')) {
                    errorMsg = 'Lỗi parse JSON - Response không hợp lệ. Có thể API yêu cầu OAuth hoặc có giới hạn với Bearer Token.';
                }
                responseDiv.innerHTML = '<div class="error">❌ Lỗi: ' + errorMsg + '</div>';
            }
        }

        function getUserInfo() {
            const username = document.getElementById('user-info-username').value;
            if (!username) {
                alert('Vui lòng nhập username!');
                return;
            }
            makeRequest('/api/user/' + username, 'Lấy thông tin user');
        }

        function getUserTweets() {
            const username = document.getElementById('user-tweets-username').value;
            const count = document.getElementById('user-tweets-count').value;
            if (!username) {
                alert('Vui lòng nhập username!');
                return;
            }
            makeRequest('/api/tweets/user/' + username + '?count=' + count, 'Lấy tweets');
        }

        function searchTweets() {
            const query = document.getElementById('search-tweets-query').value;
            const count = document.getElementById('search-tweets-count').value;
            if (!query) {
                alert('Vui lòng nhập từ khóa!');
                return;
            }
            makeRequest('/api/tweets/search?q=' + encodeURIComponent(query) + '&count=' + count, 'Tìm kiếm tweets');
        }

        function getTweetDetail() {
            const tweetId = document.getElementById('tweet-id').value;
            if (!tweetId) {
                alert('Vui lòng nhập Tweet ID!');
                return;
            }
            makeRequest('/api/tweets/' + tweetId, 'Lấy chi tiết tweet');
        }

        function getFollowing() {
            const username = document.getElementById('following-username').value;
            const count = document.getElementById('following-count').value;
            if (!username) {
                alert('Vui lòng nhập username!');
                return;
            }
            makeRequest('/api/user/' + username + '/following?count=' + count, 'Lấy danh sách following');
        }

        function getFollowers() {
            const username = document.getElementById('followers-username').value;
            const count = document.getElementById('followers-count').value;
            if (!username) {
                alert('Vui lòng nhập username!');
                return;
            }
            makeRequest('/api/user/' + username + '/followers?count=' + count, 'Lấy danh sách followers');
        }

        function getLikedTweets() {
            const username = document.getElementById('liked-username').value;
            const count = document.getElementById('liked-count').value;
            if (!username) {
                alert('Vui lòng nhập username!');
                return;
            }
            makeRequest('/api/user/' + username + '/liked?count=' + count, 'Lấy liked tweets');
        }

        function getMentions() {
            const username = document.getElementById('mentions-username').value;
            const count = document.getElementById('mentions-count').value;
            if (!username) {
                alert('Vui lòng nhập username!');
                return;
            }
            makeRequest('/api/user/' + username + '/mentions?count=' + count, 'Lấy mentions');
        }

        function searchUsers() {
            const query = document.getElementById('search-users-query').value;
            const count = document.getElementById('search-users-count').value;
            if (!query) {
                alert('Vui lòng nhập từ khóa!');
                return;
            }
            makeRequest('/api/users/search?q=' + encodeURIComponent(query) + '&count=' + count, 'Tìm kiếm users');
        }

        function listTweets() {
            const ids = document.getElementById('list-tweets-ids').value;
            if (!ids) {
                alert('Vui lòng nhập Tweet IDs!');
                return;
            }
            makeRequest('/api/tweets?ids=' + encodeURIComponent(ids), 'Lấy danh sách tweets');
        }

        function getLikingUsers() {
            const tweetId = document.getElementById('liking-users-tweet-id').value;
            const count = document.getElementById('liking-users-count').value;
            if (!tweetId) {
                alert('Vui lòng nhập Tweet ID!');
                return;
            }
            makeRequest('/api/tweets/' + tweetId + '/liking_users?count=' + count, 'Lấy liking users');
        }

        function getQuoteTweets() {
            const tweetId = document.getElementById('quote-tweets-tweet-id').value;
            const count = document.getElementById('quote-tweets-count').value;
            if (!tweetId) {
                alert('Vui lòng nhập Tweet ID!');
                return;
            }
            makeRequest('/api/tweets/' + tweetId + '/quote_tweets?count=' + count, 'Lấy quote tweets');
        }

        function getRetweetedBy() {
            const tweetId = document.getElementById('retweeted-by-tweet-id').value;
            const count = document.getElementById('retweeted-by-count').value;
            if (!tweetId) {
                alert('Vui lòng nhập Tweet ID!');
                return;
            }
            makeRequest('/api/tweets/' + tweetId + '/retweeted_by?count=' + count, 'Lấy retweeted by');
        }

        function getTweetCounts() {
            const query = document.getElementById('tweet-counts-query').value;
            const startTime = document.getElementById('tweet-counts-start').value;
            const endTime = document.getElementById('tweet-counts-end').value;
            if (!query) {
                alert('Vui lòng nhập query!');
                return;
            }
            let url = '/api/tweets/counts/recent?q=' + encodeURIComponent(query);
            if (startTime) url += '&start_time=' + encodeURIComponent(startTime);
            if (endTime) url += '&end_time=' + encodeURIComponent(endTime);
            makeRequest(url, 'Lấy tweet counts');
        }

        function listUsers() {
            const ids = document.getElementById('list-users-ids').value;
            if (!ids) {
                alert('Vui lòng nhập User IDs!');
                return;
            }
            makeRequest('/api/users?ids=' + encodeURIComponent(ids), 'Lấy danh sách users');
        }

        function getUserByID() {
            const userId = document.getElementById('user-by-id').value;
            if (!userId) {
                alert('Vui lòng nhập User ID!');
                return;
            }
            makeRequest('/api/users/' + userId, 'Lấy user theo ID');
        }

        function getMe() {
            makeRequest('/api/users/me', 'Lấy authenticated user');
        }

        function getBlockingUsers() {
            const username = document.getElementById('blocking-username').value;
            const count = document.getElementById('blocking-count').value;
            if (!username) {
                alert('Vui lòng nhập username!');
                return;
            }
            makeRequest('/api/users/' + username + '/blocking?count=' + count, 'Lấy blocking users');
        }

        function getMutingUsers() {
            const username = document.getElementById('muting-username').value;
            const count = document.getElementById('muting-count').value;
            if (!username) {
                alert('Vui lòng nhập username!');
                return;
            }
            makeRequest('/api/users/' + username + '/muting?count=' + count, 'Lấy muting users');
        }

        function getTimeline() {
            const username = document.getElementById('timeline-username').value;
            const count = document.getElementById('timeline-count').value;
            if (!username) {
                alert('Vui lòng nhập username!');
                return;
            }
            makeRequest('/api/users/' + username + '/timelines/reverse_chronological?count=' + count, 'Lấy timeline');
        }

        function getRepostsOfMe() {
            const count = document.getElementById('reposts-count').value;
            makeRequest('/api/users/reposts_of_me?count=' + count, 'Lấy reposts of me');
        }
    </script>
</body>
</html>`
}
