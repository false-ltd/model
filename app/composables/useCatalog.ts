import { watchDebounced } from "@vueuse/core";
import type { Model, CatalogMeta, ApiResponse, SelectItem } from "~/types";

const EMPTY_CATALOG = {
    data: [] as Model[],
    meta: { total: 0, page: 1, page_size: 50, total_pages: 0 } as CatalogMeta,
};

/** Stable serialization of a query object for change detection. */
const stableQuery = (q: Record<string, unknown>) =>
    JSON.stringify(
        Object.keys(q)
            .sort()
            .map((k) => {
                const v = q[k];
                return [k, Array.isArray(v) ? v.join(",") : (v ?? "")];
            }),
    );

/** Query keys that affect the catalog request; unrelated changes (e.g. the
 * compare selection's ?models=) must not trigger a refetch. */
const CATALOG_QUERY_KEYS = [
    "q", "page", "sort", "order", "providers", "reasoning", "vision", "freeOnly",
    "toolCall", "attachment", "openWeights", "structuredOutput", "temperature",
    "inputTypes", "outputTypes", "priceMin", "priceMax", "outputPriceMin",
    "outputPriceMax", "pageSize",
] as const;

const catalogQueryKey = (q: Record<string, unknown>) =>
    JSON.stringify(CATALOG_QUERY_KEYS.map((k) => [k, Array.isArray(q[k]) ? q[k].join(",") : (q[k] ?? "")]));

