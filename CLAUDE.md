# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

AI model data analysis site — browse, compare, and visualize pricing/capability data for LLMs. SPA (SSR disabled) built with Nuxt 4, Supabase, and Chart.js. Deploys to Cloudflare Pages.

## Commands

```bash
pnpm dev          # Start dev server
pnpm build        # Production build (use to verify no errors)
pnpm generate     # Static site generation
pnpm preview      # Preview production build
```

No test runner or linter configured. Use `pnpm build` to verify correctness.

## Architecture

### Data Flow

`models.dev/api.json` → `POST /api/sync` (upsert into Supabase) → Server API routes query Supabase → Client composables fetch from API routes → Vue components render

### Pages

- `index.vue` — Overview dashboard with KPI cards and Chart.js charts
- `catalog.vue` — Searchable/filterable model table (UTable) and grid view
- `providers.vue` — Provider list with search, links to filtered catalog
- `model/[id].vue` — Model detail (pricing, limits, capabilities, modalities, timeline, integration, quick start)
- `compare.vue` — Side-by-side comparison of up to 10 models

### Composables

- `useCatalog` — Catalog page state: filters, sort, pagination, URL sync
- `useCompare` — Compare selection stored in localStorage + URL params (localStorage persists across navigation, URL enables sharing)
- `useCompareData` — Field definitions and chart data for compare page
- `useProviderFilter` — Provider filter with search, grouping, URL restore
- `useModelDetail` — Single model fetch with pricing fields and quick start code

### Chart Components (`app/components/charts/`)

All use Chart.js directly (not a wrapper library). Each component registers only the Chart.js modules it needs. Pattern: `ref<HTMLCanvasElement>`, create chart in `onMounted`, destroy in `onUnmounted`, re-render on prop changes via `watch`. Use `chartColor(i)` from `app/utils/format.ts` for colors — never hardcode.

### Theming

- **Colors**: `app.config.ts` sets `primary: "rose"`, `neutral: "stone"` (warm gray)
- **CSS**: Nuxt UI semantic classes (`text-default`, `bg-elevated`, `border-default`, `text-muted`, `text-toned`, `text-primary`, `bg-accented`, etc.) — always prefer these over raw Tailwind colors
- **Dark mode**: `--ui-primary` overrides in `app/assets/css/main.css` for rose accent tuning
- **Scrollbar**: Custom `--c-scrollbar-*` vars in main.css
- **Provider logos**: `.dark img[src*="models.dev/logos"]` gets white background for visibility
- **UTable styling**: Direct `table thead th / table tbody td` selectors in main.css (UTable doesn't expose style overrides via app.config)

### i18n

English (default) + Chinese. Keys in `i18n/locales/en.json` and `zh.json`. Strategy: `prefix_except_default`. Always add keys to both files.

### Compare State (`useCompare`)

Dual persistence: localStorage (survives page navigation) + URL query params `?models=1,2,3` (for sharing). URL params take priority on load. Max 10 models.

### Supabase Schema

**providers**: id (serial), provider_id (text slug), name, npm, env (array), doc, api
**models**: id (serial), provider_id (FK → providers.provider_id), name, family, model_id, release_date, last_updated, open_weights, attachment, reasoning, tool_call, temperature, structured_output, status, knowledge, interleaved (jsonb), modalities_input/output (arrays), cost_input/output/cache_read/cache_write/input_audio/output_audio, limit_context/input/output

Integer PKs throughout (changed from string UUIDs Apr 15, 2026). Provider logos at `https://models.dev/logos/${provider_id}.svg`.

### Environment

Supabase credentials in `.env`: `NUXT_PUBLIC_SUPABASE_URL` and `NUXT_PUBLIC_SUPABASE_KEY`. Public/publishable keys, no auth in the app.

## API Schema

Data source is `models.dev/api.json`. Full schema at https://github.com/nicepkg/models.dev. Key fields:

- **Provider**: name, npm, env[], doc, api (optional, for openai-compatible)
- **Model**: name, family, attachment, reasoning, tool_call, structured_output, temperature, knowledge, release_date, last_updated, open_weights, status (alpha|beta|deprecated), interleaved, cost.{input,output,reasoning,cache_read,cache_write,input_audio,output_audio}, limit.{context,input,output}, modalities.{input,output}
