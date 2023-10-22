# Dockerfile

FROM node:16-alpine

# Copy all local files into the image
COPY frontend/ .

# install all node modules
RUN npm ci

# Build SvelteKit app
RUN npm run build

CMD ["node", "/build/index.js"]
