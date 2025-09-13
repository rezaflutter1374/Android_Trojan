#!/bin/bash

# API Test Script for C2 Server
# This script tests the API endpoints of the C2 server

# Configuration
SERVER="http://localhost:8080"
API_KEY="test-api-key"

# Colors for output
GREEN="\033[0;32m"
RED="\033[0;31m"
NC="\033[0m" # No Color

# Function to make API requests
function make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    
    echo -e "\n${GREEN}Testing $method $endpoint${NC}"
    echo "Request: $data"
    
    response=$(curl -s -X $method \
        -H "Content-Type: application/json" \
        -H "X-API-Key: $API_KEY" \
        -d "$data" \
        $SERVER$endpoint)
    
    echo "Response: $response"
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Success${NC}"
    else
        echo -e "${RED}✗ Failed${NC}"
    fi
}

# Function to test server status
function test_status() {
    echo -e "\n${GREEN}Testing GET /status${NC}"
    
    response=$(curl -s $SERVER/status)
    echo "Response: $response"
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Success${NC}"
    else
        echo -e "${RED}✗ Failed${NC}"
    fi
}

# Function to test client list endpoint
function test_clients() {
    echo -e "\n${GREEN}Testing GET /api/clients${NC}"
    
    response=$(curl -s -H "X-API-Key: $API_KEY" $SERVER/api/clients)
    echo "Response: $response"
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Success${NC}"
    else
        echo -e "${RED}✗ Failed${NC}"
    fi
}

# Function to test command endpoint
function test_command() {
    local target=$1
    local action=$2
    local payload=$3
    
    data="{\"target_id\":\"$target\",\"action\":\"$action\",\"payload\":\"$payload\"}"
    make_request "POST" "/api/command" "$data"
}

# Main test sequence
echo -e "${GREEN}Starting C2 Server API Tests${NC}"

# Test server status
test_status

# Test client list
test_clients

# Test sending command to all clients
test_command "all" "TAKE_SCREENSHOT" "resolution=high"

# Test sending command to specific client
test_command "client-123" "COLLECT_DATA" "type=contacts"

echo -e "\n${GREEN}API Tests Completed${NC}"