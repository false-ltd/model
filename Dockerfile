# Stage 1: Build Nuxt frontend
FROM node:22-alpine AS frontend
RUN corepack enable && corepack prepare pnpm@latest --activate
WORKDIR /app
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY . .
RUN pnpm run generate

# Stage 2: Build Go binary with embedded frontend
FROM golang:1.22-alpine AS backend
ARG GIT_SHA=unknown
WORKDIR /app
COPY api/go.mod api/go.sum ./
RUN go mod download
COPY api/ .
# Embed frontend build output into Go source tree
COPY --from=frontend /app/.output/public ./frontend/dist
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=${GIT_SHA}" -o /models-api ./cmd/server

# Stage 3: Minimal runtime
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
COPY --from=backend /models-api /usr/local/bin/models-api
EXPOSE 8080
CMD ["models-api"]
