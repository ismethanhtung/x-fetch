#!/bin/bash

# Monitor script - Continuously monitor Twitter accounts
# Usage: ./scripts/monitor.sh [interval_seconds]

# Configuration
INTERVAL=${1:-300}  # Default 5 minutes (300 seconds)
API_URL="${API_URL:-http://localhost:8080}"

# Accounts to monitor
ACCOUNTS=(
    "elonmusk"
    "BillGates"
    "NASA"
    "cristiano"
)

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${BLUE}🐦 Twitter Accounts Monitor${NC}"
echo "===================================="
echo "Monitoring ${#ACCOUNTS[@]} accounts"
echo "Update interval: ${INTERVAL} seconds"
echo "API: ${API_URL}"
echo "===================================="
echo ""
echo "Press Ctrl+C to stop"
echo ""

# Function to fetch and display tweets
fetch_tweets() {
    local username=$1
    local count=${2:-3}
    
    echo -e "${CYAN}📊 @${username}${NC}"
    echo "$(date '+%Y-%m-%d %H:%M:%S')"
    echo "---"
    
    response=$(curl -s "${API_URL}/api/tweets/user/${username}?count=${count}")
    
    if [ $? -eq 0 ]; then
        # Check if jq is available
        if command -v jq &> /dev/null; then
            # Get user info
            name=$(echo "$response" | jq -r '.user.name')
            followers=$(echo "$response" | jq -r '.user.metrics.followers_count')
            
            echo "Name: $name"
            echo "Followers: $(printf "%'d" $followers 2>/dev/null || echo $followers)"
            echo ""
            echo "Latest tweets:"
            
            # Display tweets
            echo "$response" | jq -r '.tweets[] | "• \(.text)\n  ❤️ \(.metrics.like_count) | 🔄 \(.metrics.retweet_count) | 💬 \(.metrics.reply_count)\n"'
        else
            # Fallback if jq not available
            echo "$response" | head -n 20
        fi
    else
        echo -e "${RED}❌ Failed to fetch tweets${NC}"
    fi
    
    echo ""
}

# Cleanup on exit
trap 'echo -e "\n${YELLOW}Monitoring stopped${NC}"; exit 0' INT TERM

# Main monitoring loop
iteration=0
while true; do
    iteration=$((iteration + 1))
    
    # Clear screen
    clear
    
    echo -e "${BLUE}🐦 Twitter Accounts Monitor - Iteration #${iteration}${NC}"
    echo "===================================="
    echo "Last update: $(date '+%Y-%m-%d %H:%M:%S')"
    echo "Next update in: ${INTERVAL} seconds"
    echo "===================================="
    echo ""
    
    # Fetch tweets for each account
    for account in "${ACCOUNTS[@]}"; do
        fetch_tweets "$account" 3
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        echo ""
    done
    
    # Wait before next iteration
    echo -e "${YELLOW}⏰ Next update in ${INTERVAL} seconds...${NC}"
    echo "Press Ctrl+C to stop"
    
    sleep $INTERVAL
done

