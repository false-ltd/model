# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

AI model data analysis site — browse, compare, and visualize pricing/capability data for LLMs. Monorepo with Nuxt 4 frontend + Go (Gin+Gorm) backend API. Deploys to k3s via GitHub Actions.

## Repository Structure

```
model/
├── app/            ← Nuxt frontend (pages, components, composables)
├── i18n/           ← Locale files (en.json, zh.json)
├── public/         ← Static assets
├── api/            ← Go backend API
│   ├── cmd/server/     Entry point
│   ├── internal/       Config, handlers, middleware, models, repositories, services, router
│   ├── migrations/     Database schema (init.sql)
│   ├── docs/swagger/   Auto-generated API docs
│   ├── api.json        Data source from models.dev
│   └── go.mod
├── Dockerfile          ← Multi-stage: Nuxt build → Go binary with embedded frontend
└── k8s.yml             ← Single deployment (Go binary + embedded frontend)
```

## Commands

### Frontend

```bash
pnpm dev          # Start dev server
pnpm build        # Production build (use to verify no errors)
pnpm generate     # Static site generation
pnpm preview      # Preview production build
```

### Backend

```bash
cd api
go build ./cmd/server     # Build (verify compilation)
go mod tidy               # Tidy dependencies
```

### Local Development

```bash
# Frontend (separate terminal)
pnpm dev

# Backend — requires local MySQL with MODEL_DATABASE_DSN or uses default
cd api
go run ./cmd/server
```

No test runner or linter configured. Use `pnpm build` to verify frontend, `go build` to verify backend.

## Architecture

### Data Flow

`models.dev/api.json` → Go API `POST /api/v1/sync` (upsert into MySQL) → Client composables fetch from Go API via `useRuntimeConfig().public.apiBase` → Vue components render

### Pages

- `index.vue` — Overview dashboard with KPI cards and Chart.js charts
- `catalog.vue` — Searchable/filterable model table (UTable) and grid view
- `providers.vue` — Provider list with search, links to filtered catalog
- `model/[id].vue` — Model detail (pricing, limits, capabilities, modalities, timeline, integration, quick start)
- `compare.vue` — Side-by-side comparison of up to 4 models

### Composables

- `useCatalog` — Catalog page state: filters, sort, pagination, URL sync. Builds API params in snake_case for Go API (`page_size`, `input_types`, `free_only`, etc.)
- `useCompare` — Compare selection stored in localStorage + URL params. Max 4 models. URL params take priority on load.
- `useCompareData` — Field definitions and chart data for compare page
- `useProviderFilter` — Provider filter with search, grouping, URL restore
- `useModelDetail` — Single model fetch with pricing fields and quick start code
- `useAutoSync` — Sync trigger + status polling
- `useProviders` — Provider list for providers page
- `useMobile` — Responsive breakpoint detection

### Catalog Table

`catalog/Table.vue` uses Nuxt UI UTable with render functions (`h()`). Boolean capability fields (reasoning, tool_call, open_weights, structured_output, attachment, temperature) are shown as colored badges under the model name in the first column — not as separate columns. API query params use snake_case; response meta uses `page_size`/`total_pages`.

### CompareFab

`CompareFab.vue` — FAB button fixed at bottom-right. Shows stacked provider avatars + count badge. Click opens popover card with model list (removable items) and "compare" link. Used on catalog and model detail pages.

### Chart Components (`app/components/charts/`)

All use Chart.js directly (not a wrapper library). Each component registers only the Chart.js modules it needs. Pattern: `ref<HTMLCanvasElement>`, create chart in `onMounted`, destroy in `onUnmounted`, re-render on prop changes via `watch`. Use `chartColor(i)` from `app/utils/format.ts` for colors — never hardcode.

### Theming

