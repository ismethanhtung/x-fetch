# 📝 Code Examples - Ví dụ Sử Dụng API

Các ví dụ thực tế để integrate với X Twitter Backend API.

## 📑 Mục Lục

1. [cURL Examples](#curl-examples)
2. [JavaScript/Node.js](#javascriptnodejs)
3. [Python](#python)
4. [Go](#go)
5. [PHP](#php)
6. [Java](#java)
7. [Shell Scripts](#shell-scripts)

---

## cURL Examples

### Basic Usage

```bash
# Health check
curl http://localhost:8080/health

# Get user info
curl http://localhost:8080/api/user/elonmusk

# Get tweets (default 10)
curl http://localhost:8080/api/tweets/user/elonmusk

# Get 20 tweets
curl "http://localhost:8080/api/tweets/user/elonmusk?count=20"

# Get following list (default 10)
curl http://localhost:8080/api/user/elonmusk/following

# Get 150 followings with pagination
curl "http://localhost:8080/api/user/elonmusk/following?count=150&pagination_token=YOUR_TOKEN"

# Pretty print JSON với jq
curl -s http://localhost:8080/api/tweets/user/elonmusk | jq '.'

# Chỉ lấy tweet text
curl -s http://localhost:8080/api/tweets/user/elonmusk | jq '.tweets[].text'

# Save to file
curl http://localhost:8080/api/tweets/user/elonmusk > tweets.json
```

### With Authentication (Future)

```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
     http://localhost:8080/api/tweets/user/elonmusk
```

---

## JavaScript/Node.js

### Using Fetch API (Browser)

```javascript
// Get user info
async function getUserInfo(username) {
  try {
    const response = await fetch(`http://localhost:8080/api/user/${username}`);
    const data = await response.json();
    
    console.log(`Name: ${data.name}`);
    console.log(`Username: @${data.username}`);
    console.log(`Followers: ${data.metrics.followers_count.toLocaleString()}`);
    
    return data;
  } catch (error) {
    console.error('Error:', error);
  }
}

// Get tweets
async function getTweets(username, count = 10) {
  try {
    const url = `http://localhost:8080/api/tweets/user/${username}?count=${count}`;
    const response = await fetch(url);
    const data = await response.json();
    
    console.log(`\nTweets from @${data.user.username}:\n`);
    
    data.tweets.forEach((tweet, index) => {
      console.log(`${index + 1}. ${tweet.text}`);
      console.log(`   ❤️ ${tweet.metrics.like_count} | 🔄 ${tweet.metrics.retweet_count}`);
      console.log('');
    });
    
    return data;
  } catch (error) {
    console.error('Error:', error);
  }
}

// Usage
getUserInfo('elonmusk');
getTweets('elonmusk', 5);
```

### Using Axios (Node.js)

```javascript
const axios = require('axios');

const API_BASE_URL = 'http://localhost:8080/api';

class TwitterAPI {
  async getUserInfo(username) {
    try {
      const { data } = await axios.get(`${API_BASE_URL}/user/${username}`);
      return data;
    } catch (error) {
      console.error('Error fetching user:', error.message);
      throw error;
    }
  }

  async getTweets(username, count = 10) {
    try {
      const { data } = await axios.get(
        `${API_BASE_URL}/tweets/user/${username}`,
        { params: { count } }
      );
      return data;
    } catch (error) {
      console.error('Error fetching tweets:', error.message);
      throw error;
    }
  }

  async getFollowing(username, count = 100, paginationToken) {
    try {
      const params = { count };
      if (paginationToken) {
        params.pagination_token = paginationToken;
      }

      const { data } = await axios.get(
        `${API_BASE_URL}/user/${username}/following`,
        { params }
      );
      return data;
    } catch (error) {
      console.error('Error fetching following list:', error.message);
      throw error;
    }
  }

  async displayUserTweets(username, count = 10) {
    const data = await this.getTweets(username, count);
    
    console.log(`\n📊 ${data.user.name} (@${data.user.username})`);
    console.log(`👥 Followers: ${data.user.metrics.followers_count.toLocaleString()}\n`);
    
    data.tweets.forEach((tweet, i) => {
      console.log(`${i + 1}. ${tweet.text}`);
      console.log(`   📅 ${new Date(tweet.created_at).toLocaleDateString()}`);
      console.log(`   ❤️ ${tweet.metrics.like_count} | 🔄 ${tweet.metrics.retweet_count} | 💬 ${tweet.metrics.reply_count}\n`);
    });
  }

  async displayFollowing(username, count = 50) {
    const data = await this.getFollowing(username, count);

    console.log(`\n👤 ${data.user.name} (@${data.user.username}) đang theo dõi:\n`);
    data.following.forEach((account, index) => {
      console.log(`${index + 1}. ${account.name} (@${account.username})`);
      console.log(`   ✅ Verified: ${account.verified ? 'Yes' : 'No'}`);
    });

    if (data.meta?.next_token) {
      console.log(`\n➡️ Next page token: ${data.meta.next_token}`);
    }
  }
}

// Usage
const api = new TwitterAPI();

(async () => {
  await api.displayUserTweets('elonmusk', 5);
  await api.displayUserTweets('BillGates', 5);
  await api.displayFollowing('elonmusk', 25);
})();
```

### React Component

```jsx
import React, { useState, useEffect } from 'react';

function TweetsFeed({ username, count = 10 }) {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchTweets = async () => {
      try {
        setLoading(true);
        const response = await fetch(
          `http://localhost:8080/api/tweets/user/${username}?count=${count}`
        );
        const json = await response.json();
        setData(json);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchTweets();
  }, [username, count]);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;
  if (!data) return null;

  return (
    <div className="tweets-feed">
      <div className="user-info">
        <h2>{data.user.name}</h2>
        <p>@{data.user.username}</p>
        <p>Followers: {data.user.metrics.followers_count.toLocaleString()}</p>
      </div>

      <div className="tweets">
        {data.tweets.map((tweet) => (
          <div key={tweet.id} className="tweet">
            <p>{tweet.text}</p>
            <div className="metrics">
              <span>❤️ {tweet.metrics.like_count}</span>
              <span>🔄 {tweet.metrics.retweet_count}</span>
              <span>💬 {tweet.metrics.reply_count}</span>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

// Usage
<TweetsFeed username="elonmusk" count={10} />
```

---

## Python

### Using Requests

```python
import requests
from datetime import datetime

API_BASE_URL = "http://localhost:8080/api"

class TwitterAPI:
    def __init__(self, base_url=API_BASE_URL):
        self.base_url = base_url
    
    def get_user_info(self, username):
        """Lấy thông tin user"""
        url = f"{self.base_url}/user/{username}"
        response = requests.get(url)
        response.raise_for_status()
        return response.json()
    
    def get_tweets(self, username, count=10):
        """Lấy tweets của user"""
        url = f"{self.base_url}/tweets/user/{username}"
        params = {"count": count}
        response = requests.get(url, params=params)
        response.raise_for_status()
        return response.json()

    def get_following(self, username, count=100, pagination_token=None):
        """Lấy danh sách accounts mà user đang theo dõi"""
        url = f"{self.base_url}/user/{username}/following"
        params = {"count": count}
        if pagination_token:
            params["pagination_token"] = pagination_token

        response = requests.get(url, params=params)
        response.raise_for_status()
        return response.json()
    
    def display_user_tweets(self, username, count=10):
        """Hiển thị tweets của user"""
        data = self.get_tweets(username, count)
        
        user = data['user']
        print(f"\n📊 {user['name']} (@{user['username']})")
        print(f"👥 Followers: {user['metrics']['followers_count']:,}\n")
        
        for i, tweet in enumerate(data['tweets'], 1):
            print(f"{i}. {tweet['text']}")
            
            # Format date
            created_at = datetime.fromisoformat(tweet['created_at'].replace('Z', '+00:00'))
            print(f"   📅 {created_at.strftime('%Y-%m-%d %H:%M')}")
            
            # Metrics
            metrics = tweet['metrics']
            print(f"   ❤️ {metrics['like_count']} | 🔄 {metrics['retweet_count']} | 💬 {metrics['reply_count']}\n")

    def display_following(self, username, count=50):
        """Hiển thị danh sách following"""
        data = self.get_following(username, count)

        print(f"\n👤 {data['user']['name']} (@{data['user']['username']}) đang theo dõi:")
        for i, account in enumerate(data['following'], 1):
            print(f"{i}. {account['name']} (@{account['username']}) - Verified: {account.get('verified', False)}")

        if data.get('meta', {}).get('next_token'):
            print(f"\n➡️ Next token: {data['meta']['next_token']}")

# Usage
if __name__ == "__main__":
    api = TwitterAPI()
    
    # Lấy tweets của nhiều users
    for username in ['elonmusk', 'BillGates', 'NASA']:
        try:
            api.display_user_tweets(username, count=5)
        except requests.exceptions.RequestException as e:
            print(f"Error fetching tweets for @{username}: {e}")

    # Lấy danh sách following của Elon Musk
    api.display_following("elonmusk", count=25)
```

### Using aiohttp (Async)

```python
import aiohttp
import asyncio

async def fetch_tweets(username, count=10):
    url = f"http://localhost:8080/api/tweets/user/{username}"
    params = {"count": count}
    
    async with aiohttp.ClientSession() as session:
        async with session.get(url, params=params) as response:
            return await response.json()

async def main():
    usernames = ['elonmusk', 'BillGates', 'NASA', 'cristiano']
    
    # Fetch all tweets concurrently
    tasks = [fetch_tweets(username, 5) for username in usernames]
    results = await asyncio.gather(*tasks)
    
    for data in results:
        user = data['user']
        print(f"\n@{user['username']}: {len(data['tweets'])} tweets")

# Run
asyncio.run(main())
```

---

## Go

```go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const APIBaseURL = "http://localhost:8080/api"

type TwitterClient struct {
	BaseURL string
	Client  *http.Client
}

type Tweet struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	CreatedAt string    `json:"created_at"`
	Metrics   struct {
		LikeCount    int `json:"like_count"`
		RetweetCount int `json:"retweet_count"`
		ReplyCount   int `json:"reply_count"`
	} `json:"metrics"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Metrics  struct {
		FollowersCount int `json:"followers_count"`
	} `json:"metrics"`
}

type TweetsResponse struct {
	Tweets []Tweet `json:"tweets"`
	User   User    `json:"user"`
}

func NewTwitterClient() *TwitterClient {
	return &TwitterClient{
		BaseURL: APIBaseURL,
		Client:  &http.Client{},
	}
}

func (c *TwitterClient) GetTweets(username string, count int) (*TweetsResponse, error) {
	endpoint := fmt.Sprintf("%s/tweets/user/%s", c.BaseURL, username)
	
	// Add query parameters
	params := url.Values{}
	params.Add("count", fmt.Sprintf("%d", count))
	
	resp, err := c.Client.Get(endpoint + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var result TweetsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	
	return &result, nil
}

func main() {
	client := NewTwitterClient()
	
	usernames := []string{"elonmusk", "BillGates", "NASA"}
	
	for _, username := range usernames {
		data, err := client.GetTweets(username, 5)
		if err != nil {
			fmt.Printf("Error fetching tweets for @%s: %v\n", username, err)
			continue
		}
		
		fmt.Printf("\n📊 %s (@%s)\n", data.User.Name, data.User.Username)
		fmt.Printf("👥 Followers: %d\n\n", data.User.Metrics.FollowersCount)
		
		for i, tweet := range data.Tweets {
			fmt.Printf("%d. %s\n", i+1, tweet.Text)
			fmt.Printf("   ❤️ %d | 🔄 %d | 💬 %d\n\n",
				tweet.Metrics.LikeCount,
				tweet.Metrics.RetweetCount,
				tweet.Metrics.ReplyCount)
		}
	}
}
```

---

## PHP

```php
<?php

class TwitterAPI {
    private $baseURL = 'http://localhost:8080/api';
    
    public function getUserInfo($username) {
        $url = "{$this->baseURL}/user/{$username}";
        return $this->makeRequest($url);
    }
    
    public function getTweets($username, $count = 10) {
        $url = "{$this->baseURL}/tweets/user/{$username}?count={$count}";
        return $this->makeRequest($url);
    }
    
    private function makeRequest($url) {
        $ch = curl_init();
        curl_setopt($ch, CURLOPT_URL, $url);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        
        $response = curl_exec($ch);
        $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        curl_close($ch);
        
        if ($httpCode !== 200) {
            throw new Exception("API request failed with code: {$httpCode}");
        }
        
        return json_decode($response, true);
    }
    
    public function displayUserTweets($username, $count = 10) {
        $data = $this->getTweets($username, $count);
        
        $user = $data['user'];
        echo "\n📊 {$user['name']} (@{$user['username']})\n";
        echo "👥 Followers: " . number_format($user['metrics']['followers_count']) . "\n\n";
        
        foreach ($data['tweets'] as $i => $tweet) {
            echo ($i + 1) . ". {$tweet['text']}\n";
            echo "   ❤️ {$tweet['metrics']['like_count']} | ";
            echo "🔄 {$tweet['metrics']['retweet_count']} | ";
            echo "💬 {$tweet['metrics']['reply_count']}\n\n";
        }
    }
}

// Usage
$api = new TwitterAPI();

try {
    $api->displayUserTweets('elonmusk', 5);
} catch (Exception $e) {
    echo "Error: " . $e->getMessage() . "\n";
}
```

---

## Java

```java
import com.google.gson.Gson;
import com.google.gson.annotations.SerializedName;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.List;

public class TwitterAPI {
    private static final String API_BASE_URL = "http://localhost:8080/api";
    private final HttpClient client;
    private final Gson gson;
    
    public TwitterAPI() {
        this.client = HttpClient.newHttpClient();
        this.gson = new Gson();
    }
    
    public TweetsResponse getTweets(String username, int count) throws Exception {
        String url = String.format("%s/tweets/user/%s?count=%d", 
            API_BASE_URL, username, count);
        
        HttpRequest request = HttpRequest.newBuilder()
            .uri(URI.create(url))
            .GET()
            .build();
        
        HttpResponse<String> response = client.send(request, 
            HttpResponse.BodyHandlers.ofString());
        
        if (response.statusCode() != 200) {
            throw new Exception("API request failed: " + response.statusCode());
        }
        
        return gson.fromJson(response.body(), TweetsResponse.class);
    }
    
    public void displayUserTweets(String username, int count) throws Exception {
        TweetsResponse data = getTweets(username, count);
        
        System.out.printf("\n📊 %s (@%s)\n", data.user.name, data.user.username);
        System.out.printf("👥 Followers: %,d\n\n", data.user.metrics.followersCount);
        
        for (int i = 0; i < data.tweets.size(); i++) {
            Tweet tweet = data.tweets.get(i);
            System.out.printf("%d. %s\n", i + 1, tweet.text);
            System.out.printf("   ❤️ %d | 🔄 %d | 💬 %d\n\n",
                tweet.metrics.likeCount,
                tweet.metrics.retweetCount,
                tweet.metrics.replyCount);
        }
    }
    
    // Data classes
    static class TweetsResponse {
        List<Tweet> tweets;
        User user;
    }
    
    static class Tweet {
        String id;
        String text;
        @SerializedName("created_at")
        String createdAt;
        Metrics metrics;
    }
    
    static class User {
        String id;
        String username;
        String name;
        UserMetrics metrics;
    }
    
    static class Metrics {
        @SerializedName("like_count")
        int likeCount;
        @SerializedName("retweet_count")
        int retweetCount;
        @SerializedName("reply_count")
        int replyCount;
    }
    
    static class UserMetrics {
        @SerializedName("followers_count")
        int followersCount;
    }
    
    public static void main(String[] args) {
        TwitterAPI api = new TwitterAPI();
        
        try {
            api.displayUserTweets("elonmusk", 5);
        } catch (Exception e) {
            System.err.println("Error: " + e.getMessage());
        }
    }
}
```

---

## Shell Scripts

### Monitor Multiple Accounts

```bash
#!/bin/bash
# monitor_accounts.sh - Monitor tweets from multiple accounts

API_URL="http://localhost:8080/api"
ACCOUNTS=("elonmusk" "BillGates" "NASA" "cristiano")
COUNT=5

echo "🐦 Twitter Accounts Monitor"
echo "======================================"

for account in "${ACCOUNTS[@]}"; do
    echo -e "\n📊 Fetching tweets from @$account..."
    
    response=$(curl -s "${API_URL}/tweets/user/${account}?count=${COUNT}")
    
    if [ $? -eq 0 ]; then
        # Parse with jq
        name=$(echo "$response" | jq -r '.user.name')
        followers=$(echo "$response" | jq -r '.user.metrics.followers_count')
        
        echo "Name: $name"
        echo "Followers: $(printf "%'d" $followers)"
        echo -e "\nRecent tweets:"
        
        echo "$response" | jq -r '.tweets[] | "- \(.text)"' | head -n 3
    else
        echo "❌ Error fetching tweets for @$account"
    fi
    
    echo "--------------------------------------"
done
```

### Continuous Monitoring

```bash
#!/bin/bash
# continuous_monitor.sh - Monitor tweets every N minutes

INTERVAL=300  # 5 minutes
USERNAME="elonmusk"

while true; do
    clear
    echo "🐦 Monitoring @$USERNAME - $(date)"
    echo "======================================"
    
    curl -s "http://localhost:8080/api/tweets/user/${USERNAME}?count=5" | \
        jq -r '.tweets[] | "\(.created_at) | \(.text) | ❤️ \(.metrics.like_count)"'
    
    echo -e "\n⏰ Next update in $INTERVAL seconds..."
    sleep $INTERVAL
done
```

### Save Tweets to CSV

```bash
#!/bin/bash
# save_to_csv.sh - Save tweets to CSV file

USERNAME=$1
COUNT=${2:-10}
OUTPUT="tweets_${USERNAME}_$(date +%Y%m%d_%H%M%S).csv"

if [ -z "$USERNAME" ]; then
    echo "Usage: $0 <username> [count]"
    exit 1
fi

echo "Fetching tweets from @$USERNAME..."

curl -s "http://localhost:8080/api/tweets/user/${USERNAME}?count=${COUNT}" | \
    jq -r '.tweets[] | [.id, .created_at, .author_id, .text, .metrics.like_count, .metrics.retweet_count, .metrics.reply_count] | @csv' > "$OUTPUT"

# Add header
sed -i '1i"Tweet ID","Created At","Author ID","Text","Likes","Retweets","Replies"' "$OUTPUT"

echo "✅ Saved to $OUTPUT"
echo "📊 Total tweets: $(wc -l < "$OUTPUT" | xargs expr -1 +)"
```

---

## 🔗 Integration Tips

### Error Handling

Luôn handle errors khi call API:

```javascript
async function safeFetchTweets(username) {
  try {
    const response = await fetch(`http://localhost:8080/api/tweets/user/${username}`);
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Error fetching tweets:', error);
    return null;
  }
}
```

### Rate Limiting

Implement rate limiting trong code:

```javascript
const rateLimit = require('express-rate-limit');

const limiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 100 // limit each IP to 100 requests per windowMs
});

app.use('/api/', limiter);
```

### Caching

Cache responses để giảm API calls:

```javascript
const NodeCache = require('node-cache');
const cache = new NodeCache({ stdTTL: 300 }); // 5 minutes

async function getCachedTweets(username, count) {
  const cacheKey = `tweets_${username}_${count}`;
  
  // Check cache
  const cached = cache.get(cacheKey);
  if (cached) {
    console.log('Returning cached data');
    return cached;
  }
  
  // Fetch from API
  const data = await fetchTweets(username, count);
  cache.set(cacheKey, data);
  
  return data;
}
```

---

**Có thêm ví dụ nào bạn cần? Mở issue hoặc PR!** 🚀

