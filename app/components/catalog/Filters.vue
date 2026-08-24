<template>
    <!-- Mobile: Title + filter toggle -->
    <template v-if="isMobile">
        <div class="flex items-center justify-between gap-3 mb-4">
            <div>
                <h1 class="text-xl font-bold text-default">{{ $t("catalog.title") }}</h1>
                <div class="text-sm text-muted">{{ total }} {{ $t("nav.catalog") }}</div>
            </div>
            <div class="flex items-center gap-2">
                <UButton
                    @click="filtersOpen = !filtersOpen"
                    icon="i-lucide-sliders-horizontal"
                    :label="$t('catalog.filters')"
                    :color="filtersOpen || activeFilterCount > 0 ? 'primary' : 'neutral'"
                    :variant="filtersOpen || activeFilterCount > 0 ? 'soft' : 'outline'"
                    size="sm"
                >
                    <template v-if="activeFilterCount > 0" #trailing>
                        <UBadge :label="String(activeFilterCount)" color="primary" variant="solid" size="xs" />
                    </template>
                </UButton>
            </div>
        </div>
        <div ref="searchWrap">
            <UInput
                :model-value="search"
                @update:model-value="(v: string) => emit('update:search', v)"
                :placeholder="$t('search.placeholder')"
                icon="i-lucide-search"
                size="lg"
                class="w-full mb-4"
            />
        </div>
    </template>

    <!-- Desktop: Row 1 — Title + Search + Column visibility -->
    <template v-else>
        <div class="flex items-center gap-3 mb-4">
            <h1 class="text-xl font-bold text-default whitespace-nowrap">
                {{ $t("catalog.title") }}
                <span class="text-sm font-normal text-muted ml-1">{{ total }} {{ $t("nav.catalog") }}</span>
            </h1>
            <div ref="searchWrap">
            <UInput
                :model-value="search"
                @update:model-value="(v: string) => emit('update:search', v)"
                @keydown.escape="emit('update:search', '')"
                :placeholder="$t('search.placeholder')"
                icon="i-lucide-search"
                size="lg"
                class="flex-1 search-input"
            >
                <template v-if="search" #trailing>
                    <UButton
                        icon="i-lucide-x"
                        color="neutral"
                        variant="ghost"
                        size="xs"
                        @click="emit('update:search', '')"
                    />
                </template>
            </UInput>
            </div>
            <UDropdownMenu
                v-if="columnMenuItems.length"
                :items="[columnMenuItems]"
                :content="{ align: 'end' }"
            >
                <UButton icon="i-lucide-columns-3-cog" color="neutral" variant="outline" size="md" />
            </UDropdownMenu>
        </div>

    <!-- Desktop: Row 2 — All filters in one line -->
    <div class="flex flex-wrap items-center gap-2 mb-4">
        <UButton
            v-for="toggle in filterToggles"
            :key="toggle.key"
            @click="$emit('toggleFilter', toggle.key)"
            :icon="toggle.icon"
            :label="toggle.label"
            :color="filters[toggle.key] ? toggle.color : 'neutral'"
            :variant="filters[toggle.key] ? 'soft' : 'outline'"
            size="sm"
        />

            <span class="text-muted text-sm select-none">|</span>

            <CatalogProviderPopover
                v-model:providerSearch="providerSearch"
                :selected-providers="selectedProviders"
                :top-providers="topProviders"
                :grouped-providers="groupedProviders"
                :filtered-providers="filteredProviders"
                :provider-trigger-label="providerTriggerLabel"
                :is-selected="isSelected"
                @toggle-provider="$emit('toggleProvider', $event)"
                @remove-provider="$emit('removeProvider', $event)"
                @clear-providers="$emit('clearProviders')"
            />

            <!-- Price Popover -->
            <UPopover>
                <UButton
                    icon="i-lucide-dollar-sign"
                    :label="$t('catalog.price')"
                    :color="priceActive ? 'primary' : 'neutral'"
                    :variant="priceActive ? 'soft' : 'outline'"
                    size="sm"
                />
                <template #content>
                    <div class="p-3 space-y-3 w-80">
                        <div>
                            <div class="flex items-center justify-between mb-1.5">
                                <span class="text-xs text-muted">{{ $t("catalog.colInputCost") }}</span>
                                <span class="text-xs text-toned font-mono tabular-nums">
                                    ${{ priceRange?.[0] ?? 0 }}–${{ priceRange?.[1] ?? 100 }}
                                </span>
                            </div>
                            <USlider v-model="priceRange" :min="0" :max="100" :step="0.1" color="primary" size="sm" />
                        </div>
                        <div>
                            <div class="flex items-center justify-between mb-1.5">
                                <span class="text-xs text-muted">{{ $t("catalog.colOutputCost") }}</span>
                                <span class="text-xs text-toned font-mono tabular-nums">
                                    ${{ outputPriceRange?.[0] ?? 0 }}–${{ outputPriceRange?.[1] ?? 100 }}
                                </span>
                            </div>
                            <USlider v-model="outputPriceRange" :min="0" :max="100" :step="0.1" color="primary" size="sm" />
                        </div>
                    </div>
                </template>
            </UPopover>

            <!-- I/O Types Popover -->
            <UPopover>
                <UButton
                    icon="i-lucide-arrow-left-right"
                    :label="$t('catalog.ioTypes')"
                    :color="ioTypeCount ? 'primary' : 'neutral'"
                    :variant="ioTypeCount ? 'soft' : 'outline'"
                    size="sm"
                >
                    <template v-if="ioTypeCount" #trailing>
                        <UBadge :label="String(ioTypeCount)" color="primary" variant="solid" size="xs" />
                    </template>
                </UButton>
                <template #content>
                    <div class="p-2 w-72">
                        <div class="grid grid-cols-2 gap-3">
                            <div>
                                <div class="text-[10px] text-muted uppercase tracking-wider font-medium mb-1.5">{{ $t('catalog.inputType') }}</div>
                                <div class="space-y-1">
                                    <button
                                        v-for="item in inputTypeItems"
                                        :key="item.value"
                                        @click="toggleInputType(item.value)"
                                        class="flex items-center gap-1.5 rounded-lg px-2 py-1.5 text-xs cursor-pointer transition-colors w-full"
                                        :class="isTypeSelected(selectedInputTypes, item.value) ? modalityBtnClass(item.value) : 'bg-default border border-default text-toned hover:border-accented'"
                                    >
                                        <UIcon :name="modalityIcon(item.value)" class="size-3.5" />
                                        <span>{{ item.label }}</span>
                                    </button>
                                </div>
                            </div>
                            <div>
                                <div class="text-[10px] text-muted uppercase tracking-wider font-medium mb-1.5">{{ $t('catalog.outputType') }}</div>
                                <div class="space-y-1">
                                    <button
                                        v-for="item in outputTypeItems"
                                        :key="item.value"
                                        @click="toggleOutputType(item.value)"
                                        class="flex items-center gap-1.5 rounded-lg px-2 py-1.5 text-xs cursor-pointer transition-colors w-full"
                                        :class="isTypeSelected(selectedOutputTypes, item.value) ? modalityBtnClass(item.value) : 'bg-default border border-default text-toned hover:border-accented'"
                                    >
                                        <UIcon :name="modalityIcon(item.value)" class="size-3.5" />
                                        <span>{{ item.label }}</span>
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>
                </template>
            </UPopover>

            <!-- Clear All (when 2+ filters active) -->
            <template v-if="activeFilterCount > 1">
                <span class="text-muted text-sm select-none">|</span>
                <UButton
                    icon="i-lucide-x"
                    :label="$t('common.clearAll')"
                    color="error"
                    variant="ghost"
                    size="xs"
                    @click="clearAllFilters"
                />
            </template>
        </div>

        <!-- Active filter chips: the current view at a glance, each removable -->
        <div v-if="activeChips.length" class="flex flex-wrap items-center gap-1.5 mb-4" data-lenis-prevent>
            <span class="text-[11px] text-muted uppercase tracking-wider font-medium mr-0.5">{{ $t("catalog.activeFilters") }}</span>
            <button
                v-for="chip in activeChips"
                :key="chip.key"
                class="group inline-flex items-center gap-1 rounded-full border border-default bg-elevated pl-2.5 pr-1 py-0.5 text-xs text-toned cursor-pointer hover:border-error/50 hover:text-default transition-colors"
                @click="chip.remove"
            >
                <UIcon v-if="chip.icon" :name="chip.icon" class="size-3 text-primary" />
                <span class="max-w-40 truncate">{{ chip.label }}</span>
                <UIcon name="i-lucide-x" class="size-3 text-muted group-hover:text-error transition-colors" />
            </button>
            <button
                v-if="activeChips.length > 1"
                class="text-xs text-muted hover:text-error cursor-pointer transition-colors ml-1"
                @click="clearAllFilters"
            >
                {{ $t("common.clearAll") }}
            </button>
        </div>
    </template>

    <!-- Mobile: Expandable filter panel -->
    <Transition name="expand" @enter="onExpandEnter" @after-enter="onExpandAfterEnter" @leave="onExpandLeave">
        <div v-if="isMobile && filtersOpen" class="mb-4 space-y-2.5">
            <div class="flex flex-wrap gap-2">
                <UButton
                    v-for="toggle in filterToggles"
                    :key="toggle.key"
                    @click="$emit('toggleFilter', toggle.key)"
                    :icon="toggle.icon"
                    :label="toggle.label"
                    :color="filters[toggle.key] ? toggle.color : 'neutral'"
                    :variant="filters[toggle.key] ? 'soft' : 'outline'"
                    size="xs"
                />
            </div>

            <CatalogProviderPopover
                v-model:providerSearch="providerSearch"
                :selected-providers="selectedProviders"
                :top-providers="topProviders"
                :grouped-providers="groupedProviders"
                :filtered-providers="filteredProviders"
                :provider-trigger-label="providerTriggerLabel"
                :is-selected="isSelected"
                @toggle-provider="$emit('toggleProvider', $event)"
                @remove-provider="$emit('removeProvider', $event)"
                @clear-providers="$emit('clearProviders')"
            />

            <div class="grid grid-cols-2 gap-3">
                <div>
                    <div class="text-[10px] text-muted uppercase tracking-wider font-medium mb-1.5">{{ $t('catalog.inputType') }}</div>
                    <div class="space-y-1">
                        <button
                            v-for="item in inputTypeItems"
                            :key="item.value"
                            @click="toggleInputType(item.value)"
                            class="flex items-center gap-1.5 rounded-lg px-2 py-1.5 text-xs cursor-pointer transition-colors w-full"
                            :class="isTypeSelected(selectedInputTypes, item.value) ? modalityBtnClass(item.value) : 'bg-default border border-default text-toned hover:border-accented'"
                        >
                            <UIcon :name="modalityIcon(item.value)" class="size-3.5" />
                            <span>{{ item.label }}</span>
                        </button>
                    </div>
                </div>
                <div>
                    <div class="text-[10px] text-muted uppercase tracking-wider font-medium mb-1.5">{{ $t('catalog.outputType') }}</div>
                    <div class="space-y-1">
                        <button
                            v-for="item in outputTypeItems"
                            :key="item.value"
                            @click="toggleOutputType(item.value)"
                            class="flex items-center gap-1.5 rounded-lg px-2 py-1.5 text-xs cursor-pointer transition-colors w-full"
                            :class="isTypeSelected(selectedOutputTypes, item.value) ? modalityBtnClass(item.value) : 'bg-default border border-default text-toned hover:border-accented'"
                        >
                            <UIcon :name="modalityIcon(item.value)" class="size-3.5" />
                            <span>{{ item.label }}</span>
                        </button>
                    </div>
                </div>
            </div>

            <div class="flex items-center gap-2 bg-default border border-default rounded-lg px-3 py-2">
                <span class="text-[10px] text-muted whitespace-nowrap">{{ $t("catalog.colInputCost") }}</span>
                <USlider v-model="priceRange" :min="0" :max="100" :step="0.1" color="primary" class="flex-1" size="sm" />
                <span class="text-[10px] text-toned font-mono whitespace-nowrap tabular-nums">
                    ${{ priceRange?.[0] ?? 0 }}–${{ priceRange?.[1] ?? 100 }}
                </span>
            </div>
            <div class="flex items-center gap-2 bg-default border border-default rounded-lg px-3 py-2">
                <span class="text-[10px] text-muted whitespace-nowrap">{{ $t("catalog.colOutputCost") }}</span>
                <USlider v-model="outputPriceRange" :min="0" :max="100" :step="0.1" color="primary" class="flex-1" size="sm" />
                <span class="text-[10px] text-toned font-mono whitespace-nowrap tabular-nums">
                    ${{ outputPriceRange?.[0] ?? 0 }}–${{ outputPriceRange?.[1] ?? 100 }}
                </span>
            </div>
        </div>
    </Transition>
