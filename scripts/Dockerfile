# syntax=docker/dockerfile:1
# Build the application from source
FROM golang:1.20-alpine3.18 AS build-stage
RUN go version
ENV CGO_ENABLED 0

COPY seed.go .


