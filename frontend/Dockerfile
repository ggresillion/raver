# stage build
FROM node:16-alpine

WORKDIR /app

# copy package.json & yarn.lock to the container
COPY package.json yarn.lock ./

# clean install all dependencies
RUN yarn install --frozen-lockfile

# copy everything to the container
COPY . .

# build SvelteKit app
RUN yarn build


# stage run
FROM node:16-alpine

WORKDIR /app

# copy dependency list
COPY package.json yarn.lock ./

# clean install dependencies, no devDependencies, no prepare script
RUN yarn install --frozen-lockfile --production --ignore-scripts

# copy built SvelteKit app to /app
COPY --from=0 /app/build ./

EXPOSE 3000
CMD ["node", "./index.js"]
