import type { Provider, SelectItem, ApiResponse } from "~/types";

export async function useProviderFilter(selectedProviders: Ref<SelectItem[]>) {
    const { t } = useI18n()
    const route = useRoute()
    const config = useRuntimeConfig()

    const { data: providers } = await useAsyncData('providers', () =>
        $fetch<ApiResponse<Provider[]>>(`${config.public.apiBase}/api/v1/providers`),
    )

    const providersList = computed(() => providers.value?.data || [])

    const providerSearch = ref('')

    const topProviders = computed(() =>
        [...providersList.value].sort((a: Provider, b: Provider) => b.modelCount - a.modelCount).slice(0, 10),
    )

    const groupedProviders = computed(() => {
        const groups: Record<string, Provider[]> = {}
        for (const p of providersList.value) {
            const letter = p.name.charAt(0).toUpperCase()
            if (!groups[letter]) groups[letter] = []
            groups[letter].push(p)
        }
        return Object.keys(groups)
            .sort()
            .map((letter) => ({ letter, items: groups[letter] }))
    })

    const filteredProviders = computed(() => {
        if (!providerSearch.value) return []
        const q = providerSearch.value.toLowerCase()
        return providersList.value.filter((p: Provider) => p.name.toLowerCase().includes(q))
    })

    const providerTriggerLabel = computed(() => {
        const count = selectedProviders.value.length
        return count === 0 ? t('catalog.providers') : `${t('catalog.providers')} (${count})`
    })

    const isSelected = (provider: { id: string }) =>
        selectedProviders.value.some((sp: SelectItem | string) => (typeof sp === 'string' ? sp : sp.value) === provider.id)

    const toggleProvider = (provider: { id: string; name: string }) => {
        if (isSelected(provider)) {
            selectedProviders.value = selectedProviders.value.filter(
                (sp: SelectItem | string) => (typeof sp === 'string' ? sp : sp.value) !== provider.id,
            )
        } else {
            selectedProviders.value.push({ label: provider.name, value: provider.id })
        }
    }

    const removeProvider = (sp: SelectItem | string) => {
        const spValue = typeof sp === 'string' ? sp : sp.value;
        selectedProviders.value = selectedProviders.value.filter(
            (p: SelectItem | string) => (typeof p === 'string' ? p : p.value) !== spValue,
        )
    }

    const clearProviders = () => {
        selectedProviders.value = []
    }

    // Restore selected providers from URL immediately
    const urlProviders = (route.query.providers as string)?.split(',').filter(Boolean) || []
    if (urlProviders.length && providersList.value.length) {
        selectedProviders.value = urlProviders
            .map((id) => {
                const p = providersList.value.find((pr: Provider) => pr.id === id)
                if (!p) return null
                return { label: p.name, value: p.id } as SelectItem
            })
            .filter((item): item is SelectItem => item !== null)
    }

    return {
        providerSearch,
        topProviders,
        groupedProviders,
        filteredProviders,
        providerTriggerLabel,
        isSelected,
        toggleProvider,
        removeProvider,
        clearProviders,
    }
}
