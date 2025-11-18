# X Twitter Backend API 🐦

Server backend chuyên nghiệp được xây dựng bằng Golang để lấy tweets và thông tin người dùng từ X/Twitter. API này cho phép bạn dễ dàng lấy các bài đăng mới nhất từ bất kỳ tài khoản Twitter/X nào.

## ✨ Tính năng

- 🔍 Lấy thông tin chi tiết về user Twitter/X
- 📱 Lấy tweets mới nhất từ bất kỳ tài khoản nào
- 📊 Hiển thị metrics (likes, retweets, replies, views)
- 👥 Lấy danh sách các tài khoản mà user đang theo dõi (following list)
- 🏷️ Trích xuất hashtags, mentions, và URLs từ tweets
- 🔐 Authentication an toàn với Twitter API v2
- 📝 Logging chi tiết và middleware
- 🚀 Performance cao với Golang
- ⚡ Graceful shutdown
- 🛡️ Error handling chuyên nghiệp
- 🌐 CORS support
- 📖 API documentation tích hợp

## 🏗️ Kiến trúc

```
x-twitter-backend/
├── config/              # Configuration management
│   └── config.go
├── handlers/            # HTTP handlers và middleware
│   ├── tweets_handler.go
│   └── middleware.go
├── models/              # Data structures
│   └── tweet.go
├── services/            # Business logic
│   └── twitter_service.go
├── main.go              # Entry point
├── go.mod               # Dependencies
├── .gitignore
├── ENV_EXAMPLE          # Environment variables template
└── README.md
```

## 📋 Yêu cầu

- **Go**: 1.21 hoặc cao hơn
- **Twitter Developer Account**: Để lấy Bearer Token
- **Twitter API v2 Access**: Free tier hoặc paid tier

## 🚀 Cài đặt

### 1. Clone hoặc download project

```bash
cd x-twitter-backend
```

### 2. Cài đặt dependencies

```bash
go mod download
```

### 3. Cấu hình Twitter API

#### Bước 1: Tạo Twitter Developer Account

