FROM node

RUN apt-get update && apt-get install -y build-essential ffmpeg

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

RUN npm run build

EXPOSE 8000

CMD node /app/dist/main.js
