<template>
    <!-- Row 1: Title + controls -->
    <div class="flex items-center justify-between gap-3 mb-3">
        <div class="shrink-0">
            <h1 class="text-xl font-bold text-default">{{ $t("catalog.title") }}</h1>
            <div class="text-sm text-muted">{{ total }} {{ $t("nav.catalog") }}</div>
        </div>
        <div class="flex items-center gap-2">
            <!-- Mobile: Filter toggle -->
            <button
                v-if="isMobile"
                @click="filtersOpen = !filtersOpen"
                class="relative flex items-center gap-1.5 rounded-lg px-2.5 py-1.5 text-xs cursor-pointer transition-colors border shrink-0"
                :class="
                    filtersOpen || activeFilterCount > 0
                        ? 'bg-primary/10 border-primary/30 text-primary'
                        : 'bg-default border-default text-toned hover:border-accented'
                "
            >
                <UIcon name="i-lucide-sliders-horizontal" class="size-3.5" />
                <span>{{ $t("catalog.filters") }}</span>
                <span
                    v-if="activeFilterCount > 0"
                    class="absolute -top-1.5 -right-1.5 bg-primary text-white text-[10px] font-bold rounded-full size-4 flex items-center justify-center"
                >
                    {{ activeFilterCount }}
                </span>
            </button>

            <!-- Search -->
            <UInput
                :model-value="search"
                @update:model-value="(v: string) => $emit('update:search', v)"
                :placeholder="$t('search.placeholder')"
                icon="i-lucide-search"
                class="w-36 sm:w-56"
                size="sm"
            />

            <!-- Desktop: View toggle + Column visibility -->
            <template v-if="!isMobile">
                <div class="flex bg-default border border-default rounded-lg p-0.5 shrink-0">
                    <button
                        @click="$emit('update:modelValue', 'grid')"
                        class="rounded-md px-2 py-1 cursor-pointer transition-colors"
                        :class="modelValue === 'grid' ? 'bg-primary text-white' : 'text-toned hover:text-default'"
                    >
                        <UIcon name="i-lucide-layout-grid" class="size-3.5" />
                    </button>
                    <button
                        @click="$emit('update:modelValue', 'table')"
                        class="rounded-md px-2 py-1 cursor-pointer transition-colors"
                        :class="modelValue === 'table' ? 'bg-primary text-white' : 'text-toned hover:text-default'"
                    >
                        <UIcon name="i-lucide-table" class="size-3.5" />
                    </button>
                </div>
                <UDropdownMenu
                    v-if="columnMenuItems.length"
                    :items="[columnMenuItems]"
                    :content="{ align: 'end', class: 'max-h-80 overflow-y-auto' }"
                >
                    <UButton icon="i-lucide-columns-3-cog" color="neutral" variant="outline" size="md" />
                </UDropdownMenu>
            </template>
        </div>
    </div>

    <!-- Desktop: Filter row -->
    <div v-if="!isMobile" class="flex items-center gap-2 mb-4 flex-wrap">
        <button
            v-for="toggle in filterToggles"
            :key="toggle.key"
            @click="$emit('toggleFilter', toggle.key)"
            class="flex items-center gap-1 rounded-lg px-2.5 py-1 text-xs cursor-pointer transition-colors shrink-0"
            :class="
                filters[toggle.key]
                    ? 'bg-primary/10 border border-primary/30 text-primary'
                    : 'bg-default border border-default text-toned hover:border-accented'
            "
        >
            <UIcon :name="toggle.icon" class="size-3.5" />
            <span>{{ toggle.label }}</span>
        </button>

        <div class="w-px h-5 bg-default shrink-0" />

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

        <USelectMenu
            v-model="selectedInputTypes"
            multiple
            clear
            :items="inputTypeItems"
            :placeholder="$t('catalog.inputType')"
            class="w-36 shrink-0"
            size="sm"
        />
        <USelectMenu
            v-model="selectedOutputTypes"
            multiple
            clear
            :items="outputTypeItems"
            :placeholder="$t('catalog.outputType')"
            class="w-36 shrink-0"
            size="sm"
        />

        <div class="flex items-center gap-2 bg-default border border-default rounded-lg px-3 py-1 shrink-0">
            <span class="text-xs text-muted whitespace-nowrap">{{ $t("catalog.priceRange") }}</span>
            <USlider v-model="priceRange" :min="0" :max="100" :step="0.1" color="primary" class="w-24" size="sm" />
            <span class="text-xs text-toned font-mono whitespace-nowrap tabular-nums">
                ${{ priceRange?.[0] ?? 0 }}–${{ priceRange?.[1] ?? 100 }}
            </span>
        </div>
    </div>

    <!-- Mobile: Expandable filter panel -->
    <Transition name="expand" @enter="onExpandEnter" @after-enter="onExpandAfterEnter" @leave="onExpandLeave">
        <div v-if="isMobile && filtersOpen" class="mb-4 space-y-2.5">
            <div class="flex items-center gap-2 flex-wrap">
                <button
                    v-for="toggle in filterToggles"
                    :key="toggle.key"
                    @click="$emit('toggleFilter', toggle.key)"
                    class="flex items-center gap-1 rounded-lg px-2.5 py-1 text-xs cursor-pointer transition-colors"
                    :class="
                        filters[toggle.key]
                            ? 'bg-primary/10 border border-primary/30 text-primary'
                            : 'bg-default border border-default text-toned hover:border-accented'
                    "
                >
                    <UIcon :name="toggle.icon" class="size-3.5" />
                    <span>{{ toggle.label }}</span>
                </button>
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

            <div class="flex items-center gap-2">
                <USelectMenu
                    v-model="selectedInputTypes"
                    multiple
                    clear
                    :items="inputTypeItems"
                    :placeholder="$t('catalog.inputType')"
                    class="flex-1"
                    size="sm"
                />
                <USelectMenu
                    v-model="selectedOutputTypes"
                    multiple
                    clear
                    :items="outputTypeItems"
                    :placeholder="$t('catalog.outputType')"
                    class="flex-1"
                    size="sm"
                />
            </div>

            <div class="flex items-center gap-2 bg-default border border-default rounded-lg px-3 py-2">
                <span class="text-xs text-muted whitespace-nowrap">{{ $t("catalog.priceRange") }}</span>
                <USlider v-model="priceRange" :min="0" :max="100" :step="0.1" color="primary" class="flex-1" size="sm" />
                <span class="text-xs text-toned font-mono whitespace-nowrap tabular-nums">
                    ${{ priceRange?.[0] ?? 0 }}–${{ priceRange?.[1] ?? 100 }}
                </span>
            </div>

            <div class="flex items-center justify-end">
                <div class="flex bg-default border border-default rounded-lg p-0.5">
                    <button
                        @click="$emit('update:modelValue', 'grid')"
                        class="rounded-md px-2 py-1 cursor-pointer transition-colors"
                        :class="modelValue === 'grid' ? 'bg-primary text-white' : 'text-toned hover:text-default'"
                    >
                        <UIcon name="i-lucide-layout-grid" class="size-3.5" />
                    </button>
                    <button
                        @click="$emit('update:modelValue', 'table')"
                        class="rounded-md px-2 py-1 cursor-pointer transition-colors"
                        :class="modelValue === 'table' ? 'bg-primary text-white' : 'text-toned hover:text-default'"
                    >
                        <UIcon name="i-lucide-table" class="size-3.5" />
                    </button>
                </div>
            </div>
        </div>
    </Transition>
