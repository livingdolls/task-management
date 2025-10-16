# Stage 1: Build
FROM golang:1.24 AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api/runner.go

# Stage 2: Runtime
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache bash

# Copy binary
COPY --from=builder /app/server .

# Copy config folder
COPY --from=builder /app/config ./config

COPY --from=builder /app/scripts/wait-for-it.sh ./wait-for-it.sh
RUN chmod +x ./wait-for-it.sh

EXPOSE 3000

CMD ["bash","./wait-for-it.sh", "mysql_db:3306", "--timeout=30", "--", "./server"]