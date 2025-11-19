# 🚀 Hướng dẫn sử dụng X/Twitter API

## 📋 Tổng quan

API Backend này cung cấp **9 endpoints miễn phí** để tương tác với X/Twitter, bao gồm:

✅ Lấy thông tin user  
✅ Lấy tweets của user  
✅ Tìm kiếm tweets  
✅ Xem chi tiết tweet  
✅ Danh sách following  
✅ Danh sách followers  
✅ Tweets đã like  
✅ Mentions  
✅ Tìm kiếm users  

---

## 🎯 Các API có sẵn

### 1. 👤 Lấy thông tin User

**Endpoint:** `GET /api/user/{username}`

**Mô tả:** Lấy thông tin chi tiết của một tài khoản Twitter/X

**Ví dụ:**
```bash
curl http://localhost:8080/api/user/elonmusk
```

**Response:**
```json
{
  "id": "44196397",
  "username": "elonmusk",
  "name": "Elon Musk",
  "description": "Tesla, SpaceX, Neuralink",
  "profile_image_url": "https://...",
  "verified": true,
  "created_at": "2009-06-02T20:12:29Z",
  "metrics": {
    "followers_count": 150000000,
    "following_count": 500,
    "tweet_count": 30000
  }
}
```

---

### 2. 📝 Lấy Tweets của User

**Endpoint:** `GET /api/tweets/user/{username}?count=10`

**Mô tả:** Lấy danh sách tweets mới nhất của một user

**Parameters:**
- `username` (required): Username của tài khoản
- `count` (optional): Số lượng tweets (default: 10, max: 100)

**Ví dụ:**
```bash
curl "http://localhost:8080/api/tweets/user/elonmusk?count=5"
```

**Response:**
```json
{
  "tweets": [
    {
      "id": "1234567890",
      "text": "Tweet content...",
      "author_id": "44196397",
      "created_at": "2024-01-15T10:30:00Z",
      "metrics": {
        "retweet_count": 1000,
        "reply_count": 500,
        "like_count": 5000,
        "quote_count": 200
      }
    }
  ],
  "user": {...},
  "meta": {
    "result_count": 5
  }
}
```

---

### 3. 🔍 Tìm kiếm Tweets

**Endpoint:** `GET /api/tweets/search?q=keyword&count=10`

**Mô tả:** Tìm kiếm tweets theo từ khóa

**Parameters:**
- `q` (required): Từ khóa tìm kiếm
- `count` (optional): Số lượng tweets (default: 10, max: 100)

**Ví dụ:**
```bash
curl "http://localhost:8080/api/tweets/search?q=golang&count=20"
```

**Search Tips:**
- Tìm chính xác: `"exact phrase"`
- Loại trừ từ: `keyword -excluded`
- Từ người dùng: `from:username`
- Hashtag: `#hashtag`
- Kết hợp: `golang OR python #programming`

---

### 4. 📄 Xem chi tiết Tweet

**Endpoint:** `GET /api/tweets/{tweet_id}`

**Mô tả:** Lấy thông tin chi tiết của một tweet cụ thể

**Parameters:**
- `tweet_id` (required): ID của tweet

**Ví dụ:**
```bash
curl http://localhost:8080/api/tweets/1234567890
```

**Response:**
```json
{
  "tweet": {
    "id": "1234567890",
    "text": "Tweet content...",
    "author_id": "44196397",
    "created_at": "2024-01-15T10:30:00Z",
    "metrics": {...},
    "entities": {
      "hashtags": [...],
      "mentions": [...],
      "urls": [...]
    }
  },
  "author": {
    "id": "44196397",
    "username": "elonmusk",
    "name": "Elon Musk",
    ...
  }
}
```

---

### 5. ➕ Danh sách Following

**Endpoint:** `GET /api/user/{username}/following?count=10`

**Mô tả:** Lấy danh sách tài khoản mà user đang theo dõi

**Parameters:**
- `username` (required): Username của tài khoản
- `count` (optional): Số lượng (default: 10, max: 1000)
- `pagination_token` (optional): Token để lấy trang tiếp theo

**Ví dụ:**
```bash
curl "http://localhost:8080/api/user/elonmusk/following?count=50"
```

**Response:**
```json
{
  "user": {...},
  "following": [
    {
      "id": "123",
      "username": "user1",
      "name": "User 1",
      ...
    }
  ],
  "meta": {
    "result_count": 50,
    "next_token": "..."
  }
}
```

---

### 6. 👥 Danh sách Followers

**Endpoint:** `GET /api/user/{username}/followers?count=10`

**Mô tả:** Lấy danh sách người đang theo dõi user

**Parameters:**
- `username` (required): Username của tài khoản
- `count` (optional): Số lượng (default: 10, max: 1000)
- `pagination_token` (optional): Token để lấy trang tiếp theo

