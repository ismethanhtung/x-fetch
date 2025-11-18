# 🚀 START HERE - Bắt Đầu Ngay!

## 👋 Chào mừng đến với X Twitter Backend API!

Đây là hướng dẫn **đầu tiên** bạn nên đọc. Chỉ mất **5 phút** để có server chạy!

---

## ⚡ Chạy Nhanh (5 phút)

### Bước 1: Lấy Twitter Bearer Token (2 phút)

1. 🔗 Truy cập: https://developer.twitter.com/
2. 🔑 Đăng nhập → Create App → Copy **Bearer Token**

> Nếu chưa có account, xem hướng dẫn chi tiết trong [TUTORIAL_VI.md](TUTORIAL_VI.md)

### Bước 2: Setup (1 phút)

```bash
# Cài dependencies
go mod download

# Tạo file config
cp ENV_EXAMPLE .env

# Mở .env và paste Bearer Token vào
# TWITTER_BEARER_TOKEN=your_token_here
```

### Bước 3: Chạy (30 giây)

```bash
go run main.go
```

Bạn sẽ thấy:
```
INFO[...] 🚀 Khởi động X Twitter Backend Server...
INFO[...] 🌐 Server đang lắng nghe...  address=localhost:8080
```

### Bước 4: Test (30 giây)

Mở browser hoặc terminal mới:

```bash
# Health check
curl http://localhost:8080/health

# Lấy tweets của Elon Musk
curl http://localhost:8080/api/tweets/user/elonmusk
```

## 🎉 Xong! Server đã chạy!

---

## 📚 Tiếp Theo - Đọc Gì?

### 🆕 Mới bắt đầu?
👉 **[QUICKSTART_VI.md](QUICKSTART_VI.md)** - Quick start guide chi tiết

### 📖 Muốn hiểu sâu?
👉 **[TUTORIAL_VI.md](TUTORIAL_VI.md)** - Tutorial từng bước, deployment guides

### 💻 Muốn code examples?
👉 **[EXAMPLES.md](EXAMPLES.md)** - Ví dụ với 7 ngôn ngữ

### 🏗️ Muốn hiểu kiến trúc?
👉 **[PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md)** - Chi tiết architecture

### 📊 Muốn overview?
👉 **[PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)** - Tổng quan toàn bộ project

### 📖 Muốn full docs?
👉 **[README.md](README.md)** - Complete documentation (đầy đủ nhất)

---

## 🎯 API Endpoints Cơ Bản

```bash
# 1. Health check
GET http://localhost:8080/health

# 2. Lấy thông tin user
GET http://localhost:8080/api/user/{username}

# 3. Lấy tweets
GET http://localhost:8080/api/tweets/user/{username}?count=10

# 4. Lấy danh sách following
GET http://localhost:8080/api/user/{username}/following?count=50

# 5. API docs
GET http://localhost:8080/api/docs
```

**Ví dụ:**
```bash
curl http://localhost:8080/api/user/elonmusk
curl http://localhost:8080/api/tweets/user/BillGates?count=5
curl http://localhost:8080/api/tweets/user/NASA?count=20
curl "http://localhost:8080/api/user/elonmusk/following?count=100"
```

---

## 🔧 Commands Hữu Ích

```bash
# Build
make build

# Run
make run

# Hot reload (cần cài air)
make dev

# Test
make test

# Help
make help

# Auto setup
./scripts/setup.sh

# Test API
./scripts/test-api.sh

# Monitor accounts
./scripts/monitor.sh
```

---

## 🐳 Chạy với Docker

```bash
# Build image
docker build -t twitter-backend .

# Run container
docker run -p 8080:8080 --env-file .env twitter-backend

# Hoặc dùng Docker Compose
docker-compose up -d
```

---

## ❗ Gặp Vấn Đề?

### Lỗi: "TWITTER_BEARER_TOKEN là bắt buộc"
✅ Check file `.env` có tồn tại và có token chưa

### Lỗi: "Port already in use"
✅ Đổi port trong `.env`: `SERVER_PORT=8081`

### Lỗi: 401 Unauthorized
✅ Check token có đúng không, regenerate nếu cần

### Lỗi khác?
✅ Xem [TUTORIAL_VI.md](TUTORIAL_VI.md) phần Troubleshooting

---

## 📁 Cấu Trúc Project (Quick View)

```
x-twitter-backend/
├── 📂 config/           # Configuration
├── 📂 handlers/         # HTTP handlers
├── 📂 models/           # Data structures
├── 📂 services/         # Business logic
├── 📂 scripts/          # Automation scripts
├── 📄 main.go           # Entry point
├── 🐳 Dockerfile        # Docker config
├── 🔧 Makefile          # Build commands
└── 📚 Docs (8 files)    # Documentation
```

---

## 🎯 Project Features

✅ **RESTful API** - 4 endpoints (tweets, user info, following list, health)
✅ **Twitter API v2** - Latest API integration
✅ **Clean Architecture** - Professional code structure
✅ **Docker Ready** - Containerized deployment
✅ **Hot Reload** - Development với Air
✅ **Comprehensive Docs** - 8 documentation files
✅ **Production Ready** - Logging, error handling, security
✅ **Easy Setup** - 5 minutes from zero to running

---

## 📞 Cần Trợ Giúp?

- 📖 **Documentation**: 8 files trong project
- 🔧 **Commands**: `make help`
- 🌐 **API Docs**: http://localhost:8080/api/docs
- 🐛 **Issues**: Tạo GitHub issue

---

## 🎓 Tài Liệu Theo Cấp Độ

### Beginner 🌱
1. START_HERE.md (you are here!)
2. QUICKSTART_VI.md
3. README.md (phần Quick Start)

### Intermediate 🌿
1. TUTORIAL_VI.md
2. EXAMPLES.md
3. README.md (full)

### Advanced 🌳
1. PROJECT_STRUCTURE.md
2. CONTRIBUTING.md
3. Source code với comments

---

## 💡 Tips

- 💾 **Save Bearer Token**: Giữ token an toàn, không commit
- 📊 **Rate Limits**: Free tier có giới hạn 300 requests/15 phút
- 🔄 **Hot Reload**: Dùng `make dev` khi development
- 🐳 **Docker**: Recommend cho production
- 📝 **Logs**: Check logs để debug

---

## 🎉 Bắt Đầu Thôi!

```bash
# 1. Setup
cp ENV_EXAMPLE .env
# Edit .env, thêm token

# 2. Run
go run main.go

# 3. Test
curl http://localhost:8080/api/tweets/user/elonmusk
curl "http://localhost:8080/api/user/elonmusk/following?count=50"

# 🎊 Done!
```

---

**Happy Coding! 🚀**

_Nếu thấy project hữu ích, hãy star ⭐ trên GitHub!_

