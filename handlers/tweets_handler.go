package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
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

func parseCount(countStr string) int {
	if countStr == "" {
		return 10
	}

	if parsedCount, err := strconv.Atoi(countStr); err == nil {
		return parsedCount
	}

	return 10
}
