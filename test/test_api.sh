#!/bin/bash
set -e

BASE_URL="http://localhost:8080"
TEST_FILE="test/sphere.obj"

echo "=== CrochetBot API Test ==="
echo ""

# Test health endpoint
echo "1. Testing health endpoint..."
curl -s "${BASE_URL}/health" | jq .
echo ""

# Test file upload
echo "2. Testing file upload..."
UPLOAD_RESPONSE=$(curl -s -X POST \
  -F "file=@${TEST_FILE}" \
  "${BASE_URL}/api/upload")

echo "$UPLOAD_RESPONSE" | jq .

# Extract filename from response
FILENAME=$(echo "$UPLOAD_RESPONSE" | jq -r '.file.filename')
FILE_ID=$(echo "$UPLOAD_RESPONSE" | jq -r '.file.id')

if [ "$FILENAME" = "null" ] || [ "$FILENAME" = "" ]; then
  echo "Error: Upload failed"
  exit 1
fi

echo ""
echo "Uploaded file: $FILENAME"
echo "File ID: $FILE_ID"
echo ""

# Test pattern generation
echo "3. Testing pattern generation..."
GENERATE_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d "{\"filename\": \"$FILENAME\", \"file_id\": \"$FILE_ID\"}" \
  "${BASE_URL}/api/generate")

echo "$GENERATE_RESPONSE" | jq .

# Extract pattern ID
PATTERN_ID=$(echo "$GENERATE_RESPONSE" | jq -r '.pattern.id')

if [ "$PATTERN_ID" = "null" ] || [ "$PATTERN_ID" = "" ]; then
  echo "Error: Pattern generation failed"
  exit 1
fi

echo ""
echo "Generated pattern ID: $PATTERN_ID"
echo ""

# Test pattern retrieval
echo "4. Testing pattern retrieval..."
curl -s "${BASE_URL}/api/pattern/${PATTERN_ID}" | jq .
echo ""

echo "=== All tests passed! ==="
