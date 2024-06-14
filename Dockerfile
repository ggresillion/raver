FROM golang:1.22.4-alpine as builder

WORKDIR /opt/app

COPY go.mod go.sum ./

RUN go mod download -x

COPY . .

RUN GOOS=linux GOARCH=arm64 go build -o raver

FROM alpine:3.19.1

WORKDIR /opt/app

COPY --from=builder /opt/app/raver /opt/app/raver

CMD ["/opt/app/raver"]
