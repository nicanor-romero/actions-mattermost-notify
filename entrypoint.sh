#!/bin/sh -l

echo 'Running entrypoint'

ls -l
ls -l /
ls -l /usr/local/go/bin

GITHUB_ACCESS_TOKEN=${1} \
MATTERMOST_ACCESS_TOKEN=${2} \
MATTERMOST_INCOMING_WEBHOOK_URL=${3} \
MATTERMOST_CHANNEL_NAME=${4} \
COMMIT_URL=${5} \
COMMIT_AUTHOR_USERNAME=${6} \
COMMIT_AUTHOR_EMAIL=${7} \
COMMIT_MESSAGE=${8} \
STATUS_CONCLUSION=${9} \
STATUS_URL=${10} \
STATUS_NAME=${11} \
STATUS_DESCRIPTION=${12} \
/usr/local/go/bin/go run /main.go

echo 'Running entrypoint done'
