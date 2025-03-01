#!/bin/bash

# Configuration
DURATION="5s"
THREADS=4
CONNECTIONS=100

# Array of routes and methods
ROUTES=(
  "POST http://localhost:8000/api/v1/polls"
  "GET http://localhost:8000/api/v1/polls"
  "POST http://localhost:8000/api/v1/polls/1/vote"
  "POST http://localhost:8000/api/v1/polls/1/skip"
)

# Payloads
PAYLOAD_POST='{"title": "title1", "options": ["op1","op2"], "tags":["tag1","tag2"]}'
PAYLOAD_VOTE_TEMPLATE='{"user_id": %d, "option_index":1}'
PAYLOAD_SKIP_TEMPLATE='{"user_id": %d}'
QUERY_PARAMS="?page=1&page_size=10&user_id=1"

IDS=()
USER_ID=1

# Function to select payload based on URL
select_payload() {
  case "$1" in
    *"vote"*) printf "$PAYLOAD_VOTE_TEMPLATE" $USER_ID ;;
    *"skip"*) printf "$PAYLOAD_SKIP_TEMPLATE" $USER_ID ;;
    *"polls"*) echo "$PAYLOAD_POST" ;;
    *) echo '{}' ;;
  esac
}

# Function to extract IDs from GET response body
extract_ids() {
  local response=$1
  IDS=($(echo "$response" | grep -o '"id":[0-9]*' | awk -F: '{print $2}'))
}

# Function to run stress test on each route
run_test() {
  local method=$1
  local url=$2
  local lua_script="request.lua"

  if [ "$method" == "GET" ]; then
    url="$url$QUERY_PARAMS"
    response=$(wrk -t1 -c1 -d1s "$url" | grep -o '{.*}')
    extract_ids "$response"
  fi

  if [[ "$url" == *"polls/1/vote"* || "$url" == *"polls/1/skip"* ]]; then
    for id in "${IDS[@]}"; do
      local dynamic_url=$(echo "$url" | sed "s/1/$id/")
      local payload=$(select_payload "$dynamic_url")
      run_wrk "$method" "$dynamic_url" "$payload"
      ((USER_ID++))
    done
    return
  fi

  local payload=$(select_payload "$url")
  run_wrk "$method" "$url" "$payload"
  ((USER_ID++))
}

run_wrk() {
  local method=$1
  local url=$2
  local payload=$3

  # Create Lua script dynamically
  cat <<EOF > request.lua
wrk.method = "$method"
wrk.headers["Content-Type"] = "application/json"
wrk.body = '$payload'
response = function(status, headers, body)
  if status ~= 200 and status ~= 201 then
    print("Unexpected status code: " .. status .. " for URL: " .. "$url")
  end
end
EOF

  echo "Running stress test on $url with $THREADS threads, $CONNECTIONS connections for $DURATION"
  wrk -t$THREADS -c$CONNECTIONS -d$DURATION -s request.lua "$url"

  # Clean up Lua script
  rm request.lua
  echo "Test completed for $url"
}

# Loop through routes
for route in "${ROUTES[@]}"; do
  method=$(echo $route | awk '{print $1}')
  url=$(echo $route | awk '{print $2}')
  run_test "$method" "$url"
done

echo "All stress tests completed"

