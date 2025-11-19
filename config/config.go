package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

// Config chứa tất cả các cấu hình của ứng dụng
type Config struct {
	// Twitter API
	TwitterBearerToken string

	// Server
	ServerPort string
	ServerHost string

	// Application
	AppEnv   string
	LogLevel string

	// Rate Limiting
	MaxTweetsPerRequest  int
	DefaultTweetsCount   int
}

var AppConfig *Config

// LoadConfig đọc và load configuration từ environment variables
func LoadConfig() (*Config, error) {
	// Load .env file nếu tồn tại (không bắt buộc trong production)
	if err := godotenv.Load(); err != nil {
		log.Warn("Không tìm thấy file .env, sẽ sử dụng environment variables từ hệ thống")
	}

	// Detect nếu đang chạy trong container
	// Kiểm tra nhiều cách để detect container environment
	isContainer := os.Getenv("CONTAINER") == "true" || 
		os.Getenv("DOCKER") == "true" ||
		os.Getenv("KUBERNETES_SERVICE_HOST") != "" ||
		fileExists("/.dockerenv")

	// Default host: 0.0.0.0 cho container, 0.0.0.0 cho tất cả (để đảm bảo hoạt động trong container)
	defaultHost := "0.0.0.0"
	if isContainer {
		log.Info("🔍 Phát hiện đang chạy trong container, đảm bảo SERVER_HOST=0.0.0.0")
	}

	serverHost := getEnv("SERVER_HOST", defaultHost)
	
	// Force 0.0.0.0 nếu đang trong container và host là localhost
	if isContainer && (serverHost == "localhost" || serverHost == "127.0.0.1") {
		log.Warnf("⚠️  Đang trong container nhưng SERVER_HOST=%s, tự động chuyển sang 0.0.0.0", serverHost)
		serverHost = "0.0.0.0"
	}

	config := &Config{
		TwitterBearerToken:  getEnv("TWITTER_BEARER_TOKEN", ""),
		ServerPort:          getEnv("SERVER_PORT", "8080"),
		ServerHost:          serverHost,
		AppEnv:              getEnv("APP_ENV", "development"),
		LogLevel:            getEnv("LOG_LEVEL", "info"),
		MaxTweetsPerRequest: getEnvAsInt("MAX_TWEETS_PER_REQUEST", 100),
		DefaultTweetsCount:  getEnvAsInt("DEFAULT_TWEETS_COUNT", 10),
	}

	// Validate required fields
	if config.TwitterBearerToken == "" {
		return nil, fmt.Errorf("TWITTER_BEARER_TOKEN là bắt buộc")
	}

	AppConfig = config
	return config, nil
}

// getEnv đọc một environment variable hoặc trả về giá trị mặc định
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt đọc environment variable dạng integer
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Warnf("Không thể parse %s thành integer, sử dụng giá trị mặc định: %d", key, defaultValue)
		return defaultValue
	}
	return value
}

// GetAddress trả về địa chỉ server đầy đủ
func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%s", c.ServerHost, c.ServerPort)
}

// fileExists kiểm tra xem file có tồn tại không
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

