# 📚 Hướng Dẫn Chi Tiết - X Twitter Backend API

Hướng dẫn từng bước để setup và sử dụng X Twitter Backend API.

## 📖 Mục Lục

1. [Cài đặt và Setup](#1-cài-đặt-và-setup)
2. [Lấy Twitter API Credentials](#2-lấy-twitter-api-credentials)
3. [Chạy Server](#3-chạy-server)
4. [Sử dụng API](#4-sử-dụng-api)
5. [Deployment](#5-deployment)
6. [Các Tính Năng Nâng Cao](#6-các-tính-năng-nâng-cao)

## 1. Cài Đặt và Setup

### Bước 1.1: Cài đặt Go

Nếu bạn chưa có Go, cài đặt từ [golang.org](https://golang.org/dl/)

Kiểm tra version:

```bash
go version
# Output: go version go1.21.x ...
```

### Bước 1.2: Clone/Download Project

```bash
cd /path/to/your/projects
# Nếu có git
git clone <repository-url>
cd x-twitter-backend

# Hoặc download và extract zip
cd x-twitter-backend
```

### Bước 1.3: Cài đặt Dependencies

```bash
# Sử dụng Go modules
go mod download

# Hoặc sử dụng Makefile
make install
```

## 2. Lấy Twitter API Credentials

### Bước 2.1: Đăng Ký Twitter Developer Account

1. **Truy cập Twitter Developer Portal**

   - Vào [developer.twitter.com](https://developer.twitter.com/)
   - Đăng nhập bằng tài khoản Twitter của bạn

2. **Sign Up cho Developer Account**

   - Click "Sign up" hoặc "Apply"
   - Điền form đăng ký:
     - Tên
     - Email
     - Country
     - Use case (chọn "Exploring the API" hoặc "Building tools")
   - Mô tả ngắn gọn về dự án của bạn
   - Đồng ý với Terms and Conditions
   - Submit application

3. **Xác Nhận Email**
   - Check email và confirm

### Bước 2.2: Tạo Twitter App

1. **Tạo Project Mới**

   - Trong Developer Portal, click "Projects & Apps"
   - Click "Create Project"
   - Đặt tên project (ví dụ: "Tweet Fetcher")
   - Chọn use case phù hợp
   - Mô tả project

2. **Tạo App trong Project**

   - Sau khi tạo project, click "Create App"
   - Đặt tên app (ví dụ: "Tweet Backend")
   - App sẽ được tạo và bạn sẽ thấy API Keys

3. **Lưu API Keys**

   - API Key
   - API Key Secret
   - **Bearer Token** (quan trọng nhất!)

   ⚠️ **LƯU Ý**: Copy và lưu Bearer Token ngay, bạn sẽ không thể xem lại!

### Bước 2.3: Setup Permissions

1. Trong app settings, vào "User authentication settings"
2. Đảm bảo app có quyền:
   - Read tweets
   - Read users
3. Save changes

### Bước 2.4: Configure Environment

Tạo file `.env` từ template:

```bash
cp ENV_EXAMPLE .env
```

Mở file `.env` và thêm Bearer Token:

```env
TWITTER_BEARER_TOKEN=AAAAAAAAAAAAAAAAAAAAABearerTokenOfYours...

SERVER_PORT=8080
SERVER_HOST=localhost
APP_ENV=development
LOG_LEVEL=info
MAX_TWEETS_PER_REQUEST=100
DEFAULT_TWEETS_COUNT=10
```

## 3. Chạy Server

### Cách 1: Chạy Trực Tiếp

```bash
go run main.go
```

Output:

```
INFO[2024-01-15 10:30:00] 🚀 Khởi động X Twitter Backend Server...
INFO[2024-01-15 10:30:00] ✅ Configuration đã được load  port=8080 host=localhost
INFO[2024-01-15 10:30:00] Twitter client đã được khởi tạo thành công
INFO[2024-01-15 10:30:00] ✅ Routes đã được thiết lập
INFO[2024-01-15 10:30:00] 🌐 Server đang lắng nghe...  address=localhost:8080
```

### Cách 2: Build và Chạy

```bash
# Build
go build -o twitter-backend

# Chạy
./twitter-backend
```

### Cách 3: Sử dụng Makefile

```bash
# Xem tất cả commands
make help

# Cài đặt dependencies
make install

# Build
make build

# Chạy
make run
```

### Cách 4: Hot Reload (Development)

```bash
# Cài đặt air
go install github.com/cosmtrek/air@latest

# Chạy với hot reload
make dev
# Hoặc
air
```

### Cách 5: Docker

```bash
# Build Docker image
docker build -t twitter-backend .

# Chạy với Docker
docker run -p 8080:8080 --env-file .env twitter-backend

# Hoặc sử dụng Docker Compose
docker-compose up -d
```

## 4. Sử dụng API

### 4.1: Test Server

Kiểm tra xem server đã chạy chưa:

```bash
curl http://localhost:8080/health
```

Response:

```json
{
  "status": "ok",
  "service": "X Twitter Backend API",
  "version": "1.0.0"
}
```

### 4.2: Lấy Thông Tin User

**Request:**

```bash
curl http://localhost:8080/api/user/elonmusk
```

**Response:**

```json
{
  "id": "44196397",
  "username": "elonmusk",
  "name": "Elon Musk",
  "description": "...",
  "profile_image_url": "https://...",
  "verified": true,
  "created_at": "2009-06-02T20:12:29Z",
  "metrics": {
    "followers_count": 168000000,
    "following_count": 500,
    "tweet_count": 35000,
    "listed_count": 120000
  }
}
```

### 4.3: Lấy Tweets của User

**Request với default count (10 tweets):**

```bash
curl http://localhost:8080/api/tweets/user/elonmusk
```

**Request với custom count:**

```bash
curl "http://localhost:8080/api/tweets/user/elonmusk?count=20"
```

**Response:**

```json
{
  "tweets": [
    {
      "id": "1234567890",
      "text": "Mars is looking good today! 🚀",
      "author_id": "44196397",
      "created_at": "2024-01-15T10:30:00Z",
      "metrics": {
        "retweet_count": 5000,
        "reply_count": 1200,
        "like_count": 50000,
        "quote_count": 800
      },
      "entities": {
        "hashtags": [],
        "mentions": [],
        "urls": []
      }
    }
    // ... more tweets
  ],
  "user": {
    "id": "44196397",
    "username": "elonmusk",
    "name": "Elon Musk",
    ...
  },
  "meta": {
    "result_count": 10
  }
}
```

### 4.4: Lấy danh sách tài khoản đang theo dõi

**Request:**

```bash
curl "http://localhost:8080/api/user/elonmusk/following?count=100"
```

**Response (rút gọn):**

```json
{
  "user": {
    "id": "44196397",
    "username": "elonmusk",
    "name": "Elon Musk",
    "metrics": {
      "following_count": 1700
    }
  },
  "following": [
    {
      "id": "20536157",
      "username": "SpaceX",
      "name": "SpaceX",
      "verified": true,
      "profile_image_url": "https://..."
    }
    // ...
  ],
  "meta": {
    "result_count": 100,
    "next_token": "7140dibdnow9c7...",
    "previous_token": ""
  }
}
```

- `count` có thể từ 1 → 1000 (Twitter API cho phép tối đa 1000)
- `pagination_token` dùng để lấy trang kế tiếp (nếu có `next_token` trong `meta`)

### 4.5: Các Ví Dụ Khác

**Lấy tweets của Bill Gates:**

```bash
curl http://localhost:8080/api/tweets/user/BillGates
```

**Lấy tweets của NASA:**

```bash
curl http://localhost:8080/api/tweets/user/NASA
```

**Lấy 50 tweets:**

```bash
curl "http://localhost:8080/api/tweets/user/cristiano?count=50"
```

### 4.6: Sử dụng với Postman

1. Mở Postman
2. Tạo request mới
3. Method: GET
4. URL: `http://localhost:8080/api/tweets/user/elonmusk`
5. Params (optional):
   - Key: `count`
   - Value: `20`
6. Send

### 4.7: Sử dụng với Browser

Mở browser và truy cập:

- `http://localhost:8080/health`
- `http://localhost:8080/api/user/elonmusk`
- `http://localhost:8080/api/tweets/user/elonmusk?count=10`
- `http://localhost:8080/api/user/elonmusk/following?count=50`

## 5. Deployment

### 5.1: Deploy lên Server Linux (VPS)

**1. Upload files lên server:**

```bash
scp -r ./x-twitter-backend user@your-server.com:/home/user/
```

**2. SSH vào server:**

```bash
ssh user@your-server.com
cd /home/user/x-twitter-backend
```

**3. Build và chạy:**

```bash
go build -o twitter-backend
./twitter-backend
```

**4. Chạy như service với systemd:**

Tạo file `/etc/systemd/system/twitter-backend.service`:

```ini
[Unit]
Description=Twitter Backend API
After=network.target

[Service]
Type=simple
User=your-user
WorkingDirectory=/home/user/x-twitter-backend
EnvironmentFile=/home/user/x-twitter-backend/.env
ExecStart=/home/user/x-twitter-backend/twitter-backend
Restart=always

[Install]
WantedBy=multi-user.target
```

Start service:

```bash
sudo systemctl daemon-reload
sudo systemctl start twitter-backend
sudo systemctl enable twitter-backend
sudo systemctl status twitter-backend
```

### 5.2: Deploy với Docker

**1. Build image:**

```bash
docker build -t twitter-backend:v1.0 .
```

**2. Run container:**

```bash
docker run -d \
  --name twitter-backend \
  -p 8080:8080 \
  --env-file .env \
  --restart unless-stopped \
  twitter-backend:v1.0
```

### 5.3: Deploy với Nginx Reverse Proxy

**Nginx config:**

```nginx
server {
    listen 80;
    server_name api.yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 5.4: Deploy lên Heroku

**1. Tạo Procfile:**

```
web: ./twitter-backend
```

**2. Deploy:**

```bash
heroku create your-app-name
heroku config:set TWITTER_BEARER_TOKEN=your_token
git push heroku main
```

## 6. Các Tính Năng Nâng Cao

### 6.1: Logging

Server log tất cả requests. Xem logs:

**Nếu chạy trực tiếp:**
Logs sẽ hiển thị trong terminal

**Nếu chạy như service:**

```bash
sudo journalctl -u twitter-backend -f
```

**Thay đổi log level:**
Trong `.env`:

```env
LOG_LEVEL=debug  # debug, info, warn, error
```

### 6.2: Monitoring

Sử dụng health check endpoint:

```bash
# Check mỗi 30 giây
watch -n 30 curl http://localhost:8080/health
```

### 6.3: Rate Limiting

API tự động tuân thủ Twitter rate limits. Nếu vượt quá:

- Wait 15 phút
- Hoặc upgrade Twitter API tier

### 6.4: Caching (Tùy chọn)

Bạn có thể thêm Redis cache để giảm calls đến Twitter API.

### 6.5: Multiple Accounts Monitoring

Tạo script để monitor nhiều accounts:

```bash
#!/bin/bash
# monitor-accounts.sh

accounts=("elonmusk" "BillGates" "NASA" "cristiano")

for account in "${accounts[@]}"; do
    echo "Fetching tweets for @$account..."
    curl -s "http://localhost:8080/api/tweets/user/$account?count=5" | jq '.tweets[].text'
    echo "---"
done
```

Chạy:

```bash
chmod +x monitor-accounts.sh
./monitor-accounts.sh
```

## 🎯 Tips & Best Practices

### Performance

1. **Caching**: Implement caching cho repeated requests
2. **Connection Pooling**: Go HTTP client đã tự động handle
3. **Concurrent Requests**: Go goroutines handle tự động

### Security

1. **Không commit .env**: Đã có trong .gitignore
2. **HTTPS in production**: Sử dụng Let's Encrypt
3. **Environment variables**: Dùng secrets manager trong production
4. **Rate limiting**: Implement application-level rate limiting

### Monitoring

1. **Health checks**: Sử dụng `/health` endpoint
2. **Metrics**: Thêm Prometheus metrics
3. **Alerts**: Setup alerts khi service down
4. **Logs**: Centralized logging với ELK stack

### Scaling

1. **Horizontal scaling**: Deploy multiple instances với load balancer
2. **Database**: Thêm database để cache tweets
3. **Queue**: Sử dụng message queue cho async processing
4. **CDN**: Cache static responses

## ❓ FAQs

**Q: API có miễn phí không?**
A: Server code miễn phí, nhưng bạn cần Twitter API access (có free tier).

**Q: Có thể lấy tweets cũ hơn không?**
A: Có, sử dụng pagination với `next_token` (cần implement).

**Q: Rate limit là bao nhiêu?**
A: Free tier: 300 user lookups, 1,500 tweets per 15 phút.

**Q: Có thể lấy tweets real-time không?**
A: Cần upgrade lên Twitter API v2 Elevated hoặc sử dụng Streaming API.

**Q: Server có thể handle bao nhiêu requests?**
A: Phụ thuộc vào server resources, nhưng Go rất performant. Bottleneck thường là Twitter API rate limits.

## 🆘 Troubleshooting Common Issues

### Issue 1: "TWITTER_BEARER_TOKEN là bắt buộc"

✅ Giải pháp:

```bash
# Kiểm tra .env file
cat .env

# Đảm bảo có TWITTER_BEARER_TOKEN
# Không có khoảng trắng quanh dấu =
TWITTER_BEARER_TOKEN=AAAAAAAAAyour_token_here
```

### Issue 2: 401 Unauthorized

✅ Giải pháp:

- Kiểm tra Bearer Token còn valid không
- Regenerate token trong Twitter Developer Portal
- Đảm bảo app có đúng permissions

### Issue 3: Port already in use

✅ Giải pháp:

```bash
# Tìm process đang dùng port 8080
lsof -i :8080

# Kill process
kill -9 <PID>

# Hoặc đổi port trong .env
SERVER_PORT=8081
```

### Issue 4: Rate limit exceeded

✅ Giải pháp:

- Đợi 15 phút
- Giảm số requests
- Implement caching
- Upgrade Twitter API tier

---

**Chúc bạn thành công! 🎉**

Nếu có vấn đề, đọc kỹ error messages trong logs và check Twitter API documentation.
