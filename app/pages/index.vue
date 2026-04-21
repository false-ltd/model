<template>
    <div>
        <!-- KPI Cards -->
        <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-6 gap-3 mb-5">
            <NuxtLink
                v-for="kpi in kpiCards"
                :key="kpi.label"
                :to="localePath(kpi.to)"
                class="bg-default border border-default rounded-xl p-4 hover:border-accented transition-colors no-underline"
            >
                <div class="text-[10px] text-muted uppercase tracking-wider mb-1">{{ kpi.label }}</div>
                <div class="text-[26px] font-bold leading-none" :style="{ color: kpi.color || 'var(--ui-text)' }">
                    {{ kpi.value }}
                </div>
                <div class="text-[10px] text-muted mt-1">{{ kpi.sub }}</div>
            </NuxtLink>
        </div>

        <!-- Donut + Modality + Provider -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-3 mb-5">
            <ChartsDonutChart
                :title="t('overview.priceDistribution')"
                :subtitle="t('overview.priceDistributionSub')"
                :segments="tierSegments"
                :total="stats.totalModels"
            />
            <ChartsModalityBreakdown
                :title="t('overview.modalityCoverage')"
                :subtitle="t('overview.modalityCoverageSub')"
                :input="modalities.input"
                :output="modalities.output"
                :total-models="stats.totalModels"
            />
            <ChartsTopKBarChart
                :title="t('overview.topProviders')"
                :subtitle="t('overview.topProvidersSub')"
                :items="providerItems"
            />
        </div>

        <!-- Open weights + Capability -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-3 mb-5">
            <ChartsDonutChart
                :title="t('overview.openWeights')"
                :subtitle="t('overview.openWeightsSub')"
                :segments="weightsSegments"
                :total="stats.totalModels"
            />
            <ChartsCapabilityBars
                :title="t('overview.capabilityCoverage')"
                :subtitle="t('overview.capabilityCoverageSub')"
                :capabilities="capabilityItems"
            />
        </div>

        <!-- Context Window Distribution -->
        <div class="mb-5">
            <ChartsContextDistribution
                :title="t('overview.contextDistribution')"
                :subtitle="t('overview.contextDistributionSub')"
                :items="contextDistribution"
            />
        </div>
    </div>
</template>

<script setup lang="ts">
    const localePath = useLocalePath();
    const { t } = useI18n();
    const { data: result } = await useAsyncData("overview-stats", () => $fetch("/api/stats"));

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

    const tierSegments = computed(() => {
        const tiers = (result.value as any)?.data?.tiers || {
            free: 0,
            budget: 0,
            standard: 0,
            premium: 0,
            enterprise: 0,
        };
        return [
            { label: t("common.free"), count: tiers.free, color: successColor() },
            { label: t("overview.budget"), count: tiers.budget, color: chartColor(0) },
            { label: t("overview.standard"), count: tiers.standard, color: chartColor(3) },
            { label: t("overview.premium"), count: tiers.premium, color: dangerColor() },
            { label: t("overview.enterprise"), count: tiers.enterprise, color: chartColor(4) },
        ];
    });

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

    const providerItems = computed(() => {
        const provs = (result.value as any)?.data?.topProviders || [];
        const maxCount = Math.max(...provs.map((p: any) => p.count), 1);
        return provs.map((p: any) => ({
            label: p.name,
            value: (p.count / maxCount) * 100,
            display: String(p.count),
            color: chartColor(2),
        }));
    });

    const modalities = computed(() => (result.value as any)?.data?.modalities || { input: [], output: [] });
    const contextDistribution = computed(() => (result.value as any)?.data?.contextDistribution || []);

    const weightsSegments = computed(() => {
        const w = (result.value as any)?.data?.weightsDistribution || { open: 0, closed: 0 };
        return [
            { label: t("overview.openWeightsLabel"), count: w.open, color: successColor() },
            { label: t("overview.closedWeightsLabel"), count: w.closed, color: chartColor(0) },
        ];
    });

    const kpiCards = computed(() => [
        { label: t("overview.providers"), value: stats.value.totalProviders, sub: "total", to: "/providers" },
        { label: t("overview.totalModels"), value: stats.value.totalModels, sub: t("overview.totalModelsSub"), to: "/catalog" },
        {
            label: t("overview.freeModels"),
            value: stats.value.freeModelsCount,
            sub: `${stats.value.freeModelsPct}% of catalog`,
            color: successColor(),
            to: "/catalog?freeOnly=true",
        },
        {
            label: t("overview.medianInputPrice"),
            value: `$${stats.value.medianInputPrice.toFixed(2)}`,
            sub: t("overview.perMTokens"),
            to: "/catalog?sort=cost_input&order=asc",
        },
        {
            label: t("overview.reasoning"),
            value: stats.value.reasoningCount,
            sub: `${stats.value.reasoningPct}% of catalog`,
            color: chartColor(3),
            to: "/catalog?reasoning=true",
        },
        {
            label: t("overview.avgContext"),
            value: formatTokens(stats.value.medianContext),
            sub: t("overview.avgContextSub"),
            to: "/catalog?sort=limit_context&order=desc",
        },
    ]);
</script>
