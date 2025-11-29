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
    <title>X/Twitter API Dashboard</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap" rel="stylesheet">
    <style>
        :root {
            --primary: #1d9bf0;
            --primary-hover: #1a8cd8;
            --bg-body: #f7f9f9;
            --bg-card: #ffffff;
            --text-main: #0f1419;
            --text-secondary: #536471;
            --border: #eff3f4;
            --success: #00ba7c;
            --error: #f91880;
            --sidebar-width: 280px;
            --header-height: 60px;
            --shadow-sm: 0 2px 8px rgba(0,0,0,0.04);
            --shadow-md: 0 8px 16px rgba(0,0,0,0.08);
        }

        * { margin: 0; padding: 0; box-sizing: border-box; }

        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
            background-color: var(--bg-body);
            color: var(--text-main);
            height: 100vh;
            overflow: hidden;
            display: flex;
        }

        /* Sidebar */
        .sidebar {
            width: var(--sidebar-width);
            background: var(--bg-card);
            border-right: 1px solid var(--border);
            display: flex;
            flex-direction: column;
            height: 100%;
            transition: transform 0.3s ease;
            z-index: 100;
        }

        .brand {
            height: var(--header-height);
            display: flex;
            align-items: center;
            padding: 0 24px;
            font-weight: 800;
            font-size: 20px;
            color: var(--primary);
            border-bottom: 1px solid var(--border);
        }

        .brand span { margin-left: 10px; color: var(--text-main); }

        .nav-menu {
            flex: 1;
            overflow-y: auto;
            padding: 20px 16px;
        }

        .nav-group { margin-bottom: 24px; }
        .nav-label {
            font-size: 12px;
            text-transform: uppercase;
            color: var(--text-secondary);
            font-weight: 600;
            padding: 0 12px;
            margin-bottom: 8px;
            letter-spacing: 0.5px;
        }

        .nav-item {
            display: flex;
            align-items: center;
            padding: 10px 12px;
            color: var(--text-main);
            text-decoration: none;
            border-radius: 8px;
            margin-bottom: 4px;
            font-size: 14px;
            font-weight: 500;
            transition: all 0.2s;
            cursor: pointer;
        }

        .nav-item:hover { background-color: var(--border); }
        .nav-item.active {
            background-color: rgba(29, 155, 240, 0.1);
            color: var(--primary);
        }
        
        .nav-item i { width: 20px; margin-right: 10px; text-align: center; }

        /* Main Content */
        .main-content {
            flex: 1;
            display: flex;
            flex-direction: column;
            height: 100%;
            overflow: hidden;
            position: relative;
        }

        .header {
            height: var(--header-height);
            background: var(--bg-card);
            border-bottom: 1px solid var(--border);
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 0 30px;
        }

        .page-title { font-size: 18px; font-weight: 700; }

        .content-scroll {
            flex: 1;
            overflow-y: auto;
            padding: 30px;
        }

        .api-section {
            display: none;
            animation: fadeIn 0.3s ease;
        }
        .api-section.active { display: block; }

        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(10px); }
            to { opacity: 1; transform: translateY(0); }
        }

        .card {
            background: var(--bg-card);
            border-radius: 16px;
            padding: 24px;
            box-shadow: var(--shadow-sm);
            border: 1px solid var(--border);
            margin-bottom: 24px;
        }

        .card-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 20px;
        }

        .card-title { font-size: 16px; font-weight: 700; display: flex; align-items: center; gap: 10px; }
        .endpoint-badge {
            background: #eff3f4;
            padding: 4px 8px;
            border-radius: 6px;
            font-family: monospace;
            font-size: 12px;
            color: var(--text-secondary);
        }

        .form-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 16px;
            margin-bottom: 20px;
        }

        .form-group label {
            display: block;
            font-size: 13px;
            font-weight: 600;
            color: var(--text-secondary);
            margin-bottom: 6px;
        }

        .form-control {
            width: 100%;
            padding: 10px 12px;
            border: 1px solid #cfd9de;
            border-radius: 8px;
            font-size: 14px;
            transition: border-color 0.2s;
            font-family: inherit;
        }

        .form-control:focus {
            outline: none;
            border-color: var(--primary);
            box-shadow: 0 0 0 2px rgba(29, 155, 240, 0.2);
        }

        .btn {
            display: inline-flex;
            align-items: center;
            justify-content: center;
            padding: 10px 20px;
            background: var(--text-main);
            color: white;
            border: none;
            border-radius: 20px;
            font-weight: 600;
            font-size: 14px;
            cursor: pointer;
            transition: all 0.2s;
        }

        .btn:hover { opacity: 0.9; transform: translateY(-1px); }
        .btn-primary { background: var(--primary); }
        .btn-primary:hover { background: var(--primary-hover); }

        /* Response Area */
        .response-area {
            position: fixed;
            right: 0;
            top: var(--header-height);
            bottom: 0;
            width: 40%;
            background: #f7f9f9;
            border-left: 1px solid var(--border);
            display: flex;
            flex-direction: column;
            transform: translateX(100%);
            transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
            z-index: 90;
            box-shadow: -5px 0 20px rgba(0,0,0,0.05);
        }

        .response-area.open { transform: translateX(0); }

        .response-header {
            padding: 16px 24px;
            background: white;
            border-bottom: 1px solid var(--border);
            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        .response-title { font-weight: 700; font-size: 15px; }
        .close-response {
            background: none;
            border: none;
            cursor: pointer;
            font-size: 20px;
            color: var(--text-secondary);
            padding: 4px;
            border-radius: 50%;
        }
        .close-response:hover { background: var(--border); }

        .response-content {
            flex: 1;
            overflow: auto;
            padding: 20px;
            font-family: 'Monaco', 'Menlo', monospace;
            font-size: 13px;
        }

        pre { margin: 0; white-space: pre-wrap; word-break: break-all; }
        
        .json-key { color: #cf222e; }
        .json-string { color: #0a3069; }
        .json-number { color: #0550ae; }
        .json-boolean { color: #0550ae; }
        .json-null { color: #953800; }

        /* Mobile Toggle */
        .menu-toggle { display: none; font-size: 20px; cursor: pointer; }

        @media (max-width: 1024px) {
            .response-area { width: 50%; }
        }

        @media (max-width: 768px) {
            .sidebar {
                position: fixed;
                left: -100%;
                height: 100%;
                box-shadow: 2px 0 10px rgba(0,0,0,0.1);
            }
            .sidebar.active { left: 0; }
            .menu-toggle { display: block; }
            .response-area { width: 100%; z-index: 200; }
        }

        /* Loading Overlay */
        .loading-overlay {
            position: fixed;
            top: 0; left: 0; right: 0; bottom: 0;
            background: rgba(255,255,255,0.8);
            display: flex;
            justify-content: center;
            align-items: center;
            z-index: 1000;
            opacity: 0;
            pointer-events: none;
            transition: opacity 0.2s;
        }
        .loading-overlay.active { opacity: 1; pointer-events: all; }
        
        .spinner {
            width: 40px; height: 40px;
            border: 3px solid var(--border);
            border-top-color: var(--primary);
            border-radius: 50%;
            animation: spin 0.8s linear infinite;
        }
        @keyframes spin { to { transform: rotate(360deg); } }

    </style>
</head>
<body>

    <!-- Sidebar -->
    <div class="sidebar" id="sidebar">
        <div class="brand">
            <svg viewBox="0 0 24 24" width="24" height="24" fill="currentColor"><g><path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"></path></g></svg>
            <span>API Tester</span>
        </div>
        <div class="nav-menu">
            <div class="nav-group">
                <div class="nav-label">User</div>
                <div class="nav-item active" onclick="showSection('user-info')">
                    <i>👤</i> Info
                </div>
                <div class="nav-item" onclick="showSection('user-following')">
                    <i>➕</i> Following
                </div>
                <div class="nav-item" onclick="showSection('user-followers')">
                    <i>👥</i> Followers
                </div>
                <div class="nav-item" onclick="showSection('user-liked')">
                    <i>❤️</i> Liked
                </div>
                <div class="nav-item" onclick="showSection('user-mentions')">
                    <i>@</i> Mentions
                </div>
            </div>
            
            <div class="nav-group">
                <div class="nav-label">Tweets</div>
                <div class="nav-item" onclick="showSection('user-tweets')">
                    <i>📝</i> User Tweets
                </div>
                <div class="nav-item" onclick="showSection('tweet-detail')">
                    <i>📄</i> Tweet Detail
                </div>
                <div class="nav-item" onclick="showSection('search-tweets')">
                    <i>🔍</i> Search
                </div>
                <div class="nav-item" onclick="showSection('tweet-counts')">
                    <i>📊</i> Counts
                </div>
            </div>

            <div class="nav-group">
                <div class="nav-label">Advanced</div>
                <div class="nav-item" onclick="showSection('liking-users')">
                    <i>👍</i> Liking Users
                </div>
                <div class="nav-item" onclick="showSection('quote-tweets')">
                    <i>💬</i> Quotes
                </div>
                <div class="nav-item" onclick="showSection('retweeted-by')">
                    <i>🔄</i> Retweets
                </div>
            </div>
        </div>
    </div>

    <!-- Main Content -->
    <div class="main-content">
        <div class="header">
            <div style="display: flex; align-items: center; gap: 15px;">
                <div class="menu-toggle" onclick="toggleSidebar()">☰</div>
                <div class="page-title" id="page-title">User Info</div>
            </div>
            <a href="/api/docs" target="_blank" class="btn" style="background: var(--text-main); font-size: 13px;">View Docs</a>
        </div>

        <div class="content-scroll">
            
            <!-- User Info -->
            <div id="user-info" class="api-section active">
                <div class="card">
                    <div class="card-header">
                        <div class="card-title">Get User Information <span class="endpoint-badge">GET /api/user/{username}</span></div>
                    </div>
                    <div class="form-grid">
                        <div class="form-group">
                            <label>Username</label>
                            <input type="text" class="form-control" id="user-info-username" value="elonmusk" placeholder="e.g. elonmusk">
                        </div>
                    </div>
                    <button class="btn btn-primary" onclick="executeRequest('getUserInfo')">Execute Request</button>
                </div>
            </div>

            <!-- User Following -->
            <div id="user-following" class="api-section">
                <div class="card">
                    <div class="card-header">
                        <div class="card-title">Get User Following <span class="endpoint-badge">GET /api/user/{username}/following</span></div>
                    </div>
                    <div class="form-grid">
                        <div class="form-group">
                            <label>Username</label>
                            <input type="text" class="form-control" id="following-username" value="elonmusk">
                        </div>
                        <div class="form-group">
                            <label>Count (1-1000)</label>
                            <input type="number" class="form-control" id="following-count" value="10">
                        </div>
                    </div>
                    <button class="btn btn-primary" onclick="executeRequest('getFollowing')">Execute Request</button>
                </div>
            </div>

            <!-- User Followers -->
            <div id="user-followers" class="api-section">
                <div class="card">
                    <div class="card-header">
                        <div class="card-title">Get User Followers <span class="endpoint-badge">GET /api/user/{username}/followers</span></div>
                    </div>
                    <div class="form-grid">
                        <div class="form-group">
                            <label>Username</label>
                            <input type="text" class="form-control" id="followers-username" value="elonmusk">
                        </div>
                        <div class="form-group">
                            <label>Count (1-1000)</label>
                            <input type="number" class="form-control" id="followers-count" value="10">
                        </div>
                    </div>
                    <button class="btn btn-primary" onclick="executeRequest('getFollowers')">Execute Request</button>
                </div>
            </div>

            <!-- User Liked -->
            <div id="user-liked" class="api-section">
                <div class="card">
                    <div class="card-header">
                        <div class="card-title">Get Liked Tweets <span class="endpoint-badge">GET /api/user/{username}/liked</span></div>
                    </div>
                    <div class="form-grid">
                        <div class="form-group">
                            <label>Username</label>
                            <input type="text" class="form-control" id="liked-username" value="elonmusk">
                        </div>
                        <div class="form-group">
                            <label>Count (1-100)</label>
                            <input type="number" class="form-control" id="liked-count" value="10">
                        </div>
                    </div>
                    <button class="btn btn-primary" onclick="executeRequest('getLikedTweets')">Execute Request</button>
                </div>
            </div>

            <!-- User Mentions -->
            <div id="user-mentions" class="api-section">
                <div class="card">
                    <div class="card-header">
                        <div class="card-title">Get User Mentions <span class="endpoint-badge">GET /api/user/{username}/mentions</span></div>
                    </div>
                    <div class="form-grid">
                        <div class="form-group">
                            <label>Username</label>
                            <input type="text" class="form-control" id="mentions-username" value="elonmusk">
                        </div>
                        <div class="form-group">
                            <label>Count (1-100)</label>
                            <input type="number" class="form-control" id="mentions-count" value="10">
                        </div>
                    </div>
                    <button class="btn btn-primary" onclick="executeRequest('getMentions')">Execute Request</button>
                </div>
            </div>

            <!-- User Tweets -->
            <div id="user-tweets" class="api-section">
                <div class="card">
                    <div class="card-header">
                        <div class="card-title">Get User Tweets <span class="endpoint-badge">GET /api/tweets/user/{username}</span></div>
                    </div>
                    <div class="form-grid">
                        <div class="form-group">
                            <label>Username</label>
                            <input type="text" class="form-control" id="user-tweets-username" value="elonmusk">
                        </div>
                        <div class="form-group">
                            <label>Count (1-100)</label>
                            <input type="number" class="form-control" id="user-tweets-count" value="10">
                        </div>
                    </div>
                    <button class="btn btn-primary" onclick="executeRequest('getUserTweets')">Execute Request</button>
                </div>
            </div>

            <!-- Tweet Detail -->
            <div id="tweet-detail" class="api-section">
                <div class="card">
                    <div class="card-header">
                        <div class="card-title">Get Tweet Detail <span class="endpoint-badge">GET /api/tweets/{id}</span></div>
                    </div>
                    <div class="form-grid">
                        <div class="form-group">
                            <label>Tweet ID</label>
                            <input type="text" class="form-control" id="tweet-id" placeholder="e.g. 1234567890">
                        </div>
                    </div>
                    <button class="btn btn-primary" onclick="executeRequest('getTweetDetail')">Execute Request</button>
                </div>
            </div>

            <!-- Search Tweets -->
            <div id="search-tweets" class="api-section">
                <div class="card">
                    <div class="card-header">
                        <div class="card-title">Search Tweets <span class="endpoint-badge">GET /api/tweets/search</span></div>
                    </div>
                    <div class="form-grid">
                        <div class="form-group">
                            <label>Query</label>
                            <input type="text" class="form-control" id="search-tweets-query" value="golang">
                        </div>
                        <div class="form-group">
                            <label>Count (1-100)</label>
                            <input type="number" class="form-control" id="search-tweets-count" value="10">
                        </div>
                    </div>
                    <button class="btn btn-primary" onclick="executeRequest('searchTweets')">Execute Request</button>
                </div>
            </div>

            <!-- Tweet Counts -->
            <div id="tweet-counts" class="api-section">
                <div class="card">
                    <div class="card-header">
                        <div class="card-title">Tweet Counts <span class="endpoint-badge">GET /api/tweets/counts/recent</span></div>
                    </div>
                    <div class="form-grid">
                        <div class="form-group">
                            <label>Query</label>
                            <input type="text" class="form-control" id="tweet-counts-query" value="golang">
                        </div>
                        <div class="form-group">
                            <label>Start Time (RFC3339)</label>
                            <input type="text" class="form-control" id="tweet-counts-start" placeholder="e.g. 2024-01-01T00:00:00Z">
                        </div>
                    </div>
                    <button class="btn btn-primary" onclick="executeRequest('getTweetCounts')">Execute Request</button>
                </div>
            </div>

            <!-- Liking Users -->
            <div id="liking-users" class="api-section">
                <div class="card">
                    <div class="card-header">
                        <div class="card-title">Liking Users <span class="endpoint-badge">GET /api/tweets/{id}/liking_users</span></div>
                    </div>
                    <div class="form-grid">
                        <div class="form-group">
                            <label>Tweet ID</label>
                            <input type="text" class="form-control" id="liking-users-tweet-id">
                        </div>
                        <div class="form-group">
                            <label>Count</label>
                            <input type="number" class="form-control" id="liking-users-count" value="10">
                        </div>
                    </div>
                    <button class="btn btn-primary" onclick="executeRequest('getLikingUsers')">Execute Request</button>
                </div>
            </div>

            <!-- Quote Tweets -->
            <div id="quote-tweets" class="api-section">
                <div class="card">
                    <div class="card-header">
                        <div class="card-title">Quote Tweets <span class="endpoint-badge">GET /api/tweets/{id}/quote_tweets</span></div>
                    </div>
                    <div class="form-grid">
                        <div class="form-group">
                            <label>Tweet ID</label>
                            <input type="text" class="form-control" id="quote-tweets-tweet-id">
                        </div>
                        <div class="form-group">
                            <label>Count</label>
                            <input type="number" class="form-control" id="quote-tweets-count" value="10">
                        </div>
                    </div>
                    <button class="btn btn-primary" onclick="executeRequest('getQuoteTweets')">Execute Request</button>
                </div>
            </div>

            <!-- Retweeted By -->
            <div id="retweeted-by" class="api-section">
                <div class="card">
                    <div class="card-header">
                        <div class="card-title">Retweeted By <span class="endpoint-badge">GET /api/tweets/{id}/retweeted_by</span></div>
                    </div>
                    <div class="form-grid">
                        <div class="form-group">
                            <label>Tweet ID</label>
                            <input type="text" class="form-control" id="retweeted-by-tweet-id">
                        </div>
                        <div class="form-group">
                            <label>Count</label>
                            <input type="number" class="form-control" id="retweeted-by-count" value="10">
                        </div>
                    </div>
                    <button class="btn btn-primary" onclick="executeRequest('getRetweetedBy')">Execute Request</button>
                </div>
            </div>

        </div>
    </div>

    <!-- Response Sidebar -->
    <div class="response-area" id="response-area">
        <div class="response-header">
            <div class="response-title">Response Output</div>
            <button class="close-response" onclick="closeResponse()">×</button>
        </div>
        <div class="response-content">
            <pre id="response-output">Select an API and execute request to see response here...</pre>
        </div>
    </div>

    <!-- Loading -->
    <div class="loading-overlay" id="loading">
        <div class="spinner"></div>
    </div>

    <script>
        // Navigation
        function showSection(id) {
            document.querySelectorAll('.api-section').forEach(el => el.classList.remove('active'));
            document.getElementById(id).classList.add('active');
            
            document.querySelectorAll('.nav-item').forEach(el => el.classList.remove('active'));
            event.currentTarget.classList.add('active');

            // Update title
            const title = event.currentTarget.innerText.trim();
            document.getElementById('page-title').innerText = title;

            // Mobile: close sidebar
            if (window.innerWidth <= 768) {
                document.getElementById('sidebar').classList.remove('active');
            }
        }

        function toggleSidebar() {
            document.getElementById('sidebar').classList.toggle('active');
        }

        function closeResponse() {
            document.getElementById('response-area').classList.remove('open');
        }

        // API Execution
        async function executeRequest(type) {
            const loading = document.getElementById('loading');
            const responseArea = document.getElementById('response-area');
            const output = document.getElementById('response-output');
            
            loading.classList.add('active');
            responseArea.classList.add('open');
            output.innerHTML = 'Loading...';

            let url = '';
            
            try {
                switch(type) {
                    case 'getUserInfo':
                        const uName = document.getElementById('user-info-username').value;
                        if(!uName) throw new Error('Username is required');
                        url = '/api/user/' + uName;
                        break;
                    case 'getFollowing':
                        url = '/api/user/' + document.getElementById('following-username').value + '/following?count=' + document.getElementById('following-count').value;
                        break;
                    case 'getFollowers':
                        url = '/api/user/' + document.getElementById('followers-username').value + '/followers?count=' + document.getElementById('followers-count').value;
                        break;
                    case 'getLikedTweets':
                        url = '/api/user/' + document.getElementById('liked-username').value + '/liked?count=' + document.getElementById('liked-count').value;
                        break;
                    case 'getMentions':
                        url = '/api/user/' + document.getElementById('mentions-username').value + '/mentions?count=' + document.getElementById('mentions-count').value;
                        break;
                    case 'getUserTweets':
                        url = '/api/tweets/user/' + document.getElementById('user-tweets-username').value + '?count=' + document.getElementById('user-tweets-count').value;
                        break;
                    case 'getTweetDetail':
                        const tId = document.getElementById('tweet-id').value;
                        if(!tId) throw new Error('Tweet ID is required');
                        url = '/api/tweets/' + tId;
                        break;
                    case 'searchTweets':
                        url = '/api/tweets/search?q=' + encodeURIComponent(document.getElementById('search-tweets-query').value) + '&count=' + document.getElementById('search-tweets-count').value;
                        break;
                    case 'getTweetCounts':
                        let q = '/api/tweets/counts/recent?query=' + encodeURIComponent(document.getElementById('tweet-counts-query').value);
                        const start = document.getElementById('tweet-counts-start').value;
                        if(start) q += '&start_time=' + start;
                        url = q;
                        break;
                    case 'getLikingUsers':
                        url = '/api/tweets/' + document.getElementById('liking-users-tweet-id').value + '/liking_users?count=' + document.getElementById('liking-users-count').value;
                        break;
                    case 'getQuoteTweets':
                        url = '/api/tweets/' + document.getElementById('quote-tweets-tweet-id').value + '/quote_tweets?count=' + document.getElementById('quote-tweets-count').value;
                        break;
                    case 'getRetweetedBy':
                        url = '/api/tweets/' + document.getElementById('retweeted-by-tweet-id').value + '/retweeted_by?count=' + document.getElementById('retweeted-by-count').value;
                        break;
                }

                const res = await fetch(url);
                const data = await res.json();
                output.innerHTML = syntaxHighlight(data);

            } catch (err) {
                output.innerHTML = '<span style="color:red">Error: ' + err.message + '</span>';
            } finally {
                loading.classList.remove('active');
            }
        }

        function syntaxHighlight(json) {
            if (typeof json != 'string') {
                json = JSON.stringify(json, undefined, 2);
            }
            json = json.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
            return json.replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, function (match) {
                var cls = 'json-number';
                if (/^"/.test(match)) {
                    if (/:$/.test(match)) {
                        cls = 'json-key';
                    } else {
                        cls = 'json-string';
                    }
                } else if (/true|false/.test(match)) {
                    cls = 'json-boolean';
                } else if (/null/.test(match)) {
                    cls = 'json-null';
                }
                return '<span class="' + cls + '">' + match + '</span>';
            });
        }
    </script>
</body>
</html>`
}
