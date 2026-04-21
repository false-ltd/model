FROM node:22-alpine AS build
ARG NUXT_PUBLIC_SUPABASE_URL
ARG NUXT_PUBLIC_SUPABASE_KEY
ENV NUXT_PUBLIC_SUPABASE_URL=$NUXT_PUBLIC_SUPABASE_URL
ENV NUXT_PUBLIC_SUPABASE_KEY=$NUXT_PUBLIC_SUPABASE_KEY
RUN corepack enable && corepack prepare pnpm@latest --activate
WORKDIR /app
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY . .
RUN pnpm run generate

FROM nginx:alpine
COPY --from=build /app/.output/public /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
