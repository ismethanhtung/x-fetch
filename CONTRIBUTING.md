# 🤝 Contributing to X Twitter Backend

Cảm ơn bạn đã quan tâm đến việc đóng góp cho project! Mọi contribution đều được hoan nghênh.

## 📋 Table of Contents

1. [Code of Conduct](#code-of-conduct)
2. [Getting Started](#getting-started)
3. [How to Contribute](#how-to-contribute)
4. [Development Guidelines](#development-guidelines)
5. [Pull Request Process](#pull-request-process)
6. [Coding Standards](#coding-standards)
7. [Testing](#testing)
8. [Documentation](#documentation)

## 📜 Code of Conduct

### Our Pledge

Chúng tôi cam kết tạo một môi trường thân thiện, chuyên nghiệp và tôn trọng cho tất cả mọi người.

### Expected Behavior

- ✅ Sử dụng ngôn ngữ thân thiện và bao dung
- ✅ Tôn trọng quan điểm và kinh nghiệm khác nhau
- ✅ Chấp nhận phê bình mang tính xây dựng
- ✅ Tập trung vào điều tốt nhất cho cộng đồng
- ✅ Thể hiện sự đồng cảm với các thành viên khác

### Unacceptable Behavior

- ❌ Ngôn ngữ hoặc hình ảnh mang tính khiêu dâm
- ❌ Trolling, bình luận xúc phạm hoặc công격 cá nhân
- ❌ Quấy rối công khai hoặc riêng tư
- ❌ Công bố thông tin riêng tư của người khác
- ❌ Hành vi không chuyên nghiệp khác

## 🚀 Getting Started

### Prerequisites

- Go 1.21 hoặc cao hơn
- Git
- Twitter Developer Account (để test)
- Code editor (VS Code, GoLand, etc.)

### Setup Development Environment

1. **Fork repository**
   ```bash
   # Click "Fork" button trên GitHub
   ```

2. **Clone fork của bạn**
   ```bash
   git clone https://github.com/YOUR_USERNAME/x-twitter-backend.git
   cd x-twitter-backend
   ```

3. **Add upstream remote**
   ```bash
   git remote add upstream https://github.com/ORIGINAL_OWNER/x-twitter-backend.git
   ```

4. **Install dependencies**
   ```bash
   go mod download
   ```

5. **Setup environment**
   ```bash
   cp ENV_EXAMPLE .env
   # Edit .env và thêm Twitter Bearer Token
   ```

6. **Run server**
   ```bash
   go run main.go
   ```

## 🔨 How to Contribute

### Types of Contributions

1. **🐛 Bug Reports**
   - Tìm thấy bug? Tạo issue với label `bug`
   - Bao gồm: steps to reproduce, expected vs actual behavior, screenshots

2. **✨ Feature Requests**
   - Có ý tưởng mới? Tạo issue với label `enhancement`
   - Mô tả feature, use cases, và potential implementation

3. **📝 Documentation**
   - Cải thiện docs, fix typos, add examples
   - Documentation rất quan trọng!

4. **🔧 Code Contributions**
   - Fix bugs
   - Implement new features
   - Refactor code
   - Improve performance

### Before You Start

1. **Check existing issues**
   - Tìm xem issue đã tồn tại chưa
   - Comment nếu bạn muốn work on it

2. **Create or comment on issue**
   - Discuss approach trước khi code
   - Get feedback từ maintainers

3. **Keep PRs focused**
   - One feature/fix per PR
   - Easier to review và merge

## 💻 Development Guidelines

### Branch Naming

```
feature/add-caching        # New features
bugfix/fix-cors-issue      # Bug fixes
docs/update-readme         # Documentation
refactor/improve-service   # Code refactoring
test/add-unit-tests        # Tests
```

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add Redis caching layer
fix: resolve CORS preflight issue
docs: update API examples
refactor: simplify error handling
test: add service layer tests
chore: update dependencies
```

**Format:**
```
<type>: <subject>

<body (optional)>

<footer (optional)>
```

**Types:**
- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation
- `style` - Formatting
- `refactor` - Code restructuring
- `test` - Tests
- `chore` - Maintenance

### Code Style

#### Go Code Style

1. **Follow Go standards**
   ```bash
   gofmt -s -w .
   go vet ./...
   ```

2. **Use meaningful names**
   ```go
   // ✅ Good
   func GetUserTweets(username string, count int) (*TweetsResponse, error)
   
   // ❌ Bad
   func gut(u string, c int) (*TR, error)
   ```

3. **Add comments**
   ```go
   // GetUserTweets retrieves recent tweets from a specific user
   // It fetches user info first, then retrieves tweets with metrics
   func (s *TwitterService) GetUserTweets(ctx context.Context, username string, maxResults int) (*models.TweetsResponse, error) {
       // Implementation
   }
   ```

4. **Handle errors properly**
   ```go
   // ✅ Good
   if err != nil {
       log.WithError(err).Error("Failed to fetch tweets")
       return nil, fmt.Errorf("failed to fetch tweets: %w", err)
   }
   
   // ❌ Bad
   if err != nil {
       return nil, err
   }
   ```

5. **Use contexts**
   ```go
   // ✅ Good
   func (s *Service) FetchData(ctx context.Context) error {
       // Use ctx for timeouts, cancellation
   }
   ```

#### Project Structure

- Put models trong `models/`
- Put business logic trong `services/`
- Put HTTP handlers trong `handlers/`
- Put configuration trong `config/`
- Keep `main.go` minimal

## 🔄 Pull Request Process

### 1. Update Your Fork

```bash
git checkout main
git fetch upstream
git merge upstream/main
git push origin main
```

### 2. Create Feature Branch

```bash
git checkout -b feature/your-feature-name
```

### 3. Make Changes

- Write code
- Follow coding standards
- Add tests
- Update documentation

### 4. Test Your Changes

```bash
# Run tests
go test ./...

# Run linter
golangci-lint run

# Test manually
go run main.go
```

### 5. Commit Changes

```bash
git add .
git commit -m "feat: add your feature"
```

### 6. Push to Fork

```bash
git push origin feature/your-feature-name
```

### 7. Create Pull Request

1. Go to your fork trên GitHub
2. Click "New Pull Request"
3. Fill in template:
   - Description of changes
   - Related issues
   - Screenshots (if UI)
   - Testing done

### 8. Address Review Comments

- Respond to feedback
- Make requested changes
- Push updates

### 9. Merge

- Maintainer sẽ merge sau khi approve
- Celebrate! 🎉

## ✅ Pull Request Checklist

- [ ] Code follows project style
- [ ] Self-review completed
- [ ] Comments added for complex code
- [ ] Documentation updated
- [ ] Tests added/updated
- [ ] All tests pass
- [ ] No linter errors
- [ ] Commit messages follow convention
- [ ] PR description is clear

## 🧪 Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./services/...

# Verbose output
go test -v ./...
```

### Writing Tests

```go
package services

import (
    "context"
    "testing"
)

func TestGetUserTweets(t *testing.T) {
    // Setup
    service := NewTwitterService(mockConfig)
    
    // Execute
    result, err := service.GetUserTweets(context.Background(), "testuser", 10)
    
    // Assert
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    
    if len(result.Tweets) == 0 {
        t.Error("Expected tweets, got none")
    }
}
```

### Test Coverage Goals

- Aim for > 80% coverage
- Focus on critical paths
- Test error cases
- Test edge cases

## 📚 Documentation

### What to Document

1. **Code Comments**
   - Public functions
   - Complex logic
   - Non-obvious decisions

2. **README Updates**
   - New features
   - Changed behavior
   - New dependencies

3. **API Documentation**
   - New endpoints
   - Changed parameters
   - New responses

4. **Examples**
   - Usage examples
   - Integration examples
   - Common patterns

### Documentation Style

```go
// GetUserTweets retrieves recent tweets from a Twitter user.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - username: Twitter username (without @)
//   - maxResults: Maximum number of tweets to retrieve (1-100)
//
// Returns:
//   - *TweetsResponse: Contains tweets, user info, and metadata
//   - error: Error if request fails
//
// Example:
//   response, err := service.GetUserTweets(ctx, "elonmusk", 10)
//   if err != nil {
//       log.Fatal(err)
//   }
func (s *TwitterService) GetUserTweets(ctx context.Context, username string, maxResults int) (*models.TweetsResponse, error) {
    // Implementation
}
```

## 🐛 Bug Report Template

```markdown
**Bug Description**
A clear and concise description of the bug.

**Steps to Reproduce**
1. Go to '...'
2. Call API with '...'
3. See error

**Expected Behavior**
What you expected to happen.

**Actual Behavior**
What actually happened.

**Environment**
- OS: [e.g., macOS 14.0]
- Go version: [e.g., 1.21]
- Server version: [e.g., 1.0.0]

**Logs**
```
Paste relevant logs here
```

**Screenshots**
If applicable, add screenshots.
```

## ✨ Feature Request Template

```markdown
**Feature Description**
Clear description of the feature.

**Use Case**
Why is this feature needed? What problem does it solve?

**Proposed Solution**
How you think it should work.

**Alternatives Considered**
Other solutions you've thought about.

**Additional Context**
Any other context, screenshots, or examples.
```

## 💡 Tips for Contributors

### Quick Wins

- Fix typos trong documentation
- Improve error messages
- Add code comments
- Add examples trong EXAMPLES.md
- Improve logging

### Good First Issues

Look for issues labeled:
- `good first issue`
- `help wanted`
- `documentation`

### Communication

- Be patient and respectful
- Ask questions if unclear
- Provide context trong comments
- Update issues with progress

## 🏆 Recognition

Contributors sẽ được:
- Listed trong README
- Mentioned trong release notes
- Our eternal gratitude! 🙏

## 📞 Getting Help

- **Questions**: Open a GitHub Discussion
- **Bugs**: Open a GitHub Issue
- **Feature Ideas**: Open a GitHub Issue
- **Chat**: [Discord/Slack link] (if available)

## 📜 License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

**Thank you for contributing! 🎉**

Every contribution, no matter how small, makes this project better.

