# build stage
FROM golang:1.18.1-alpine3.15 AS build

ENV GOOS=linux
ENV GOARCH=arm

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /home_app

#runtime stage
FROM --platform=linux/arm/v6 arm32v6/alpine:3.15

WORKDIR /app

COPY --from=build /home_app .
COPY start.sh .
