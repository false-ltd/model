import type { Provider, ApiResponse } from "~/types";

export function useProviders() {
    const config = useRuntimeConfig();

    // Session-scoped so the view preference survives page navigation.
    const viewMode = useState<"grid" | "table">("providers-view", () => "grid");
    const search = ref("");

    const { data: result, status, error, execute } = useAsyncData("providers-list", () =>
        $fetch<ApiResponse<Provider[]>>(`${config.public.apiBase}/api/v1/providers`),
    );

    const loading = computed(() => status.value === "pending");
    const loadError = computed(() => error.value != null);

    const providersList = computed(() => result.value?.data || []);

    const filtered = computed(() => {
        if (!search.value.trim()) return providersList.value;
        const q = search.value.toLowerCase();
        return providersList.value.filter(
            (p: Provider) =>
                p.name.toLowerCase().includes(q) ||
                p.id.toLowerCase().includes(q) ||
                p.npm?.toLowerCase().includes(q),
        );
    });

    return { viewMode, search, providersList, filtered, loading, loadError, retry: () => execute() };
}
