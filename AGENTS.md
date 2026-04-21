# AGENTS.md

## Verification Commands (not `pnpm dev`)

```bash
pnpm build   # Frontend: production build to verify no errors
go build ./cmd/server   # Backend: must be in api/ directory or use full path
```

No test runner or linter configured. These are the only verification steps.

## Project Structure

- `app/` — Nuxt 4 frontend (pages, components, composables)
- `api/` — Go backend (Gin+Gorm). Entry point: `api/cmd/server/main.go`
- `migrations/init.sql` — MySQL schema (providers, models tables)
- `api/api.json` — Data source from models.dev

## Go API

- Runs on `model.false.ltd`
- Response format: `{ code, data, message, meta? }`
- Query params: **snake_case** (`page_size`, `input_types`, `free_only`)
- Auth: only `POST /api/v1/sync` requires Bearer token (API_KEYS env)

## Frontend Patterns

### UTable (catalog table)

Boolean fields (reasoning, tool_call, open_weights, etc.) render as **colored badges under model name in first column** — not as separate columns.

### Chart Components

Located in `app/components/charts/`. Use Chart.js **directly** (not vue-chartjs wrapper). Pattern:

```ts
const canvasRef = ref<HTMLCanvasElement>()
onMounted(() => { new Chart(canvasRef.value, {...}) })
onUnmounted(() => { chart.destroy() })
watch(props, () => { chart.update() })
```

Use `chartColor(i)` from `app/utils/format.ts` — never hardcode colors.

### CSS / Theming

- Use **Nuxt UI semantic classes** (`text-default`, `bg-elevated`, `border-default`, `text-muted`) — never raw Tailwind colors
- Dark mode accent: `--ui-primary` in `app/assets/css/main.css`
- Provider logos: `.dark img[src*="models.dev/logos"]` needs white background

### i18n

- Strategy: `prefix_except_default` (English default, Chinese prefixed)
- **Always add keys to both files** (`i18n/locales/en.json`, `zh.json`)

## Data Flow

`models.dev/api.json` → `POST /api/v1/sync` → MySQL → Go API → `useRuntimeConfig().public.apiBase` → Vue

## Local Dev

```bash
docker-compose up   # Go API + MySQL
pnpm dev            # Frontend (separate terminal)
```

`NUXT_PUBLIC_API_BASE` is build-time embedded — must rebuild frontend to change API URL.
