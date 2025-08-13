FROM golang:1.23.0-alpine AS builder

WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o s3 ./cmd

FROM alpine:latest
WORKDIR /appbuild
COPY --from=builder /app/s3 .
COPY images ./images
COPY .env .
CMD ["./s3"]