1. Truy cập [Twitter Developer Portal](https://developer.twitter.com/)
2. Đăng ký tài khoản developer (miễn phí)
3. Tạo một project mới
4. Tạo một app trong project đó

#### Bước 2: Lấy Bearer Token

1. Trong Twitter Developer Portal, vào phần "Keys and tokens"
2. Tạo hoặc copy **Bearer Token**
3. Lưu token này để sử dụng

### 4. Thiết lập Environment Variables

Tạo file `.env` từ template:

```bash
cp ENV_EXAMPLE .env
```

Chỉnh sửa file `.env` và thêm Bearer Token của bạn:

```env
# Twitter API Configuration
TWITTER_BEARER_TOKEN=your_actual_bearer_token_here

# Server Configuration
SERVER_PORT=8080
SERVER_HOST=localhost

# Application Configuration
APP_ENV=development
LOG_LEVEL=info

# Rate Limiting
MAX_TWEETS_PER_REQUEST=100
DEFAULT_TWEETS_COUNT=10
```

### 5. Chạy server

```bash
go run main.go
```

Hoặc build và chạy:

```bash
go build -o twitter-backend
./twitter-backend
```

Server sẽ chạy tại `http://localhost:8080`

## 📖 API Documentation

### 1. Health Check

Kiểm tra trạng thái server.

**Endpoint:** `GET /health`

**Response:**

```json
{
  "status": "ok",
  "service": "X Twitter Backend API",
  "version": "1.0.0"
}
```

**Ví dụ:**

```bash
curl http://localhost:8080/health
```

### 2. Lấy thông tin User

Lấy thông tin chi tiết về một user Twitter/X.

**Endpoint:** `GET /api/user/{username}`

**Parameters:**

- `username` (path): Username của tài khoản Twitter/X

**Response:**

```json
{
  "id": "44196397",
  "username": "elonmusk",
  "name": "Elon Musk",
  "description": "Tesla, SpaceX, Twitter",
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

**Ví dụ:**

```bash
# Lấy thông tin Elon Musk
curl http://localhost:8080/api/user/elonmusk

# Lấy thông tin Bill Gates
curl http://localhost:8080/api/user/BillGates
```

### 3. Lấy Tweets của User

Lấy các tweets mới nhất từ một user.

**Endpoint:** `GET /api/tweets/user/{username}`

**Parameters:**

- `username` (path): Username của tài khoản Twitter/X
- `count` (query, optional): Số lượng tweets cần lấy (default: 10, max: 100)

**Response:**

```json
{
  "tweets": [
    {
      "id": "1234567890",
      "text": "This is a tweet...",
      "author_id": "44196397",
      "created_at": "2024-01-15T10:30:00Z",
      "metrics": {
        "retweet_count": 1000,
        "reply_count": 500,
        "like_count": 5000,
        "quote_count": 200
      },
      "entities": {
        "hashtags": [
          {"tag": "AI"}
        ],
        "mentions": [
          {"username": "someone", "id": "123"}
        ],
        "urls": [
          {
            "url": "https://t.co/xyz",
            "expanded_url": "https://example.com",
            "display_url": "example.com"
          }
        ]
      }
    }
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

**Ví dụ:**

```bash
# Lấy 10 tweets mới nhất của Elon Musk (default)
curl http://localhost:8080/api/tweets/user/elonmusk

# Lấy 20 tweets mới nhất
curl http://localhost:8080/api/tweets/user/elonmusk?count=20

# Lấy 50 tweets của Bill Gates
curl "http://localhost:8080/api/tweets/user/BillGates?count=50"
```

### 4. Lấy danh sách tài khoản đang theo dõi

Lấy toàn bộ danh sách accounts mà user đang theo dõi (following list).

**Endpoint:** `GET /api/user/{username}/following`

**Parameters:**

- `username` (path): Username của tài khoản Twitter/X
- `count` (query, optional): Số lượng accounts cần lấy (default: 10, max: 1000)
- `pagination_token` (query, optional): Token để lấy trang kế tiếp (nếu kết quả nhiều hơn giới hạn)

**Response:**

```json
{
  "user": {
    "id": "44196397",
    "username": "elonmusk",
    "name": "Elon Musk",
    "metrics": {
      "following_count": 170
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
  ],
  "meta": {
    "result_count": 100,
    "next_token": "7140dibdnow9c7obbb09mjwo7xgns86sp4l83vr0b8npg",
    "previous_token": ""
  }
}
```

**Ví dụ:**

```bash
# Lấy 100 accounts mà Elon Musk đang theo dõi
curl "http://localhost:8080/api/user/elonmusk/following?count=100"

# Lấy trang tiếp theo bằng pagination token
curl "http://localhost:8080/api/user/elonmusk/following?count=100&pagination_token=YOUR_TOKEN"
```

### 5. API Documentation

Xem tài liệu API tích hợp.

**Endpoint:** `GET /api/docs`

**Ví dụ:**

```bash
curl http://localhost:8080/api/docs
```

## 🔧 Configuration

### Environment Variables

| Variable                 | Mô tả                                    | Default     | Required |
| ------------------------ | ---------------------------------------- | ----------- | -------- |
| `TWITTER_BEARER_TOKEN`   | Bearer token từ Twitter Developer Portal | -           | ✅ Yes   |
| `SERVER_PORT`            | Port để chạy server                      | 8080        | No       |
| `SERVER_HOST`            | Host để bind server                      | localhost   | No       |
| `APP_ENV`                | Environment (development/production)     | development | No       |
| `LOG_LEVEL`              | Log level (debug/info/warn/error)        | info        | No       |
| `MAX_TWEETS_PER_REQUEST` | Số lượng tweets tối đa mỗi request       | 100         | No       |
| `DEFAULT_TWEETS_COUNT`   | Số lượng tweets mặc định                 | 10          | No       |

## 🛠️ Development

### Build

```bash
go build -o twitter-backend
```

### Run với hot reload (sử dụng air)

```bash
# Cài đặt air
go install github.com/cosmtrek/air@latest

# Chạy với hot reload
air
```

### Testing

```bash
go test ./...
```

### Linting

```bash
# Cài đặt golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
```

## 📊 Dependencies

- **github.com/gorilla/mux** - HTTP router và URL matcher
- **github.com/michimani/gotwi** - Twitter API v2 client cho Go
- **github.com/joho/godotenv** - Load environment variables từ .env
- **github.com/sirupsen/logrus** - Structured logger

## 🔐 Security

### Best Practices

1. **Không commit Bearer Token**: File `.env` đã được thêm vào `.gitignore`
2. **Rate Limiting**: API tuân thủ rate limits của Twitter
3. **Error Handling**: Không expose sensitive information trong errors
4. **CORS**: Có thể configure CORS cho production
5. **HTTPS**: Nên sử dụng HTTPS trong production

### Production Deployment

Khi deploy lên production:

1. Sử dụng environment variables thay vì file `.env`
2. Enable HTTPS
3. Configure proper CORS origins
4. Set up monitoring và logging
5. Sử dụng reverse proxy (nginx, Caddy)
6. Implement rate limiting ở application level

## 🚦 Rate Limits

Twitter API có các rate limits sau (Free tier):

- **User lookup**: 300 requests / 15 phút
- **User tweets**: 1,500 requests / 15 phút

Server này tuân thủ các rate limits của Twitter API. Nếu bạn vượt quá rate limit, API sẽ trả về error từ Twitter.

## 📝 Examples

### Ví dụ với cURL

```bash
# 1. Health check
curl http://localhost:8080/health

# 2. Lấy thông tin Elon Musk
curl http://localhost:8080/api/user/elonmusk

# 3. Lấy 15 tweets mới nhất của Elon Musk
curl "http://localhost:8080/api/tweets/user/elonmusk?count=15"

# 4. Lấy tweets của nhiều users khác
curl http://localhost:8080/api/tweets/user/BillGates
curl http://localhost:8080/api/tweets/user/NASA
curl http://localhost:8080/api/tweets/user/cristiano
```

### Ví dụ với JavaScript (fetch)

```javascript
// Lấy tweets của Elon Musk
async function getElonTweets() {
  try {
    const response = await fetch(
      "http://localhost:8080/api/tweets/user/elonmusk?count=10"
    );
    const data = await response.json();

    console.log(`User: ${data.user.name} (@${data.user.username})`);
    console.log(`Followers: ${data.user.metrics.followers_count}`);
    console.log(`\nTweets:`);

    data.tweets.forEach((tweet, index) => {
      console.log(`\n${index + 1}. ${tweet.text}`);
      console.log(
        `   ❤️ ${tweet.metrics.like_count} | 🔄 ${tweet.metrics.retweet_count}`
      );
    });
  } catch (error) {
    console.error("Error:", error);
  }
}

getElonTweets();
```

### Ví dụ với Python

```python
import requests

def get_user_tweets(username, count=10):
    url = f"http://localhost:8080/api/tweets/user/{username}"
    params = {"count": count}

    response = requests.get(url, params=params)
    data = response.json()

    print(f"User: {data['user']['name']} (@{data['user']['username']})")
    print(f"Followers: {data['user']['metrics']['followers_count']:,}")
    print(f"\nTweets:")

    for i, tweet in enumerate(data['tweets'], 1):
        print(f"\n{i}. {tweet['text']}")
        print(f"   ❤️ {tweet['metrics']['like_count']} | 🔄 {tweet['metrics']['retweet_count']}")

# Sử dụng
get_user_tweets("elonmusk", count=5)
```

## ❗ Troubleshooting

### Lỗi: "TWITTER_BEARER_TOKEN là bắt buộc"

**Giải pháp**: Đảm bảo bạn đã:

1. Tạo file `.env`
2. Thêm Bearer Token vào file `.env`
3. Bearer Token hợp lệ và chưa expire

### Lỗi: "Unauthorized" hoặc 401

**Giải pháp**:

1. Kiểm tra Bearer Token có đúng không
2. Đảm bảo Twitter app của bạn có quyền truy cập API v2
3. Kiểm tra xem token có bị revoke không

### Lỗi: Rate limit exceeded

**Giải pháp**:

1. Đợi 15 phút để rate limit reset
2. Cân nhắc upgrade lên paid tier của Twitter API
3. Implement caching ở application level

### Server không start được

**Giải pháp**:

1. Kiểm tra port 8080 có bị sử dụng không
2. Thay đổi `SERVER_PORT` trong `.env`
3. Kiểm tra logs để xem lỗi cụ thể

## 🤝 Contributing

Mọi đóng góp đều được hoan nghênh! Vui lòng:

1. Fork project
2. Tạo feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Mở Pull Request

## 📄 License

Project này được phát triển cho mục đích học tập và sử dụng cá nhân.

## 👤 Author

Được tạo bởi AI Assistant với ❤️ cho người dùng.

## 🙏 Acknowledgments

- [Twitter API v2](https://developer.twitter.com/en/docs/twitter-api)
- [gotwi](https://github.com/michimani/gotwi) - Thư viện Twitter API v2 cho Go
- [Gorilla Mux](https://github.com/gorilla/mux)
- [Logrus](https://github.com/sirupsen/logrus)

## 📞 Support

Nếu bạn gặp vấn đề hoặc có câu hỏi, vui lòng:

1. Đọc phần Troubleshooting
2. Kiểm tra Twitter API documentation
3. Tạo issue trên GitHub

---

**Happy Coding! 🚀**
# x-fetch