</template>

<script setup lang="ts">
    const props = defineProps<{
        modelValue: "grid" | "table";
        search: string;
        total: number;
        columnMenuItems: any[];
        filterToggles: { key: "freeOnly" | "reasoning" | "vision"; label: string; icon: string }[];
        filters: { freeOnly: boolean; reasoning: boolean; vision: boolean };
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

    defineEmits<{
        "update:modelValue": [value: "grid" | "table"];
        "update:search": [value: string];
        toggleFilter: [key: "freeOnly" | "reasoning" | "vision"];
        toggleProvider: [provider: any];
        removeProvider: [provider: any];
        clearProviders: [];
    }>();

    // Mobile detection
    const isMobile = ref(false);
    const filtersOpen = ref(false);
    onMounted(() => {
        isMobile.value = window.innerWidth < 768;
        const mq = window.matchMedia("(min-width: 768px)");
        const handler = (e: MediaQueryListEvent) => {
            isMobile.value = !e.matches;
            if (!isMobile.value) filtersOpen.value = false;
        };
        mq.addEventListener("change", handler);
        onUnmounted(() => mq.removeEventListener("change", handler));
    });

    // Active filter count for mobile badge
    const activeFilterCount = computed(() => {
        let count = 0;
        if (props.filters.freeOnly) count++;
        if (props.filters.reasoning) count++;
        if (props.filters.vision) count++;
        if (props.selectedProviders?.length) count++;
        if (selectedInputTypes.value?.length) count++;
        if (selectedOutputTypes.value?.length) count++;
        const pr = priceRange.value || [0, 100];
        if (pr[0] !== 0 || pr[1] !== 100) count++;
        return count;
    });

    // Height animation helpers
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
