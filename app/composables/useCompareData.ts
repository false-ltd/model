export function useCompareData(compareModels: Ref<any[]>) {
    const { t } = useI18n();

    const identityFields = computed(() => [
        { key: "model_id", label: t("compare.modelId") },
        { key: "name", label: t("compare.name") },
        { key: "family", label: t("compare.family") },
        { key: "open_weights", label: t("compare.openWeights"), format: (v: boolean) => (v ? "✓" : "—") },
    ]);

    const pricingFields = computed(() => [
        { key: "cost_input", label: t("compare.input"), format: (v: number) => (v != null ? `$${v}` : "—") },
        { key: "cost_output", label: t("compare.output"), format: (v: number) => (v != null ? `$${v}` : "—") },
        { key: "cost_cache_read", label: t("compare.cacheRead"), format: (v: number) => (v != null ? `$${v}` : "—") },
        { key: "cost_cache_write", label: t("compare.cacheWrite"), format: (v: number) => (v != null ? `$${v}` : "—") },
        { key: "cost_reasoning", label: t("compare.reasoningCost"), format: (v: number) => (v != null ? `$${v}` : "—") },
        { key: "cost_input_audio", label: t("compare.inputAudio"), format: (v: number) => (v != null ? `$${v}` : "—") },
        { key: "cost_output_audio", label: t("compare.outputAudio"), format: (v: number) => (v != null ? `$${v}` : "—") },
    ]);

    const limitFields = computed(() => [
        { key: "limit_context", label: t("compare.contextWindow") },
        { key: "limit_input", label: t("compare.maxInput") },
        { key: "limit_output", label: t("compare.maxOutput") },
    ]);

    const capabilityFields = computed(() => [
        { key: "reasoning", label: t("detail.reasoning") },
        { key: "tool_call", label: t("compare.toolCall") },
        { key: "attachment", label: t("compare.attachment") },
        { key: "structured_output", label: t("compare.structuredOutput") },
        { key: "temperature", label: t("compare.temperature") },
    ]);

    const timelineFields = computed(() => [
        { key: "release_date", label: t("compare.releaseDate") },
        { key: "last_updated", label: t("compare.lastUpdated") },
        { key: "knowledge", label: t("compare.knowledge") },
        { key: "status", label: t("compare.status") },
    ]);

    const integrationFields = computed(() => [
        { key: "providers.name", label: t("detail.provider") },
        { key: "providers.npm", label: t("detail.npmPackage") },
        { key: "providers.api", label: t("detail.apiEndpoint") },
        { key: "providers.env.0", label: t("detail.envVariable") },
    ]);

    const hasPaidModels = computed(() =>
        compareModels.value.some((m) => (m.cost_input ?? 0) > 0 || (m.cost_output ?? 0) > 0),
    );

    const isCheapest = (m: any, field: string) => {
        const val = m[field];
        if (val == null) return false;
        const values = compareModels.value.map((m2) => m2[field]).filter((v) => v != null && v > 0);
        return values.length > 0 && val === Math.min(...values);
    };

    const isBest = (m: any, field: string) => {
        const val = m[field];
        if (val == null) return false;
        const values = compareModels.value.map((m2) => m2[field]).filter((v) => v != null);
        return values.length > 0 && val === Math.max(...values);
    };

    const pricingTabs = computed(() => [
        { key: "all", label: t("compare.allPricing") },
        { key: "input", label: t("compare.input") },
        { key: "output", label: t("compare.output") },
        { key: "cache", label: t("compare.cache") },
    ]);

    const pricingDatasets = (tab: string) => {
        const datasets = [];
        if (tab !== "output" && tab !== "cache")
            datasets.push({ label: t("compare.input"), data: compareModels.value.map((m) => m.cost_input || 0), backgroundColor: chartColor(0) });
        if (tab !== "input" && tab !== "cache")
            datasets.push({ label: t("compare.output"), data: compareModels.value.map((m) => m.cost_output || 0), backgroundColor: chartColor(1) });
        if (tab === "all" || tab === "cache") {
            datasets.push({ label: t("compare.cacheRead"), data: compareModels.value.map((m) => m.cost_cache_read || 0), backgroundColor: chartColor(2) });
            datasets.push({ label: t("compare.cacheWrite"), data: compareModels.value.map((m) => m.cost_cache_write || 0), backgroundColor: chartColor(4) });
        }
        return datasets;
    };

    const limitDatasets = computed(() => [
        {
            label: t("compare.context"),
            data: compareModels.value.map((m) => m.limit_context || 0),
            backgroundColor: chartColor(1),
        },
        {
            label: t("compare.maxOutput"),
            data: compareModels.value.map((m) => m.limit_output || 0),
            backgroundColor: chartColor(1) + "88",
        },
    ]);

    return {
        identityFields,
        pricingFields,
        limitFields,
        capabilityFields,
        timelineFields,
        integrationFields,
        hasPaidModels,
        isCheapest,
        isBest,
        pricingTabs,
        pricingDatasets,
        limitDatasets,
    };
}
