FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/url-shortener .
EXPOSE 8080
CMD ["./url-shortener"]
