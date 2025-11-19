package models

import "time"

// Tweet đại diện cho một tweet từ Twitter/X
type Tweet struct {
	ID               string          `json:"id"`
	Text             string          `json:"text"`
	AuthorID         string          `json:"author_id"`
	CreatedAt        time.Time       `json:"created_at"`
	Metrics          *TweetMetrics   `json:"metrics,omitempty"`
	Entities         *TweetEntities  `json:"entities,omitempty"`
	ReferencedTweets []ReferencedTweet `json:"referenced_tweets,omitempty"`
}

// TweetMetrics chứa các số liệu của tweet
type TweetMetrics struct {
	RetweetCount int `json:"retweet_count"`
	ReplyCount   int `json:"reply_count"`
	LikeCount    int `json:"like_count"`
	QuoteCount   int `json:"quote_count"`
	ViewCount    int `json:"view_count,omitempty"`
}

// TweetEntities chứa các entities trong tweet (hashtags, mentions, urls, etc.)
type TweetEntities struct {
	Hashtags []Hashtag `json:"hashtags,omitempty"`
	Mentions []Mention `json:"mentions,omitempty"`
	URLs     []URL     `json:"urls,omitempty"`
}

// Hashtag đại diện cho một hashtag trong tweet
type Hashtag struct {
	Tag string `json:"tag"`
}

// Mention đại diện cho một mention (@username) trong tweet
type Mention struct {
	Username string `json:"username"`
	ID       string `json:"id"`
}

// URL đại diện cho một URL trong tweet
type URL struct {
	URL         string `json:"url"`
	ExpandedURL string `json:"expanded_url"`
	DisplayURL  string `json:"display_url"`
}

// ReferencedTweet đại diện cho tweet được tham chiếu (reply, retweet, quote)
type ReferencedTweet struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// User đại diện cho thông tin user Twitter/X
type User struct {
	ID              string      `json:"id"`
	Username        string      `json:"username"`
	Name            string      `json:"name"`
	Description     string      `json:"description,omitempty"`
	ProfileImageURL string      `json:"profile_image_url,omitempty"`
	Verified        bool        `json:"verified"`
	CreatedAt       time.Time   `json:"created_at"`
	Metrics         *UserMetrics `json:"metrics,omitempty"`
}

// UserMetrics chứa các số liệu của user
type UserMetrics struct {
	FollowersCount int `json:"followers_count"`
	FollowingCount int `json:"following_count"`
	TweetCount     int `json:"tweet_count"`
	ListedCount    int `json:"listed_count"`
}

// TweetsResponse là response structure cho API lấy tweets
type TweetsResponse struct {
	Tweets []Tweet `json:"tweets"`
	User   *User   `json:"user,omitempty"`
	Meta   *Meta   `json:"meta,omitempty"`
}

// Meta chứa metadata của response
type Meta struct {
	ResultCount   int    `json:"result_count"`
	NextToken     string `json:"next_token,omitempty"`
	PreviousToken string `json:"previous_token,omitempty"`
}

// FollowingResponse là response structure cho API lấy danh sách followings
type FollowingResponse struct {
	User      *User  `json:"user"`
	Following []User `json:"following"`
	Meta      *Meta  `json:"meta,omitempty"`
}

// ErrorResponse là response structure cho lỗi
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// SearchTweetsResponse là response structure cho API tìm kiếm tweets
type SearchTweetsResponse struct {
	Tweets []Tweet `json:"tweets"`
	Meta   *Meta   `json:"meta,omitempty"`
}

// TweetDetailResponse là response structure cho API lấy chi tiết tweet
type TweetDetailResponse struct {
	Tweet  Tweet `json:"tweet"`
	Author *User `json:"author,omitempty"`
}

// FollowersResponse là response structure cho API lấy danh sách followers
type FollowersResponse struct {
	User      *User  `json:"user"`
	Followers []User `json:"followers"`
	Meta      *Meta  `json:"meta,omitempty"`
}

// LikedTweetsResponse là response structure cho API lấy liked tweets
type LikedTweetsResponse struct {
	User   *User   `json:"user"`
	Tweets []Tweet `json:"tweets"`
	Meta   *Meta   `json:"meta,omitempty"`
}

// SearchUsersResponse là response structure cho API tìm kiếm users
type SearchUsersResponse struct {
	Users []User `json:"users"`
	Meta  *Meta  `json:"meta,omitempty"`
}

// MentionsResponse là response structure cho API lấy mentions
type MentionsResponse struct {
	User   *User   `json:"user"`
	Tweets []Tweet `json:"tweets"`
	Meta   *Meta   `json:"meta,omitempty"`
}

// LikingUsersResponse là response structure cho API lấy users đã like tweet
type LikingUsersResponse struct {
	TweetID string `json:"tweet_id"`
	Users   []User `json:"users"`
	Meta    *Meta  `json:"meta,omitempty"`
}

// QuoteTweetsResponse là response structure cho API lấy quote tweets
type QuoteTweetsResponse struct {
	TweetID string  `json:"tweet_id"`
	Tweets  []Tweet `json:"tweets"`
	Meta    *Meta   `json:"meta,omitempty"`
}

// RetweetedByResponse là response structure cho API lấy users đã retweet
type RetweetedByResponse struct {
	TweetID string `json:"tweet_id"`
	Users   []User `json:"users"`
	Meta    *Meta  `json:"meta,omitempty"`
}

// TweetCountsResponse là response structure cho API lấy tweet counts
type TweetCountsResponse struct {
	Query       string         `json:"query"`
	Counts     []TweetCount   `json:"counts"`
	Meta       *Meta          `json:"meta,omitempty"`
}

// TweetCount chứa thông tin count của tweet trong một khoảng thời gian
type TweetCount struct {
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
	TweetCount int       `json:"tweet_count"`
}

// UsersListResponse là response structure cho API lấy danh sách users
type UsersListResponse struct {
	Users []User `json:"users"`
	Meta  *Meta  `json:"meta,omitempty"`
}

// BlockingUsersResponse là response structure cho API lấy blocking users
type BlockingUsersResponse struct {
	User   *User  `json:"user"`
	Users  []User `json:"users"`
	Meta   *Meta  `json:"meta,omitempty"`
}

// MutingUsersResponse là response structure cho API lấy muting users
type MutingUsersResponse struct {
	User   *User  `json:"user"`
	Users  []User `json:"users"`
	Meta   *Meta  `json:"meta,omitempty"`
}

// HideTweetResponse là response structure cho API hide/unhide tweet
type HideTweetResponse struct {
	TweetID string `json:"tweet_id"`
	Hidden  bool   `json:"hidden"`
	Message string `json:"message,omitempty"`
}

// RepostsResponse là response structure cho API lấy reposts
type RepostsResponse struct {
	User   *User   `json:"user"`
	Tweets []Tweet `json:"tweets"`
	Meta   *Meta   `json:"meta,omitempty"`
}

