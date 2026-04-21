# Model — AI Model Catalog

Browse, compare, and visualize pricing and capability data for LLMs. Built with [Nuxt 4](https://nuxt.com), [Go (Gin+Gorm)](https://gin-gonic.com), and [Chart.js](https://www.chartjs.org).

Data sourced from [models.dev](https://models.dev).

## Features

- **Catalog** — Searchable, sortable table with 25+ columns, filterable by provider, input/output types, price range, and capabilities (free, reasoning, vision)
- **Model Detail** — Pricing breakdown, context/output limits, capabilities, modalities, timeline, integration info, and quick start code
- **Compare** — Side-by-side comparison of up to 10 models with charts for pricing and limits
- **Providers** — Provider directory with model counts, npm packages, and API endpoints
- **Overview Dashboard** — KPI cards, release timeline, price distribution, modality coverage, capability stats, and context window distribution charts
- **Dark Mode** — System preference with manual toggle
- **i18n** — English and Chinese
- **Shareable URLs** — Filters and compare selections are reflected in URL query params

## Tech Stack

| Layer | Technology |
|---|---|
| Framework | Nuxt 4 (SPA mode, SSR disabled) |
| UI | Nuxt UI v4 (Rose + Stone theme) |
| Charts | Chart.js |
| API | Go (Gin + Gorm) |
| Styling | Tailwind CSS v4 |
| Deployment | Cloudflare Pages |

## Getting Started

### Prerequisites

- Node.js 18+
- pnpm
- Go API server running (see API docs in `swagger.json`)

### Setup

```bash
# Install dependencies
pnpm install

# Configure environment
cp .env.example .env
# Edit .env with your Go API base URL
```

### Environment Variables

```
NUXT_PUBLIC_API_BASE=http://localhost:8080
```

### Development

```bash
pnpm dev          # Start dev server
pnpm build        # Production build
pnpm generate     # Static site generation
pnpm preview      # Preview production build
```

## Data Flow

```
models.dev/api.json
       ↓
Go API POST /api/v1/sync  (upsert into database)
       ↓
Client composables  (fetch from Go API)
       ↓
Vue components  (render)
```

## Project Structure

```
app/
├── pages/              # Routes (catalog, compare, providers, model/[id])
├── components/
│   ├── charts/         # Chart.js components (donut, bar, timeline, etc.)
│   ├── catalog/        # Catalog page sub-components (Table, Filters, CompareBar, ProviderPopover)
│   └── model/          # Model card, badge, limit bar
├── composables/        # State management (useCatalog, useCompare, useCompareData, etc.)
├── utils/format.ts     # Shared formatters and chart colors
└── assets/css/main.css # Global styles and UTable overrides

i18n/locales/           # English and Chinese translations
swagger.json            # Go API documentation
```

## License

MIT
