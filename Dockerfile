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

COPY . /go/src/github.com/american-factory-os/glowplug/
COPY .git/ /go/src/github.com/american-factory-os/glowplug/.git/

WORKDIR /go/src/github.com/american-factory-os/glowplug/

RUN set -Eeux && \
    go mod download && \
    go mod verify

ARG GOOS=linux
ARG GOARCH=amd64
ARG GOCGO_ENABLED=0
ARG VERSION=0.0.0
ARG REVISION=development
    
RUN GOOS=${GOOS} \
    GOARCH=${GOARCH} \
    GOCGO_ENABLED=${GOCGO_ENABLED} \
    echo "building version ${VERSION} revision ${REVISION}" && \
     go build \
     -trimpath \
     -ldflags="-w -s -X 'main.Revision=${REVISION}' -X 'main.Version=${VERSION}'" \
     -o glowplug main.go

RUN go test -cover -v ./...

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