export function useCatalog() {
    const route = useRoute();
    const router = useRouter();
    const config = useRuntimeConfig();

    const sortField = ref((route.query.sort as string) || "last_updated");
    const sortOrder = ref((route.query.order as string) === "asc" ? "asc" : "desc");
    const selectedProviders = ref<SelectItem[]>([]);
    const selectedInputTypes = ref<SelectItem[]>([]);
    const selectedOutputTypes = ref<SelectItem[]>([]);
    const priceRange = ref([Number(route.query.priceMin) || 0, Number(route.query.priceMax) || 100]);
    const outputPriceRange = ref([Number(route.query.outputPriceMin) || 0, Number(route.query.outputPriceMax) || 100]);
    const searchQuery = ref((route.query.q as string) || "");
    const currentPage = ref(Number(route.query.page) || 1);
    const pageSize = ref(Number(route.query.pageSize) || 50);
    const filters = reactive({
        freeOnly: route.query.freeOnly === "true",
        reasoning: route.query.reasoning === "true",
        vision: route.query.vision === "true",
        toolCall: route.query.toolCall === "true",
        attachment: route.query.attachment === "true",
        openWeights: route.query.openWeights === "true",
        structuredOutput: route.query.structuredOutput === "true",
        temperature: route.query.temperature === "true",
    });

    const { t } = useI18n();

    const filterToggles = computed(() => [
        {
            key: "freeOnly" as const,
            label: t("catalog.free"),
            icon: "i-lucide-badge-dollar-sign",
            color: "success",
        },
        {
            key: "reasoning" as const,
            label: t("catalog.reasoning"),
            icon: "i-lucide-brain",
            color: "warning",
        },
        {
            key: "toolCall" as const,
            label: t("catalog.colToolCall"),
            icon: "i-lucide-wrench",
            color: "info",
        },
        { key: "vision" as const, label: t("catalog.vision"), icon: "i-lucide-eye", color: "primary" },
        {
            key: "openWeights" as const,
            label: t("catalog.colWeights"),
            icon: "i-lucide-unlock",
            color: "secondary",
        },
        { key: "attachment" as const, label: t("catalog.colAttachment"), icon: "i-lucide-paperclip", color: "error" },
        {
            key: "structuredOutput" as const,
            label: t("catalog.colStructured"),
            icon: "i-lucide-braces",
            color: "info",
        },
        {
            key: "temperature" as const,
            label: t("catalog.colTemperature"),
            icon: "i-lucide-thermometer",
            color: "secondary",
        },
    ]);

    const buildQuery = (): Record<string, string> => {
        const query: Record<string, string> = {};
        if (searchQuery.value.trim()) query.q = searchQuery.value.trim();
        if (currentPage.value > 1) query.page = String(currentPage.value);
        if (sortField.value && sortField.value !== "last_updated") query.sort = sortField.value;
        if (sortOrder.value === "asc") query.order = "asc";
        if (selectedProviders.value.length)
            query.providers = selectedProviders.value.map((p: SelectItem | string) => typeof p === 'string' ? p : p.value).join(",");
        if (filters.reasoning) query.reasoning = "true";
        if (filters.vision) query.vision = "true";
        if (filters.freeOnly) query.freeOnly = "true";
        if (filters.toolCall) query.toolCall = "true";
        if (filters.attachment) query.attachment = "true";
        if (filters.openWeights) query.openWeights = "true";
        if (filters.structuredOutput) query.structuredOutput = "true";
        if (filters.temperature) query.temperature = "true";
        if (selectedInputTypes.value.length)
            query.inputTypes = selectedInputTypes.value.map((t: SelectItem | string) => typeof t === 'string' ? t : t.value).join(",");
        if (selectedOutputTypes.value.length)
            query.outputTypes = selectedOutputTypes.value.map((t: SelectItem | string) => typeof t === 'string' ? t : t.value).join(",");
        if (priceRange.value[0]! > 0) query.priceMin = String(priceRange.value[0]);
        if (priceRange.value[1]! < 100) query.priceMax = String(priceRange.value[1]);
        if (outputPriceRange.value[0]! > 0) query.outputPriceMin = String(outputPriceRange.value[0]);
        if (outputPriceRange.value[1]! < 100) query.outputPriceMax = String(outputPriceRange.value[1]);
        if (pageSize.value !== 50) query.pageSize = String(pageSize.value);
        return query;
    };

    const buildApiParams = () => {
        const params: Record<string, string> = {
            page: String(currentPage.value),
            page_size: String(pageSize.value),
            sort: sortField.value,
        };
        if (sortOrder.value === "desc") params.order = "desc";
        if (searchQuery.value.trim()) params.q = searchQuery.value.trim();
        if (selectedProviders.value.length)
            params.providers = selectedProviders.value.map((p: SelectItem | string) => typeof p === 'string' ? p : p.value).join(",");
        if (filters.reasoning) params.reasoning = "true";
        if (filters.vision) params.vision = "true";
        if (filters.freeOnly) params.free_only = "true";
        if (filters.toolCall) params.tool_call = "true";
        if (filters.attachment) params.attachment = "true";
        if (filters.openWeights) params.open_weights = "true";
        if (filters.structuredOutput) params.structured_output = "true";
        if (filters.temperature) params.temperature = "true";
        if (selectedInputTypes.value.length)
            params.input_types = selectedInputTypes.value.map((t: SelectItem | string) => typeof t === 'string' ? t : t.value).join(",");
        if (selectedOutputTypes.value.length)
            params.output_types = selectedOutputTypes.value.map((t: SelectItem | string) => typeof t === 'string' ? t : t.value).join(",");
        if (priceRange.value[0]! > 0 || priceRange.value[1]! < 100) {
            params.price_min = String(priceRange.value[0]);
            params.price_max = String(priceRange.value[1]);
        }
        if (outputPriceRange.value[0]! > 0 || outputPriceRange.value[1]! < 100) {
            params.price_output_min = String(outputPriceRange.value[0]);
            params.price_output_max = String(outputPriceRange.value[1]);
        }
        return params;
    };

    // Only touch the URL when the query actually changed — a redundant
    // router.replace still fires the route watcher below and would trigger
    // a duplicate fetch.
    const syncToUrl = () => {
        const query = buildQuery();
        if (stableQuery(query) === stableQuery(route.query as Record<string, unknown>)) return;
        router.replace({ query });
    };

    const toggleSort = (field: string) => {
        if (sortField.value === field) {
            sortOrder.value = sortOrder.value === "asc" ? "desc" : "asc";
        } else {
            sortField.value = field;
            sortOrder.value = "asc";
        }
        currentPage.value = 1;
        syncToUrl();
    };

    const toggleFilter = (key: keyof typeof filters) => {
        filters[key] = !filters[key];
        // No explicit syncToUrl here: the filters watcher handles it.
    };

    const { data: catalogData, pending: loading, execute: fetchCatalog } = useAsyncData(
        "catalog",
        async () => {
            const res = await $fetch<ApiResponse<Model[]>>(`${config.public.apiBase}/api/v1/models`, { params: buildApiParams() });
            return { data: res.data, meta: res.meta ?? EMPTY_CATALOG.meta };
        },
        {
            immediate: false,
            default: () => EMPTY_CATALOG,
        },
    );

    const models = computed((): Model[] => catalogData.value.data);

    // Filter changes flow through the URL; when the current state already
    // matches the URL (e.g. providers restored from ?providers= on mount)
    // there is nothing to do — and importantly no page reset.
    watch(
        [selectedProviders, selectedInputTypes, selectedOutputTypes, filters],
        () => {
            if (stableQuery(buildQuery()) === stableQuery(route.query as Record<string, unknown>)) return;
            currentPage.value = 1;
            syncToUrl();
        },
        { deep: true },
    );

    watchDebounced(
        [priceRange, outputPriceRange],
        () => {
            if (stableQuery(buildQuery()) === stableQuery(route.query as Record<string, unknown>)) return;
            currentPage.value = 1;
            syncToUrl();
        },
        { debounce: 300, deep: true },
    );

    watchDebounced(searchQuery, () => {
        if (stableQuery(buildQuery()) === stableQuery(route.query as Record<string, unknown>)) return;
        currentPage.value = 1;
        syncToUrl();
    }, { debounce: 300 });

    let lastCatalogKey = catalogQueryKey(route.query as Record<string, unknown>);

    watch(
        () => route.query,
        () => {
            const key = catalogQueryKey(route.query as Record<string, unknown>);
            if (key === lastCatalogKey) return;
            lastCatalogKey = key;
            currentPage.value = Number(route.query.page) || 1;
            fetchCatalog();
        },
    );

    const goToPage = (page: number) => {
        if (page < 1 || page > catalogData.value.meta.total_pages) return;
        currentPage.value = page;
        syncToUrl();
        if (import.meta.client) window.scrollTo({ top: 0, behavior: "smooth" });
    };

    const changePageSize = (size: number) => {
        pageSize.value = size;
        currentPage.value = 1;
        syncToUrl();
    };

    const inputTypeItems = ["text", "image", "pdf", "video", "audio"].map((v) => ({
        label: v.charAt(0).toUpperCase() + v.slice(1),
        value: v,
    }));

    const outputTypeItems = ["text", "image", "video", "audio"].map((v) => ({
        label: v.charAt(0).toUpperCase() + v.slice(1),
        value: v,
    }));

    return {
        searchQuery,
        sortField,
        sortOrder,
        selectedProviders,
        selectedInputTypes,
        selectedOutputTypes,
        priceRange,
        outputPriceRange,
        currentPage,
        pageSize,
        filters,
        filterToggles,
        toggleFilter,
        catalogData,
        models,
        loading,
        fetchCatalog,
        toggleSort,
        goToPage,
        changePageSize,
        inputTypeItems,
        outputTypeItems,
    };
}
