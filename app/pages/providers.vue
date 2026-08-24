<template>
    <div>
        <!-- Header with catalog summary -->
        <div class="flex flex-col sm:flex-row sm:items-end justify-between gap-3 mb-5">
            <div>
                <h1 class="font-display text-2xl font-bold text-default tracking-tight">{{ t("providers.title") }}</h1>
                <div class="text-sm text-muted mt-1 tabular">
                    {{ providersList.length }} {{ t("providers.total") }}
                    <template v-if="summary.models">
                        · {{ summary.models.toLocaleString() }} {{ t("nav.catalog") }}
                        · {{ summary.openPct }}% {{ t("providers.openShare") }}
                    </template>
                </div>
            </div>
            <div class="flex items-center gap-2">
                <UInput
                    v-model="search"
                    :placeholder="t('providers.search')"
                    icon="i-lucide-search"
                    class="w-full sm:w-56"
                    size="sm"
                />
                <!-- View toggle (desktop) -->
                <div class="hidden sm:flex bg-default border border-default rounded-lg p-0.5 shrink-0">
                    <button
                        @click="viewMode = 'grid'"
                        class="rounded-md px-2 py-1 cursor-pointer transition-colors"
                        :class="viewMode === 'grid' ? 'bg-primary text-white' : 'text-toned hover:text-default'"
                    >
                        <UIcon name="i-lucide-layout-grid" class="size-3.5" />
                    </button>
                    <button
                        @click="viewMode = 'table'"
                        class="rounded-md px-2 py-1 cursor-pointer transition-colors"
                        :class="viewMode === 'table' ? 'bg-primary text-white' : 'text-toned hover:text-default'"
                    >
                        <UIcon name="i-lucide-table" class="size-3.5" />
                    </button>
                </div>
            </div>
        </div>

        <!-- Error state -->
        <div v-if="loadError" class="flex flex-col items-center justify-center py-24 text-center">
            <div class="size-16 rounded-full bg-default border border-default flex items-center justify-center mb-4">
                <UIcon name="i-lucide-cloud-off" class="size-7 text-muted" />
            </div>
            <div class="text-base font-medium text-default mb-1">{{ t("common.loadError") }}</div>
            <UButton :label="t('common.retry')" icon="i-lucide-refresh-cw" color="primary" variant="soft" size="sm" class="mt-3" @click="retry" />
        </div>

        <!-- Loading skeletons -->
        <div v-else-if="loading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
            <div v-for="i in 6" :key="i" class="bg-default border border-default rounded-xl p-4">
                <div class="flex items-center gap-3 mb-3">
                    <USkeleton class="size-9 rounded-lg" />
                    <div class="flex-1">
                        <USkeleton class="h-4 w-28 rounded mb-1.5" />
                        <USkeleton class="h-3 w-16 rounded" />
                    </div>
                </div>
                <USkeleton class="h-3 w-full rounded mb-2" />
                <USkeleton class="h-3 w-2/3 rounded" />
            </div>
        </div>

        <!-- Grid view -->
        <template v-else-if="viewMode === 'grid'">
            <!-- Searching: flat results -->
            <template v-if="search.trim()">
                <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
                    <ProviderCard v-for="p in filtered" :key="p.id" :provider="p" :enrich="enrich.get(p.id)" />
                </div>
                <div v-if="!filtered.length" class="text-center py-16 text-sm text-muted">
                    {{ t("providers.noResults") }}
                </div>
            </template>

            <!-- Browsing: featured + A-Z -->
            <template v-else>
                <!-- Featured -->
                <div class="flex items-center gap-2 mb-3">
                    <UIcon name="i-lucide-award" class="size-4 text-primary" />
                    <span class="text-xs font-semibold text-muted uppercase tracking-wider">{{ t("providers.featured") }}</span>
                </div>
                <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 mb-8">
                    <ProviderCard v-for="p in featured" :key="p.id" :provider="p" :enrich="enrich.get(p.id)" featured />
                </div>

                <!-- A-Z jump rail -->
                <div class="sticky top-14 z-10 -mx-1 px-1 py-2 mb-4">
                    <div class="flex flex-wrap items-center gap-x-2 gap-y-1.5 bg-elevated/70 backdrop-blur-xl border border-default rounded-xl shadow-sm px-3 py-2.5">
                        <span class="inline-flex items-center gap-1.5 text-sm font-semibold text-default mr-1">
                            <UIcon name="i-lucide-arrow-down-a-z" class="size-4 text-primary" />
                            {{ t("providers.aToZ") }}
                        </span>
                        <button
                            v-for="g in aToZGroups"
                            :key="g.letter"
                            @click="jumpTo(g.letter)"
                            class="size-7 sm:size-8 rounded-md text-[13px] sm:text-sm font-semibold cursor-pointer transition-all tabular"
                            :class="letterClass(g.letter, g.items.length)"
                            :disabled="!g.items.length"
                        >
                            {{ g.letter }}
                        </button>
                    </div>
                </div>

                <!-- Letter sections -->
                <section v-for="g in aToZGroups.filter((x) => x.items.length)" :key="g.letter" :id="`prov-letter-${g.letter}`" class="mb-6 scroll-mt-24">
                    <div class="sticky top-[4.6rem] z-[5] bg-default/95 backdrop-blur-sm py-1.5 mb-2 flex items-baseline gap-2">
                        <span class="font-display font-bold text-lg text-primary">{{ g.letter }}</span>
                        <span class="text-[11px] text-muted tabular">{{ g.items.length }}</span>
                    </div>
                    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-2">
                        <ProviderCard v-for="p in g.items" :key="p.id" :provider="p" :enrich="enrich.get(p.id)" compact />
                    </div>
                </section>
            </template>
        </template>

        <!-- Table view -->
        <div v-else class="border border-default rounded-xl overflow-hidden overflow-x-auto">
            <table class="w-full text-sm">
                <thead>
                    <tr class="border-b border-default bg-elevated/50">
                        <th class="text-left px-4 py-2.5 text-xs font-medium text-muted uppercase tracking-wider">
                            <button class="inline-flex items-center gap-1 cursor-pointer hover:text-default transition-colors" @click="setSort('name')">
                                {{ t("providers.title") }}
                                <UIcon :name="sortIcon('name')" class="size-3" :class="sortKey === 'name' ? 'text-primary' : 'opacity-40'" />
                            </button>
                        </th>
                        <th class="text-left px-4 py-2.5 text-xs font-medium text-muted uppercase tracking-wider">ID</th>
                        <th class="text-center px-4 py-2.5 text-xs font-medium text-muted uppercase tracking-wider">
                            <button class="inline-flex items-center gap-1 cursor-pointer hover:text-default transition-colors" @click="setSort('count')">
                                {{ t("providers.models") }}
                                <UIcon :name="sortIcon('count')" class="size-3" :class="sortKey === 'count' ? 'text-primary' : 'opacity-40'" />
                            </button>
                        </th>
                        <th class="text-center px-4 py-2.5 text-xs font-medium text-muted uppercase tracking-wider">{{ t("providers.openShare") }}</th>
                        <th class="hidden md:table-cell text-left px-4 py-2.5 text-xs font-medium text-muted uppercase tracking-wider">NPM</th>
                        <th class="hidden lg:table-cell text-left px-4 py-2.5 text-xs font-medium text-muted uppercase tracking-wider">API</th>
                        <th class="text-center px-4 py-2.5 text-xs font-medium text-muted uppercase tracking-wider">{{ t("providers.docs") }}</th>
                    </tr>
                </thead>
                <tbody>
                    <NuxtLink
                        v-for="p in sortedTable"
                        :key="p.id"
                        :to="localePath(`/catalog?providers=${p.id}`)"
                        custom
                    >
                        <template #default="{ navigate }">
                            <tr
                                class="border-b border-default last:border-b-0 hover:bg-elevated/50 transition-colors cursor-pointer"
                                @click="navigate"
                            >
                                <td class="px-4 py-3">
                                    <div class="flex items-center gap-2">
                                        <ProviderLogo :provider-id="p.id" cls="w-7 h-7 rounded-lg p-0.5 shrink-0" />
                                        <span class="font-medium text-default truncate">{{ p.name }}</span>
                                    </div>
                                </td>
                                <td class="px-4 py-3 text-toned font-mono text-xs">{{ p.id }}</td>
                                <td class="px-4 py-3 text-center">
                                    <span class="text-xs font-semibold text-toned bg-elevated px-2 py-0.5 rounded tabular">
                                        {{ p.model_count }}
                                    </span>
                                </td>
                                <td class="px-4 py-3 text-center text-xs tabular" :class="openToneClass(enrich.get(p.id)?.openPct)">
                                    {{ enrich.get(p.id)?.openPct ?? 0 }}%
                                </td>
                                <td class="hidden md:table-cell px-4 py-3 text-toned font-mono text-xs truncate max-w-40">{{ p.npm || "—" }}</td>
                                <td class="hidden lg:table-cell px-4 py-3 text-toned font-mono text-xs truncate max-w-60">{{ p.api_url || "—" }}</td>
                                <td class="px-4 py-3 text-center">
                                    <a
                                        v-if="p.doc_url"
                                        :href="p.doc_url"
                                        target="_blank"
                                        class="text-primary hover:underline text-xs"
                                        @click.stop
                                    >
                                        <span class="inline-flex items-center gap-0.5">{{ t("providers.docs") }}<UIcon name="i-lucide-arrow-up-right" class="size-3" /></span>
                                    </a>
                                    <span v-else class="text-muted">—</span>
                                </td>
                            </tr>
                        </template>
                    </NuxtLink>
                </tbody>
            </table>
            <div v-if="!sortedTable.length" class="text-center py-16 text-sm text-muted">
                {{ t("providers.noResults") }}
            </div>
        </div>

        <!-- Empty state (grid, not searching) -->
        <div v-if="viewMode === 'grid' && !search.trim() && !loading && !loadError && !providersList.length" class="text-center py-16 text-sm text-muted">
            {{ t("providers.noResults") }}
        </div>
    </div>
