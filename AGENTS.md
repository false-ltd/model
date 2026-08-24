# AGENTS.md

本文件是 AI 编码助手（Claude Code / Codex / Cursor 等）在本仓库工作的唯一指令来源。CLAUDE.md 是本文件的符号链接。

## 验证命令（不要用 `pnpm dev` 验证）

```bash
pnpm build                          # 前端：生产构建验证无错误（仓库根目录执行）
cd api && go build ./cmd/server     # 后端：编译验证
cd api && go test ./...             # 后端：单元测试（internal/cache）
```

无 lint 配置，以上是仅有的验证步骤。

本地编译 Go 需先确保 `api/frontend/dist/` 存在（go:embed 要求，目录被 gitignore）：

```bash
mkdir -p api/frontend/dist && printf '<!doctype html>' > api/frontend/dist/index.html
```

Docker 构建时会用真实前端产物覆盖该目录；本地占位文件不会进入 git。

## 项目结构

- `app/` — Nuxt 4 前端（`ssr: false` SPA，单二进制部署是刻意架构，勿建议切 SSR）
- `api/` — Go 后端（Gin + Gorm），入口 `api/cmd/server/main.go`
- `api/migrations/` — MySQL schema（`init.sql` 全量 + `002_add_provider_index.sql` 增量）
- `api/api.json` — models.dev 数据快照（本地联调可将其用 http.server 做为 mock 上游）
- `i18n/locales/` — en.json / zh.json，**两文件 key 必须严格对等**（校验习惯：扫描 `t("…")` 引用比对）
- `docs/` — 参考资料（非代码）

## 后端要点

- 响应统一 `{ code, data, message, meta? }`；用 `model/response.go` 的构造器（`SuccessResponse` / `PagedSuccessResponse` / `CountSuccessResponse`），不要手写 `gin.H`
- 查询参数 **snake_case**（`page_size`、`input_types`、`free_only`）
- 关联字段命名：模型上的服务商对象是 `provider`（单数）；服务商计数是 `model_count`
- `api/internal/cache` — TTL 内存缓存 + singleflight，覆盖 stats/providers/sitemap；**sync 成功后自动失效**。改数据写入路径时记得评估缓存失效
- sync（`service/sync.go`）：互斥锁防并发（409）、单事务全量 upsert、按 `(provider_id, model_id)` 复合键清理、空上游拒绝同步；`MODEL_API_KEYS` 为空时 sync 端点 fail-closed（503）
- 错误处理：500 响应一律用 `handler.internalError`（记日志 + 通用文案），不得把 `err.Error()` 回给客户端
- 排序白名单在 `repository/model.go` 的 `mapSortColumn`（防注入）；新增可排序列要同时维护 `nullableSortColumns`

## 前端模式

### UTable（目录表）

布尔能力字段（reasoning、tool_call、open_weights 等）渲染为**第一列模型名下方的彩色徽章**，不单独成列。列渲染用 `h()` 渲染函数（`app/components/catalog/Table.vue`）。

### 图表

位于 `app/components/charts/`，直接用 Chart.js（不用 vue-chartjs 包装）。统一通过 `app/composables/useChart.ts`：

```ts
const { canvasRef } = useChart(getConfig, () => props.items)
```

useChart 已处理：销毁/重建、数据变化重绘（不重放动画）、**color mode 切换重绘**。颜色一律 `chartColor(i)`（`app/utils/format.ts`），禁止硬编码。

### 图标与符号

全站使用 Lucide 图标（`i-lucide-*`）。**界面禁止 emoji 和 Unicode 符号**（✓ → ↗ 之类一律换成图标组件；空值占位的 em-dash "—" 除外）。

### 设计系统

- 字体本地打包（`@fontsource-variable/inter` + `space-grotesk`），大标题/大数字用 `font-display` 工具类
- 颜色只用 Nuxt UI 语义类（`text-default` / `bg-elevated` / `border-default` / `text-muted` / `text-toned`），不用原生 Tailwind 色板（图表/徽章的历史半透明色类除外）
- 动效：令牌 `--ease-out-quint`；滚动显现用 `v-reveal` 指令（`app/plugins/reveal.client.ts`，值 = 交错延迟 ms）；**所有动画必须尊重 `prefers-reduced-motion`**（main.css 已有全局兜底）
- 深色模式 provider logo 兼容：`.dark img[src*="models.dev/logos"]` 白底规则在 main.css

### 数据获取

- `useAsyncData` 的 key 是跨组件共享缓存的唯一凭据：providers 一律用 `"providers-list"`（useProviders 与 useProviderFilter 共用），compare 全量数据用 `"compare-models"`
- CompareFab 渲染本地快照（`useCompare` 存 `{id, name, provider_id}`），**不额外请求**；仅 URL 恢复缺名字时补水一次
- catalog 的筛选流向是单向的：filter 状态 → `syncToUrl`（内容不变时跳过 replace）→ route.query watcher → fetch；改筛选逻辑时不要打破这条链
- ⌘K 命令面板（`app/components/AppCommandPalette.vue`）：Nuxt UI v4 的 group items 只接受静态数组，异步搜索通过 `v-model:search-term` + watchDebounced 驱动（注意 `@vueuse/core` 需显式 import）

### i18n

策略 `prefix_except_default`（英文无前缀，中文 `/zh`）。**新 key 必须同时加 en.json 和 zh.json**。

## 数据流

```
models.dev/api.json → POST /api/v1/sync（事务）→ MySQL → Go API（cache）→ apiBase → Vue
```

`NUXT_PUBLIC_API_BASE` 构建时嵌入，换 API 地址必须重新 build 前端。

## 行为准则（精简版）

1. **先想后写**：假设要显式说出；有多种解读时列出来，不要默默选一个
2. **最小改动**：只动必须动的；不顺手"改进"无关代码；匹配既有风格
3. **可验证**：每个改动要能对应一条验证方式（build / test / curl）；完成前跑验证命令
4. **清理自己产生的孤儿**：自己的改动导致的无用 import/变量要删；项目里既有的死代码只报告不动手

## GitNexus

本仓库已被 GitNexus 索引（`.gitnexus/`，gitignored）。CLI skill 位于 `.claude/skills/gitnexus/`（exploring / impact-analysis / debugging / refactoring / guide / cli）。

- 索引过期时先跑 `npx gitnexus analyze`
- 若会话中 GitNexus MCP 工具（`gitnexus_impact` 等）可用：改符号前先做影响分析，提交前跑 `gitnexus_detect_changes()`，HIGH/CRITICAL 风险要先告知用户
- MCP 工具不可用时，以人工调用图分析替代（本仓库规模下 grep + 读代码足够）

## 本地联调

```bash
# API（本地 MySQL，DSN 按需覆盖）
cd api && go run ./cmd/server

# 前端（默认 3000；被占用时 --port 指定，如 3300）
pnpm dev
```
