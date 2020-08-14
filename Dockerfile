FROM node:12.18.3-alpine3.11 

RUN apk add --no-cache build-base python ffmpeg

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

RUN npm run build

CMD node /app/dist/main.js
