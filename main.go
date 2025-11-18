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

	// Tweets routes
	api.HandleFunc("/tweets/user/{username}", tweetsHandler.GetUserTweets).Methods("GET")

	// API documentation endpoint
	api.HandleFunc("/docs", handleAPIDocs).Methods("GET")

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
  "version": "1.0.0",
  "description": "API để lấy tweets và thông tin user từ X/Twitter",
  "endpoints": [
    {
      "path": "/health",
      "method": "GET",
      "description": "Health check endpoint",
      "response": {
        "status": "ok",
        "service": "X Twitter Backend API",
        "version": "1.0.0"
      }
    },
    {
      "path": "/api/user/{username}",
      "method": "GET",
      "description": "Lấy thông tin user theo username",
      "parameters": {
        "username": "Username của tài khoản Twitter/X (ví dụ: elonmusk)"
      },
      "example": "/api/user/elonmusk"
    },
    {
      "path": "/api/tweets/user/{username}",
      "method": "GET",
      "description": "Lấy tweets mới nhất của một user",
      "parameters": {
        "username": "Username của tài khoản Twitter/X",
        "count": "Số lượng tweets cần lấy (optional, default: 10, max: 100)"
      },
      "example": "/api/tweets/user/elonmusk?count=20"
    },
    {
      "path": "/api/user/{username}/following",
      "method": "GET",
      "description": "Lấy danh sách tài khoản mà user đang theo dõi",
      "parameters": {
        "username": "Username của tài khoản Twitter/X",
        "count": "Số lượng accounts cần lấy (optional, default: 10, max: 1000)",
        "pagination_token": "Token để lấy trang tiếp theo (optional)"
      },
      "example": "/api/user/elonmusk/following?count=100"
    }
  ],
  "authentication": "Yêu cầu TWITTER_BEARER_TOKEN trong environment variables",
  "notes": [
    "API tuân thủ rate limits của Twitter API",
    "Tất cả responses trả về dạng JSON",
    "Errors được trả về với format chuẩn: {error, message, code}"
  ]
}`

	w.Write([]byte(docs))
}
