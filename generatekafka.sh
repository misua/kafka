#!/bin/bash

# This script generates and sends 100 random order POST requests to a specified endpoint.

# The API endpoint URL
URL="http://localhost:8080/order"

# A list of possible statuses for the orders
STATUSES=("pending" "processing" "shipped" "delivered" "cancelled")

echo "Sending 100 POST requests to $URL..."

# Loop 100 times to create and send requests
for i in {1..100}; do
  # Generate a unique order ID using uuidgen
  ORDER_ID="ORD-$(uuidgen | head -c 10 | tr '[:lower:]' '[:upper:]')"
  
  # Generate a random user ID between 100 and 999
  USER_ID="user-$(shuf -i 100-999 -n 1)"

  # Generate a random floating-point number for the total
  # This uses awk to generate a number with two decimal places
  TOTAL=$(awk "BEGIN{printf \"%.2f\", (rand() * 490) + 10}")

  # Select a random status from the array
  RANDOM_STATUS_INDEX=$(( RANDOM % ${#STATUSES[@]} ))
  STATUS=${STATUSES[$RANDOM_STATUS_INDEX]}

  # Generate a recent timestamp in ISO 8601 format
  TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

  # Construct the JSON payload
  PAYLOAD="{
    \"id\": \"$ORDER_ID\",
    \"user_id\": \"$USER_ID\",
    \"total\": $TOTAL,
    \"status\": \"$STATUS\",
    \"timestamp\": \"$TIMESTAMP\"
  }"

  # Send the POST request using curl and capture the status code
  HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST \
    -H 'Content-Type: application/json' \
    -d "$PAYLOAD" "$URL")

  echo "Request $i/100 - Status Code: $HTTP_STATUS"

  # Optional: Print the payload for verification
  # echo "Payload: $PAYLOAD"

done

echo ""
echo "Script finished."