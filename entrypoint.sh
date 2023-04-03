#!/bin/sh -l

echo "Starting entrypoint"
echo "1: $1"
echo "2: $2"
echo "3: $3"
echo "4: $4"
echo "5: $5"
echo "6: $6"

MATTERMOST_PERSONAL_ACCESS_TOKEN="$1" \
COMMIT_URL="$2" \
COMMIT_AUTHOR_USERNAME="$3" \
COMMIT_AUTHOR_EMAIL="$4" \
COMMIT_MESSAGE="$5" \
TEST_JOB_OUTPUT="\"$6\"" \
go run /main.go
