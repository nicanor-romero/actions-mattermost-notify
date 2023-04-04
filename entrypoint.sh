#!/bin/sh -l

echo 'Running entrypoint'

MATTERMOST_ACCESS_TOKEN='$1' \
REPO_URL='$2' \
COMMIT_SHA='$3' \
CHECK_RUN_NAME='$4' \
CHECK_RUN_CONCLUSION='$5' \
CHECK_RUN_URL='$6' \
CHECK_RUN_OUTPUT_TITLE='$7' \
CHECK_RUN_OUTPUT_TEXT='$8' \
CHECK_RUN_OUTPUT_SUMMARY='$9' \
go run /main.go

echo 'Running entrypoint done'