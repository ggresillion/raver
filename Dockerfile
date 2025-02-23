FROM golang:1.23-alpine AS builder

WORKDIR /opt/app

COPY go.mod go.sum ./

RUN go mod download -x

COPY . .

RUN go build -o raver

FROM alpine:3.20.2

RUN apk add python3 py3-pip yt-dlp

WORKDIR /opt/app

COPY --from=builder /opt/app/raver /opt/app/raver

CMD ["/opt/app/raver"]
