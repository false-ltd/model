<template>
    <div class="flex items-center gap-2 mb-4 flex-nowrap overflow-x-auto">
        <!-- Search -->
        <div class="relative shrink-0">
            <UIcon name="i-lucide-search" class="size-3.5 absolute left-2.5 top-1/2 -translate-y-1/2 text-muted" />
            <input
                :value="search"
                @input="($event) => $emit('update:search', ($event.target as HTMLInputElement).value)"
                :placeholder="$t('search.placeholder')"
                class="bg-elevated rounded-lg pl-8 pr-3 py-1.5 text-sm text-toned w-48 border border-default focus:border-primary focus:outline-none transition-colors"
            />
        </div>

        <!-- Separator -->
        <div class="w-px h-5 bg-default shrink-0" />

        <!-- Toggle pills -->
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

        <!-- Provider picker -->
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

        <!-- Input / output type selects -->
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

        <!-- Price range -->
        <div class="flex items-center gap-2 bg-default border border-default rounded-lg px-3 py-1 shrink-0">
            <span class="text-xs text-muted whitespace-nowrap">{{ $t("catalog.priceRange") }}</span>
            <USlider v-model="priceRange" :min="0" :max="100" :step="0.1" color="primary" class="w-24" size="sm" />
            <span class="text-xs text-toned font-mono whitespace-nowrap tabular-nums">
                ${{ priceRange[0] }}–${{ priceRange[1] }}
            </span>
        </div>

        <!-- Count -->
        <span class="text-md px-2 text-muted shrink-0">{{ total }} </span>

        <!-- Spacer -->
        <div class="flex-1 min-w-2" />

        <!-- View toggle -->
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

        <!-- Column visibility -->
        <UDropdownMenu
            v-if="columnMenuItems.length"
            :items="[columnMenuItems]"
            :content="{ align: 'end', class: 'max-h-80 overflow-y-auto' }"
        >
            <UButton icon="i-lucide-columns-3-cog" color="neutral" variant="outline" size="md" />
        </UDropdownMenu>
    </div>
</template>

<script setup lang="ts">
    defineProps<{
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
    const priceRange = defineModel<number[]>("priceRange", { required: true });

    defineEmits<{
        "update:modelValue": [value: "grid" | "table"];
        "update:search": [value: string];
        toggleFilter: [key: "freeOnly" | "reasoning" | "vision"];
        toggleProvider: [provider: any];
        removeProvider: [provider: any];
        clearProviders: [];
    }>();
</script>
