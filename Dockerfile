FROM node:12.18.3-alpine3.11 as build

RUN apk add --no-cache build-base python

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

RUN npm run build

FROM node:12.18.3-alpine3.11

RUN apk add --no-cache ffmpeg

WORKDIR /app

COPY --from=build /app/node_modules /app/node_modules

COPY --from=build /app/dist /app

RUN ls /app > test

CMD cat test