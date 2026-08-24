<template>
    <UModal v-model:open="open" :ui="{ content: 'max-w-xl top-[15%] translate-y-0' }">
        <template #content>
            <div data-lenis-prevent>
                <UCommandPalette
                    v-model:search-term="query"
                    :groups="groups"
                    :placeholder="t('search.placeholder')"
                    :ui="{ input: 'h-12 text-base' }"
                    @update:model-value="onSelect"
                />
            </div>
        </template>
    </UModal>
</template>

<script setup lang="ts">
    import { watchDebounced } from "@vueuse/core";

    /**
     * Global ⌘K / Ctrl+K command palette: page navigation, quick actions
     * (theme toggle, saved filters), provider search, and live model
     * search. Nuxt UI v4 groups only take static item arrays, so async
     * search is driven by watching the palette's own search term.
     */
    const open = defineModel<boolean>("open", { default: false });

    const { t } = useI18n();
    const localePath = useLocalePath();
    const config = useRuntimeConfig();
    const colorMode = useColorMode();

    const query = ref("");
    const searchResults = ref<any[]>([]);
    let searchSeq = 0;

    const navGroup = computed(() => ({
        id: "pages",
        label: t("palette.pages"),
        items: [
            { label: t("nav.overview"), icon: "i-lucide-layout-dashboard", to: localePath("/") },
            { label: t("nav.catalog"), icon: "i-lucide-library", to: localePath("/catalog") },
            { label: t("nav.providers"), icon: "i-lucide-building-2", to: localePath("/providers") },
            { label: t("nav.compare"), icon: "i-lucide-git-compare", to: localePath("/compare") },
        ],
    }));

    const actionGroup = computed(() => ({
        id: "actions",
        label: t("palette.actions"),
        items: [
            {
                label: t("palette.toggleTheme"),
                icon: "i-lucide-moon-star",
                onSelect: () => {
                    colorMode.preference = colorMode.value === "dark" ? "light" : "dark";
                },
            },
            {
                label: t("palette.freeModels"),
                icon: "i-lucide-badge-dollar-sign",
                to: localePath("/catalog?freeOnly=true"),
            },
            {
                label: t("palette.openWeightsModels"),
                icon: "i-lucide-unlock",
                to: localePath("/catalog?openWeights=true"),
            },
        ],
    }));

    const { data: providersRaw } = useAsyncData("providers-list", () =>
        $fetch<any>(`${config.public.apiBase}/api/v1/providers`),
    );

    const providerGroup = computed(() => {
        const list: any[] = providersRaw.value?.data || [];
        return {
            id: "providers",
            label: t("palette.providers"),
            items: list.slice(0, 40).map((p: any) => ({
                label: p.name,
                suffix: String(p.model_count),
                icon: "i-lucide-building-2",
                to: localePath(`/catalog?providers=${p.id}`),
            })),
        };
    });

    const modelGroup = computed(() => ({
        id: "models",
        label: t("palette.models"),
        ignoreFilter: true,
        highlightedIcon: "i-lucide-arrow-right",
        // Results are already filtered server-side by q.
        items: searchResults.value.map((m) => ({
            label: m.name,
            suffix: m.model_id,
            icon: "i-lucide-box",
            to: localePath(`/model/${m.id}`),
        })),
    }));

    const groups = computed(() => {
        const g = [actionGroup.value, navGroup.value];
        if (providerGroup.value.items.length) g.push(providerGroup.value);
        if (modelGroup.value.items.length) g.push(modelGroup.value);
        return g;
    });

    watchDebounced(
        query,
        async (q) => {
            const term = q.trim();
            if (term.length < 2) {
                searchSeq++;
                searchResults.value = [];
                return;
            }
            const seq = ++searchSeq;
            try {
                const res = await $fetch<{ data: any[] }>(`${config.public.apiBase}/api/v1/models`, {
                    params: { q: term, page_size: 8 },
                });
                if (seq !== searchSeq) return; // stale response
                searchResults.value = res.data || [];
            } catch {
                if (seq === searchSeq) searchResults.value = [];
            }
        },
        { debounce: 250 },
    );

    const onSelect = () => {
        open.value = false;
    };

    watch(open, (v) => {
        if (!v) {
            query.value = "";
        }
    });

    onMounted(() => {
        const onKeydown = (e: KeyboardEvent) => {
            if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === "k") {
                e.preventDefault();
                open.value = !open.value;
            }
        };
        window.addEventListener("keydown", onKeydown);
        onUnmounted(() => window.removeEventListener("keydown", onKeydown));
    });
</script>
