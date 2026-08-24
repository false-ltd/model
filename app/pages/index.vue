<template>
    <div>
        <!-- Error state -->
        <div v-if="loadError" class="flex flex-col items-center justify-center py-24 text-center">
            <div class="size-16 rounded-full bg-default border border-default flex items-center justify-center mb-4">
                <UIcon name="i-lucide-cloud-off" class="size-7 text-muted" />
            </div>
            <div class="text-base font-medium text-default mb-1">{{ t("common.loadError") }}</div>
            <UButton :label="t('common.retry')" icon="i-lucide-refresh-cw" color="primary" variant="soft" size="sm" class="mt-3" @click="retry" />
        </div>

        <!-- Loading skeletons -->
        <div v-else-if="loading">
            <div class="py-16 sm:py-24 text-center">
                <USkeleton class="h-3 w-44 rounded mx-auto mb-6" />
                <USkeleton class="h-12 sm:h-16 w-3/4 max-w-xl rounded mx-auto mb-3" />
                <USkeleton class="h-12 sm:h-16 w-1/2 max-w-sm rounded mx-auto mb-6" />
                <USkeleton class="h-4 w-2/3 max-w-lg rounded mx-auto" />
            </div>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-px bg-default border border-default rounded-2xl overflow-hidden mb-5">
                <div v-for="i in 4" :key="i" class="p-6 space-y-3">
                    <USkeleton class="h-3 w-20 rounded" />
                    <USkeleton class="h-10 w-24 rounded" />
                </div>
            </div>
        </div>

        <template v-else>
            <!-- Hero -->
            <section class="relative py-14 sm:py-20 lg:py-24">
                <AmbientBackground />
                <div class="max-w-4xl">
                    <div v-reveal class="flex items-center gap-2 mb-5">
                        <span class="relative flex size-2">
                            <span class="absolute inline-flex h-full w-full rounded-full bg-success opacity-60 animate-ping" />
                            <span class="relative inline-flex rounded-full size-2 bg-success" />
                        </span>
                        <span class="text-[11px] font-medium text-muted uppercase tracking-[0.18em]">
                            {{ t("overview.heroEyebrow") }}
                        </span>
                    </div>

                    <h1
                        v-reveal="80"
                        class="font-display font-bold text-default leading-[0.98] tracking-tight text-[clamp(2.75rem,7vw,5.5rem)]"
                    >
                        {{ t("overview.heroLine1") }}<br />
                        <span class="text-transparent bg-clip-text bg-linear-to-r from-amber-500 via-orange-500 to-rose-500">
                            {{ t("overview.heroLine2") }}
                        </span>
                    </h1>

                    <p v-reveal="160" class="mt-6 max-w-xl text-base sm:text-lg text-muted leading-relaxed">
                        {{ t("overview.heroSub", { models: stats.totalModels.toLocaleString(), providers: stats.totalProviders.toLocaleString() }) }}
                    </p>

                    <div v-reveal="240" class="mt-8 flex flex-wrap items-center gap-3">
                        <div v-magnetic class="inline-block rounded-xl">
                            <UButton
                                :to="localePath('/catalog')"
                                :label="t('overview.heroCtaCatalog')"
                                icon="i-lucide-arrow-right"
                                trailing-icon
                                size="lg"
                                class="rounded-xl"
                            />
                        </div>
                        <div v-magnetic class="inline-block rounded-xl">
                            <UButton
                                :to="localePath('/compare')"
                                :label="t('overview.heroCtaCompare')"
                                icon="i-lucide-git-compare"
                                color="neutral"
                                variant="outline"
                                size="lg"
                                class="rounded-xl"
                            />
                        </div>
                    </div>
                </div>
            </section>

            <!-- Stats strip -->
            <section class="mb-5 grid grid-cols-2 md:grid-cols-4 gap-px bg-default border border-default rounded-2xl overflow-hidden">
                <NuxtLink
                    v-for="(stat, i) in statStrip"
                    :key="stat.label"
                    v-reveal="i * 70"
                    :to="localePath(stat.to)"
                    class="group relative bg-default p-5 sm:p-6 no-underline transition-colors hover:bg-elevated/60"
                >
                    <div class="text-[10px] font-medium text-muted uppercase tracking-[0.16em] mb-2.5 flex items-center gap-1.5">
                        <UIcon :name="stat.icon" class="size-3 text-primary" />
                        {{ stat.label }}
                    </div>
                    <div class="font-display tabular font-bold text-default text-4xl sm:text-5xl leading-none group-hover:text-primary transition-colors duration-300">
                        {{ stat.value }}
                    </div>
                    <div class="text-[11px] text-muted mt-2">{{ stat.sub }}</div>
                </NuxtLink>
            </section>

            <!-- The Atlas — every model as a star -->
            <section ref="atlasSectionRef" v-reveal class="atlas-section mb-5 bg-default border border-default rounded-2xl overflow-hidden">
                <div class="flex flex-wrap items-center justify-between gap-3 px-5 py-4 border-b border-default">
                    <div>
                        <h2 class="font-display font-bold text-lg text-default tracking-tight flex items-center gap-2">
                            <UIcon name="i-lucide-sparkles" class="size-4 text-primary" />
                            {{ t("atlas.title") }}
                        </h2>
                        <p class="text-xs text-muted mt-0.5">{{ t("atlas.subtitle") }}</p>
                    </div>
                    <div class="flex items-center gap-2">
                        <span class="hidden sm:inline text-[11px] text-muted">{{ t("atlas.hint") }}</span>
                        <UButton
                            icon="i-lucide-maximize"
                            color="neutral"
                            variant="ghost"
                            size="xs"
                            :aria-label="t('atlas.reset')"
                            @click="atlasRef?.resetView()"
                        />
                        <UButton
                            :icon="atlasLite ? 'i-lucide-sparkle' : 'i-lucide-sparkles'"
                            color="neutral"
                            variant="ghost"
                            size="xs"
                            :aria-label="t('atlas.toggleLite')"
                            :title="t('atlas.toggleLite')"
                            :class="atlasLite ? 'text-primary' : ''"
                            @click="atlasLite = !atlasLite"
                        />
                        <UButton
                            v-if="fullscreenSupported"
                            :icon="isFullscreen ? 'i-lucide-minimize' : 'i-lucide-expand'"
                            color="neutral"
                            variant="ghost"
                            size="xs"
                            :aria-label="isFullscreen ? t('atlas.exitFullscreen') : t('atlas.fullscreen')"
                            @click="toggleFullscreen"
                        />
                    </div>
                </div>

                <div class="atlas-body relative h-[420px] sm:h-[480px]">
                    <AtlasCanvas ref="atlasRef" :marginal="true" />

                    <!-- Zoom controls -->
                    <div class="absolute bottom-3 right-3 flex flex-col bg-default/70 backdrop-blur-sm border border-default rounded-lg overflow-hidden">
                        <button
                            :aria-label="t('atlas.toggle3D')"
                            :title="t('atlas.toggle3D')"
                            class="size-8 flex items-center justify-center transition-colors cursor-pointer"
                            :class="view3D ? 'text-primary bg-primary/10' : 'text-toned hover:text-primary hover:bg-accented'"
                            @click="toggle3D"
                        >
                            <UIcon :name="view3D ? 'i-lucide-box' : 'i-lucide-square'" class="size-4" />
                        </button>
                        <div class="h-px bg-default" />
                        <button
                            :aria-label="t('atlas.zoomIn')"
                            :title="t('atlas.zoomIn')"
                            class="size-8 flex items-center justify-center text-toned hover:text-primary hover:bg-accented transition-colors cursor-pointer"
                            @click="atlasRef?.zoomIn()"
                        >
                            <UIcon name="i-lucide-plus" class="size-4" />
                        </button>
                        <div class="h-px bg-default" />
                        <button
                            :aria-label="t('atlas.zoomOut')"
                            :title="t('atlas.zoomOut')"
                            class="size-8 flex items-center justify-center text-toned hover:text-primary hover:bg-accented transition-colors cursor-pointer"
                            @click="atlasRef?.zoomOut()"
                        >
                            <UIcon name="i-lucide-minus" class="size-4" />
                        </button>
                    </div>

                    <!-- Legend: top providers, hover to isolate (top-right so
                         it stays clear of the dense low-price cluster) -->
                    <div class="absolute top-3 right-3 flex flex-col gap-1 bg-default/70 backdrop-blur-sm border border-default rounded-lg px-2.5 py-2 max-w-44">
                        <div class="text-[9px] text-muted uppercase tracking-wider mb-0.5">{{ t("atlas.legend") }}</div>
                        <button
                            v-for="pv in legendProviders"
                            :key="pv.id"
                            class="flex items-center gap-1.5 text-left cursor-pointer group"
                            @mouseenter="atlasRef && (atlasRef.highlightProvider = pv.id)"
                            @mouseleave="atlasRef && (atlasRef.highlightProvider = null)"
                        >
                            <span
                                class="size-2 rounded-full shrink-0 transition-transform group-hover:scale-125"
                                :style="{ backgroundColor: pv.color }"
                            />
                            <span
                                class="text-[11px] truncate transition-colors"
                                :class="atlasRef?.highlightProvider === pv.id ? 'text-default font-medium' : 'text-muted group-hover:text-toned'"
                            >{{ pv.name }}</span>
                        </button>
                    </div>
                </div>

                <!-- Axis captions -->
                <div class="flex justify-between px-5 py-2.5 border-t border-default text-[10px] text-muted tabular">
                    <span>{{ t("atlas.axisX") }}</span>
                    <span>{{ t("atlas.axisY") }}</span>
                </div>
            </section>

            <!-- Decision charts: what capability costs, and the cheapest
                 model per context tier. Both derive from the atlas point
                 cloud client-side — zero extra requests. -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-3 mb-3">
                <div v-reveal>
                    <ChartsCapabilityPriceChart
                        :title="t('overview.capabilityPrice')"
                        :subtitle="t('overview.capabilityPriceSub')"
                        :points="capabilityCurve"
                        :x-label="t('overview.capabilityAxis')"
                    />
                </div>
                <div v-reveal="80">
                    <ChartsFrontierChart
                        :title="t('overview.frontier')"
                        :subtitle="t('overview.frontierSub')"
                        :items="frontierItems"
                        :free-label="t('overview.frontierFree')"
                    />
                </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-3 mb-3">
                <div v-reveal>
                    <ChartsProvidersValueChart
                        :title="t('overview.topProviders')"
                        :subtitle="t('overview.providersValueSub')"
                        :items="providerValue"
                        :median-label="t('overview.medianLabel')"
                    />
                </div>
                <div v-reveal="80">
                    <ChartsCapabilityBars
                        :title="t('overview.capabilityCoverage')"
                        :subtitle="t('overview.capabilityCoverageSub')"
                        :capabilities="capabilityItems"
                    />
                </div>
            </div>

            <!-- Compact catalog profile: open-weights share + modality legend -->
            <div v-reveal class="bg-default border border-default rounded-xl px-5 py-4 mb-3 flex flex-wrap items-center gap-x-8 gap-y-4">
                <div class="flex items-center gap-3">
                    <div class="font-display tabular font-bold text-3xl text-default leading-none">{{ openPct }}<span class="text-lg text-muted">%</span></div>
                    <div>
                        <div class="text-xs font-medium text-default">{{ t("overview.profileOpen") }}</div>
                        <div class="text-[11px] text-muted">{{ openCount.toLocaleString() }} / {{ stats.totalModels.toLocaleString() }}</div>
                    </div>
                </div>
                <div class="w-40 h-1.5 rounded-full bg-elevated overflow-hidden">
                    <div class="h-full rounded-full bg-primary" :style="{ width: openPct + '%' }" />
                </div>
                <div class="h-8 w-px bg-default hidden sm:block" />
                <div class="flex flex-wrap items-center gap-3">
                    <div
                        v-for="m in modalityLegend"
                        :key="m.type"
                        class="flex items-center gap-1.5 text-xs text-toned"
                    >
                        <UIcon :name="modalityIcon(m.type)" class="size-3.5 text-primary" />
                        <span class="capitalize">{{ m.type }}</span>
                        <span class="text-muted tabular">{{ m.count.toLocaleString() }}</span>
                    </div>
                </div>
            </div>

            <!-- Recently indexed -->
            <div v-reveal>
                <RecentTicker :title="t('overview.recentTitle')" :items="recentItems" />
            </div>
        </template>
    </div>
