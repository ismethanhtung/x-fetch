package services

import (
	"context"
	"fmt"
	"x-twitter-backend/config"
	"x-twitter-backend/models"

	"strings"
	"time"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/fields"
	"github.com/michimani/gotwi/resources"
	"github.com/michimani/gotwi/tweet/timeline"
	timelineTypes "github.com/michimani/gotwi/tweet/timeline/types"
	"github.com/michimani/gotwi/tweet/like"
	likeTypes "github.com/michimani/gotwi/tweet/like/types"
	"github.com/michimani/gotwi/tweet/quotetweet"
	quotetweetTypes "github.com/michimani/gotwi/tweet/quotetweet/types"
	"github.com/michimani/gotwi/tweet/retweet"
	retweetTypes "github.com/michimani/gotwi/tweet/retweet/types"
	"github.com/michimani/gotwi/tweet/searchtweet"
	searchTypes "github.com/michimani/gotwi/tweet/searchtweet/types"
	"github.com/michimani/gotwi/tweet/tweetcount"
	tweetcountTypes "github.com/michimani/gotwi/tweet/tweetcount/types"
	"github.com/michimani/gotwi/tweet/tweetlookup"
	lookupTypes "github.com/michimani/gotwi/tweet/tweetlookup/types"
	"github.com/michimani/gotwi/user/follow"
	followTypes "github.com/michimani/gotwi/user/follow/types"
	"github.com/michimani/gotwi/user/userlookup"
	userlookupTypes "github.com/michimani/gotwi/user/userlookup/types"
	log "github.com/sirupsen/logrus"
)

// TwitterService xử lý tất cả các tương tác với Twitter API
type TwitterService struct {
	client *gotwi.Client
	config *config.Config
}

// NewTwitterService tạo một instance mới của TwitterService
func NewTwitterService(cfg *config.Config) (*TwitterService, error) {
	// Khởi tạo Twitter client với Bearer Token
	client, err := gotwi.NewClientWithAccessToken(&gotwi.NewClientWithAccessTokenInput{
		AccessToken: cfg.TwitterBearerToken,
	})
	if err != nil {
		return nil, fmt.Errorf("không thể khởi tạo Twitter client: %w", err)
	}

	log.Info("Twitter client đã được khởi tạo thành công")

	return &TwitterService{
		client: client,
		config: cfg,
	}, nil
}

// GetUserByUsername lấy thông tin user theo username
func (s *TwitterService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	log.WithField("username", username).Info("Đang lấy thông tin user")

	params := &userlookupTypes.GetByUsernameInput{
		Username: username,
		UserFields: fields.UserFieldList{
			fields.UserFieldID,
			fields.UserFieldName,
			fields.UserFieldUsername,
			fields.UserFieldDescription,
			fields.UserFieldProfileImageUrl,
			fields.UserFieldVerified,
			fields.UserFieldCreatedAt,
			fields.UserFieldPublicMetrics,
		},
	}

	resp, err := userlookup.GetByUsername(ctx, s.client, params)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy thông tin user: %w", err)
	}

	if resp.Data.ID == nil || gotwi.StringValue(resp.Data.ID) == "" {
		return nil, fmt.Errorf("không tìm thấy user với username: %s", username)
	}

	user := s.convertToUser(&resp.Data)
	log.WithFields(log.Fields{
		"user_id":  user.ID,
		"username": user.Username,
	}).Info("Đã lấy thông tin user thành công")

	return user, nil
}

// GetUserTweets lấy tweets của một user theo username
func (s *TwitterService) GetUserTweets(ctx context.Context, username string, maxResults int) (*models.TweetsResponse, error) {
	log.WithFields(log.Fields{
		"username":    username,
		"max_results": maxResults,
	}).Info("Đang lấy tweets của user")

	// Đầu tiên, lấy thông tin user để có user ID
	user, err := s.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	// Validate và điều chỉnh maxResults
	if maxResults <= 0 {
		maxResults = s.config.DefaultTweetsCount
	}
	if maxResults > s.config.MaxTweetsPerRequest {
		maxResults = s.config.MaxTweetsPerRequest
	}

	// Lấy tweets của user
	params := &timelineTypes.ListTweetsInput{
		ID:         user.ID,
		MaxResults: timelineTypes.ListMaxResults(maxResults),
		TweetFields: fields.TweetFieldList{
			fields.TweetFieldID,
			fields.TweetFieldText,
			fields.TweetFieldAuthorID,
			fields.TweetFieldCreatedAt,
			fields.TweetFieldPublicMetrics,
			fields.TweetFieldEntities,
			fields.TweetFieldReferencedTweets,
		},
		UserFields: fields.UserFieldList{
			fields.UserFieldID,
			fields.UserFieldName,
			fields.UserFieldUsername,
			fields.UserFieldProfileImageUrl,
		},
	}

	resp, err := timeline.ListTweets(ctx, s.client, params)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy tweets: %w", err)
	}

	// Convert response sang models
	tweets := make([]models.Tweet, 0)
	if len(resp.Data) > 0 {
		for i := range resp.Data {
			tweet := s.convertToTweet(&resp.Data[i])
			tweets = append(tweets, tweet)
		}
	}

	result := &models.TweetsResponse{
		Tweets: tweets,
		User:   user,
		Meta:   buildMetaFromTimeline(resp.Meta, len(tweets)),
	}

	log.WithFields(log.Fields{
		"username":     username,
		"tweets_count": len(tweets),
	}).Info("Đã lấy tweets thành công")

	return result, nil
}

