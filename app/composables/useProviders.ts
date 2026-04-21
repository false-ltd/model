import type { Provider, ApiResponse } from "~/types";

export function useProviders() {
    const { t } = useI18n();
    const config = useRuntimeConfig();

    const viewMode = ref<"grid" | "table">("grid");
    const search = ref("");

    const { data: result } = useAsyncData("providers-list", () =>
        $fetch<ApiResponse<Provider[]>>(`${config.public.apiBase}/api/v1/providers`),
    );

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

    return { viewMode, search, providersList, filtered };
}
