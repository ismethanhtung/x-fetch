package services

import (
	"context"
	"fmt"
	"x-twitter-backend/config"
	"x-twitter-backend/models"

	"github.com/michimani/gotwi"
	"github.com/michimani/gotwi/fields"
	"github.com/michimani/gotwi/resources"
	"github.com/michimani/gotwi/tweet/timeline"
	timelineTypes "github.com/michimani/gotwi/tweet/timeline/types"
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
