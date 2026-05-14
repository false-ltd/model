# Stage 1: Build Nuxt frontend (native amd64, no QEMU)
FROM --platform=$BUILDPLATFORM node:22-alpine AS frontend
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable && corepack prepare pnpm@10 --activate
WORKDIR /app
COPY package.json pnpm-lock.yaml ./
RUN --mount=type=cache,target=/root/.local/share/pnpm/store pnpm install --frozen-lockfile
COPY . .
RUN pnpm run generate

# Stage 2: Build Go binary (native amd64, cross-compile to arm64)
FROM --platform=$BUILDPLATFORM golang:alpine AS backend
ARG GIT_SHA=unknown
ARG TARGETARCH
WORKDIR /app
COPY api/go.mod api/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY api/ .
COPY --from=frontend /app/.output/public ./frontend/dist
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w -X main.version=${GIT_SHA}" -o /model ./cmd/server

# Stage 3: Scratch runtime (no RUN = no QEMU)
FROM scratch
ENV TZ=Asia/Shanghai
COPY --from=backend /model /model
COPY --from=backend /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
ENTRYPOINT ["/model"]
