FROM node:latest

RUN apt-get update && apt-get install ffmpeg -y

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

RUN npm run build

CMD node dist/main.js
