# syntax=docker/dockerfile:1

# We use a multi-stage build setup.
# (https://docs.docker.com/build/building/multi-stage/)

###############################################################################
# Stage 1 (to create a "build" image, ~850MB)                                 #
###############################################################################
# Image from https://hub.docker.com/_/golang
FROM golang:1.22.3 AS builder
# smoke test to verify if golang is available
RUN go version

ARG PROJECT_VERSION

COPY . /go/src/github.com/american-factory-os/glowplug/
WORKDIR /go/src/github.com/american-factory-os/glowplug/
RUN set -Eeux && \
    go mod download && \
    go mod verify

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build \
    -trimpath \
    -ldflags="-w -s -X 'main.Revision=$(git rev-parse --short HEAD)'" -X 'main.Version=$(git describe --abbrev=0 --tags)'" \
    -o glowplug main.go
# RUN go test -cover -v ./...

###############################################################################
# Stage 2 (to create a downsized "container executable", ~5MB)                #
###############################################################################

FROM golang:alpine

RUN apk --no-cache add ca-certificates

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/github.com/american-factory-os/glowplug/glowplug .

# EXPOSE 80
ENTRYPOINT [ "/app/glowplug" ]
