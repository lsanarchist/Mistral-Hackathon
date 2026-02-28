#!/bin/bash

# Enhanced Load Generation Script for TriageProf Demo
# Generates realistic traffic patterns to demonstrate performance issues

SERVER_URL=${1:-http://localhost:6060}
DURATION=${2:-30}  # Duration in seconds
CONCURRENCY=${3:-10} # Concurrent requests

echo "🚀 Starting enhanced load generation on $SERVER_URL"
echo "⏱  Duration: ${DURATION}s | 👥 Concurrency: ${CONCURRENCY}"
echo "📊 Generating realistic traffic patterns..."
echo ""

END_TIME=$((SECONDS + DURATION))
REQUEST_COUNT=0

# Function to make requests with random endpoints
make_request() {
    while [ $SECONDS -lt $END_TIME ]; do
        # Randomly select an endpoint
        ENDPOINT_NUM=$((RANDOM % 6))
        
        case $ENDPOINT_NUM in
            0) curl -s "$SERVER_URL/api/users" > /dev/null ;;
            1) curl -s "$SERVER_URL/api/search?q=test" > /dev/null ;;
            2) curl -s "$SERVER_URL/api/analytics" > /dev/null ;;
            3) curl -s "$SERVER_URL/api/export" > /dev/null ;;
            4) curl -s -X POST "$SERVER_URL/api/process" -H "Content-Type: application/json" -d '{"data":"test"}' > /dev/null ;;
            5) curl -s "$SERVER_URL/api/strings" > /dev/null ;;
            6) curl -s "$SERVER_URL/api/nocache" > /dev/null ;;
            7) curl -s "$SERVER_URL/api/iobound" > /dev/null ;;
        esac
        
        REQUEST_COUNT=$((REQUEST_COUNT + 1))
        
        # Random delay between requests (50-500ms)
        sleep $(echo "scale=3; $RANDOM / 32768 * 0.45 + 0.05" | bc)
    done
}

# Start concurrent workers
echo "🔥 Starting $CONCURRENCY concurrent workers..."
for i in $(seq 1 $CONCURRENCY); do
    make_request &
done

# Show progress
while [ $SECONDS -lt $END_TIME ]; do
    REMAINING=$((END_TIME - SECONDS))
    echo -ne "📊 Progress: $((100 - (REMAINING * 100 / DURATION)))% | Requests: $REQUEST_COUNT | Remaining: ${REMAINING}s\r"
    sleep 1
done

echo -e "\n✅ Load generation completed!"
echo "📈 Total requests generated: $REQUEST_COUNT"
echo "🎯 Load profile: Mixed read/write operations with realistic timing"
echo ""
echo "💡 Performance issues that should be visible:"
echo "  • JSON serialization overhead in /api/users"
echo "  • Database lock contention in /api/search"
echo "  • CPU-bound processing in /api/analytics"
echo "  • Memory allocation patterns in /api/export"
echo "  • Mutex contention in /api/process"
echo "  • Inefficient string operations in /api/strings"
echo "  • Lack of caching in /api/nocache"
echo "  • I/O bottlenecks in /api/iobound"
echo ""
echo "🚀 Ready for TriageProf analysis!"