// GetUserFollowing lấy danh sách tài khoản mà user đang theo dõi
func (s *TwitterService) GetUserFollowing(ctx context.Context, username string, maxResults int, paginationToken string) (*models.FollowingResponse, error) {
	log.WithFields(log.Fields{
		"username":    username,
		"max_results": maxResults,
		"page_token":  paginationToken,
	}).Info("Đang lấy danh sách following của user")

	user, err := s.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if maxResults <= 0 {
		maxResults = s.config.DefaultTweetsCount
	}
	if maxResults > 1000 {
		maxResults = 1000
	}

	params := &followTypes.ListFollowingsInput{
		ID:         user.ID,
		MaxResults: followTypes.ListMaxResults(maxResults),
		UserFields: fields.UserFieldList{
			fields.UserFieldID,
			fields.UserFieldName,
			fields.UserFieldUsername,
			fields.UserFieldDescription,
			fields.UserFieldProfileImageUrl,
			fields.UserFieldVerified,
			fields.UserFieldCreatedAt,
			fields.UserFieldPublicMetrics,
		},
	}

	if paginationToken != "" {
		params.PaginationToken = paginationToken
	}

	resp, err := follow.ListFollowings(ctx, s.client, params)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy danh sách following: %w", err)
	}

	following := make([]models.User, 0, len(resp.Data))
	for i := range resp.Data {
		following = append(following, *s.convertToUser(&resp.Data[i]))
	}

	result := &models.FollowingResponse{
		User:      user,
		Following: following,
		Meta:      buildMetaFromPagination(resp.Meta, len(following)),
	}

	log.WithFields(log.Fields{
		"username":        username,
		"following_count": len(following),
	}).Info("Đã lấy danh sách following thành công")

	return result, nil
}

// GetTweetsByUserID lấy tweets theo user ID trực tiếp
func (s *TwitterService) GetTweetsByUserID(ctx context.Context, userID string, maxResults int) ([]models.Tweet, error) {
	log.WithFields(log.Fields{
		"user_id":     userID,
		"max_results": maxResults,
	}).Info("Đang lấy tweets theo user ID")

	// Validate và điều chỉnh maxResults
	if maxResults <= 0 {
		maxResults = s.config.DefaultTweetsCount
	}
	if maxResults > s.config.MaxTweetsPerRequest {
		maxResults = s.config.MaxTweetsPerRequest
	}

	params := &timelineTypes.ListTweetsInput{
		ID:         userID,
		MaxResults: timelineTypes.ListMaxResults(maxResults),
		TweetFields: fields.TweetFieldList{
			fields.TweetFieldID,
			fields.TweetFieldText,
			fields.TweetFieldAuthorID,
			fields.TweetFieldCreatedAt,
			fields.TweetFieldPublicMetrics,
			fields.TweetFieldEntities,
			fields.TweetFieldReferencedTweets,
		},
	}

	resp, err := timeline.ListTweets(ctx, s.client, params)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy tweets: %w", err)
	}

	tweets := make([]models.Tweet, 0)
	if len(resp.Data) > 0 {
		for i := range resp.Data {
			tweet := s.convertToTweet(&resp.Data[i])
			tweets = append(tweets, tweet)
		}
	}

	log.WithFields(log.Fields{
		"user_id":      userID,
		"tweets_count": len(tweets),
	}).Info("Đã lấy tweets thành công")

	return tweets, nil
}

func buildMetaFromTimeline(meta resources.TweetTimelineMeta, fallbackCount int) *models.Meta {
	resultCount := fallbackCount
	if meta.ResultCount != nil {
		resultCount = gotwi.IntValue(meta.ResultCount)
	}

	return &models.Meta{
		ResultCount: resultCount,
		NextToken:   gotwi.StringValue(meta.NextToken),
	}
}

func buildMetaFromPagination(meta resources.PaginationMeta, fallbackCount int) *models.Meta {
	resultCount := fallbackCount
	if meta.ResultCount != nil {
		resultCount = gotwi.IntValue(meta.ResultCount)
	}

	out := &models.Meta{
		ResultCount: resultCount,
		NextToken:   gotwi.StringValue(meta.NextToken),
	}

	if meta.PreviousToken != nil {
		out.PreviousToken = gotwi.StringValue(meta.PreviousToken)
	}

	return out
}

// buildMetaFromQuoteTweets chuyển đổi QuoteTweetsMeta sang models.Meta
func buildMetaFromQuoteTweets(meta resources.QuoteTweetsMeta, fallbackCount int) *models.Meta {
	resultCount := fallbackCount
	if meta.ResultCount != nil {
		resultCount = gotwi.IntValue(meta.ResultCount)
	}

	return &models.Meta{
		ResultCount: resultCount,
		NextToken:   gotwi.StringValue(meta.NextToken),
	}
}

// contains kiểm tra xem string có chứa substring không (case-insensitive)
func contains(s, substr string) bool {
	sLower := strings.ToLower(s)
	substrLower := strings.ToLower(substr)
	return strings.Contains(sLower, substrLower)
}

