# Stage 1: build
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/server ./cmd/server

# Stage 2: final image
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 50051

ENTRYPOINT ["/app/server"]
