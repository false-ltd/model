export async function useProviderFilter(selectedProviders: Ref<any[]>) {
    const { t } = useI18n()
    const route = useRoute()

    const { data: providers } = await useAsyncData('providers', () =>
        $fetch<{ data: { id: string; name: string; modelCount: number }[] }>('/api/providers'),
    )

    const providersList = computed(() => (providers.value as any)?.data || [])

    const providerSearch = ref('')

    const topProviders = computed(() =>
        [...providersList.value].sort((a: any, b: any) => b.modelCount - a.modelCount).slice(0, 10),
    )

    const groupedProviders = computed(() => {
        const groups: Record<string, any[]> = {}
        for (const p of providersList.value) {
            const letter = (p as any).name.charAt(0).toUpperCase()
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
        return providersList.value.filter((p: any) => p.name.toLowerCase().includes(q))
    })

    const providerTriggerLabel = computed(() => {
        const count = selectedProviders.value.length
        return count === 0 ? t('catalog.providers') : `${t('catalog.providers')} (${count})`
    })

    const isSelected = (provider: { id: string }) =>
        selectedProviders.value.some((sp: any) => (sp.value ?? sp) === provider.id)

    const toggleProvider = (provider: { id: string; name: string }) => {
        if (isSelected(provider)) {
            selectedProviders.value = selectedProviders.value.filter(
                (sp: any) => (sp.value ?? sp) !== provider.id,
            )
        } else {
            selectedProviders.value.push({ label: provider.name, value: provider.id })
        }
    }

    const removeProvider = (sp: any) => {
        selectedProviders.value = selectedProviders.value.filter(
            (p: any) => (p.value ?? p) !== (sp.value ?? sp),
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
                const p = providersList.value.find((pr: any) => pr.id === id)
                if (!p) return null
                return { label: (p as any).name, value: (p as any).id }
            })
            .filter(Boolean) as any[]
    }

    return {
        providersList,
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
