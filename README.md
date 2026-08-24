# Model — AI Model Atlas

Browse, compare, and visualize pricing and capability data for LLMs. Data synced from [models.dev](https://models.dev), served by a self-contained Go binary with an embedded frontend.

线上地址: [model.false.ltd](https://model.false.ltd)

## Features

- **Overview** — Hero dashboard with animated counters, price tier distribution, modality coverage, capability stats, and context window charts
- **Catalog** — Searchable, sortable table; filters for provider, I/O modalities, price ranges, and capabilities (free, reasoning, vision, open weights, …); column visibility persisted per device
- **Model Detail** — Pricing gauges, limits, capabilities, modalities, timeline, integration info, and Quick Start code (AI SDK / cURL / Python)
- **Compare** — Side-by-side comparison of up to 4 models with pricing/limit charts; selection persists in URL and localStorage
- **Providers** — Provider directory with model counts, npm packages, API endpoints, docs links
- **⌘K Command Palette** — Global fuzzy page navigation + live model search
- **Dark Mode** — System preference with manual toggle; charts re-theme instantly
- **i18n** — English + Chinese（中文路由带 `/zh` 前缀）
- **Shareable URLs** — Filters and compare selections are reflected in URL query params

## Tech Stack

| Layer | Technology |
|---|---|
| Frontend | Nuxt 4 (SPA mode, SSR disabled) |
| UI | Nuxt UI v4 (Amber + Stone theme), self-hosted Space Grotesk / Inter variable fonts |
| Charts | Chart.js (tree-shaken per component) |
| Backend | Go (Gin + Gorm) |
| Database | MySQL 8.0 |
| Deployment | Docker (scratch image, single binary) + k8s + GitHub Actions |

## Getting Started

### Prerequisites

- Node.js 22+ & pnpm
- Go 1.25+
- MySQL 8.0

### Database Setup

```sql
CREATE DATABASE models DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- 建表
SOURCE api/migrations/init.sql;
-- 已有旧库的环境执行增量索引迁移
SOURCE api/migrations/002_add_provider_index.sql;
```

### Frontend

```bash
pnpm install
pnpm dev          # http://localhost:3000
```

`NUXT_PUBLIC_API_BASE`（默认取 `.env`，指向 `http://localhost:8080`）在构建时嵌入，改动后需重新 build。

### Backend

```bash
cd api
mkdir -p frontend/dist && printf '<!doctype html>' > frontend/dist/index.html   # go:embed 占位，本地编译必需（gitignored）
go run ./cmd/server   # http://localhost:8080
```

Docker 构建时会用真实前端产物覆盖 `api/frontend/dist`，本地占位文件不影响产物。

## Environment Variables

所有后端变量使用 `MODEL_` 前缀：

| Variable | Default | Description |
|----------|---------|-------------|
| `MODEL_DATABASE_DSN` | *(必填，无默认)* | MySQL DSN，缺失时启动即失败 |
| `MODEL_API_KEYS` | *(empty)* | `POST /api/v1/sync` 的 Bearer key。**为空时该端点直接拒绝（fail-closed）**，需配置后才可手动触发同步 |
| `MODEL_TRUSTED_PROXIES` | *(empty)* | 受信反代 CIDR 列表（逗号分隔）。为空 = 不信任任何代理；在 nginx ingress 后部署时需设置，否则限流器只看到一个 IP |
| `MODEL_CORS_ALLOWED_ORIGINS` | `*` | 允许的前端源。`*` 时不携带凭证；显式列出源时启用凭证 |
| `MODEL_SERVER_PORT` | `8080` | HTTP 监听端口 |
| `MODEL_SYNC_CRON_MINUTES` | `60` | 定时同步间隔；`0` 关闭后台同步 |
| `MODEL_SYNC_COOLDOWN_MINUTES` | `10` | 手动同步最小间隔 |
| `MODEL_MODELS_DEV_URL` | `https://models.dev/api.json` | 数据源 URL |
| `MODEL_GIN_MODE` | `release` | Gin 模式（`debug`/`release`） |
| `MODEL_SITE_URL` | `https://model.false.ltd` | 站点地址（sitemap 用） |

## API Endpoints

| Method | Path | Auth | Rate limit | Description |
|--------|------|------|-----------|-------------|
| GET | `/health` | No | — | Health check |
| GET | `/api/v1/models` | No | 60/min | List models — pagination + filters（snake_case 参数；布尔筛选支持 `true`/`false`；价格区间过滤不含无价模型） |
| GET | `/api/v1/models/:id` | No | 60/min | Model detail（404 与服务端错误严格区分） |
| GET | `/api/v1/providers` | No | 60/min | Providers with `model_count` |
| GET | `/api/v1/stats` | No | 60/min | Aggregated statistics（内存缓存，同步时失效） |
| GET | `/api/v1/compare` | No | 60/min | Compare models（`?ids=1,2,3`，保持请求顺序并去重） |
| POST | `/api/v1/sync` | Bearer | 10/min | Trigger sync（并发触发返回 409；空 API_KEYS 返回 503） |
| GET | `/api/v1/sync/status` | No | 10/min | Last sync time |
| GET | `/sitemap.xml` | No | — | Cached sitemap（en/zh hreflang） |

响应统一为 `{ code, data, message, meta? }`；分页 meta 为 `{ total, page, page_size, total_pages }`。模型的关联字段为 `provider`（单数对象），服务商计数为 `model_count`。

## Architecture

```
models.dev/api.json
       ↓  POST /api/v1/sync（单事务全量 upsert + 按 provider 复合键清理，互斥锁防并发）
   MySQL (providers, models)          api/internal/cache — stats/providers/sitemap
       ↓  Go API（Gin + Gorm）            内存缓存（singleflight，同步成功即失效）
Client composables（useAsyncData 共享 key 去重）
       ↓
Vue 3 组件（图表 useChart 统一生命周期 + 主题重绘）
```

**同步语义**：全量拉取 → 单事务内 upsert providers/models → 按 `(provider_id, model_id)` 复合键清理上游已删数据 → 更新 `synced_at` → 失效缓存。上游返回空数据时拒绝同步（防止清库）。

## Deployment

Single Docker image (multi-stage: Nuxt build → scratch Go binary with embedded frontend). Push to `main` triggers GitHub Actions:

1. Build multi-arch image (amd64 + arm64) with Git SHA version
2. Push to GHCR with `sha-<short>` and `latest` tags
3. Apply `k8s.yml` with SHA-pinned image
4. Wait for rollout, auto-rollback on failure

## Project Structure

```
app/
├── pages/              # 路由（index, catalog, compare, providers, model/[modelId]）
├── components/         # charts/ catalog/ compare/ model/ + 全局组件（AppHeader、AppCommandPalette…）
├── composables/        # useCatalog, useCompare, useChart, useCountUp, …
├── plugins/            # reveal.client.ts（v-reveal 滚动显现指令）
├── utils/              # format / chart / badge
└── assets/css/main.css # 设计令牌、动效、字体、UTable 覆盖

api/
├── cmd/server/         # 入口（HTTP 超时 + 优雅关停）
├── internal/
│   ├── cache/          # TTL 内存缓存 + singleflight（含单元测试）
│   ├── config/         # MODEL_* 环境变量
│   ├── handler/        # HTTP handlers（统一响应构造、错误不泄露内部细节）
│   ├── middleware/     # CORS / gzip / 限流 / 认证（fail-closed）/ 日志 / RequestID
│   ├── model/          # 数据模型与过滤类型
│   ├── repository/     # GORM 查询（WithTx 事务支持）
│   ├── service/        # 业务逻辑（stats 聚合、sync 事务）
│   └── router/         # 路由 + SPA 静态服务（白名单回退）
├── migrations/         # init.sql + 002_add_provider_index.sql
└── api.json            # models.dev 数据快照（本地测试用）

docs/                   # 参考资料（benchlm 数据源等）
i18n/locales/           # en.json / zh.json（双语言 key 严格对等）
```

## License

MIT — see [LICENSE](LICENSE).
