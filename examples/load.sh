#!/bin/bash

# Load script to generate traffic on demo server
# Usage: ./load.sh [server_url]

SERVER_URL=${1:-http://localhost:6060}

echo "Generating load on $SERVER_URL..."

# Hit CPU hotspot endpoint
curl -s "$SERVER_URL/cpu-hotspot" > /dev/null &

# Hit allocation heavy endpoint  
curl -s "$SERVER_URL/alloc-heavy" > /dev/null &

# Hit mutex contention endpoint
curl -s "$SERVER_URL/mutex-contention" > /dev/null &

# Wait for all requests to complete
wait

echo "Load generation completed."