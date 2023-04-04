#!/bin/sh -l

echo 'Running entrypoint'

MATTERMOST_ACCESS_TOKEN=${1} \
COMMIT_URL=${2} \
COMMIT_AUTHOR_USERNAME=${3} \
COMMIT_AUTHOR_EMAIL=${4} \
COMMIT_MESSAGE=${5} \
STATUS_CONCLUSION=${6} \
STATUS_URL=${7} \
STATUS_NAME=${8} \
STATUS_DESCRIPTION=${9} \
go run /main.go

echo 'Running entrypoint done'