</template>

<style scoped>
.search-input :deep(input:focus) {
    box-shadow: 0 0 0 2px var(--ui-bg-elevated), 0 0 0 4px var(--ui-border-accented);
    border-color: var(--ui-border-accented);
}
</style>

<script setup lang="ts">
    const props = defineProps<{
        search: string;
        total: number;
        columnMenuItems: any[];
        filterToggles: { key: string; label: string; icon: string; color: string }[];
        filters: Record<string, boolean>;
        selectedProviders: any[];
        topProviders: any[];
        groupedProviders: { letter: string; items?: any[] }[];
        filteredProviders: any[];
        providerTriggerLabel: string;
        isSelected: (p: { id: string }) => boolean;
        inputTypeItems: { label: string; value: string }[];
        outputTypeItems: { label: string; value: string }[];
    }>();

    const providerSearch = defineModel<string>("providerSearch", { required: true });
    const selectedInputTypes = defineModel<any[]>("selectedInputTypes", { required: true });
    const selectedOutputTypes = defineModel<any[]>("selectedOutputTypes", { required: true });
    const priceRange = defineModel<number[]>("priceRange", { default: () => [0, 100] });
    const outputPriceRange = defineModel<number[]>("outputPriceRange", { default: () => [0, 100] });

    const emit = defineEmits<{
        "update:search": [value: string];
        toggleFilter: [key: string];
        toggleProvider: [provider: any];
        removeProvider: [provider: any];
        clearProviders: [];
    }>();

    const { isMobile } = useMobile();
    const { t: t2 } = useI18n();
    const filtersOpen = ref(false);

    const modalityBtnClass = (type: string) => modalityClass(type);
    const isTypeSelected = (list: any[], value: string) => list?.some((t: any) => (t.value ?? t) === value);
    const toggleInputType = (value: string) => {
        const exists = selectedInputTypes.value?.some((t: any) => (t.value ?? t) === value);
        if (exists) {
            selectedInputTypes.value = selectedInputTypes.value.filter((t: any) => (t.value ?? t) !== value);
        } else {
            selectedInputTypes.value = [...(selectedInputTypes.value || []), { label: value.charAt(0).toUpperCase() + value.slice(1), value }];
        }
    };
    const toggleOutputType = (value: string) => {
        const exists = selectedOutputTypes.value?.some((t: any) => (t.value ?? t) === value);
        if (exists) {
            selectedOutputTypes.value = selectedOutputTypes.value.filter((t: any) => (t.value ?? t) !== value);
        } else {
            selectedOutputTypes.value = [...(selectedOutputTypes.value || []), { label: value.charAt(0).toUpperCase() + value.slice(1), value }];
        }
    };

    const priceActive = computed(() => {
        const pr = priceRange.value || [0, 100];
        const opr = outputPriceRange.value || [0, 100];
        return pr[0] !== 0 || pr[1] !== 100 || opr[0] !== 0 || opr[1] !== 100;
    });
    const ioTypeCount = computed(() => (selectedInputTypes.value?.length || 0) + (selectedOutputTypes.value?.length || 0));

    const activeFilterCount = computed(() => {
        let count = props.search ? 1 : 0;
        count += Object.values(props.filters).filter(Boolean).length;
        count += props.selectedProviders?.length || 0;
        count += (selectedInputTypes.value?.length || 0);
        count += (selectedOutputTypes.value?.length || 0);
        if (priceActive.value) count++;
        return count;
    });

    interface FilterChip {
        key: string;
        icon?: string;
        label: string;
        remove: () => void;
    }

    const activeChips = computed<FilterChip[]>(() => {
        const chips: FilterChip[] = [];
        if (props.search.trim()) {
            chips.push({
                key: "search",
                icon: "i-lucide-search",
                label: props.search.trim(),
                remove: () => emit("update:search", ""),
            });
        }
        for (const toggle of props.filterToggles) {
            if (props.filters[toggle.key]) {
                chips.push({
                    key: `f-${toggle.key}`,
                    icon: toggle.icon,
                    label: toggle.label,
                    remove: () => emit("toggleFilter", toggle.key),
                });
            }
        }
        for (const sp of props.selectedProviders) {
            chips.push({
                key: `p-${typeof sp === "string" ? sp : sp.value}`,
                icon: "i-lucide-building-2",
                label: typeof sp === "string" ? sp : sp.label,
                remove: () => emit("removeProvider", sp),
            });
        }
        for (const it of selectedInputTypes.value || []) {
            const v = typeof it === "string" ? it : it.value;
            chips.push({
                key: `in-${v}`,
                icon: modalityIcon(v),
                label: `${t2("catalog.inputType")}: ${v}`,
                remove: () => toggleInputType(v),
            });
        }
        for (const ot of selectedOutputTypes.value || []) {
            const v = typeof ot === "string" ? ot : ot.value;
            chips.push({
                key: `out-${v}`,
                icon: modalityIcon(v),
                label: `${t2("catalog.outputType")}: ${v}`,
                remove: () => toggleOutputType(v),
            });
        }
        const pr = priceRange.value || [0, 100];
        if (pr[0] !== 0 || pr[1] !== 100) {
            chips.push({
                key: "price-in",
                icon: "i-lucide-dollar-sign",
                label: `$${pr[0]}–$${pr[1]}`,
                remove: () => (priceRange.value = [0, 100]),
            });
        }
        const opr = outputPriceRange.value || [0, 100];
        if (opr[0] !== 0 || opr[1] !== 100) {
            chips.push({
                key: "price-out",
                icon: "i-lucide-coins",
                label: `${t2("catalog.colOutputCost")} $${opr[0]}–$${opr[1]}`,
                remove: () => (outputPriceRange.value = [0, 100]),
            });
        }
        return chips;
    });

    // "/" focuses the catalog search from anywhere on the page.
    const searchWrap = useTemplateRef<HTMLElement>("searchWrap");
    onMounted(() => {
        const onKeydown = (e: KeyboardEvent) => {
            if (e.key !== "/" || e.metaKey || e.ctrlKey || e.altKey) return;
            const el = e.target as HTMLElement;
            if (el && (el.tagName === "INPUT" || el.tagName === "TEXTAREA" || el.isContentEditable)) return;
            e.preventDefault();
            searchWrap.value?.querySelector("input")?.focus();
        };
        window.addEventListener("keydown", onKeydown);
        onUnmounted(() => window.removeEventListener("keydown", onKeydown));
    });

    const clearAllFilters = () => {
        emit("update:search", "");
        (Object.keys(props.filters) as (keyof typeof props.filters)[]).forEach((key) => {
            if (props.filters[key]) emit("toggleFilter", key);
        });
        emit("clearProviders");
        selectedInputTypes.value = [];
        selectedOutputTypes.value = [];
        priceRange.value = [0, 100];
        outputPriceRange.value = [0, 100];
    };

    onMounted(() => {
        watch(isMobile, (mobile) => {
            if (!mobile) filtersOpen.value = false;
        });
    });

    function onExpandEnter(el: Element) {
        const htmlEl = el as HTMLElement;
        htmlEl.style.height = "0";
        htmlEl.style.overflow = "hidden";
        void htmlEl.offsetHeight;
        htmlEl.style.height = htmlEl.scrollHeight + "px";
    }
    function onExpandAfterEnter(el: Element) {
        const htmlEl = el as HTMLElement;
        htmlEl.style.height = "auto";
        htmlEl.style.overflow = "";
    }
    function onExpandLeave(el: Element) {
        const htmlEl = el as HTMLElement;
        htmlEl.style.height = htmlEl.scrollHeight + "px";
        htmlEl.style.overflow = "hidden";
        void htmlEl.offsetHeight;
        htmlEl.style.height = "0";
    }
</script>
