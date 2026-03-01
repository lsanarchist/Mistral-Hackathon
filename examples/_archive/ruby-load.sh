#!/bin/bash

# Load script for Ruby demo server
# Usage: ./ruby-load.sh [server_url]

SERVER_URL=${1:-http://localhost:4567}

echo "Generating load on Ruby demo server at $SERVER_URL..."

# Hit CPU intensive endpoint
curl -s "$SERVER_URL/cpu-intensive" > /dev/null &

# Hit memory heavy endpoint  
curl -s "$SERVER_URL/memory-heavy" > /dev/null &

# Hit object creation endpoint
curl -s "$SERVER_URL/object-creation" > /dev/null &

# Hit JSON processing endpoint
curl -s "$SERVER_URL/json-processing" > /dev/null &

# Hit database queries endpoint
curl -s "$SERVER_URL/database-queries" > /dev/null &

# Wait for all requests to complete
wait

echo "Load generation completed for Ruby demo server."