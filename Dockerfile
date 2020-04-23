FROM node:latest as frontend

WORKDIR /app

COPY client/package*.json ./

RUN npm install

COPY client/ .

RUN npm run build -- --prod

FROM node:latest

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

RUN npm run build

COPY --from=frontend /app/dist/DiscordSoundBoard /app/dist/client

RUN apt-get update && apt-get install ffmpeg -y

CMD node dist/main.js