- **CSS**: Nuxt UI semantic classes (`text-default`, `bg-elevated`, `border-default`, `text-muted`, `text-toned`, `text-primary`, `bg-accented`, etc.) — always prefer these over raw Tailwind colors
- **Dark mode**: `--ui-primary` overrides in `app/assets/css/main.css` for rose accent tuning
- **Scrollbar**: Custom `--c-scrollbar-*` vars in main.css
- **Provider logos**: `.dark img[src*="models.dev/logos"]` gets white background for visibility
- **UTable styling**: Direct `table thead th / table tbody td` selectors in main.css (UTable doesn't expose style overrides via app.config)

### i18n

English (default) + Chinese. Keys in `i18n/locales/en.json` and `zh.json`. Strategy: `prefix_except_default`. Always add keys to both files.

## Go API

### Routes

| Method | Path                  | Auth          | Description                                          |
| ------ | --------------------- | ------------- | ---------------------------------------------------- |
| GET    | `/health`             | No            | Health check                                         |
| GET    | `/api/v1/models`      | No            | List models (paginated, filtered, snake_case params) |
| GET    | `/api/v1/models/:id`  | No            | Get model by ID                                      |
| GET    | `/api/v1/providers`   | No            | List providers with model counts                     |
| GET    | `/api/v1/stats`       | No            | Aggregated statistics                                |
| GET    | `/api/v1/compare`     | No            | Compare multiple models (`?ids=1,2,3`)               |
| POST   | `/api/v1/sync`        | Yes (API key) | Trigger data sync                                    |
| GET    | `/api/v1/sync/status` | No            | Sync status                                          |

### Response Format

All responses wrapped in `{ code, data, message, meta? }`. Paginated endpoints include `meta: { page, page_size, total, total_pages }`.

### Middleware

- **CORS**: Configured via `MODEL_CORS_ALLOWED_ORIGINS`
- **Rate limiting**: 60 req/min (general), 10 req/min (sync)
- **Auth**: Bearer token for sync endpoint only
- **Request ID**: UUID generation

### Database

MySQL 8.0. Schema in `api/migrations/init.sql`. Tables: `providers`, `models`.

### Environment Variables

All app-specific env vars use `MODEL_` prefix:

| Variable | Default | Description |
|----------|---------|-------------|
| `MODEL_DATABASE_DSN` | `root:123456@tcp(127.0.0.1:3306)/models?charset=utf8mb4&parseTime=True&loc=Local` | MySQL connection DSN |
| `MODEL_API_KEYS` | *(empty)* | Comma-separated auth keys for sync endpoint |
| `MODEL_CORS_ALLOWED_ORIGINS` | `*` | Allowed frontend origins |
| `MODEL_SERVER_PORT` | `8080` | HTTP listen port |
| `MODEL_SYNC_CRON_MINUTES` | `60` | Auto-sync interval in minutes |
| `MODEL_SYNC_COOLDOWN_MINUTES` | `10` | Minimum gap between manual syncs |
| `MODEL_MODELS_DEV_URL` | `https://models.dev/api.json` | Data source URL |
| `MODEL_GIN_MODE` | `debug` | Gin framework mode (`debug`/`release`) |

Frontend build-time: `NUXT_PUBLIC_API_BASE` — Go API URL (embedded at build time).

### Database Setup

MySQL 8.0, schema in `api/migrations/init.sql`. Two tables: `providers`, `models`. Composite unique key on `models` is `(model_id, provider_id)` — sync uses upsert on this pair.

```sql
CREATE DATABASE model DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'model'@'10.0.0.%' IDENTIFIED BY 'PASSWORD';
GRANT ALL PRIVILEGES ON model.* TO 'model'@'10.0.0.%';
FLUSH PRIVILEGES;
```

## Deployment

Single `k8s.yml` with one `model` deployment (Go binary with embedded frontend). GitHub Actions builds multi-arch Docker image (multi-stage: Nuxt → Go binary) and deploys on push to main. Image tagged with `sha-<short>` for rollbacks.

## API Schema

Data source is `api/api.json` (from models.dev). Full schema at https://github.com/nicepkg/models.dev. Key fields:

- **Provider**: name, npm, env[], doc, api (optional, for openai-compatible)
- **Model**: name, family, attachment, reasoning, tool_call, structured_output, temperature, knowledge, release_date, last_updated, open_weights, status (alpha|beta|deprecated), interleaved, cost.{input,output,reasoning,cache_read,cache_write,input_audio,output_audio}, limit.{context,input,output}, modalities.{input,output}

Behavioral guidelines to reduce common LLM coding mistakes. Merge with project-specific instructions as needed.

**Tradeoff:** These guidelines bias toward caution over speed. For trivial tasks, use judgment.

## 1. Think Before Coding

**Don't assume. Don't hide confusion. Surface tradeoffs.**

Before implementing:

- State your assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them - don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.

## 2. Simplicity First

**Minimum code that solves the problem. Nothing speculative.**

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.

Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

## 3. Surgical Changes

**Touch only what you must. Clean up only your own mess.**

When editing existing code:

- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it - don't delete it.

When your changes create orphans:

- Remove imports/variables/functions that YOUR changes made unused.
- Don't remove pre-existing dead code unless asked.

The test: Every changed line should trace directly to the user's request.

## 4. Goal-Driven Execution

**Define success criteria. Loop until verified.**

Transform tasks into verifiable goals:

- "Add validation" → "Write tests for invalid inputs, then make them pass"
- "Fix the bug" → "Write a test that reproduces it, then make it pass"
- "Refactor X" → "Ensure tests pass before and after"

For multi-step tasks, state a brief plan:

1. [Step] → verify: [check]
2. [Step] → verify: [check]
3. [Step] → verify: [check]

Strong success criteria let you loop independently. Weak criteria ("make it work") require constant clarification.

---

**These guidelines are working if:** fewer unnecessary changes in diffs, fewer rewrites due to overcomplication, and clarifying questions come before implementation rather than after mistakes.
