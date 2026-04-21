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
                @update:model-value="(v: string) => emit('update:search', v)"
                :placeholder="$t('search.placeholder')"
                icon="i-lucide-search"
                class="w-40 sm:w-60"
            />

            <!-- Desktop: Capabilities Popover -->
            <UPopover v-if="!isMobile">
                <button :class="triggerClass(!!capabilityCount)">
                    <UIcon name="i-lucide-sliders-horizontal" class="size-3.5" />
                    <span>{{ $t('catalog.features') }}</span>
                    <span
                        v-if="capabilityCount"
                        class="inline-flex items-center justify-center bg-primary text-white text-[10px] font-bold rounded-full size-4"
                    >
                        {{ capabilityCount }}
                    </span>
                </button>
                <template #content>
                    <div class="p-2">
                        <div class="grid grid-cols-2 gap-1">
                            <button
                                v-for="toggle in filterToggles"
                                :key="toggle.key"
                                @click="$emit('toggleFilter', toggle.key)"
                                class="flex items-center gap-1.5 rounded-lg px-2.5 py-1.5 text-xs cursor-pointer transition-colors"
                                :class="filterBtnClass(toggle)"
                            >
                                <UIcon :name="toggle.icon" class="size-3.5" />
                                <span>{{ toggle.label }}</span>
                            </button>
                        </div>
                    </div>
                </template>
            </UPopover>

            <!-- Desktop: Price Popover -->
            <UPopover v-if="!isMobile">
                <button :class="triggerClass(priceActive)">
                    <UIcon name="i-lucide-dollar-sign" class="size-3.5" />
                    <span>{{ $t('catalog.price') }}</span>
                </button>
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

            <!-- Desktop: I/O Types Popover -->
            <UPopover v-if="!isMobile">
                <button :class="triggerClass(!!ioTypeCount)">
                    <UIcon name="i-lucide-arrow-left-right" class="size-3.5" />
                    <span>{{ $t('catalog.ioTypes') }}</span>
                    <span
                        v-if="ioTypeCount"
                        class="inline-flex items-center justify-center bg-primary text-white text-[10px] font-bold rounded-full size-4"
                    >
                        {{ ioTypeCount }}
                    </span>
                </button>
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

            <!-- Desktop: Provider Popover -->
            <CatalogProviderPopover
                v-if="!isMobile"
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

            <!-- Desktop: Column visibility -->
            <UDropdownMenu
                v-if="!isMobile && columnMenuItems.length"
                :items="[columnMenuItems]"
                :content="{ align: 'end' }"
            >
                <UButton icon="i-lucide-columns-3-cog" color="neutral" variant="outline" size="md" />
            </UDropdownMenu>
        </div>
    </div>

    <!-- Active filter pills -->
    <TransitionGroup v-if="activePills.length" name="fade" tag="div" class="flex flex-wrap items-center gap-1.5 mb-4">
        <span
            v-for="(pill, i) in activePills"
            :key="i"
            class="inline-flex items-center gap-1 rounded-full pl-2.5 pr-1.5 py-0.5 bg-primary/10 text-primary text-xs font-medium"
        >
            {{ pill.label }}
            <button @click="pill.remove" class="hover:text-error cursor-pointer transition-colors">
                <UIcon name="i-lucide-x" class="size-3" />
            </button>
        </span>
        <button v-if="activePills.length > 1" @click="clearAllFilters" class="text-xs text-muted hover:text-error cursor-pointer transition-colors">
            {{ $t("common.clearAll") }}
        </button>
    </TransitionGroup>

    <!-- Mobile: Expandable filter panel -->
    <Transition name="expand" @enter="onExpandEnter" @after-enter="onExpandAfterEnter" @leave="onExpandLeave">
        <div v-if="isMobile && filtersOpen" class="mb-4 space-y-2.5">
            <div class="flex items-center gap-2 flex-wrap">
                <button
                    v-for="toggle in filterToggles"
                    :key="toggle.key"
                    @click="$emit('toggleFilter', toggle.key)"
                    class="flex items-center gap-1 rounded-lg px-2.5 py-1 text-xs cursor-pointer transition-colors"
                    :class="filterBtnClass(toggle)"
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
.fade-enter-active,
.fade-leave-active {
    transition: all 0.2s ease;
}
.fade-enter-from {
    opacity: 0;
    transform: translateY(4px);
}
.fade-leave-to {
    opacity: 0;
    transform: translateY(-4px);
}
</style>

