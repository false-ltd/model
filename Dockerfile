FROM node:lts-alpine3.22 AS builder
#node vsersion is 24
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM node:lts-alpine3.22
#node vsersion is 24
WORKDIR /app
COPY --from=builder /app/.output ./output
ENV HOST=0.0.0.0
ENV PORT=3000
ENV NODE_ENV=production
EXPOSE 3000
CMD ["node", "./output/server/index.mjs"]