// convertToUser chuyển đổi Twitter user data sang models.User
func (s *TwitterService) convertToUser(data *resources.User) *models.User {
	user := &models.User{
		ID:       gotwi.StringValue(data.ID),
		Username: gotwi.StringValue(data.Username),
		Name:     gotwi.StringValue(data.Name),
	}

	if data.Description != nil {
		user.Description = *data.Description
	}

	if data.ProfileImageURL != nil {
		user.ProfileImageURL = *data.ProfileImageURL
	}

	if data.Verified != nil {
		user.Verified = *data.Verified
	}

	if data.CreatedAt != nil {
		user.CreatedAt = gotwi.TimeValue(data.CreatedAt)
	}

	if data.PublicMetrics != nil {
		user.Metrics = &models.UserMetrics{
			FollowersCount: gotwi.IntValue(data.PublicMetrics.FollowersCount),
			FollowingCount: gotwi.IntValue(data.PublicMetrics.FollowingCount),
			TweetCount:     gotwi.IntValue(data.PublicMetrics.TweetCount),
			ListedCount:    gotwi.IntValue(data.PublicMetrics.ListedCount),
		}
	}

	return user
}

// convertToTweet chuyển đổi Twitter tweet data sang models.Tweet
func (s *TwitterService) convertToTweet(data *resources.Tweet) models.Tweet {
	tweet := models.Tweet{
		ID:       gotwi.StringValue(data.ID),
		Text:     gotwi.StringValue(data.Text),
		AuthorID: gotwi.StringValue(data.AuthorID),
	}

	if data.CreatedAt != nil {
		tweet.CreatedAt = gotwi.TimeValue(data.CreatedAt)
	}

	if data.PublicMetrics != nil {
		tweet.Metrics = &models.TweetMetrics{
			RetweetCount: gotwi.IntValue(data.PublicMetrics.RetweetCount),
			ReplyCount:   gotwi.IntValue(data.PublicMetrics.ReplyCount),
			LikeCount:    gotwi.IntValue(data.PublicMetrics.LikeCount),
			QuoteCount:   gotwi.IntValue(data.PublicMetrics.QuoteCount),
		}
	}

	if data.Entities != nil {
		tweet.Entities = &models.TweetEntities{}

		if len(data.Entities.HashTags) > 0 {
			tweet.Entities.Hashtags = make([]models.Hashtag, 0, len(data.Entities.HashTags))
			for _, ht := range data.Entities.HashTags {
				tweet.Entities.Hashtags = append(tweet.Entities.Hashtags, models.Hashtag{
					Tag: gotwi.StringValue(ht.Tag),
				})
			}
		}

		if len(data.Entities.Mentions) > 0 {
			tweet.Entities.Mentions = make([]models.Mention, 0, len(data.Entities.Mentions))
			for _, mention := range data.Entities.Mentions {
				tweet.Entities.Mentions = append(tweet.Entities.Mentions, models.Mention{
					Username: gotwi.StringValue(mention.Tag),
				})
			}
		}

		if len(data.Entities.URLs) > 0 {
			tweet.Entities.URLs = make([]models.URL, 0, len(data.Entities.URLs))
			for _, url := range data.Entities.URLs {
				tweet.Entities.URLs = append(tweet.Entities.URLs, models.URL{
					URL:         gotwi.StringValue(url.URL),
					ExpandedURL: gotwi.StringValue(url.ExpandedURL),
					DisplayURL:  gotwi.StringValue(url.DisplayURL),
				})
			}
		}
	}

	if data.ReferencedTweets != nil {
		tweet.ReferencedTweets = make([]models.ReferencedTweet, 0)
		for _, ref := range data.ReferencedTweets {
			tweet.ReferencedTweets = append(tweet.ReferencedTweets, models.ReferencedTweet{
				Type: gotwi.StringValue(ref.Type),
				ID:   gotwi.StringValue(ref.ID),
			})
		}
	}

	return tweet
}

// GetUserFollowers lấy danh sách người theo dõi (followers) của user
func (s *TwitterService) GetUserFollowers(ctx context.Context, username string, maxResults int, paginationToken string) (*models.FollowersResponse, error) {
	log.WithFields(log.Fields{
		"username":    username,
		"max_results": maxResults,
		"page_token":  paginationToken,
	}).Info("Đang lấy danh sách followers của user")

	user, err := s.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if maxResults <= 0 {
		maxResults = s.config.DefaultTweetsCount
	}
	if maxResults > 1000 {
		maxResults = 1000
	}

	params := &followTypes.ListFollowersInput{
		ID:         user.ID,
		MaxResults: followTypes.ListMaxResults(maxResults),
		UserFields: fields.UserFieldList{
			fields.UserFieldID,
			fields.UserFieldName,
			fields.UserFieldUsername,
			fields.UserFieldDescription,
			fields.UserFieldProfileImageUrl,
			fields.UserFieldVerified,
			fields.UserFieldCreatedAt,
			fields.UserFieldPublicMetrics,
		},
	}

	if paginationToken != "" {
		params.PaginationToken = paginationToken
	}

	resp, err := follow.ListFollowers(ctx, s.client, params)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy danh sách followers: %w", err)
	}

	followers := make([]models.User, 0, len(resp.Data))
	for i := range resp.Data {
		followers = append(followers, *s.convertToUser(&resp.Data[i]))
	}

	result := &models.FollowersResponse{
		User:      user,
		Followers: followers,
		Meta:      buildMetaFromPagination(resp.Meta, len(followers)),
	}

	log.WithFields(log.Fields{
		"username":        username,
		"followers_count": len(followers),
	}).Info("Đã lấy danh sách followers thành công")

	return result, nil
}

