#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🔨 Building Weather Microservice...${NC}"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed. Please install Go 1.21 or higher.${NC}"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.21"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo -e "${RED}❌ Go version $REQUIRED_VERSION or higher is required. Current version: $GO_VERSION${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Go version: $GO_VERSION${NC}"

# Create bin directory if it doesn't exist
mkdir -p bin

# Download dependencies
echo -e "${YELLOW}📦 Downloading dependencies...${NC}"
go mod download
go mod tidy

# Run tests
echo -e "${YELLOW}🧪 Running tests...${NC}"
go test -v ./... || {
    echo -e "${RED}❌ Tests failed!${NC}"
    exit 1
}

# Build for current platform
echo -e "${YELLOW}🔨 Building for current platform...${NC}"
go build -o bin/weather-server -ldflags="-s -w" ./cmd/server/main.go

echo -e "${GREEN}✅ Build successful!${NC}"
echo -e "${GREEN}📦 Binary: bin/weather-server${NC}"

# Build for multiple platforms
echo -e "${YELLOW}🔨 Building for multiple platforms...${NC}"

# Linux
GOOS=linux GOARCH=amd64 go build -o bin/weather-server-linux-amd64 -ldflags="-s -w" ./cmd/server/main.go
echo -e "${GREEN}✅ Built: bin/weather-server-linux-amd64${NC}"

# macOS
GOOS=darwin GOARCH=amd64 go build -o bin/weather-server-darwin-amd64 -ldflags="-s -w" ./cmd/server/main.go
echo -e "${GREEN}✅ Built: bin/weather-server-darwin-amd64${NC}"

# Windows
GOOS=windows GOARCH=amd64 go build -o bin/weather-server-windows-amd64.exe -ldflags="-s -w" ./cmd/server/main.go
echo -e "${GREEN}✅ Built: bin/weather-server-windows-amd64.exe${NC}"

# ARM64 for Raspberry Pi / ARM servers
GOOS=linux GOARCH=arm64 go build -o bin/weather-server-linux-arm64 -ldflags="-s -w" ./cmd/server/main.go
echo -e "${GREEN}✅ Built: bin/weather-server-linux-arm64${NC}"

echo ""
echo -e "${GREEN}🎉 All builds completed successfully!${NC}"
echo -e "${GREEN}📦 Binaries are in the bin/ directory${NC}"
echo ""
echo -e "${YELLOW}To run the server:${NC}"
echo -e "  ${GREEN}./bin/weather-server${NC}"