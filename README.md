# Model — AI Model Data Analysis

Browse, compare, and visualize pricing and capability data for LLMs. Data sourced from [models.dev](https://models.dev).

## Features

- **Overview Dashboard** — KPI cards, price distribution, modality coverage, capability stats, context window distribution charts
- **Catalog** — Searchable, sortable table with filters for provider, input/output types, price range, and capabilities (free, reasoning, vision)
- **Model Detail** — Pricing breakdown, context/output limits, capabilities, modalities, timeline, integration info, and quick start code
- **Compare** — Side-by-side comparison of up to 4 models with charts for pricing and limits
- **Providers** — Provider directory with model counts, npm packages, and API endpoints
- **Dark Mode** — System preference with manual toggle
- **i18n** — English and Chinese
- **Shareable URLs** — Filters and compare selections are reflected in URL query params

## Tech Stack

| Layer | Technology |
|---|---|
| Frontend | Nuxt 4 (SPA mode, SSR disabled) |
| UI | Nuxt UI v4 (Amber + Stone theme) |
| Charts | Chart.js |
| Backend | Go (Gin + Gorm) |
| Database | MySQL 8.0 |
| Deployment | Docker + k3s + GitHub Actions |

## Getting Started

### Prerequisites

- Node.js 22+ & pnpm
- Go 1.25+
- MySQL 8.0

### Database Setup

```sql
CREATE DATABASE model DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'model'@'10.0.0.%' IDENTIFIED BY 'PASSWORD';
GRANT ALL PRIVILEGES ON model.* TO 'model'@'10.0.0.%';
FLUSH PRIVILEGES;
```

Schema: `api/migrations/init.sql` (tables: `providers`, `models`)

### Frontend

```bash
pnpm install
pnpm dev          # http://localhost:3000
```

### Backend

```bash
cd api
go run ./cmd/server   # http://localhost:8080
```

## Environment Variables

All app-specific env vars use `MODEL_` prefix:

| Variable | Default | Description |
|----------|---------|-------------|
| `MODEL_DATABASE_DSN` | `root:123456@tcp(127.0.0.1:3306)/models?charset=utf8mb4&parseTime=True&loc=Local` | MySQL connection DSN |
| `MODEL_API_KEYS` | *(empty)* | Comma-separated auth keys for sync endpoint |
| `MODEL_CORS_ALLOWED_ORIGINS` | `*` | Allowed frontend origins |
| `MODEL_SERVER_PORT` | `8080` | HTTP listen port |
| `MODEL_SYNC_CRON_MINUTES` | `60` | Auto-sync interval |
| `MODEL_SYNC_COOLDOWN_MINUTES` | `10` | Minimum gap between manual syncs |
| `MODEL_MODELS_DEV_URL` | `https://models.dev/api.json` | Data source URL |
| `MODEL_GIN_MODE` | `debug` | Gin framework mode (`debug`/`release`) |

Frontend build-time: `NUXT_PUBLIC_API_BASE` — Go API URL (embedded at build time).

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/health` | No | Health check |
| GET | `/api/v1/models` | No | List models (paginated, filtered) |
| GET | `/api/v1/models/:id` | No | Get model by ID |
| GET | `/api/v1/providers` | No | List providers with model counts |
| GET | `/api/v1/stats` | No | Aggregated statistics |
| GET | `/api/v1/compare` | No | Compare models (`?ids=1,2,3`) |
| POST | `/api/v1/sync` | Yes | Trigger data sync |
| GET | `/api/v1/sync/status` | No | Sync status |

## Data Flow

```
models.dev/api.json
       ↓
Go API POST /api/v1/sync  (upsert into MySQL, composite key: model_id + provider_id)
       ↓
Client composables  (fetch from Go API)
       ↓
Vue components  (render)
```

## Deployment

Single Docker image (multi-stage: Nuxt build → Go binary with embedded frontend). Push to `main` triggers GitHub Actions:

1. Build multi-arch image (amd64 + arm64) with Git SHA version
2. Push to GHCR with `sha-<short>` and `latest` tags
3. Apply `k8s.yml` with SHA-pinned image
4. Wait for rollout, auto-rollback on failure

## Project Structure

```
app/
├── pages/              # Routes (index, catalog, compare, providers, model/[id])
├── components/
│   ├── charts/         # Chart.js components (donut, bar, capability, etc.)
│   ├── catalog/        # Catalog page sub-components (Table, Filters)
│   ├── compare/        # Compare page (Chart, MobileFieldList, ModelPicker)
│   └── model/          # Model card, badge, price/limit gauge, quick start
├── composables/        # State management (useCatalog, useCompare, useCompareData, etc.)
├── utils/              # Shared formatters, chart config, badge helpers
└── assets/css/main.css # Global styles and UTable overrides

api/
├── cmd/server/         # Entry point
├── internal/
│   ├── config/         # Env config (MODEL_* prefix)
│   ├── handler/        # HTTP handlers
│   ├── middleware/      # CORS, rate limit, auth, request ID
│   ├── model/          # Data models & filter types
│   ├── repository/     # Database queries (upsert on model_id+provider_id)
│   ├── service/        # Business logic & sync
│   └── router/         # Route registration & SPA serving
├── migrations/          # Database schema (init.sql)
└── api.json             # Data source snapshot from models.dev
```

## License

Private