// SearchTweets tìm kiếm tweets theo keyword
func (s *TwitterService) SearchTweets(ctx context.Context, query string, maxResults int) (*models.SearchTweetsResponse, error) {
	log.WithFields(log.Fields{
		"query":       query,
		"max_results": maxResults,
	}).Info("Đang tìm kiếm tweets")

	if maxResults <= 0 {
		maxResults = s.config.DefaultTweetsCount
	}
	if maxResults > 100 {
		maxResults = 100
	}

	params := &searchTypes.ListRecentInput{
		Query:      query,
		MaxResults: searchTypes.ListMaxResults(maxResults),
		TweetFields: fields.TweetFieldList{
			fields.TweetFieldID,
			fields.TweetFieldText,
			fields.TweetFieldAuthorID,
			fields.TweetFieldCreatedAt,
			fields.TweetFieldPublicMetrics,
			fields.TweetFieldEntities,
			fields.TweetFieldReferencedTweets,
		},
	}

	resp, err := searchtweet.ListRecent(ctx, s.client, params)
	if err != nil {
		return nil, fmt.Errorf("không thể tìm kiếm tweets: %w", err)
	}

	tweets := make([]models.Tweet, 0, len(resp.Data))
	for i := range resp.Data {
		tweet := s.convertToTweet(&resp.Data[i])
		tweets = append(tweets, tweet)
	}

	result := &models.SearchTweetsResponse{
		Tweets: tweets,
		Meta:   buildMetaFromPagination(resp.Meta, len(tweets)),
	}

	log.WithFields(log.Fields{
		"query":        query,
		"tweets_count": len(tweets),
	}).Info("Đã tìm kiếm tweets thành công")

	return result, nil
}

// GetTweetByID lấy chi tiết tweet theo ID
func (s *TwitterService) GetTweetByID(ctx context.Context, tweetID string) (*models.TweetDetailResponse, error) {
	log.WithField("tweet_id", tweetID).Info("Đang lấy chi tiết tweet")

	params := &lookupTypes.GetInput{
		ID: tweetID,
		TweetFields: fields.TweetFieldList{
			fields.TweetFieldID,
			fields.TweetFieldText,
			fields.TweetFieldAuthorID,
			fields.TweetFieldCreatedAt,
			fields.TweetFieldPublicMetrics,
			fields.TweetFieldEntities,
			fields.TweetFieldReferencedTweets,
		},
		Expansions: fields.ExpansionList{
			fields.ExpansionAuthorID,
		},
		UserFields: fields.UserFieldList{
			fields.UserFieldID,
			fields.UserFieldName,
			fields.UserFieldUsername,
			fields.UserFieldProfileImageUrl,
			fields.UserFieldVerified,
			fields.UserFieldPublicMetrics,
		},
	}

	resp, err := tweetlookup.Get(ctx, s.client, params)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy tweet: %w", err)
	}

	tweet := s.convertToTweet(&resp.Data)

	result := &models.TweetDetailResponse{
		Tweet: tweet,
	}

	// Nếu có thông tin author trong response
	if len(resp.Includes.Users) > 0 {
		result.Author = s.convertToUser(&resp.Includes.Users[0])
	}

	log.WithField("tweet_id", tweetID).Info("Đã lấy chi tiết tweet thành công")

	return result, nil
}

// GetLikedTweets lấy danh sách tweets mà user đã like
func (s *TwitterService) GetLikedTweets(ctx context.Context, username string, maxResults int) (*models.LikedTweetsResponse, error) {
	log.WithFields(log.Fields{
		"username":    username,
		"max_results": maxResults,
	}).Info("Đang lấy liked tweets")

	user, err := s.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if maxResults <= 0 {
		maxResults = s.config.DefaultTweetsCount
	}
	if maxResults > 100 {
		maxResults = 100
	}

	params := &likeTypes.ListInput{
		ID:         user.ID,
		MaxResults: likeTypes.ListMaxResults(maxResults),
		TweetFields: fields.TweetFieldList{
			fields.TweetFieldID,
			fields.TweetFieldText,
			fields.TweetFieldAuthorID,
			fields.TweetFieldCreatedAt,
			fields.TweetFieldPublicMetrics,
			fields.TweetFieldEntities,
		},
	}

	resp, err := like.List(ctx, s.client, params)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy liked tweets: %w", err)
	}

	tweets := make([]models.Tweet, 0, len(resp.Data))
	for i := range resp.Data {
		tweet := s.convertToTweet(&resp.Data[i])
		tweets = append(tweets, tweet)
	}

	result := &models.LikedTweetsResponse{
		User:   user,
		Tweets: tweets,
		Meta:   buildMetaFromPagination(resp.Meta, len(tweets)),
	}

	log.WithFields(log.Fields{
		"username":     username,
		"tweets_count": len(tweets),
	}).Info("Đã lấy liked tweets thành công")

	return result, nil
}

