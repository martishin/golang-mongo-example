# Build stage
FROM golang:1.23-bullseye AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
RUN CGO_ENABLED=0 go build -o mongodb-example

# Run stage
FROM alpine:3.20

RUN apk add --no-cache bash coreutils

WORKDIR /app
COPY wait-for-it.sh wait-for-it.sh
COPY --from=builder /app/mongodb-example /app/mongodb-example

CMD ./wait-for-it.sh -s -t 30 $MONGODB_SERVER -- ./mongodb-example
