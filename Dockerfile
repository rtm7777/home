# build stage
FROM --platform=linux/arm/v6 golang:1.17-alpine AS build

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=arm go build -o /home_app

#runtime stage
FROM --platform=linux/arm/v6 arm32v6/alpine:3.15

WORKDIR /app

COPY --from=build /home_app .
COPY start.sh .