// SearchUsers tìm kiếm users theo query
func (s *TwitterService) SearchUsers(ctx context.Context, query string, maxResults int) (*models.SearchUsersResponse, error) {
	log.WithFields(log.Fields{
		"query":       query,
		"max_results": maxResults,
	}).Info("Đang tìm kiếm users")

	if maxResults <= 0 {
		maxResults = s.config.DefaultTweetsCount
	}
	if maxResults > 100 {
		maxResults = 100
	}

	// Twitter API v2 không hỗ trợ user search trực tiếp với Bearer token
	// Thay vào đó, chúng ta sẽ tìm kiếm tweets với query và lấy unique authors
	searchParams := &searchTypes.ListRecentInput{
		Query:      query,
		MaxResults: searchTypes.ListMaxResults(maxResults),
		TweetFields: fields.TweetFieldList{
			fields.TweetFieldAuthorID,
		},
		Expansions: fields.ExpansionList{
			fields.ExpansionAuthorID,
		},
		UserFields: fields.UserFieldList{
			fields.UserFieldID,
			fields.UserFieldName,
			fields.UserFieldUsername,
			fields.UserFieldDescription,
			fields.UserFieldProfileImageUrl,
			fields.UserFieldVerified,
			fields.UserFieldCreatedAt,
			fields.UserFieldPublicMetrics,
		},
	}

	resp, err := searchtweet.ListRecent(ctx, s.client, searchParams)
	if err != nil {
		return nil, fmt.Errorf("không thể tìm kiếm users: %w", err)
	}

	// Extract unique users from includes
	usersMap := make(map[string]*models.User)
	if len(resp.Includes.Users) > 0 {
		for i := range resp.Includes.Users {
			user := s.convertToUser(&resp.Includes.Users[i])
			usersMap[user.ID] = user
		}
	}

	users := make([]models.User, 0, len(usersMap))
	for _, user := range usersMap {
		users = append(users, *user)
	}

	result := &models.SearchUsersResponse{
		Users: users,
		Meta: &models.Meta{
			ResultCount: len(users),
		},
	}

	log.WithFields(log.Fields{
		"query":       query,
		"users_count": len(users),
	}).Info("Đã tìm kiếm users thành công")

	return result, nil
}

// GetUserMentions lấy danh sách tweets có mention đến user
func (s *TwitterService) GetUserMentions(ctx context.Context, username string, maxResults int) (*models.MentionsResponse, error) {
	log.WithFields(log.Fields{
		"username":    username,
		"max_results": maxResults,
	}).Info("Đang lấy mentions của user")

	user, err := s.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if maxResults <= 0 {
		maxResults = s.config.DefaultTweetsCount
	}
	if maxResults > 100 {
		maxResults = 100
	}

	params := &timelineTypes.ListMentionsInput{
		ID:         user.ID,
		MaxResults: timelineTypes.ListMaxResults(maxResults),
		TweetFields: fields.TweetFieldList{
			fields.TweetFieldID,
			fields.TweetFieldText,
			fields.TweetFieldAuthorID,
			fields.TweetFieldCreatedAt,
			fields.TweetFieldPublicMetrics,
			fields.TweetFieldEntities,
		},
	}

	resp, err := timeline.ListMentions(ctx, s.client, params)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy mentions: %w", err)
	}

	tweets := make([]models.Tweet, 0, len(resp.Data))
	for i := range resp.Data {
		tweet := s.convertToTweet(&resp.Data[i])
		tweets = append(tweets, tweet)
	}

	result := &models.MentionsResponse{
		User:   user,
		Tweets: tweets,
		Meta:   buildMetaFromTimeline(resp.Meta, len(tweets)),
	}

	log.WithFields(log.Fields{
		"username":     username,
		"tweets_count": len(tweets),
	}).Info("Đã lấy mentions thành công")

	return result, nil
}

// ListTweets lấy danh sách tweets theo IDs
func (s *TwitterService) ListTweets(ctx context.Context, tweetIDs []string) (*models.SearchTweetsResponse, error) {
	log.WithField("count", len(tweetIDs)).Info("Đang lấy danh sách tweets")

	if len(tweetIDs) == 0 {
		return &models.SearchTweetsResponse{Tweets: []models.Tweet{}}, nil
	}

	// Twitter API chỉ cho phép tối đa 100 IDs mỗi request
	if len(tweetIDs) > 100 {
		tweetIDs = tweetIDs[:100]
	}

	params := &lookupTypes.ListInput{
		IDs: tweetIDs,
		TweetFields: fields.TweetFieldList{
			fields.TweetFieldID,
			fields.TweetFieldText,
			fields.TweetFieldAuthorID,
			fields.TweetFieldCreatedAt,
			fields.TweetFieldPublicMetrics,
			fields.TweetFieldEntities,
			fields.TweetFieldReferencedTweets,
		},
	}

	resp, err := tweetlookup.List(ctx, s.client, params)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy danh sách tweets: %w", err)
	}

	tweets := make([]models.Tweet, 0, len(resp.Data))
	for i := range resp.Data {
		tweet := s.convertToTweet(&resp.Data[i])
		tweets = append(tweets, tweet)
	}

	result := &models.SearchTweetsResponse{
		Tweets: tweets,
		Meta: &models.Meta{
			ResultCount: len(tweets),
		},
	}

	log.WithField("tweets_count", len(tweets)).Info("Đã lấy danh sách tweets thành công")
	return result, nil
}

