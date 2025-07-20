#!/bin/bash

# GoCart Memory Cleanup Script
# This script stops all services and frees up memory

echo "🧹 Starting GoCart cleanup..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Step 1: Stop Docker Compose services
print_status "Stopping Docker Compose services..."
if docker-compose down 2>/dev/null; then
    print_success "Docker Compose services stopped"
else
    print_warning "No Docker Compose services running or docker-compose.yml not found"
fi

# Step 2: Stop any remaining Docker containers
print_status "Stopping any remaining Docker containers..."
RUNNING_CONTAINERS=$(docker ps -q)
if [ -n "$RUNNING_CONTAINERS" ]; then
    docker stop $RUNNING_CONTAINERS
    print_success "Stopped running Docker containers"
else
    print_success "No running Docker containers found"
fi

# Step 3: Kill processes on common development ports
print_status "Killing processes on development ports..."

PORTS=(3000 8080 8081 8082 5432)
KILLED_PROCESSES=0

for PORT in "${PORTS[@]}"; do
    PIDS=$(lsof -ti:$PORT 2>/dev/null)
    if [ -n "$PIDS" ]; then
        echo "  Killing processes on port $PORT: $PIDS"
        kill -9 $PIDS 2>/dev/null
        KILLED_PROCESSES=$((KILLED_PROCESSES + 1))
    fi
done

if [ $KILLED_PROCESSES -gt 0 ]; then
    print_success "Killed processes on $KILLED_PROCESSES port(s)"
else
    print_success "No processes found on development ports"
fi

# Step 4: Kill any localhost processes (additional cleanup)
print_status "Killing any remaining localhost processes..."
pkill -f "localhost:808" 2>/dev/null || true
pkill -f "localhost:3000" 2>/dev/null || true
pkill -f "localhost:5432" 2>/dev/null || true

# Step 5: Docker system cleanup
print_status "Cleaning up Docker resources..."
CLEANUP_OUTPUT=$(docker system prune -f 2>&1)
if echo "$CLEANUP_OUTPUT" | grep -q "Total reclaimed space"; then
    RECLAIMED=$(echo "$CLEANUP_OUTPUT" | grep "Total reclaimed space" | awk '{print $4 " " $5}')
    print_success "Docker cleanup completed - Reclaimed: $RECLAIMED"
else
    print_success "Docker cleanup completed - No resources to reclaim"
fi

# Step 6: Final verification
print_status "Verifying cleanup..."
REMAINING_PROCESSES=$(lsof -i :8080 -i :8081 -i :8082 -i :3000 -i :5432 2>/dev/null | wc -l)
if [ $REMAINING_PROCESSES -eq 0 ]; then
    print_success "All target ports are now free"
else
    print_warning "Some processes may still be running on target ports"
    lsof -i :8080 -i :8081 -i :8082 -i :3000 -i :5432 2>/dev/null || true
fi

# Step 7: Show memory status (optional)
print_status "Current memory usage:"
if command -v free >/dev/null 2>&1; then
    free -h
elif command -v vm_stat >/dev/null 2>&1; then
    # macOS memory info
    echo "  $(vm_stat | head -4)"
fi

echo ""
print_success "🎉 Cleanup completed! Memory has been freed up."
print_status "You can now restart services with: docker-compose up -d" 