**Ví dụ:**
```bash
curl "http://localhost:8080/api/user/elonmusk/followers?count=100"
```

---

### 7. ❤️ Tweets đã Like

**Endpoint:** `GET /api/user/{username}/liked?count=10`

**Mô tả:** Lấy danh sách tweets mà user đã like

**Parameters:**
- `username` (required): Username của tài khoản
- `count` (optional): Số lượng (default: 10, max: 100)

**Ví dụ:**
```bash
curl "http://localhost:8080/api/user/elonmusk/liked?count=20"
```

**Response:**
```json
{
  "user": {...},
  "tweets": [...],
  "meta": {
    "result_count": 20
  }
}
```

---

### 8. @ Mentions

**Endpoint:** `GET /api/user/{username}/mentions?count=10`

**Mô tả:** Lấy danh sách tweets có mention đến user

**Parameters:**
- `username` (required): Username của tài khoản
- `count` (optional): Số lượng (default: 10, max: 100)

**Ví dụ:**
```bash
curl "http://localhost:8080/api/user/elonmusk/mentions?count=15"
```

---

### 9. 🔎 Tìm kiếm Users

**Endpoint:** `GET /api/users/search?q=keyword&count=10`

**Mô tả:** Tìm kiếm users theo từ khóa

**Parameters:**
- `q` (required): Từ khóa tìm kiếm
- `count` (optional): Số lượng (default: 10, max: 100)

**Ví dụ:**
```bash
curl "http://localhost:8080/api/users/search?q=elon&count=10"
```

**Lưu ý:** API này tìm kiếm users thông qua tweets có chứa từ khóa, vì Twitter API v2 không hỗ trợ user search trực tiếp với Bearer token.

---

## 🧪 Trang Test

Truy cập **http://localhost:8080** hoặc **http://localhost:8080/test** để sử dụng giao diện test đẹp mắt với:

✨ 9 API cards với form nhập liệu  
✨ Hiển thị response JSON đẹp mắt  
✨ Loading animation  
✨ Error handling  
✨ Responsive design  

---

## 🚀 Cách chạy

### 1. Cấu hình

Tạo file `.env`:
```bash
cp ENV_EXAMPLE .env
```

Chỉnh sửa `.env`:
```
TWITTER_BEARER_TOKEN=your_bearer_token_here
SERVER_HOST=localhost
SERVER_PORT=8080
APP_ENV=development
LOG_LEVEL=info
```

### 2. Chạy server

```bash
# Build
go build -o bin/x-twitter-backend

# Chạy
./bin/x-twitter-backend
```

Hoặc dùng Makefile:
```bash
make run
```

### 3. Test API

Truy cập: http://localhost:8080

---

## 📊 API Documentation

Xem full API documentation tại:
```
GET http://localhost:8080/api/docs
```

---

## ⚠️ Lưu ý

1. **Rate Limits**: API tuân thủ rate limits của Twitter API v2
2. **Bearer Token**: Cần có Twitter Bearer Token hợp lệ
3. **Public Data**: Chỉ truy cập được dữ liệu public
4. **Free Tier**: Tất cả API đều miễn phí với Twitter API v2 Basic access

---

## 🛠️ Tech Stack

- **Language:** Go 1.21+
- **Framework:** Gorilla Mux
- **Twitter Library:** github.com/michimani/gotwi v0.14.0
- **Logging:** Logrus
- **Config:** godotenv

---

## 📝 Examples với JavaScript

### Fetch API
```javascript
// Lấy thông tin user
fetch('http://localhost:8080/api/user/elonmusk')
  .then(res => res.json())
  .then(data => console.log(data));

// Tìm kiếm tweets
fetch('http://localhost:8080/api/tweets/search?q=golang&count=20')
  .then(res => res.json())
  .then(data => console.log(data));
```

### Axios
```javascript
import axios from 'axios';

// Lấy tweets
const tweets = await axios.get('http://localhost:8080/api/tweets/user/elonmusk', {
  params: { count: 10 }
});

// Tìm kiếm users
const users = await axios.get('http://localhost:8080/api/users/search', {
  params: { q: 'elon', count: 5 }
});
```

---

## 🎓 Error Handling

Tất cả API trả về error với format:
```json
{
  "error": "ERROR_CODE",
  "message": "Chi tiết lỗi",
  "code": 400
}
```

**Common Error Codes:**
- `400` - Bad Request (thiếu parameters)
- `404` - Not Found (không tìm thấy user/tweet)
- `500` - Internal Server Error
- `429` - Rate Limit Exceeded

---

## 📧 Support

Nếu có vấn đề, vui lòng tạo issue trên GitHub hoặc liên hệ team.

**Happy Coding! 🎉**

