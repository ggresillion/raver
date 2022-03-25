FROM node:17-alpine as frontend

WORKDIR /app

COPY ./frontend/package.json ./frontend/yarn.lock /app/

RUN yarn install

COPY ./frontend /app/

RUN yarn build

FROM golang:1.18.0-alpine as backend

RUN apk add build-base

WORKDIR /app

COPY ./backend/go.mod ./backend/go.sum /app/

RUN go mod download

COPY ./backend /app/

RUN go build -o discordsoundboard ./cmd/discordsoundboard/main.go

FROM alpine:3.15.2

WORKDIR /app

COPY --from=frontend /app/dist /app/static

COPY --from=backend /app/discordsoundboard /app/discordsoundboard

CMD ./discordsoundboard