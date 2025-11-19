package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"x-twitter-backend/models"
	"x-twitter-backend/services"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// TweetsHandler xử lý các HTTP requests liên quan đến tweets
type TweetsHandler struct {
	twitterService *services.TwitterService
}

// NewTweetsHandler tạo một instance mới của TweetsHandler
func NewTweetsHandler(twitterService *services.TwitterService) *TweetsHandler {
	return &TweetsHandler{
		twitterService: twitterService,
	}
}

// GetUserTweets xử lý request lấy tweets của một user
// GET /api/tweets/user/{username}?count=10
func (h *TweetsHandler) GetUserTweets(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	if username == "" {
		h.respondWithError(w, http.StatusBadRequest, "Username là bắt buộc", "MISSING_USERNAME")
		return
	}

	count := parseCount(r.URL.Query().Get("count"))

	log.WithFields(log.Fields{
		"username": username,
		"count":    count,
		"ip":       r.RemoteAddr,
	}).Info("Nhận request lấy tweets")

	// Gọi service để lấy tweets
	response, err := h.twitterService.GetUserTweets(r.Context(), username, count)
	if err != nil {
		log.WithError(err).Error("Lỗi khi lấy tweets")
		h.respondWithError(w, http.StatusInternalServerError, "Không thể lấy tweets: "+err.Error(), "FETCH_ERROR")
		return
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

// GetUserInfo xử lý request lấy thông tin user
// GET /api/user/{username}
func (h *TweetsHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	if username == "" {
		h.respondWithError(w, http.StatusBadRequest, "Username là bắt buộc", "MISSING_USERNAME")
		return
	}

	log.WithFields(log.Fields{
		"username": username,
		"ip":       r.RemoteAddr,
	}).Info("Nhận request lấy thông tin user")

	// Gọi service để lấy thông tin user
	user, err := h.twitterService.GetUserByUsername(r.Context(), username)
	if err != nil {
		log.WithError(err).Error("Lỗi khi lấy thông tin user")
		h.respondWithError(w, http.StatusInternalServerError, "Không thể lấy thông tin user: "+err.Error(), "FETCH_ERROR")
		return
	}

	h.respondWithJSON(w, http.StatusOK, user)
}

// GetUserFollowing xử lý request lấy danh sách tài khoản mà user đang theo dõi
// GET /api/user/{username}/following?count=50&pagination_token=xxx
func (h *TweetsHandler) GetUserFollowing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	if username == "" {
		h.respondWithError(w, http.StatusBadRequest, "Username là bắt buộc", "MISSING_USERNAME")
		return
	}

	count := parseCount(r.URL.Query().Get("count"))
	paginationToken := r.URL.Query().Get("pagination_token")

	log.WithFields(log.Fields{
		"username":   username,
		"count":      count,
		"page_token": paginationToken,
		"ip":         r.RemoteAddr,
	}).Info("Nhận request lấy danh sách following")

	response, err := h.twitterService.GetUserFollowing(r.Context(), username, count, paginationToken)
	if err != nil {
		log.WithError(err).Error("Lỗi khi lấy danh sách following")
		h.respondWithError(w, http.StatusInternalServerError, "Không thể lấy danh sách following: "+err.Error(), "FETCH_ERROR")
		return
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

// HealthCheck xử lý health check request
// GET /health
func (h *TweetsHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":  "ok",
		"service": "X Twitter Backend API",
		"version": "1.0.0",
	}
	h.respondWithJSON(w, http.StatusOK, response)
}

// respondWithJSON gửi JSON response
func (h *TweetsHandler) respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.WithError(err).Error("Lỗi khi encode JSON response")
	}
}

// respondWithError gửi error response
func (h *TweetsHandler) respondWithError(w http.ResponseWriter, statusCode int, message, code string) {
	errorResponse := models.ErrorResponse{
		Error:   code,
		Message: message,
		Code:    statusCode,
	}
	h.respondWithJSON(w, statusCode, errorResponse)
}

