FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /opt/app

COPY go.mod go.sum ./
RUN go mod download -x

COPY . .

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o raver

FROM alpine:3.20.2

RUN apk add --no-cache yt-dlp ffmpeg

WORKDIR /opt/app

COPY --from=builder /opt/app/raver /opt/app/raver

CMD ["/opt/app/raver"]
