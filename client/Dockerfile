FROM node:23-alpine

WORKDIR /svelte-app

COPY client/package.json client/package-lock.json client/tailwind.config.cjs ./

RUN npm ci

COPY client/ ./
EXPOSE 5173

CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
