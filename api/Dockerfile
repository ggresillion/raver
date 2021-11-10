FROM node:alpine as builder

RUN apk update && apk add ffmpeg make gcc g++ libtool autoconf automake python3 && rm -rf /var/cache/apk/*

WORKDIR /app

COPY package.json yarn.lock ./

RUN yarn install

COPY . .

RUN yarn build

RUN yarn install --production

FROM node:alpine

RUN apk update && apk add ffmpeg && rm -rf /var/cache/apk/*

WORKDIR /app

COPY --from=builder /app/node_modules ./node_modules

COPY --from=builder /app/dist ./dist

EXPOSE 8000

CMD ["node", "./dist/main.js"]
