# Container image that runs your code
FROM golang:alpine

# Copies your code file from your action repository to the filesystem path `/` of the container
COPY main.go /main.go
COPY go.mod /go.mod
COPY go.sum /go.sum

ENV MATTERMOST_PERSONAL_ACCESS_TOKEN=$1
ENV COMMIT_URL=$2
ENV COMMIT_AUTHOR_USERNAME=$3
ENV COMMIT_AUTHOR_EMAIL=$4
ENV COMMIT_MESSAGE=$5
ENV TEST_JOB_OUTPUT=$6

RUN go run /main.go