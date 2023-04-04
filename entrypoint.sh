#!/bin/sh -l

echo 'Running entrypoint'

GITHUB_ACCESS_TOKEN=${1} \
MATTERMOST_ACCESS_TOKEN=${2} \
REPO_URL=${3} \
COMMIT_SHA=${4} \
CHECK_RUN_NAME="${5}" \
CHECK_RUN_CONCLUSION="${6}" \
CHECK_RUN_URL=${7} \
CHECK_RUN_OUTPUT_TITLE="${8}" \
CHECK_RUN_OUTPUT_TEXT="${9}" \
CHECK_RUN_OUTPUT_SUMMARY="${10}" \
go run /main.go

echo 'Running entrypoint done'