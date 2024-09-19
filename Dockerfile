# Build stage
FROM golang:1.23.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o imagu .

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