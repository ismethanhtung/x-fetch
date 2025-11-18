#!/bin/bash

# Setup script for X Twitter Backend
# Tự động cài đặt và cấu hình project

set -e  # Exit on error

echo "🚀 X Twitter Backend - Setup Script"
echo "===================================="

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check Go installation
echo -e "\n${YELLOW}[1/5]${NC} Checking Go installation..."
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed!${NC}"
    echo "Please install Go from: https://golang.org/dl/"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo -e "${GREEN}✅ Go is installed: ${GO_VERSION}${NC}"

# Check Go version
REQUIRED_GO_VERSION="go1.21"
if [[ "$GO_VERSION" < "$REQUIRED_GO_VERSION" ]]; then
    echo -e "${YELLOW}⚠️  Warning: Go 1.21 or higher is recommended${NC}"
fi

# Install dependencies
echo -e "\n${YELLOW}[2/5]${NC} Installing dependencies..."
go mod download
go mod tidy
echo -e "${GREEN}✅ Dependencies installed${NC}"

# Create .env file
echo -e "\n${YELLOW}[3/5]${NC} Setting up environment file..."
if [ -f .env ]; then
    echo -e "${YELLOW}⚠️  .env file already exists, skipping...${NC}"
else
    cp ENV_EXAMPLE .env
    echo -e "${GREEN}✅ Created .env file${NC}"
    echo -e "${YELLOW}⚠️  Please edit .env and add your TWITTER_BEARER_TOKEN${NC}"
fi

# Build application
echo -e "\n${YELLOW}[4/5]${NC} Building application..."
go build -o twitter-backend main.go
echo -e "${GREEN}✅ Application built successfully${NC}"

# Test run
echo -e "\n${YELLOW}[5/5]${NC} Testing configuration..."
if [ -z "$(grep 'TWITTER_BEARER_TOKEN=.*[^[:space:]]' .env | grep -v '^#')" ]; then
    echo -e "${YELLOW}⚠️  TWITTER_BEARER_TOKEN is not set in .env${NC}"
    echo -e "${YELLOW}   Please add your token before running the server${NC}"
else
    echo -e "${GREEN}✅ Configuration looks good${NC}"
fi

# Summary
echo ""
echo "===================================="
echo -e "${GREEN}🎉 Setup completed successfully!${NC}"
echo ""
echo "Next steps:"
echo "  1. Edit .env and add your TWITTER_BEARER_TOKEN"
echo "  2. Run the server:"
echo "     • Using Go: go run main.go"
echo "     • Using binary: ./twitter-backend"
echo "     • Using Make: make run"
echo ""
echo "Documentation:"
echo "  • Quick Start: QUICKSTART_VI.md"
echo "  • Full Tutorial: TUTORIAL_VI.md"
echo "  • Examples: EXAMPLES.md"
echo ""
echo "Happy coding! 🚀"
echo "===================================="

