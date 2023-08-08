# syntax=docker/dockerfile:1
# Build the application from source
FROM golang:1.20-alpine3.18 AS build-stage
RUN go version
ENV CGO_ENABLED 0

WORKDIR /app

COPY . .

RUN go mod download


RUN  go build  -o /api


# Deploy the application binary into a lean image
FROM alpine:3.18 AS build-release-stage

WORKDIR /
COPY --from=build-stage /api /api
COPY .env.dev .
COPY .env.compose .
EXPOSE 3000