// GetLikingUsers lấy danh sách users đã like một tweet
// Lưu ý: API này có thể yêu cầu OAuth 1.0a hoặc có giới hạn với Bearer Token
func (s *TwitterService) GetLikingUsers(ctx context.Context, tweetID string, maxResults int, paginationToken string) (*models.LikingUsersResponse, error) {
	log.WithFields(log.Fields{
		"tweet_id":     tweetID,
		"max_results":  maxResults,
		"page_token":   paginationToken,
	}).Info("Đang lấy danh sách users đã like tweet")

	if maxResults <= 0 {
		maxResults = s.config.DefaultTweetsCount
	}
	if maxResults > 100 {
		maxResults = 100
	}

	params := &likeTypes.ListUsersInput{
		ID:         tweetID,
		MaxResults: likeTypes.ListUsersMaxResults(maxResults),
		UserFields: fields.UserFieldList{
			fields.UserFieldID,
			fields.UserFieldName,
			fields.UserFieldUsername,
			fields.UserFieldDescription,
			fields.UserFieldProfileImageUrl,
			fields.UserFieldVerified,
			fields.UserFieldCreatedAt,
			fields.UserFieldPublicMetrics,
		},
	}

	if paginationToken != "" {
		params.PaginationToken = paginationToken
	}

	resp, err := like.ListUsers(ctx, s.client, params)
	if err != nil {
		// Kiểm tra nếu là lỗi 403, trả về thông báo rõ ràng hơn
		if errStr := err.Error(); contains(errStr, "403") || contains(errStr, "Forbidden") {
			return nil, fmt.Errorf("API liking users có thể yêu cầu OAuth 1.0a hoặc tweet không public. Bearer Token có giới hạn với API này. Lỗi: %w", err)
		}
		return nil, fmt.Errorf("không thể lấy danh sách liking users: %w", err)
	}

	users := make([]models.User, 0, len(resp.Data))
	for i := range resp.Data {
		users = append(users, *s.convertToUser(&resp.Data[i]))
	}

	result := &models.LikingUsersResponse{
		TweetID: tweetID,
		Users:   users,
		Meta: &models.Meta{
			ResultCount: len(users),
		},
	}

	log.WithFields(log.Fields{
		"tweet_id":      tweetID,
		"users_count":   len(users),
	}).Info("Đã lấy danh sách liking users thành công")

	return result, nil
}

// GetQuoteTweets lấy danh sách quote tweets của một tweet
func (s *TwitterService) GetQuoteTweets(ctx context.Context, tweetID string, maxResults int) (*models.QuoteTweetsResponse, error) {
	log.WithFields(log.Fields{
		"tweet_id":    tweetID,
		"max_results": maxResults,
	}).Info("Đang lấy danh sách quote tweets")

	if maxResults <= 0 {
		maxResults = s.config.DefaultTweetsCount
	}
	if maxResults > 100 {
		maxResults = 100
	}

	params := &quotetweetTypes.ListInput{
		ID:         tweetID,
		MaxResults: quotetweetTypes.ListMaxResults(maxResults),
		TweetFields: fields.TweetFieldList{
			fields.TweetFieldID,
			fields.TweetFieldText,
			fields.TweetFieldAuthorID,
			fields.TweetFieldCreatedAt,
			fields.TweetFieldPublicMetrics,
			fields.TweetFieldEntities,
			fields.TweetFieldReferencedTweets,
		},
	}

	resp, err := quotetweet.List(ctx, s.client, params)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy danh sách quote tweets: %w", err)
	}

	tweets := make([]models.Tweet, 0, len(resp.Data))
	for i := range resp.Data {
		tweet := s.convertToTweet(&resp.Data[i])
		tweets = append(tweets, tweet)
	}

	result := &models.QuoteTweetsResponse{
		TweetID: tweetID,
		Tweets:  tweets,
		Meta:    buildMetaFromQuoteTweets(resp.Meta, len(tweets)),
	}

	log.WithFields(log.Fields{
		"tweet_id":     tweetID,
		"tweets_count": len(tweets),
	}).Info("Đã lấy danh sách quote tweets thành công")

	return result, nil
}

// GetRetweetedBy lấy danh sách users đã retweet một tweet
func (s *TwitterService) GetRetweetedBy(ctx context.Context, tweetID string, maxResults int, paginationToken string) (*models.RetweetedByResponse, error) {
	log.WithFields(log.Fields{
		"tweet_id":     tweetID,
		"max_results":  maxResults,
		"page_token":   paginationToken,
	}).Info("Đang lấy danh sách users đã retweet")

	if maxResults <= 0 {
		maxResults = s.config.DefaultTweetsCount
	}
	if maxResults > 100 {
		maxResults = 100
	}

	params := &retweetTypes.ListUsersInput{
		ID:         tweetID,
		MaxResults: retweetTypes.ListUsersMaxResults(maxResults),
		UserFields: fields.UserFieldList{
			fields.UserFieldID,
			fields.UserFieldName,
			fields.UserFieldUsername,
			fields.UserFieldDescription,
			fields.UserFieldProfileImageUrl,
			fields.UserFieldVerified,
			fields.UserFieldCreatedAt,
			fields.UserFieldPublicMetrics,
		},
	}

	if paginationToken != "" {
		params.PaginationToken = paginationToken
	}

	resp, err := retweet.ListUsers(ctx, s.client, params)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy danh sách retweeted by: %w", err)
	}

	users := make([]models.User, 0, len(resp.Data))
	for i := range resp.Data {
		users = append(users, *s.convertToUser(&resp.Data[i]))
	}

	result := &models.RetweetedByResponse{
		TweetID: tweetID,
		Users:   users,
		Meta: &models.Meta{
			ResultCount: len(users),
		},
	}

	log.WithFields(log.Fields{
		"tweet_id":    tweetID,
		"users_count": len(users),
	}).Info("Đã lấy danh sách retweeted by thành công")

	return result, nil
}

