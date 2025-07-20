#!/bin/bash

# GoCart Quick Start Script
# This script starts all GoCart services

echo "🚀 Starting GoCart services..."

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# Start Docker Compose services
print_status "Starting Docker Compose services..."
if docker-compose up --build -d; then
    print_success "All services started successfully!"
    echo ""
    echo "📋 Service URLs:"
    echo "  🛍️  Product Service:  http://localhost:8080"
    echo "  👥 User Service:     http://localhost:8081" 
    echo "  🛒 Order Service:    http://localhost:8082"
    echo "  🌐 Frontend:         http://localhost:3000"
    echo "  🗄️  PostgreSQL:      localhost:5432"
    echo ""
    print_status "Use 'docker-compose logs -f' to view live logs"
    print_status "Use './cleanup.sh' to stop all services and free memory"
else
    echo "❌ Failed to start services"
    exit 1
fi 