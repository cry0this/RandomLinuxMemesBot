# Build Stage
FROM golang:alpine AS build

WORKDIR /src/RandomLinuxMemesBot

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /app/RandomLinuxMemesBot

# Final Stage
FROM alpine:latest

# Add curl for health checks
RUN apk add curl

WORKDIR /app
COPY --from=build /app/RandomLinuxMemesBot /app/

EXPOSE 8090

CMD ./RandomLinuxMemesBot