// GetTweetCounts lấy số lượng tweets theo query và time range
func (s *TwitterService) GetTweetCounts(ctx context.Context, query string, startTime, endTime string) (*models.TweetCountsResponse, error) {
	log.WithFields(log.Fields{
		"query":      query,
		"start_time": startTime,
		"end_time":   endTime,
	}).Info("Đang lấy tweet counts")

	params := &tweetcountTypes.ListRecentInput{
		Query: query,
	}

	if startTime != "" {
		if t, err := time.Parse(time.RFC3339, startTime); err == nil {
			params.StartTime = &t
		}
	}
	if endTime != "" {
		if t, err := time.Parse(time.RFC3339, endTime); err == nil {
			params.EndTime = &t
		}
	}

	resp, err := tweetcount.ListRecent(ctx, s.client, params)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy tweet counts: %w", err)
	}

	counts := make([]models.TweetCount, 0, len(resp.Data))
	for i := range resp.Data {
		count := models.TweetCount{
			TweetCount: gotwi.IntValue(resp.Data[i].TweetCount),
		}
		if resp.Data[i].Start != nil {
			count.Start = gotwi.TimeValue(resp.Data[i].Start)
		}
		if resp.Data[i].End != nil {
			count.End = gotwi.TimeValue(resp.Data[i].End)
		}
		counts = append(counts, count)
	}

	result := &models.TweetCountsResponse{
		Query:   query,
		Counts:  counts,
		Meta: &models.Meta{
			ResultCount: len(counts),
		},
	}

	log.WithFields(log.Fields{
		"query":        query,
		"counts_count": len(counts),
	}).Info("Đã lấy tweet counts thành công")

	return result, nil
}

// GetUserByID lấy thông tin user theo ID
func (s *TwitterService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	log.WithField("user_id", userID).Info("Đang lấy thông tin user theo ID")

	params := &userlookupTypes.GetInput{
		ID: userID,
		UserFields: fields.UserFieldList{
			fields.UserFieldID,
			fields.UserFieldName,
			fields.UserFieldUsername,
			fields.UserFieldDescription,
			fields.UserFieldProfileImageUrl,
			fields.UserFieldVerified,
			fields.UserFieldCreatedAt,
			fields.UserFieldPublicMetrics,
		},
	}

	resp, err := userlookup.Get(ctx, s.client, params)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy thông tin user: %w", err)
	}

	if resp.Data.ID == nil || gotwi.StringValue(resp.Data.ID) == "" {
		return nil, fmt.Errorf("không tìm thấy user với ID: %s", userID)
	}

	user := s.convertToUser(&resp.Data)
	log.WithFields(log.Fields{
		"user_id":  user.ID,
		"username": user.Username,
	}).Info("Đã lấy thông tin user thành công")

	return user, nil
}

// ListUsers lấy danh sách users theo IDs
func (s *TwitterService) ListUsers(ctx context.Context, userIDs []string) (*models.UsersListResponse, error) {
	log.WithField("count", len(userIDs)).Info("Đang lấy danh sách users")

	if len(userIDs) == 0 {
		return &models.UsersListResponse{Users: []models.User{}}, nil
	}

	// Twitter API chỉ cho phép tối đa 100 IDs mỗi request
	if len(userIDs) > 100 {
		userIDs = userIDs[:100]
	}

	params := &userlookupTypes.ListInput{
		IDs: userIDs,
		UserFields: fields.UserFieldList{
			fields.UserFieldID,
			fields.UserFieldName,
			fields.UserFieldUsername,
			fields.UserFieldDescription,
			fields.UserFieldProfileImageUrl,
			fields.UserFieldVerified,
			fields.UserFieldCreatedAt,
			fields.UserFieldPublicMetrics,
		},
	}

	resp, err := userlookup.List(ctx, s.client, params)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy danh sách users: %w", err)
	}

	users := make([]models.User, 0, len(resp.Data))
	for i := range resp.Data {
		users = append(users, *s.convertToUser(&resp.Data[i]))
	}

	result := &models.UsersListResponse{
		Users: users,
		Meta: &models.Meta{
			ResultCount: len(users),
		},
	}

	log.WithField("users_count", len(users)).Info("Đã lấy danh sách users thành công")
	return result, nil
}

