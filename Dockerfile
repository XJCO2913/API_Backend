FROM golang:1.20 AS builder

ENV GOPROXY=https://goproxy.cn

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY ./config/config.yml /root/config/config.yml
CMD ["./main"]