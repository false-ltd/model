export function useCatalog() {
    const route = useRoute();
    const router = useRouter();

    const viewMode = ref<"grid" | "table">("table");
    const sortField = ref((route.query.sort as string) || "last_updated");
    const sortOrder = ref((route.query.order as string) === "asc" ? "asc" : "desc");
    const selectedProviders = ref<any[]>([]);
    const selectedInputTypes = ref<any[]>([]);
    const selectedOutputTypes = ref<any[]>([]);
    const priceRange = ref([Number(route.query.priceMin) || 0, Number(route.query.priceMax) || 100]);
    const searchQuery = ref((route.query.q as string) || "");
    const currentPage = ref(Number(route.query.page) || 1);
    const filters = reactive({
        freeOnly: route.query.freeOnly === "true",
        reasoning: route.query.reasoning === "true",
        vision: route.query.vision === "true",
    });

    const { t } = useI18n();

    const filterToggles = computed(() => [
        { key: "freeOnly" as const, label: t("catalog.free"), icon: "i-lucide-badge-dollar-sign" },
        { key: "reasoning" as const, label: t("catalog.reasoning"), icon: "i-lucide-brain" },
        { key: "vision" as const, label: t("catalog.vision"), icon: "i-lucide-eye" },
    ]);

    const toggleFilter = (key: keyof typeof filters) => {
        filters[key] = !filters[key];
        currentPage.value = 1;
        syncToUrl();
    };

    const catalogData = ref<{
        data: any[];
        meta: { total: number; page: number; pageSize: number; totalPages: number };
    }>({
        data: [],
        meta: { total: 0, page: 1, pageSize: 100, totalPages: 0 },
    });

    const models = computed(() => catalogData.value.data);

    const buildQuery = (): Record<string, string> => {
        const query: Record<string, string> = {};
        if (searchQuery.value.trim()) query.q = searchQuery.value.trim();
        if (currentPage.value > 1) query.page = String(currentPage.value);
        if (sortField.value && sortField.value !== "last_updated") query.sort = sortField.value;
        if (sortOrder.value === "asc") query.order = "asc";
        if (selectedProviders.value.length)
            query.providers = selectedProviders.value.map((p: any) => p.value ?? p).join(",");
        if (filters.reasoning) query.reasoning = "true";
        if (filters.vision) query.vision = "true";
        if (filters.freeOnly) query.freeOnly = "true";
        if (selectedInputTypes.value.length)
            query.inputTypes = selectedInputTypes.value.map((t: any) => t.value ?? t).join(",");
        if (selectedOutputTypes.value.length)
            query.outputTypes = selectedOutputTypes.value.map((t: any) => t.value ?? t).join(",");
        if (priceRange.value[0]! > 0) query.priceMin = String(priceRange.value[0]);
        if (priceRange.value[1]! < 100) query.priceMax = String(priceRange.value[1]);
        return query;
    };

    const buildApiParams = () => {
        const params: Record<string, string> = {
            page: String(currentPage.value),
            pageSize: "50",
            sort: sortField.value,
        };
        if (sortOrder.value === "desc") params.order = "desc";
        if (searchQuery.value.trim()) params.q = searchQuery.value.trim();
        if (selectedProviders.value.length)
            params.providers = selectedProviders.value.map((p: any) => p.value ?? p).join(",");
        if (filters.reasoning) params.reasoning = "true";
        if (filters.vision) params.vision = "true";
        if (filters.freeOnly) params.freeOnly = "true";
        if (selectedInputTypes.value.length)
            params.inputTypes = selectedInputTypes.value.map((t: any) => t.value ?? t).join(",");
        if (selectedOutputTypes.value.length)
            params.outputTypes = selectedOutputTypes.value.map((t: any) => t.value ?? t).join(",");
        if (priceRange.value[0]! > 0 || priceRange.value[1]! < 100) {
            params.priceMin = String(priceRange.value[0]);
            params.priceMax = String(priceRange.value[1]);
        }
        return params;
    };

    const syncToUrl = () => {
        router.replace({ query: buildQuery() });
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

    const fetchCatalog = async () => {
        const res = await $fetch("/api/models", { params: buildApiParams() });
        catalogData.value = res as any;
    };

    // Filter/sort changes → sync to URL
    watch(
        [selectedProviders, selectedInputTypes, selectedOutputTypes, priceRange, filters],
        () => {
            currentPage.value = 1;
            syncToUrl();
        },
        { deep: true },
    );

    // Search query changes → sync to URL (debounced)
    let searchTimer: ReturnType<typeof setTimeout> | null = null;
    watch(searchQuery, () => {
        if (searchTimer) clearTimeout(searchTimer);
        searchTimer = setTimeout(() => {
            currentPage.value = 1;
            syncToUrl();
        }, 300);
    });

    // URL query changes → fetch data
    watch(
        () => route.query,
        () => {
            currentPage.value = Number(route.query.page) || 1;
            fetchCatalog();
        },
    );

    const goToPage = (page: number) => {
        if (page < 1 || page > catalogData.value.meta.totalPages) return;
        currentPage.value = page;
        syncToUrl();
        window.scrollTo({ top: 0, behavior: "smooth" });
    };

    const inputTypeItems = ["text", "image", "pdf", "video", "audio"].map((v) => ({
        label: v.charAt(0).toUpperCase() + v.slice(1),
        value: v,
    }));

    const outputTypeItems = ["text", "image", "video"].map((v) => ({
        label: v.charAt(0).toUpperCase() + v.slice(1),
        value: v,
    }));

    return {
        viewMode,
        searchQuery,
        sortField,
        sortOrder,
        selectedProviders,
        selectedInputTypes,
        selectedOutputTypes,
        priceRange,
        currentPage,
        filters,
        filterToggles,
        toggleFilter,
        catalogData,
        models,
        fetchCatalog,
        syncToUrl,
        toggleSort,
        goToPage,
        inputTypeItems,
        outputTypeItems,
    };
}
