#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

APP_NAME="weather-microservice"
DOCKER_IMAGE="weather-microservice:latest"

echo -e "${BLUE}🚀 Weather Microservice Deployment Script${NC}"
echo ""

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to deploy to Docker
deploy_docker() {
    echo -e "${YELLOW}🐳 Deploying with Docker...${NC}"
    
    # Build image
    echo -e "${YELLOW}Building Docker image...${NC}"
    docker build -t $DOCKER_IMAGE .
    
    # Stop existing container
    if [ "$(docker ps -aq -f name=$APP_NAME)" ]; then
        echo -e "${YELLOW}Stopping existing container...${NC}"
        docker stop $APP_NAME || true
        docker rm $APP_NAME || true
    fi
    
    # Run new container
    echo -e "${YELLOW}Starting new container...${NC}"
    docker run -d \
        --name $APP_NAME \
        -p 8080:8080 \
        -v $(pwd)/.apiConfig:/app/.apiConfig:ro \
        --restart unless-stopped \
        $DOCKER_IMAGE
    
    echo -e "${GREEN}✅ Docker deployment complete!${NC}"
    echo -e "${GREEN}🌐 Service available at: http://localhost:8080${NC}"
}

# Function to deploy with Docker Compose
deploy_docker_compose() {
    echo -e "${YELLOW}🐳 Deploying with Docker Compose...${NC}"
    
    docker-compose down
    docker-compose build
    docker-compose up -d
    
    echo -e "${GREEN}✅ Docker Compose deployment complete!${NC}"
    echo -e "${GREEN}🌐 Weather API: http://localhost:8080${NC}"
    echo -e "${GREEN}📊 Prometheus: http://localhost:9090${NC}"
    echo -e "${GREEN}📈 Grafana: http://localhost:3000${NC}"
}

# Function to deploy to Kubernetes
deploy_kubernetes() {
    echo -e "${YELLOW}☸️  Deploying to Kubernetes...${NC}"
    
    if ! command_exists kubectl; then
        echo -e "${RED}❌ kubectl is not installed${NC}"
        exit 1
    fi
    
    # Apply configurations
    kubectl apply -f k8s/deployment.yaml
    
    # Wait for deployment
    echo -e "${YELLOW}Waiting for deployment to be ready...${NC}"
    kubectl wait --for=condition=available --timeout=300s \
        deployment/weather-microservice -n weather-microservice
    
    echo -e "${GREEN}✅ Kubernetes deployment complete!${NC}"
    echo -e "${YELLOW}Getting service information...${NC}"
    kubectl get svc -n weather-microservice
}

# Function to deploy to cloud platform
deploy_cloud() {
    echo -e "${YELLOW}☁️  Cloud Deployment Options:${NC}"
    echo ""
    echo "1. AWS (ECS/EKS)"
    echo "2. GCP (Cloud Run/GKE)"
    echo "3. Azure (Container Instances/AKS)"
    echo "4. Heroku"
    echo "5. Render"
    echo ""
    read -p "Select platform (1-5): " choice
    
    case $choice in
        1)
            echo -e "${YELLOW}🌩️  AWS Deployment${NC}"
            echo "Please use AWS CLI or Console to deploy"
            echo "Docker image: $DOCKER_IMAGE"
            ;;
        2)
            echo -e "${YELLOW}☁️  GCP Deployment${NC}"
            echo "Run: gcloud run deploy $APP_NAME --image $DOCKER_IMAGE --platform managed"
            ;;
        3)
            echo -e "${YELLOW}☁️  Azure Deployment${NC}"
            echo "Please use Azure CLI or Portal to deploy"
            ;;
        4)
            echo -e "${YELLOW}🟣 Heroku Deployment${NC}"
            echo "Run: heroku container:push web && heroku container:release web"
            ;;
        5)
            echo -e "${YELLOW}🎨 Render Deployment${NC}"
            echo "Connect your GitHub repo to Render and deploy"
            ;;
        *)
            echo -e "${RED}Invalid choice${NC}"
            ;;
    esac
}

# Main menu
echo "Select deployment option:"
echo "1. Docker"
echo "2. Docker Compose (with monitoring)"
echo "3. Kubernetes"
echo "4. Cloud Platforms"
echo "5. Build only (no deployment)"
echo ""
read -p "Enter choice (1-5): " choice

case $choice in
    1)
        if ! command_exists docker; then
            echo -e "${RED}❌ Docker is not installed${NC}"
            exit 1
        fi
        deploy_docker
        ;;
    2)
        if ! command_exists docker-compose; then
            echo -e "${RED}❌ Docker Compose is not installed${NC}"
            exit 1
        fi
        deploy_docker_compose
        ;;
    3)
        deploy_kubernetes
        ;;
    4)
        deploy_cloud
        ;;
    5)
        echo -e "${YELLOW}🔨 Building only...${NC}"
        ./scripts/build.sh
        ;;
    *)
        echo -e "${RED}❌ Invalid choice${NC}"
        exit 1
        ;;
esac

echo ""
echo -e "${GREEN}🎉 Deployment process completed!${NC}"