</template>

<script setup lang="ts">
    const localePath = useLocalePath();
    const { t } = useI18n();
    const { viewMode, search, providersList, filtered, loading, loadError, retry } = useProviders();

    useSeoMeta({
        title: t("seo.providersTitle"),
        ogTitle: t("seo.providersTitle"),
        description: t("seo.providersDescription"),
        ogDescription: t("seo.providersDescription"),
        keywords: t("seo.providersKeywords"),
        twitterCard: "summary",
    });

    // ---- enrichment from the shared atlas cloud (open share, median price) ----
    const { data: atlasRaw } = useAtlasData();

    const enrich = computed(() => {
        const acc = new Map<string, { count: number; open: number; prices: number[] }>();
        for (const p of atlasRaw.value?.data?.points || []) {
            let e = acc.get(p.p);
            if (!e) {
                e = { count: 0, open: 0, prices: [] };
                acc.set(p.p, e);
            }
            e.count++;
            if (p.ow) e.open++;
            if (p.ci != null && p.ci > 0) e.prices.push(p.ci);
        }
        const out = new Map<string, { count: number; openPct: number; median: number | null }>();
        for (const [id, e] of acc) {
            const s = [...e.prices].sort((a, b) => a - b);
            const mid = Math.floor(s.length / 2);
            out.set(id, {
                count: e.count,
                openPct: e.count ? Math.round((e.open / e.count) * 100) : 0,
                median: s.length ? (s.length % 2 ? s[mid] : (s[mid - 1] + s[mid]) / 2) : null,
            });
        }
        return out;
    });

    const summary = computed(() => {
        let models = 0;
        let open = 0;
        for (const e of enrich.value.values()) {
            models += e.count;
            open += Math.round((e.openPct / 100) * e.count);
        }
        return { models, openPct: models ? Math.round((open / models) * 100) : 0 };
    });

    // ---- featured + A-Z sections (grid browsing) ----
    // 9 = exactly 3 rows on desktop, keeping the A-Z jump rail inside the
    // first viewport.
    const FEATURED_COUNT = 9;
    const featured = computed(() =>
        [...providersList.value]
            .sort((a, b) => (enrich.value.get(b.id)?.count ?? b.model_count) - (enrich.value.get(a.id)?.count ?? a.model_count))
            .slice(0, FEATURED_COUNT),
    );

    const aToZGroups = computed(() => {
        const featuredIds = new Set(featured.value.map((p) => p.id));
        const groups = new Map<string, any[]>();
        for (const p of providersList.value) {
            if (featuredIds.has(p.id)) continue;
            const letter = /[a-z]/i.test(p.name.charAt(0)) ? p.name.charAt(0).toUpperCase() : "#";
            if (!groups.has(letter)) groups.set(letter, []);
            groups.get(letter)!.push(p);
        }
        return [...groups.entries()]
            .sort((a, b) => (a[0] === "#" ? 1 : b[0] === "#" ? -1 : a[0].localeCompare(b[0])))
            .map(([letter, items]) => ({ letter, items }));
    });

    const jumpTo = (letter: string) => {
        document.getElementById(`prov-letter-${letter}`)?.scrollIntoView({ behavior: "smooth", block: "start" });
    };

    // Scroll-spy: the letter you are currently browsing lights up solid.
    // Sections render after data arrives, so (re)bind when groups appear.
    const activeLetter = ref<string | null>(null);
    const letterClass = (letter: string, hasItems: boolean) => {
        if (!hasItems) return "text-muted/40 cursor-default";
        if (activeLetter.value === letter) return "bg-primary text-white border-primary shadow-sm scale-110";
        return "bg-primary/10 text-primary border-primary/30 hover:bg-primary hover:text-white hover:border-primary";
    };
    let spy: IntersectionObserver | null = null;
    const bindSpy = () => {
        spy?.disconnect();
        const sections = Array.from(document.querySelectorAll<HTMLElement>("section[id^=prov-letter-]"));
        if (!sections.length) return;
        spy = new IntersectionObserver(
            (entries) => {
                for (const entry of entries) {
                    if (entry.isIntersecting) {
                        activeLetter.value = entry.target.id.replace("prov-letter-", "");
                    }
                }
            },
            { rootMargin: "-15% 0px -70% 0px", threshold: 0 },
        );
        sections.forEach((s) => spy!.observe(s));
    };
    onMounted(() => {
        watch(
            () => aToZGroups.value.filter((g) => g.items.length).length,
            () => nextTick(bindSpy),
            { immediate: true },
        );
        onUnmounted(() => spy?.disconnect());
    });

    // ---- table sorting ----
    const sortKey = ref<"name" | "count">("name");
    const sortDir = ref<"asc" | "desc">("asc");
    const setSort = (key: "name" | "count") => {
        if (sortKey.value === key) sortDir.value = sortDir.value === "asc" ? "desc" : "asc";
        else {
            sortKey.value = key;
            sortDir.value = key === "count" ? "desc" : "asc";
        }
    };
    const sortIcon = (key: "name" | "count") =>
        sortKey.value === key
            ? sortDir.value === "asc"
                ? "i-lucide-arrow-up-narrow-wide"
                : "i-lucide-arrow-down-wide-narrow"
            : "i-lucide-arrow-up-down";

    const sortedTable = computed(() => {
        const arr = [...filtered.value];
        const dir = sortDir.value === "asc" ? 1 : -1;
        if (sortKey.value === "count") {
            arr.sort((a, b) => ((enrich.value.get(a.id)?.count ?? a.model_count) - (enrich.value.get(b.id)?.count ?? b.model_count)) * dir);
        } else {
            arr.sort((a, b) => a.name.localeCompare(b.name) * dir);
        }
        return arr;
    });

    const openToneClass = (pct: number | undefined) =>
        pct == null ? "text-muted" : pct >= 60 ? "text-success font-semibold" : pct >= 20 ? "text-toned" : "text-muted";
</script>
