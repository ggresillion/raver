FROM node:latest
WORKDIR /app
COPY package*.json /app/
RUN npm install
COPY ./ /app/
RUN npm run build -- --prod --output-path=./dist/out
RUN rm -rf /var/www/* && cp -r dist/out /var/www
