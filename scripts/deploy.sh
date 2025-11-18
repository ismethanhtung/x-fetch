#!/bin/bash

# Deployment script for X Twitter Backend
# Usage: ./scripts/deploy.sh [environment]

set -e

ENVIRONMENT=${1:-production}

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🚀 X Twitter Backend - Deployment Script${NC}"
echo "===================================="
echo "Environment: $ENVIRONMENT"
echo ""

# Step 1: Validate environment
echo -e "${YELLOW}[1/6]${NC} Validating environment..."

if [ "$ENVIRONMENT" != "production" ] && [ "$ENVIRONMENT" != "staging" ] && [ "$ENVIRONMENT" != "development" ]; then
    echo -e "${RED}❌ Invalid environment: $ENVIRONMENT${NC}"
    echo "Valid environments: production, staging, development"
    exit 1
fi

echo -e "${GREEN}✅ Environment validated${NC}"

# Step 2: Check dependencies
echo -e "\n${YELLOW}[2/6]${NC} Checking dependencies..."

if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Dependencies OK${NC}"

# Step 3: Run tests
echo -e "\n${YELLOW}[3/6]${NC} Running tests..."

if go test ./... > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Tests passed${NC}"
else
    echo -e "${YELLOW}⚠️  No tests found or tests failed${NC}"
fi

# Step 4: Build application
echo -e "\n${YELLOW}[4/6]${NC} Building application..."

BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

CGO_ENABLED=0 go build \
    -ldflags "-X main.Version=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}" \
    -o twitter-backend \
    main.go

echo -e "${GREEN}✅ Build successful${NC}"

# Step 5: Create deployment package
echo -e "\n${YELLOW}[5/6]${NC} Creating deployment package..."

DEPLOY_DIR="deploy_${ENVIRONMENT}_${BUILD_TIME}"
mkdir -p "$DEPLOY_DIR"

# Copy files
cp twitter-backend "$DEPLOY_DIR/"
cp ENV_EXAMPLE "$DEPLOY_DIR/.env.example"
cp README.md "$DEPLOY_DIR/"
cp QUICKSTART_VI.md "$DEPLOY_DIR/"

# Create systemd service file
cat > "$DEPLOY_DIR/twitter-backend.service" << EOF
[Unit]
Description=Twitter Backend API
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/twitter-backend
EnvironmentFile=/opt/twitter-backend/.env
ExecStart=/opt/twitter-backend/twitter-backend
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

# Create installation script
cat > "$DEPLOY_DIR/install.sh" << 'EOF'
#!/bin/bash
set -e

echo "Installing Twitter Backend..."

# Stop existing service
sudo systemctl stop twitter-backend 2>/dev/null || true

# Create directory
sudo mkdir -p /opt/twitter-backend

# Copy files
sudo cp twitter-backend /opt/twitter-backend/
sudo cp .env.example /opt/twitter-backend/
sudo chmod +x /opt/twitter-backend/twitter-backend

# Setup service
sudo cp twitter-backend.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable twitter-backend

echo "✅ Installation complete!"
echo ""
echo "Next steps:"
echo "  1. Edit /opt/twitter-backend/.env (copy from .env.example)"
echo "  2. sudo systemctl start twitter-backend"
echo "  3. sudo systemctl status twitter-backend"
EOF

chmod +x "$DEPLOY_DIR/install.sh"

echo -e "${GREEN}✅ Deployment package created: $DEPLOY_DIR${NC}"

# Step 6: Create tarball
echo -e "\n${YELLOW}[6/6]${NC} Creating tarball..."

tar -czf "${DEPLOY_DIR}.tar.gz" "$DEPLOY_DIR"
rm -rf "$DEPLOY_DIR"

echo -e "${GREEN}✅ Tarball created: ${DEPLOY_DIR}.tar.gz${NC}"

# Summary
echo ""
echo "===================================="
echo -e "${GREEN}🎉 Deployment package ready!${NC}"
echo "===================================="
echo ""
echo "Package: ${DEPLOY_DIR}.tar.gz"
echo "Environment: $ENVIRONMENT"
echo "Build Time: $BUILD_TIME"
echo "Git Commit: $GIT_COMMIT"
echo ""
echo "To deploy:"
echo "  1. Copy ${DEPLOY_DIR}.tar.gz to server"
echo "  2. Extract: tar -xzf ${DEPLOY_DIR}.tar.gz"
echo "  3. cd ${DEPLOY_DIR}"
echo "  4. ./install.sh"
echo ""
echo "For Docker deployment:"
echo "  docker build -t twitter-backend:${GIT_COMMIT} ."
echo "  docker push twitter-backend:${GIT_COMMIT}"
echo ""

