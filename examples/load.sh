#!/bin/bash

# Enhanced Load Generation Script for TriageProf Demo
# Generates realistic traffic patterns to demonstrate performance issues

SERVER_URL=${1:-http://localhost:6060}
DURATION=${2:-60}  # Duration in seconds (increased from 30 to 60)
CONCURRENCY=${3:-15} # Concurrent requests (increased from 10 to 15)

echo "🚀 Starting enhanced load generation on $SERVER_URL"
echo "⏱  Duration: ${DURATION}s | 👥 Concurrency: ${CONCURRENCY}"
echo "📊 Generating realistic traffic patterns..."
echo ""

END_TIME=$((SECONDS + DURATION))
REQUEST_COUNT_FILE=$(mktemp)
echo "0" > "$REQUEST_COUNT_FILE"

# Function to make requests with random endpoints
make_request() {
    while [ $SECONDS -lt $END_TIME ]; do
        # Randomly select an endpoint
        ENDPOINT_NUM=$((RANDOM % 11))
        
        case $ENDPOINT_NUM in
            0) curl -s "$SERVER_URL/api/users" > /dev/null ;;
            1) curl -s "$SERVER_URL/api/search?q=test" > /dev/null ;;
            2) curl -s "$SERVER_URL/api/analytics" > /dev/null ;;
            3) curl -s "$SERVER_URL/api/export" > /dev/null ;;
            4) curl -s -X POST "$SERVER_URL/api/process" -H "Content-Type: application/json" -d '{"data":"test"}' > /dev/null ;;
            5) curl -s "$SERVER_URL/api/strings" > /dev/null ;;
            6) curl -s "$SERVER_URL/api/nocache" > /dev/null ;;
            7) curl -s "$SERVER_URL/api/iobound" > /dev/null ;;
            8) curl -s "$SERVER_URL/api/leak" > /dev/null ;;
            9) curl -s "$SERVER_URL/api/blocking" > /dev/null ;;
            10) curl -s "$SERVER_URL/api/goroutine" > /dev/null ;;
        esac
        
        # Increment request count using file
        CURRENT_COUNT=$(cat "$REQUEST_COUNT_FILE")
        echo $((CURRENT_COUNT + 1)) > "$REQUEST_COUNT_FILE"
        
        # Random delay between requests (20-300ms for faster pace)
        sleep $(echo "scale=3; $RANDOM / 32768 * 0.28 + 0.02" | bc)
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
    CURRENT_COUNT=$(cat "$REQUEST_COUNT_FILE")
    echo -ne "📊 Progress: $((100 - (REMAINING * 100 / DURATION)))% | Requests: $CURRENT_COUNT | Remaining: ${REMAINING}s\r"
    sleep 1
done

CURRENT_COUNT=$(cat "$REQUEST_COUNT_FILE")
echo -e "\n✅ Load generation completed!"
echo "📈 Total requests generated: $CURRENT_COUNT"
echo "🎯 Load profile: Mixed read/write operations with realistic timing"

# Clean up temp file
rm -f "$REQUEST_COUNT_FILE"
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
echo "  • Memory leaks in /api/leak"
echo "  • Blocking I/O operations in /api/blocking"
echo "  • Goroutine leaks in /api/goroutine"
echo ""
echo "🚀 Ready for TriageProf analysis!"
