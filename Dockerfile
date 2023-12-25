# Build Stage
FROM golang:1.21.5 AS builder
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go

# Run Stage
FROM zc1185230223/alpine:3.18

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 9090
CMD ["/app/main", "-debug=false"]