# syntax=docker/dockerfile:1.2

FROM golang:1.17-buster AS base

LABEL maintainer="Shinto C V <cshintov@gmail.com>"

WORKDIR /app

# Creating user shinto (same as host) with home directory
# For avoiding permission issues for files created in container
RUN groupadd -g 1001 shinto
RUN useradd --create-home --shell /bin/bash -u 1001 -g 1001 shinto
RUN usermod -aG sudo shinto
USER shinto

# Change gopath from /go to
ENV GOPATH=/home/shinto/go
ENV PATH="${GOPATH}/bin:${PATH}"

# Install dev tools
RUN go install github.com/pilu/fresh@latest && \
    go install golang.org/x/tools/cmd/godoc@latest

COPY go.* .
RUN --mount=type=cache,target=/home/shinto/go/pkg/mod,uid=1001,gid=1001 \
    go mod download

ENV GO_ENABLED=0

EXPOSE 3000