<script setup lang="ts">
    const props = defineProps<{
        search: string;
        total: number;
        columnMenuItems: any[];
        filterToggles: { key: string; label: string; icon: string; color: string; default?: boolean }[];
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
    const filtersOpen = ref(false);

    const colorMap: Record<string, { bg: string; border: string; text: string }> = {
        emerald: { bg: "bg-emerald-500/10", border: "border-emerald-500/30", text: "text-emerald-500" },
        amber:   { bg: "bg-amber-500/10",   border: "border-amber-500/30",   text: "text-amber-500" },
        blue:    { bg: "bg-blue-500/10",     border: "border-blue-500/30",     text: "text-blue-500" },
        violet:  { bg: "bg-violet-500/10",   border: "border-violet-500/30",   text: "text-violet-500" },
        rose:    { bg: "bg-rose-500/10",      border: "border-rose-500/30",     text: "text-rose-500" },
        purple:  { bg: "bg-purple-500/10",    border: "border-purple-500/30",   text: "text-purple-500" },
        stone:   { bg: "bg-stone-500/10",     border: "border-stone-500/30",    text: "text-stone-500" },
        primary: { bg: "bg-primary/10",       border: "border-primary/30",       text: "text-primary" },
    };

    const filterBtnClass = (toggle: { key: string; color: string }) => {
        if (props.filters[toggle.key]) {
            const c = colorMap[toggle.color] ?? colorMap.primary!;
            return `${c.bg} border ${c.border} ${c.text}`;
        }
        return "bg-default border border-default text-toned hover:border-accented";
    };

    const triggerClass = (active: boolean) => {
        const base = "flex items-center gap-1.5 rounded-lg px-2.5 py-1.5 text-xs cursor-pointer transition-colors border shrink-0";
        if (active) return `${base} bg-primary/10 border-primary/30 text-primary`;
        return `${base} bg-default border-default text-toned hover:border-accented`;
    };

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

    const capabilityCount = computed(() => Object.values(props.filters).filter(Boolean).length);
    const priceActive = computed(() => {
        const pr = priceRange.value || [0, 100];
        const opr = outputPriceRange.value || [0, 100];
        return pr[0] !== 0 || pr[1] !== 100 || opr[0] !== 0 || opr[1] !== 100;
    });
    const ioTypeCount = computed(() => (selectedInputTypes.value?.length || 0) + (selectedOutputTypes.value?.length || 0));

    const activeFilterCount = computed(() => activePills.value.length);

    const activePills = computed(() => {
        const pills: { label: string; remove: () => void }[] = [];
        if (props.search) pills.push({ label: `"${props.search}"`, remove: () => emit("update:search", "") });
        props.filterToggles.forEach((toggle) => {
            if ((props.filters as any)[toggle.key]) pills.push({ label: toggle.label, remove: () => emit("toggleFilter", toggle.key) });
        });
        props.selectedProviders?.forEach((p: any) => {
            const name = p.label || p.name || p.value;
            pills.push({ label: name, remove: () => emit("removeProvider", p) });
        });
        selectedInputTypes.value?.forEach((item: any) => {
            pills.push({ label: item.label || item.value, remove: () => { selectedInputTypes.value = selectedInputTypes.value.filter((x: any) => (x.value ?? x) !== (item.value ?? item)); } });
        });
        selectedOutputTypes.value?.forEach((item: any) => {
            pills.push({ label: item.label || item.value, remove: () => { selectedOutputTypes.value = selectedOutputTypes.value.filter((x: any) => (x.value ?? x) !== (item.value ?? item)); } });
        });
        const pr = priceRange.value || [0, 100];
        if (pr[0] !== 0 || pr[1] !== 100) pills.push({ label: `In $${pr[0]}–${pr[1]}`, remove: () => { priceRange.value = [0, 100]; } });
        const opr = outputPriceRange.value || [0, 100];
        if (opr[0] !== 0 || opr[1] !== 100) pills.push({ label: `Out $${opr[0]}–${opr[1]}`, remove: () => { outputPriceRange.value = [0, 100]; } });
        return pills;
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
