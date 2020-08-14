FROM node:12.18.3-alpine3.11 as build
WORKDIR /app
COPY package*.json /app/
RUN npm install
COPY . /app
ARG configuration=production
RUN npm run build -- --outputPath=./dist/out --configuration $configuration

FROM nginx:alpine
COPY --from=build /app/dist/out/ /usr/share/nginx/html
COPY /nginx.conf /etc/nginx/conf.d/default.conf
