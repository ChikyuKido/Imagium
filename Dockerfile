# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o iamgu .

# Runtime stage
FROM debian:bookworm-slim

WORKDIR /root/


RUN apt-get update && \
    apt-get install -y imagemagick && \
    rm -rf /var/lib/apt/lists/*


RUN ln -s /usr/bin/convert /usr/bin/magick

COPY --from=builder /app/iamgu .
COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./iamgu"]
