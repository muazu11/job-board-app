FROM node:23-alpine3.20 AS builder 

WORKDIR /app

COPY package*.json .

# install all node modules
RUN npm ci

# Copy all local files into the image
COPY . .

# Build SvelteKit app
RUN npm run build

FROM node:23-alpine3.20

WORKDIR /app

COPY --from=builder --chown=node:node /app/build ./build
COPY --from=builder --chown=node:node /app/package.json ./
COPY --from=builder --chown=node:node /app/package-lock.json ./

RUN npm ci --omit=dev

EXPOSE 80

USER node

CMD ["node", "build/index.js"]