// GetMe lấy thông tin user hiện tại (authenticated user)
func (s *TwitterService) GetMe(ctx context.Context) (*models.User, error) {
	log.Info("Đang lấy thông tin authenticated user")

	params := &userlookupTypes.GetMeInput{
		UserFields: fields.UserFieldList{
			fields.UserFieldID,
			fields.UserFieldName,
			fields.UserFieldUsername,
			fields.UserFieldDescription,
			fields.UserFieldProfileImageUrl,
			fields.UserFieldVerified,
			fields.UserFieldCreatedAt,
			fields.UserFieldPublicMetrics,
		},
	}

	resp, err := userlookup.GetMe(ctx, s.client, params)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy thông tin authenticated user: %w", err)
	}

	if resp.Data.ID == nil || gotwi.StringValue(resp.Data.ID) == "" {
		return nil, fmt.Errorf("không thể lấy thông tin authenticated user")
	}

	user := s.convertToUser(&resp.Data)
	log.WithFields(log.Fields{
		"user_id":  user.ID,
		"username": user.Username,
	}).Info("Đã lấy thông tin authenticated user thành công")

	return user, nil
}

// GetBlockingUsers lấy danh sách users bị block
// Lưu ý: API này yêu cầu OAuth 1.0a với authenticated user context
// Chỉ có thể xem blocking list của chính authenticated user, không thể xem của user khác
func (s *TwitterService) GetBlockingUsers(ctx context.Context, username string, maxResults int, paginationToken string) (*models.BlockingUsersResponse, error) {
	log.WithFields(log.Fields{
		"username":    username,
		"max_results": maxResults,
		"page_token":  paginationToken,
	}).Info("Đang lấy danh sách blocking users")

	// Blocking API chỉ hoạt động với authenticated user (OAuth 1.0a)
	// Với Bearer Token, chúng ta chỉ có thể lấy blocking list của authenticated user
	// Tạm thời trả về lỗi thông báo rõ ràng
	return nil, fmt.Errorf("API blocking users yêu cầu OAuth 1.0a với authenticated user context. Bearer Token chỉ hỗ trợ xem blocking list của chính authenticated user. Vui lòng sử dụng OAuth 1.0a để truy cập API này")
}

// GetMutingUsers lấy danh sách users bị mute
// Lưu ý: API này yêu cầu OAuth 1.0a với authenticated user context
// Chỉ có thể xem muting list của chính authenticated user, không thể xem của user khác
func (s *TwitterService) GetMutingUsers(ctx context.Context, username string, maxResults int, paginationToken string) (*models.MutingUsersResponse, error) {
	log.WithFields(log.Fields{
		"username":    username,
		"max_results": maxResults,
		"page_token":  paginationToken,
	}).Info("Đang lấy danh sách muting users")

	// Muting API chỉ hoạt động với authenticated user (OAuth 1.0a)
	// Với Bearer Token, chúng ta chỉ có thể lấy muting list của authenticated user
	// Tạm thời trả về lỗi thông báo rõ ràng
	return nil, fmt.Errorf("API muting users yêu cầu OAuth 1.0a với authenticated user context. Bearer Token chỉ hỗ trợ xem muting list của chính authenticated user. Vui lòng sử dụng OAuth 1.0a để truy cập API này")
}

// HideTweet ẩn/hiện một tweet (chỉ áp dụng cho authenticated user's tweets)
// Lưu ý: Twitter API v2 với Bearer Token không hỗ trợ hide tweet trực tiếp
// Endpoint này yêu cầu OAuth 1.0a với write permissions
func (s *TwitterService) HideTweet(ctx context.Context, tweetID string, hidden bool) (*models.HideTweetResponse, error) {
	log.WithFields(log.Fields{
		"tweet_id": tweetID,
		"hidden":   hidden,
	}).Info("Đang thay đổi trạng thái hidden của tweet")

	// Twitter API v2 không hỗ trợ hide tweet với Bearer Token
	// Cần OAuth 1.0a với write permissions
	// Tạm thời trả về thông báo lỗi
	return nil, fmt.Errorf("API hide tweet không được hỗ trợ với Bearer Token. Cần OAuth 1.0a với write permissions")
}

// GetUserTimelineReverseChronological lấy timeline reverse chronological của user
// Lưu ý: API này yêu cầu OAuth 1.0a với authenticated user context cho reverse chronological
// Với Bearer Token, chúng ta sử dụng GetUserTweets (cũng là reverse chronological)
func (s *TwitterService) GetUserTimelineReverseChronological(ctx context.Context, username string, maxResults int) (*models.TweetsResponse, error) {
	log.WithFields(log.Fields{
		"username":    username,
		"max_results": maxResults,
	}).Info("Đang lấy timeline reverse chronological")

	// Sử dụng cùng method với GetUserTweets vì đó là reverse chronological
	// Timeline reverse chronological endpoint yêu cầu OAuth 1.0a
	// Nhưng GetUserTweets cũng trả về reverse chronological nên có thể dùng được
	return s.GetUserTweets(ctx, username, maxResults)
}

// GetRepostsOfMe lấy danh sách reposts của authenticated user
// Lưu ý: API này yêu cầu OAuth 1.0a với authenticated user context
func (s *TwitterService) GetRepostsOfMe(ctx context.Context, maxResults int) (*models.RepostsResponse, error) {
	log.WithField("max_results", maxResults).Info("Đang lấy reposts của authenticated user")

	// Reposts Of Me API yêu cầu OAuth 1.0a với authenticated user context
	// Bearer Token không hỗ trợ API này
	return nil, fmt.Errorf("API reposts_of_me yêu cầu OAuth 1.0a với authenticated user context. Bearer Token không hỗ trợ API này. Vui lòng sử dụng OAuth 1.0a để truy cập")
}
