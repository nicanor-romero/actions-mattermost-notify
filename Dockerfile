# Container image that runs your code
FROM golang:alpine

# Copies your code file from your action repository to the filesystem path `/` of the container
COPY main.go /main.go
COPY go.mod /go.mod
COPY go.sum /go.sum
COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]