</template>

<script setup lang="ts">
    const localePath = useLocalePath();
    const { t } = useI18n();
    const config = useRuntimeConfig();
    const colorMode = useColorMode();

    useSeoMeta({
        title: t("seo.indexTitle"),
        ogTitle: t("seo.indexTitle"),
        description: t("seo.indexDescription"),
        ogDescription: t("seo.indexDescription"),
        keywords: t("seo.indexKeywords"),
        twitterCard: "summary_large_image",
    });

    // Not awaited so skeletons render while the request is in flight.
    const { data: result, status, error, execute } = useAsyncData("overview-stats", () =>
        $fetch(`${config.public.apiBase}/api/v1/stats`),
    );
    const loading = computed(() => status.value === "pending");
    const loadError = computed(() => error.value != null);
    const retry = () => execute();

    const stats = computed(
        () =>
            (result.value as any)?.data?.stats || {
                totalProviders: 0,
                totalModels: 0,
                freeModelsCount: 0,
                freeModelsPct: 0,
                medianInputPrice: 0,
                reasoningCount: 0,
                reasoningPct: 0,
                medianContext: 0,
            },
    );

    // Count-up numerals for the stats strip.
    const modelsCount = useCountUp(computed(() => stats.value.totalModels));
    const providersCount = useCountUp(computed(() => stats.value.totalProviders));
    const freeCount = useCountUp(computed(() => stats.value.freeModelsCount));

    const statStrip = computed(() => [
        {
            label: t("overview.totalModels"),
            icon: "i-lucide-layers",
            value: modelsCount.value.toLocaleString(),
            sub: t("overview.totalModelsSub"),
            to: "/catalog",
        },
        {
            label: t("overview.providers"),
            icon: "i-lucide-building-2",
            value: providersCount.value.toLocaleString(),
            sub: "total",
            to: "/providers",
        },
        {
            label: t("overview.freeModels"),
            icon: "i-lucide-badge-dollar-sign",
            value: freeCount.value.toLocaleString(),
            sub: `${stats.value.freeModelsPct}% of catalog`,
            to: "/catalog?freeOnly=true",
        },
        {
            label: t("overview.avgContext"),
            icon: "i-lucide-expand",
            value: formatTokens(stats.value.medianContext),
            sub: t("overview.avgContextSub"),
            to: "/catalog?sort=limit_context&order=desc",
        },
    ]);

    const capabilityItems = computed(() => {
        const caps = (result.value as any)?.data?.capabilities || {};
        const colors: Record<string, string> = {
            reasoning: chartColor(3),
            tool_call: successColor(),
            vision: chartColor(0),
            attachment: chartColor(5),
            temperature: chartColor(1),
        };
        return Object.entries(caps).map(([key, val]: [string, any]) => ({
            label: key.replace("_", " ").replace(/\b\w/g, (l) => l.toUpperCase()),
            count: val.count,
            total: val.total,
            pct: val.total > 0 ? Math.round((val.count / val.total) * 100) : 0,
            color: colors[key] || chartColor(0),
        }));
    });

    // Atlas legend: top 8 providers by model count, colored with the same
    // palette order the canvas uses (provider order = color index).
    const { data: atlasRaw } = useAtlasData();
    const atlasProviders = computed<any[]>(() => atlasRaw.value?.data?.providers || []);
    const atlasPoints = computed<any[]>(() => atlasRaw.value?.data?.points || []);
    const legendProviders = computed(() => {
        const counts = new Map<string, number>();
        for (const p of atlasPoints.value) counts.set(p.p, (counts.get(p.p) ?? 0) + 1);
        return [...counts.entries()]
            .sort((a, b) => b[1] - a[1])
            .slice(0, 8)
            .map(([id, count]) => ({
                id,
                name: atlasProviders.value.find((p: any) => p.id === id)?.name || id,
                count,
                color: providerColor(id, colorMode.value === "dark"),
            }));
    });

    // ---- decision-chart data, all derived from the shared atlas cloud ----
    const medianOf = (nums: number[]): number | null => {
        if (!nums.length) return null;
        const sorted = [...nums].sort((a, b) => a - b);
        const mid = Math.floor(sorted.length / 2);
        return sorted.length % 2 ? sorted[mid] : (sorted[mid - 1] + sorted[mid]) / 2;
    };

    const capabilityCurve = computed(() => {
        const groups = new Map<number, number[]>();
        for (const p of atlasPoints.value) {
            if (p.ci == null) continue;
            if (!groups.has(p.k)) groups.set(p.k, []);
            groups.get(p.k)!.push(p.ci);
        }
        return [...groups.entries()]
            .sort((a, b) => a[0] - b[0])
            .map(([k, arr]) => ({ k, median: medianOf(arr) ?? 0 }));
    });

    const FRONTIER_TIERS = [8192, 32768, 131072, 262144, 1048576, 2097152];
    const frontierItems = computed(() => {
        const labels = ["8K", "32K", "128K", "256K", "1M", "2M"];
        return FRONTIER_TIERS.map((tier, i) => {
            const inTier = atlasPoints.value.filter((p: any) => p.ctx != null && p.ctx >= tier);
            const paid = inTier.filter((p: any) => p.ci != null && p.ci > 0);
            let price: number | null = null;
            let model: string | null = null;
            for (const p of paid) {
                if (price == null || p.ci < price) {
                    price = p.ci;
                    model = p.n;
                }
            }
            const free = inTier.some((p: any) => p.ci === 0);
            return { tier: labels[i], price, model, free };
        });
    });

    const providerValue = computed(() => {
        const groups = new Map<string, { open: number; closed: number; prices: number[] }>();
        for (const p of atlasPoints.value) {
            let g = groups.get(p.p);
            if (!g) {
                g = { open: 0, closed: 0, prices: [] };
                groups.set(p.p, g);
            }
            if (p.ow) g.open++;
            else g.closed++;
            if (p.ci != null && p.ci > 0) g.prices.push(p.ci);
        }
        return [...groups.entries()]
            .sort((a, b) => b[1].open + b[1].closed - (a[1].open + a[1].closed))
            .slice(0, 8)
            .map(([id, g]) => ({
                name: atlasProviders.value.find((p: any) => p.id === id)?.name || id,
                open: g.open,
                closed: g.closed,
                median: medianOf(g.prices),
            }));
    });

    const openPct = computed(() => (result.value as any)?.data?.weightsDistribution?.openPct ?? 0);
    const openCount = computed(() => (result.value as any)?.data?.weightsDistribution?.open ?? 0);
    const modalityLegend = computed(() => ((result.value as any)?.data?.modalities?.input || []).slice(0, 6));
    const recentItems = computed(() => (result.value as any)?.data?.recent || []);

    const atlasRef = useTemplateRef<any>("atlasRef");
    const view3D = ref(true);
    const atlasLite = useAtlasLite();
    const toggle3D = () => {
        view3D.value = !view3D.value;
        atlasRef.value?.set3D(view3D.value);
    };

    // ---- fullscreen ----
    const atlasSectionRef = useTemplateRef<HTMLElement>("atlasSectionRef");
    const isFullscreen = ref(false);
    // ssr:false — safe to read at render time; avoids mount-order issues.
    const fullscreenSupported = computed(() => !!document.fullscreenEnabled);
    const toggleFullscreen = async () => {
        try {
            if (!document.fullscreenElement) await atlasSectionRef.value?.requestFullscreen?.();
            else await document.exitFullscreen();
        } catch {}
    };
    onMounted(() => {
        const onFs = () => (isFullscreen.value = document.fullscreenElement === atlasSectionRef.value);
        document.addEventListener("fullscreenchange", onFs);
        onUnmounted(() => document.removeEventListener("fullscreenchange", onFs));
    });
</script>

<style scoped>
/* Fullscreen atlas: header + canvas + axis captions fill the screen.
   The bg-default utility keeps working in fullscreen. */
.atlas-section:fullscreen {
    display: flex;
    flex-direction: column;
    border: none;
    border-radius: 0;
}
.atlas-section:fullscreen .atlas-body {
    flex: 1;
    height: auto;
}
</style>
