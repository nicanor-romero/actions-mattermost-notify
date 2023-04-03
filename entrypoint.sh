#!/bin/sh -l

echo "Starting"
echo "1: $1"
echo "2: $2"
echo "3: $3"
echo "4: $4"
echo "5: $5"
echo "6: $6"
echo "7: $7"

setenv MATTERMOST_PERSONAL_ACCESS_TOKEN: "$1"
setenv COMMIT_URL: "$2"
setenv COMMIT_AUTHOR_USERNAME: "$3"
setenv COMMIT_AUTHOR_EMAIL: "$4"
setenv COMMIT_MESSAGE: $5
setenv TEST_JOB_OUTPUT: "$6"
setenv TEST_ENV: "$7"

go run /main.go
