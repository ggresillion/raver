FROM node:latest
WORKDIR /app
COPY package*.json /app/
RUN npm install
COPY ./ /app/
RUN npm run build -- --output-path=./dist/out
RUN cp -r dist/out /var/www
