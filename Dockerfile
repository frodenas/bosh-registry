FROM golang
MAINTAINER Ferran Rodenas <frodenas@gmail.com>

# Set environment variables
ENV CGO_ENABLED 0
ENV GOARCH      amd64
ENV GOARM       5
ENV GOOS        linux

# Build BOSH Registry
RUN go get -a -installsuffix cgo -ldflags '-s' github.com/frodenas/bosh-registry/main

# Add files
ADD Dockerfile.final /go/bin/Dockerfile
ADD bosh-registry.json /go/bin/bosh-registry.json

# Command to run
CMD docker build -t frodenas/bosh-registry /go/bin
