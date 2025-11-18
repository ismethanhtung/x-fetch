# ⚡ Quick Start - Chạy Nhanh trong 5 Phút

Hướng dẫn nhanh nhất để chạy X Twitter Backend API.

## 🎯 Bước 1: Lấy Twitter Bearer Token (2 phút)

1. Vào [developer.twitter.com](https://developer.twitter.com/)
2. Đăng nhập và tạo app (nếu chưa có)
3. Vào "Keys and tokens" → Copy **Bearer Token**

## 🔧 Bước 2: Setup Project (1 phút)

```bash
# Di chuyển vào thư mục project
cd x-twitter-backend

# Cài đặt dependencies
go mod download

# Tạo file .env
cp ENV_EXAMPLE .env
```

Mở file `.env` và paste Bearer Token:

```env
TWITTER_BEARER_TOKEN=paste_your_bearer_token_here
```

## 🚀 Bước 3: Chạy Server (1 phút)

```bash
go run main.go
```

Bạn sẽ thấy:

```
INFO[...] 🚀 Khởi động X Twitter Backend Server...
INFO[...] 🌐 Server đang lắng nghe...  address=localhost:8080
```

## ✅ Bước 4: Test API (1 phút)

Mở browser hoặc terminal mới:

**Test health check:**

```bash
curl http://localhost:8080/health
```

**Lấy tweets của Elon Musk:**

```bash
curl http://localhost:8080/api/tweets/user/elonmusk
```

**Lấy danh sách followings của Elon Musk:**

```bash
curl "http://localhost:8080/api/user/elonmusk/following?count=50"
```

Hoặc mở trong browser:

- `http://localhost:8080/api/user/elonmusk`
- `http://localhost:8080/api/tweets/user/elonmusk?count=10`

## 🎉 Xong!

API của bạn đã chạy! Các endpoints có sẵn:

- `GET /health` - Health check
- `GET /api/user/{username}` - Thông tin user
- `GET /api/tweets/user/{username}?count=10` - Tweets của user
- `GET /api/user/{username}/following?count=50` - Accounts mà user đang theo dõi

## 📚 Tiếp Theo

- Đọc [README.md](README.md) để biết full documentation
- Đọc [TUTORIAL_VI.md](TUTORIAL_VI.md) cho hướng dẫn chi tiết
- Xem [API Documentation](http://localhost:8080/api/docs)

## ⚠️ Lưu Ý Nhanh

- **Rate Limits**: Free tier có giới hạn requests
- **Bearer Token**: Giữ bí mật, không commit lên git
- **Port**: Mặc định 8080, đổi trong `.env` nếu cần

## 🐳 Chạy với Docker (Tùy Chọn)

```bash
# Build
docker build -t twitter-backend .

# Run
docker run -p 8080:8080 --env-file .env twitter-backend
```

---

**Cần trợ giúp?** Xem [TUTORIAL_VI.md](TUTORIAL_VI.md) hoặc [README.md](README.md)
