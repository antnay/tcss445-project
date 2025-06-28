# FROM node:20-alpine AS development-dependencies-env
# WORKDIR /app
# COPY client/ ./
# RUN npm ci

# FROM node:20-alpine AS production-dependencies-env
# COPY client/package.json client/package-lock.json client/vite.config.ts /app/
# WORKDIR /app
# COPY client/ ./
# RUN npm ci --omit=dev

# FROM node:20-alpine AS build-env
# WORKDIR /app
# COPY client/ ./
# COPY --from=development-dependencies-env /app/node_modules /app/node_modules
# RUN npm run build

# FROM node:20-alpine
# COPY client/package.json client/package-lock.json /app/
# COPY --from=production-dependencies-env /app/node_modules client/vite.config.ts /app/node_modules
# COPY --from=build-env /app/build /app/build
# WORKDIR /app

# CMD ["npm", "run", "start"]

FROM node:20-alpine

WORKDIR /react-app

COPY client/package.json client/package-lock.json client/vite.config.ts ./

RUN npm install

COPY client ./

EXPOSE 5173

CMD ["npm", "run", "dev"]
