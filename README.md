# Model — AI Model Catalog

Browse, compare, and visualize pricing and capability data for LLMs. Built with [Nuxt 4](https://nuxt.com), [Supabase](https://supabase.com), and [Chart.js](https://www.chartjs.org).

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
| Database | Supabase (PostgreSQL) |
| Styling | Tailwind CSS v4 |
| Deployment | Cloudflare Pages |

## Getting Started

### Prerequisites

- Node.js 18+
- pnpm

### Setup

```bash
# Install dependencies
pnpm install

# Configure environment
cp .env.example .env
# Edit .env with your Supabase credentials
```

### Environment Variables

```
NUXT_PUBLIC_SUPABASE_URL=https://your-project.supabase.co
NUXT_PUBLIC_SUPABASE_KEY=your-public-anon-key
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
POST /api/sync  (upsert into Supabase)
       ↓
Server API routes  (query Supabase)
       ↓
Client composables  (fetch from API)
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
├── assets/css/main.css # Global styles and UTable overrides
└── types/              # Auto-generated Supabase types

server/api/             # REST endpoints (models, compare, providers, stats, sync)
supabase/migrations/    # Database schema
i18n/locales/           # English and Chinese translations
```

## License

MIT
