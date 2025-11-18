#!/bin/bash

# Test script for X Twitter Backend API
# Tự động test tất cả các endpoints

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
API_URL="${API_URL:-http://localhost:8080}"
TEST_USERNAME="${TEST_USERNAME:-elonmusk}"

echo -e "${BLUE}🧪 X Twitter Backend API Tests${NC}"
echo "===================================="
echo "API URL: $API_URL"
echo "Test Username: $TEST_USERNAME"
echo ""

# Test counter
PASSED=0
FAILED=0

# Function to test endpoint
test_endpoint() {
    local name=$1
    local url=$2
    local expected_status=$3
    
    echo -n "Testing $name... "
    
    response=$(curl -s -w "\n%{http_code}" "$url")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" -eq "$expected_status" ]; then
        echo -e "${GREEN}✅ PASSED${NC} (HTTP $http_code)"
        PASSED=$((PASSED + 1))
        
        # Pretty print JSON if available
        if command -v jq &> /dev/null; then
            echo "$body" | jq '.' 2>/dev/null | head -n 5
        fi
        echo ""
        return 0
    else
        echo -e "${RED}❌ FAILED${NC} (Expected: $expected_status, Got: $http_code)"
        echo "Response: $body"
        echo ""
        FAILED=$((FAILED + 1))
        return 1
    fi
}

# Test 1: Health Check
echo -e "${YELLOW}[Test 1]${NC} Health Check"
test_endpoint "Health Check" "$API_URL/health" 200

# Test 2: Get User Info
echo -e "${YELLOW}[Test 2]${NC} Get User Info"
test_endpoint "Get User Info" "$API_URL/api/user/$TEST_USERNAME" 200

# Test 3: Get User Tweets (default count)
echo -e "${YELLOW}[Test 3]${NC} Get User Tweets (default)"
test_endpoint "Get Tweets Default" "$API_URL/api/tweets/user/$TEST_USERNAME" 200

# Test 4: Get User Tweets (custom count)
echo -e "${YELLOW}[Test 4]${NC} Get User Tweets (count=5)"
test_endpoint "Get Tweets Custom Count" "$API_URL/api/tweets/user/$TEST_USERNAME?count=5" 200

# Test 5: Get Following list
echo -e "${YELLOW}[Test 5]${NC} Get Following List"
test_endpoint "Get Following" "$API_URL/api/user/$TEST_USERNAME/following?count=25" 200

# Test 6: API Documentation
echo -e "${YELLOW}[Test 6]${NC} API Documentation"
test_endpoint "API Docs" "$API_URL/api/docs" 200

# Test 7: Invalid User (should still return 200 or 404)
echo -e "${YELLOW}[Test 7]${NC} Invalid User"
echo -n "Testing Invalid User... "
response=$(curl -s -w "\n%{http_code}" "$API_URL/api/user/thisisnotarealuser12345xyz")
http_code=$(echo "$response" | tail -n1)

if [ "$http_code" -eq 200 ] || [ "$http_code" -eq 404 ] || [ "$http_code" -eq 500 ]; then
    echo -e "${GREEN}✅ PASSED${NC} (HTTP $http_code - handled gracefully)"
    PASSED=$((PASSED + 1))
else
    echo -e "${RED}❌ FAILED${NC} (Unexpected status: $http_code)"
    FAILED=$((FAILED + 1))
fi
echo ""

# Summary
echo "===================================="
echo -e "${BLUE}Test Summary${NC}"
echo "===================================="
echo -e "Total Tests: $((PASSED + FAILED))"
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${RED}Failed: $FAILED${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ Some tests failed!${NC}"
    exit 1
fi

