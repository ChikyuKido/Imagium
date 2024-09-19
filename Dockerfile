# Build stage
FROM golang:1.23.1-alpine AS builder

WORKDIR /app
ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64
RUN apk add --no-cache upx

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" imagu

RUN upx --best --lzma imagu

# Runtime stage
FROM alpine:latest

WORKDIR /app/


RUN apk add --no-cache imagemagick

COPY --from=builder /app/imagu .
COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./imagu"]
