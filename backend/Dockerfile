FROM golang:1.18.0-alpine as builder

RUN apk add build-base

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

COPY . /app/

RUN go build -o discordsoundboard ./cmd/discordsoundboard/main.go

FROM alpine:3.15.2

WORKDIR /app

RUN apk add ffmpeg

COPY --from=builder /app/discordsoundboard /app/discordsoundboard

CMD ./discordsoundboard
