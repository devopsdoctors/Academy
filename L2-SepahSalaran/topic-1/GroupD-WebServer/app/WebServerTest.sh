#!/bin/bash
# Comprehensive test script for the Go web server
# Usage: ./test_server.sh [host] [port]

# Default values
HOST=${1:-"localhost"}
PORT=${2:-8080}
BASE_URL="http://$HOST:$PORT"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counter
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

# Function to run tests
run_test() {
  local name=$1
  local command=$2
  local expected_status=$3
  local expected_content=$4
  
  echo -e "${YELLOW}Running test: $name${NC}"
  echo -e "${BLUE}Command:${NC} $command"
  
  # Run the command and capture output and status
  output=$(eval $command 2>&1)
  status=$?
  
  # Get HTTP status code if present in output
  http_status=$(echo "$output" | grep -oP "HTTP/\d\.\d \K\d+")
  
  # Increment test counter
  TESTS_RUN=$((TESTS_RUN + 1))
  
  # Check if test passed
  status_check=true
  if [[ -n "$expected_status" && "$http_status" != "$expected_status" ]]; then
    status_check=false
  fi
  
  # Check content if provided
  content_check=true
  if [[ -n "$expected_content" && ! "$output" =~ $expected_content ]]; then
    content_check=false
  fi
  
  if [[ $status -eq 0 && "$status_check" == true && "$content_check" == true ]]; then
    echo -e "${GREEN}✓ Test passed${NC}"
    TESTS_PASSED=$((TESTS_PASSED + 1))
  else
    echo -e "${RED}✗ Test failed${NC}"
    if [[ "$status_check" == false ]]; then
      echo -e "${RED}Expected status: $expected_status, Got: $http_status${NC}"
    fi
    if [[ "$content_check" == false ]]; then
      echo -e "${RED}Expected content not found${NC}"
    fi
    echo -e "Output: $output"
    TESTS_FAILED=$((TESTS_FAILED + 1))
  fi
  echo ""
}

echo "Starting tests for web server at $BASE_URL"
echo "==========================================="

# Basic functionality tests
echo -e "${BLUE}Basic Functionality Tests${NC}"
echo "-------------------------------------------"

# Test 1: Root endpoint
run_test "Root endpoint" "curl -s -i $BASE_URL/" "200" "Welcome to the Go Web Server"

# Test 2: Echo endpoint
run_test "Echo endpoint" "curl -s -i $BASE_URL/echo/hello-world" "200" "hello-world"

# Test 3: User-agent endpoint
run_test "User-agent endpoint" "curl -s -i $BASE_URL/user-agent -H 'User-Agent: CurlTestScript'" "200" "CurlTestScript"

# Test 4: API status endpoint
run_test "API status endpoint" "curl -s -i $BASE_URL/api/status" "200" "\"status\":\"ok\""

# Test 5: API time endpoint
run_test "API time endpoint" "curl -s -i $BASE_URL/api/time" "200" "\"time\":"

# Test 6: API echo endpoint
run_test "API echo endpoint" "curl -s -i -X POST $BASE_URL/api/echo -d '{\"test\":\"data\"}' -H 'Content-Type: application/json'" "200" "\"test\":\"data\""

# File operations tests
echo -e "${BLUE}File Operations Tests${NC}"
echo "-------------------------------------------"

# Test 7: Create a test file
run_test "Create file" "curl -s -i -X POST $BASE_URL/files/test.txt -d 'This is a test file'" "201" "File created"

# Test 8: Get the created file
run_test "Get file" "curl -s -i $BASE_URL/files/test.txt" "200" "This is a test file"

# Test 9: Delete the test file
run_test "Delete file" "curl -s -i -X DELETE $BASE_URL/files/test.txt" "200" "File deleted"

# Test 10: Try to get non-existent file
run_test "Get non-existent file" "curl -s -i $BASE_URL/files/nonexistent.txt" "404" "File not found"

# Security tests
echo -e "${BLUE}Security Tests${NC}"
echo "-------------------------------------------"

# Test 11: Path traversal attempt
#run_test "Path traversal attempt" "curl -s -i $BASE_URL/files/../../../etc/passwd" "403" "Path traversal not allowed"

# Test 12: Another path traversal variant
run_test "Path traversal variant" "curl -s -i $BASE_URL/files/%2e%2e/%2e%2e/etc/passwd" "403" "Path traversal not allowed"

# Session tests
echo -e "${BLUE}Session Tests${NC}"
echo "-------------------------------------------"

# Test 13: Test API session endpoint
run_test "API session endpoint" "curl -s -i $BASE_URL/api/session -c cookies.txt" "200" "\"session_id\":"

# Test 14: Test session persistence
run_test "Session persistence" "curl -s -i $BASE_URL/api/session -b cookies.txt" "200" "\"session_id\":"

# Performance and feature tests
echo -e "${BLUE}Performance and Feature Tests${NC}"
echo "-------------------------------------------"

# Test 15: Gzip encoding
run_test "Gzip encoding" "curl -s -i $BASE_URL/ --compressed -H 'Accept-Encoding: gzip'" "200" "Content-Encoding: gzip"

# Test 16: Directory listing
run_test "Directory listing" "curl -s -i $BASE_URL/files/" "200" "Directory Listing"

# Test 17: Method not allowed
#run_test "Method not allowed" "curl -s -i -X PUT $BASE_URL/user-agent" "404" "Not Found"

# Test 18: Large request body
run_test "Large request body" "dd if=/dev/zero bs=1024 count=100 2>/dev/null | curl -s -i -X POST $BASE_URL/files/large.bin --data-binary @-" "201" "File created"

# Test 19: Clean up large file
run_test "Delete large file" "curl -s -i -X DELETE $BASE_URL/files/large.bin" "200" "File deleted"

# Test 20: Security headers
run_test "Security headers" "curl -s -i $BASE_URL/" "200" "X-Content-Type-Options: nosniff"

# Test 21: Multiple concurrent requests
echo -e "${YELLOW}Running multiple concurrent requests...${NC}"
for i in {1..10}; do
  curl -s $BASE_URL/ &>/dev/null &
done
wait
echo -e "${GREEN}Concurrent requests completed${NC}"
echo ""

# Test 22: Very long URL
long_url=$(printf "%0.s$" {1..500})
run_test "Very long URL" "curl -s -i \"$BASE_URL/echo/$long_url\"" "200"

# Test 23: Long header
run_test "Long header" "curl -s -i $BASE_URL/ -H \"X-Custom-Header: $(printf '%0.s$' {1..500})\"" "200" "Welcome to the Go Web Server"

# Test 24: Non-existent path
run_test "Non-existent path" "curl -s -i $BASE_URL/notfound" "404" "Not Found"

# Test 25: Verify files endpoint methods
run_test "PUT method not allowed" "curl -s -i -X PUT $BASE_URL/files/test.txt -d 'content'" "405" "Method not allowed"

# Summary
echo "==========================================="
echo "Test Summary:"
echo "Tests run: $TESTS_RUN"
echo -e "${GREEN}Tests passed: $TESTS_PASSED${NC}"
echo -e "${RED}Tests failed: $TESTS_FAILED${NC}"

# Clean up cookies file
rm -f cookies.txt

# Exit with error code if any test failed
if [ $TESTS_FAILED -gt 0 ]; then
  exit 1
fi

exit 0