// GetUserFollowers xử lý request lấy danh sách followers
// GET /api/user/{username}/followers?count=50&pagination_token=xxx
func (h *TweetsHandler) GetUserFollowers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	if username == "" {
		h.respondWithError(w, http.StatusBadRequest, "Username là bắt buộc", "MISSING_USERNAME")
		return
	}

	count := parseCount(r.URL.Query().Get("count"))
	paginationToken := r.URL.Query().Get("pagination_token")

	log.WithFields(log.Fields{
		"username":   username,
		"count":      count,
		"page_token": paginationToken,
		"ip":         r.RemoteAddr,
	}).Info("Nhận request lấy danh sách followers")

	response, err := h.twitterService.GetUserFollowers(r.Context(), username, count, paginationToken)
	if err != nil {
		log.WithError(err).Error("Lỗi khi lấy danh sách followers")
		h.respondWithError(w, http.StatusInternalServerError, "Không thể lấy danh sách followers: "+err.Error(), "FETCH_ERROR")
		return
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

// SearchTweets xử lý request tìm kiếm tweets
// GET /api/tweets/search?q=golang&count=20
func (h *TweetsHandler) SearchTweets(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	if query == "" {
		h.respondWithError(w, http.StatusBadRequest, "Query là bắt buộc", "MISSING_QUERY")
		return
	}

	count := parseCount(r.URL.Query().Get("count"))

	log.WithFields(log.Fields{
		"query": query,
		"count": count,
		"ip":    r.RemoteAddr,
	}).Info("Nhận request tìm kiếm tweets")

	response, err := h.twitterService.SearchTweets(r.Context(), query, count)
	if err != nil {
		log.WithError(err).Error("Lỗi khi tìm kiếm tweets")
		h.respondWithError(w, http.StatusInternalServerError, "Không thể tìm kiếm tweets: "+err.Error(), "SEARCH_ERROR")
		return
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

// GetTweetByID xử lý request lấy chi tiết tweet
// GET /api/tweets/{tweet_id}
func (h *TweetsHandler) GetTweetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tweetID := vars["tweet_id"]

	if tweetID == "" {
		h.respondWithError(w, http.StatusBadRequest, "Tweet ID là bắt buộc", "MISSING_TWEET_ID")
		return
	}

	log.WithFields(log.Fields{
		"tweet_id": tweetID,
		"ip":       r.RemoteAddr,
	}).Info("Nhận request lấy chi tiết tweet")

	response, err := h.twitterService.GetTweetByID(r.Context(), tweetID)
	if err != nil {
		log.WithError(err).Error("Lỗi khi lấy chi tiết tweet")
		h.respondWithError(w, http.StatusInternalServerError, "Không thể lấy chi tiết tweet: "+err.Error(), "FETCH_ERROR")
		return
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

// GetLikedTweets xử lý request lấy liked tweets
// GET /api/user/{username}/liked?count=20
func (h *TweetsHandler) GetLikedTweets(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	if username == "" {
		h.respondWithError(w, http.StatusBadRequest, "Username là bắt buộc", "MISSING_USERNAME")
		return
	}

	count := parseCount(r.URL.Query().Get("count"))

	log.WithFields(log.Fields{
		"username": username,
		"count":    count,
		"ip":       r.RemoteAddr,
	}).Info("Nhận request lấy liked tweets")

	response, err := h.twitterService.GetLikedTweets(r.Context(), username, count)
	if err != nil {
		log.WithError(err).Error("Lỗi khi lấy liked tweets")
		h.respondWithError(w, http.StatusInternalServerError, "Không thể lấy liked tweets: "+err.Error(), "FETCH_ERROR")
		return
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

// SearchUsers xử lý request tìm kiếm users
// GET /api/users/search?q=elon&count=10
func (h *TweetsHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	if query == "" {
		h.respondWithError(w, http.StatusBadRequest, "Query là bắt buộc", "MISSING_QUERY")
		return
	}

	count := parseCount(r.URL.Query().Get("count"))

	log.WithFields(log.Fields{
		"query": query,
		"count": count,
		"ip":    r.RemoteAddr,
	}).Info("Nhận request tìm kiếm users")

	response, err := h.twitterService.SearchUsers(r.Context(), query, count)
	if err != nil {
		log.WithError(err).Error("Lỗi khi tìm kiếm users")
		h.respondWithError(w, http.StatusInternalServerError, "Không thể tìm kiếm users: "+err.Error(), "SEARCH_ERROR")
		return
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

// GetUserMentions xử lý request lấy mentions
// GET /api/user/{username}/mentions?count=20
func (h *TweetsHandler) GetUserMentions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	if username == "" {
		h.respondWithError(w, http.StatusBadRequest, "Username là bắt buộc", "MISSING_USERNAME")
		return
	}

	count := parseCount(r.URL.Query().Get("count"))

	log.WithFields(log.Fields{
		"username": username,
		"count":    count,
		"ip":       r.RemoteAddr,
	}).Info("Nhận request lấy mentions")

	response, err := h.twitterService.GetUserMentions(r.Context(), username, count)
	if err != nil {
		log.WithError(err).Error("Lỗi khi lấy mentions")
		h.respondWithError(w, http.StatusInternalServerError, "Không thể lấy mentions: "+err.Error(), "FETCH_ERROR")
		return
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

// ListTweets xử lý request lấy danh sách tweets theo IDs
// GET /api/tweets?ids=123,456,789
func (h *TweetsHandler) ListTweets(w http.ResponseWriter, r *http.Request) {
	idsParam := r.URL.Query().Get("ids")
	if idsParam == "" {
		h.respondWithError(w, http.StatusBadRequest, "ids là bắt buộc (comma-separated)", "MISSING_IDS")
		return
	}

	// Parse IDs từ comma-separated string
	ids := []string{}
	for _, id := range splitComma(idsParam) {
		if id != "" {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		h.respondWithError(w, http.StatusBadRequest, "Không có ID hợp lệ", "INVALID_IDS")
		return
	}

	log.WithFields(log.Fields{
		"ids_count": len(ids),
		"ip":        r.RemoteAddr,
	}).Info("Nhận request lấy danh sách tweets")

	response, err := h.twitterService.ListTweets(r.Context(), ids)
	if err != nil {
		log.WithError(err).Error("Lỗi khi lấy danh sách tweets")
		h.respondWithError(w, http.StatusInternalServerError, "Không thể lấy danh sách tweets: "+err.Error(), "FETCH_ERROR")
		return
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

// GetLikingUsers xử lý request lấy users đã like tweet
// GET /api/tweets/{tweet_id}/liking_users?count=10
func (h *TweetsHandler) GetLikingUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tweetID := vars["tweet_id"]

	if tweetID == "" {
		h.respondWithError(w, http.StatusBadRequest, "Tweet ID là bắt buộc", "MISSING_TWEET_ID")
		return
	}

	count := parseCount(r.URL.Query().Get("count"))
	paginationToken := r.URL.Query().Get("pagination_token")

	log.WithFields(log.Fields{
		"tweet_id": tweetID,
		"count":    count,
		"ip":       r.RemoteAddr,
	}).Info("Nhận request lấy liking users")

	response, err := h.twitterService.GetLikingUsers(r.Context(), tweetID, count, paginationToken)
	if err != nil {
		log.WithError(err).Error("Lỗi khi lấy liking users")
		h.respondWithError(w, http.StatusInternalServerError, "Không thể lấy liking users: "+err.Error(), "FETCH_ERROR")
		return
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

// GetQuoteTweets xử lý request lấy quote tweets
// GET /api/tweets/{tweet_id}/quote_tweets?count=10
func (h *TweetsHandler) GetQuoteTweets(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tweetID := vars["tweet_id"]

	if tweetID == "" {
		h.respondWithError(w, http.StatusBadRequest, "Tweet ID là bắt buộc", "MISSING_TWEET_ID")
		return
	}

	count := parseCount(r.URL.Query().Get("count"))

	log.WithFields(log.Fields{
		"tweet_id": tweetID,
		"count":    count,
		"ip":       r.RemoteAddr,
	}).Info("Nhận request lấy quote tweets")

	response, err := h.twitterService.GetQuoteTweets(r.Context(), tweetID, count)
	if err != nil {
		log.WithError(err).Error("Lỗi khi lấy quote tweets")
		h.respondWithError(w, http.StatusInternalServerError, "Không thể lấy quote tweets: "+err.Error(), "FETCH_ERROR")
		return
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

// GetRetweetedBy xử lý request lấy users đã retweet
// GET /api/tweets/{tweet_id}/retweeted_by?count=10
func (h *TweetsHandler) GetRetweetedBy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tweetID := vars["tweet_id"]

	if tweetID == "" {
		h.respondWithError(w, http.StatusBadRequest, "Tweet ID là bắt buộc", "MISSING_TWEET_ID")
		return
	}

	count := parseCount(r.URL.Query().Get("count"))
	paginationToken := r.URL.Query().Get("pagination_token")

	log.WithFields(log.Fields{
		"tweet_id": tweetID,
		"count":    count,
		"ip":       r.RemoteAddr,
	}).Info("Nhận request lấy retweeted by")

	response, err := h.twitterService.GetRetweetedBy(r.Context(), tweetID, count, paginationToken)
	if err != nil {
		log.WithError(err).Error("Lỗi khi lấy retweeted by")
		h.respondWithError(w, http.StatusInternalServerError, "Không thể lấy retweeted by: "+err.Error(), "FETCH_ERROR")
		return
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

// GetTweetCounts xử lý request lấy tweet counts
// GET /api/tweets/counts/recent?q=keyword&start_time=2024-01-01T00:00:00Z&end_time=2024-01-02T00:00:00Z
func (h *TweetsHandler) GetTweetCounts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		h.respondWithError(w, http.StatusBadRequest, "Query là bắt buộc", "MISSING_QUERY")
		return
	}

	startTime := r.URL.Query().Get("start_time")
	endTime := r.URL.Query().Get("end_time")

	log.WithFields(log.Fields{
		"query":      query,
		"start_time": startTime,
		"end_time":   endTime,
		"ip":         r.RemoteAddr,
	}).Info("Nhận request lấy tweet counts")

	response, err := h.twitterService.GetTweetCounts(r.Context(), query, startTime, endTime)
	if err != nil {
		log.WithError(err).Error("Lỗi khi lấy tweet counts")
		h.respondWithError(w, http.StatusInternalServerError, "Không thể lấy tweet counts: "+err.Error(), "FETCH_ERROR")
		return
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

// GetUserByID xử lý request lấy user theo ID
// GET /api/users/{user_id}
func (h *TweetsHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if userID == "" {
		h.respondWithError(w, http.StatusBadRequest, "User ID là bắt buộc", "MISSING_USER_ID")
		return
	}

	log.WithFields(log.Fields{
		"user_id": userID,
		"ip":      r.RemoteAddr,
	}).Info("Nhận request lấy user theo ID")

	user, err := h.twitterService.GetUserByID(r.Context(), userID)
	if err != nil {
		log.WithError(err).Error("Lỗi khi lấy user theo ID")
		h.respondWithError(w, http.StatusInternalServerError, "Không thể lấy user: "+err.Error(), "FETCH_ERROR")
		return
	}

	h.respondWithJSON(w, http.StatusOK, user)
}

// ListUsers xử lý request lấy danh sách users theo IDs
// GET /api/users?ids=123,456,789
func (h *TweetsHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	idsParam := r.URL.Query().Get("ids")
	if idsParam == "" {
		h.respondWithError(w, http.StatusBadRequest, "ids là bắt buộc (comma-separated)", "MISSING_IDS")
		return
	}

	ids := []string{}
	for _, id := range splitComma(idsParam) {
		if id != "" {
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		h.respondWithError(w, http.StatusBadRequest, "Không có ID hợp lệ", "INVALID_IDS")
		return
	}

	log.WithFields(log.Fields{
		"ids_count": len(ids),
		"ip":        r.RemoteAddr,
	}).Info("Nhận request lấy danh sách users")

	response, err := h.twitterService.ListUsers(r.Context(), ids)
	if err != nil {
		log.WithError(err).Error("Lỗi khi lấy danh sách users")
		h.respondWithError(w, http.StatusInternalServerError, "Không thể lấy danh sách users: "+err.Error(), "FETCH_ERROR")
		return
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

// GetMe xử lý request lấy thông tin authenticated user
// GET /api/users/me
func (h *TweetsHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	log.WithField("ip", r.RemoteAddr).Info("Nhận request lấy authenticated user")

	user, err := h.twitterService.GetMe(r.Context())
	if err != nil {
		log.WithError(err).Error("Lỗi khi lấy authenticated user")
		h.respondWithError(w, http.StatusInternalServerError, "Không thể lấy authenticated user: "+err.Error(), "FETCH_ERROR")
		return
	}

	h.respondWithJSON(w, http.StatusOK, user)
}

// GetBlockingUsers xử lý request lấy blocking users
// GET /api/users/{username}/blocking?count=10
func (h *TweetsHandler) GetBlockingUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	if username == "" {
		h.respondWithError(w, http.StatusBadRequest, "Username là bắt buộc", "MISSING_USERNAME")
		return
	}

	count := parseCount(r.URL.Query().Get("count"))
	paginationToken := r.URL.Query().Get("pagination_token")

	log.WithFields(log.Fields{
		"username": username,
		"count":    count,
		"ip":       r.RemoteAddr,
	}).Info("Nhận request lấy blocking users")

	response, err := h.twitterService.GetBlockingUsers(r.Context(), username, count, paginationToken)
	if err != nil {
		log.WithError(err).Error("Lỗi khi lấy blocking users")
		h.respondWithError(w, http.StatusInternalServerError, "Không thể lấy blocking users: "+err.Error(), "FETCH_ERROR")
		return
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

// GetMutingUsers xử lý request lấy muting users
// GET /api/users/{username}/muting?count=10
func (h *TweetsHandler) GetMutingUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	if username == "" {
		h.respondWithError(w, http.StatusBadRequest, "Username là bắt buộc", "MISSING_USERNAME")
		return
	}

	count := parseCount(r.URL.Query().Get("count"))
	paginationToken := r.URL.Query().Get("pagination_token")

	log.WithFields(log.Fields{
		"username": username,
		"count":    count,
		"ip":       r.RemoteAddr,
	}).Info("Nhận request lấy muting users")

	response, err := h.twitterService.GetMutingUsers(r.Context(), username, count, paginationToken)
	if err != nil {
		log.WithError(err).Error("Lỗi khi lấy muting users")
		h.respondWithError(w, http.StatusInternalServerError, "Không thể lấy muting users: "+err.Error(), "FETCH_ERROR")
		return
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

// HideTweet xử lý request hide/unhide tweet
// PUT /api/tweets/{tweet_id}/hidden
func (h *TweetsHandler) HideTweet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tweetID := vars["tweet_id"]

	if tweetID == "" {
		h.respondWithError(w, http.StatusBadRequest, "Tweet ID là bắt buộc", "MISSING_TWEET_ID")
		return
	}

	hiddenParam := r.URL.Query().Get("hidden")
	if hiddenParam == "" {
		h.respondWithError(w, http.StatusBadRequest, "hidden parameter là bắt buộc (true/false)", "MISSING_HIDDEN")
		return
	}

	hidden := hiddenParam == "true"

	log.WithFields(log.Fields{
		"tweet_id": tweetID,
		"hidden":   hidden,
		"ip":       r.RemoteAddr,
	}).Info("Nhận request hide tweet")

	response, err := h.twitterService.HideTweet(r.Context(), tweetID, hidden)
	if err != nil {
		log.WithError(err).Error("Lỗi khi hide tweet")
		h.respondWithError(w, http.StatusInternalServerError, "Không thể hide tweet: "+err.Error(), "FETCH_ERROR")
		return
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

// GetUserTimelineReverseChronological xử lý request lấy timeline reverse chronological
// GET /api/users/{username}/timelines/reverse_chronological?count=10
func (h *TweetsHandler) GetUserTimelineReverseChronological(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	if username == "" {
		h.respondWithError(w, http.StatusBadRequest, "Username là bắt buộc", "MISSING_USERNAME")
		return
	}

	count := parseCount(r.URL.Query().Get("count"))

	log.WithFields(log.Fields{
		"username": username,
		"count":    count,
		"ip":       r.RemoteAddr,
	}).Info("Nhận request lấy timeline reverse chronological")

	response, err := h.twitterService.GetUserTimelineReverseChronological(r.Context(), username, count)
	if err != nil {
		log.WithError(err).Error("Lỗi khi lấy timeline")
		h.respondWithError(w, http.StatusInternalServerError, "Không thể lấy timeline: "+err.Error(), "FETCH_ERROR")
		return
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

// GetRepostsOfMe xử lý request lấy reposts của authenticated user
// GET /api/users/reposts_of_me?count=10
func (h *TweetsHandler) GetRepostsOfMe(w http.ResponseWriter, r *http.Request) {
	count := parseCount(r.URL.Query().Get("count"))

	log.WithFields(log.Fields{
		"count": count,
		"ip":    r.RemoteAddr,
	}).Info("Nhận request lấy reposts of me")

	response, err := h.twitterService.GetRepostsOfMe(r.Context(), count)
	if err != nil {
		log.WithError(err).Error("Lỗi khi lấy reposts")
		h.respondWithError(w, http.StatusInternalServerError, "Không thể lấy reposts: "+err.Error(), "FETCH_ERROR")
		return
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

func parseCount(countStr string) int {
	if countStr == "" {
		return 10
	}

	if parsedCount, err := strconv.Atoi(countStr); err == nil {
		return parsedCount
	}

	return 10
}

// splitComma tách chuỗi comma-separated thành slice
func splitComma(s string) []string {
	if s == "" {
		return []string{}
	}
	
	result := []string{}
	parts := strings.Split(s, ",")
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
