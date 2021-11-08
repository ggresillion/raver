FROM node:alpine

RUN apk update && apk add ffmpeg make gcc g++ libtool autoconf automake python3

WORKDIR /app

COPY package*.json ./

RUN yarn install

COPY . .

RUN yarn build

EXPOSE 8000

CMD node /app/dist/main.js
