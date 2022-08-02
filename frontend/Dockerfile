# stage build
FROM node:16-alpine

WORKDIR /app

# copy package.json & package-lock.json to the container
COPY package.json package-lock.json ./

# clean install all dependencies
RUN npm ci

# copy everything to the container
COPY . .

# build SvelteKit app
RUN npm run build


# stage run
FROM node:16-alpine

WORKDIR /app

# copy dependency list
COPY package.json package-lock.json ./

# clean install dependencies, no devDependencies, no prepare script
RUN npm ci --omit=dev

# copy built SvelteKit app to /app
COPY --from=0 /app/build ./

EXPOSE 3000
CMD ["node", "./index.js"]
