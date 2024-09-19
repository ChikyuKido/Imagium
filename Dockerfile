# Build stage
FROM golang:1.23.1-alpine AS builder

WORKDIR /app
ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64
RUN apk add --no-cache upx gcc musl-dev

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" imagu

RUN upx --best --lzma imagu

# Runtime stage
FROM debian:bookworm-slim

WORKDIR /app/


RUN apt-get update && \
    apt-get install -y imagemagick && \
    rm -rf /var/lib/apt/lists/*


RUN ln -s /usr/bin/convert /usr/bin/magick

COPY --from=builder /app/imagu .
COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./imagu"]