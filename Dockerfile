# Container image that runs your code
FROM golang:alpine

# Copies your code file from your action repository to the filesystem path `/` of the container
COPY main.go /main.go
COPY go.mod /go.mod
COPY go.sum /go.sum
COPY entrypoint.sh /entrypoint.sh


ENV TEST_ENV_2='This is a test env 2 variable'
ENV TEST_ENV_3=$2

ENTRYPOINT ["/entrypoint.sh"]