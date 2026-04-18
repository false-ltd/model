<template>
    <UPopover modal>
        <UButton label="Open" color="neutral" variant="subtle">
            <UIcon name="i-lucide-building-2" class="size-3.5" /> {{ providerTriggerLabel }}
        </UButton>

        <template #content>
            <div class="w-96 max-h-105 overflow-y-auto bg-default border border-default rounded-xl shadow-lg">
                <div class="sticky top-0 bg-default px-2 pt-2 pb-1.5 z-10">
                    <UInput
                        v-model="providerSearch"
                        :placeholder="$t('catalog.searchProviders')"
                        icon="i-lucide-search"
                        size="sm"
                        class="w-full"
                    />
                </div>

                <!-- Selected providers (always on top) -->
                <template v-if="selectedProviders.length && !providerSearch">
                    <div class="px-3 pb-1.5">
                        <div class="flex items-center justify-between mb-1">
                            <div class="text-xs uppercase tracking-wider text-primary font-medium">
                                {{ $t("catalog.selected") }}
                            </div>
                            <button
                                @click="$emit('clearProviders')"
                                class="text-xs text-muted hover:text-error cursor-pointer transition-colors"
                            >
                                {{ $t("catalog.clearProviders") }}
                            </button>
                        </div>
                        <div class="grid grid-cols-2 gap-0.5">
                            <button
                                v-for="sp in selectedProviders"
                                :key="sp.value"
                                @click="$emit('removeProvider', sp)"
                                class="flex items-center gap-2 px-2 py-1.5 rounded-md cursor-pointer transition-colors text-left bg-primary/10 hover:bg-primary/20"
                            >
                                <img
                                    :src="`https://models.dev/logos/${sp.value}.svg`"
                                    class="w-5 h-5 rounded shrink-0"
                                    @error="($event.target as HTMLImageElement).style.display = 'none'"
                                />
                                <span class="flex-1 text-sm truncate text-primary font-medium">
                                    {{ sp.label }}
                                </span>
                                <UIcon name="i-lucide-x" class="size-3.5 text-primary shrink-0" />
                            </button>
                        </div>
                    </div>
                    <div class="border-t border-default mx-3" />
                </template>

                <!-- Top Providers (empty search only) -->
                <template v-if="!providerSearch">
                    <div class="px-3 pb-1.5">
                        <div class="text-xs uppercase tracking-wider text-muted font-medium mb-1">
                            {{ $t("catalog.topProviders") }}
                        </div>
                        <div class="grid grid-cols-2 gap-0.5">
                            <button
                                v-for="p in topProviders"
                                :key="p.id"
                                @click="$emit('toggleProvider', p)"
                                class="flex items-center gap-2 px-2 py-1.5 rounded-md cursor-pointer transition-colors text-left"
                                :class="isSelected(p) ? 'bg-primary/10' : 'hover:bg-accented'"
                            >
                                <img
                                    :src="`https://models.dev/logos/${(p as any).id}.svg`"
                                    class="w-5 h-5 rounded shrink-0"
                                    @error="($event.target as HTMLImageElement).style.display = 'none'"
                                />
                                <span
                                    class="flex-1 text-sm truncate"
                                    :class="isSelected(p) ? 'text-primary font-medium' : 'text-toned'"
                                >
                                    {{ (p as any).name }}
                                </span>
                                <span class="text-xs text-muted tabular-nums">{{ (p as any).modelCount }}</span>
                                <UIcon v-if="isSelected(p)" name="i-lucide-check" class="size-3.5 text-primary shrink-0" />
                            </button>
                        </div>
                    </div>

                    <div class="border-t border-default mx-3" />

                    <!-- All Providers A-Z -->
                    <div class="px-3 py-1.5">
                        <div class="text-xs uppercase tracking-wider text-muted font-medium mb-1">
                            {{ $t("catalog.allProviders") }}
                        </div>
                        <template v-for="group in groupedProviders" :key="group.letter">
                            <div class="text-xs font-semibold text-muted uppercase px-2 py-1 sticky top-10 bg-default">
                                {{ group.letter }}
                            </div>
                            <button
                                v-for="p in group.items"
                                :key="(p as any).id"
                                @click="$emit('toggleProvider', p)"
                                class="flex items-center gap-2 w-full px-2 py-1.5 rounded-md cursor-pointer transition-colors text-left"
                                :class="isSelected(p) ? 'bg-primary/10' : 'hover:bg-accented'"
                            >
                                <img
                                    :src="`https://models.dev/logos/${(p as any).id}.svg`"
                                    class="w-4 h-4 rounded shrink-0"
                                    @error="($event.target as HTMLImageElement).style.display = 'none'"
                                />
                                <span
                                    class="flex-1 text-sm truncate"
                                    :class="isSelected(p) ? 'text-primary font-medium' : 'text-toned'"
                                >
                                    {{ (p as any).name }}
                                </span>
                                <span class="text-xs text-muted tabular-nums">{{ (p as any).modelCount }}</span>
                                <UIcon v-if="isSelected(p)" name="i-lucide-check" class="size-3.5 text-primary shrink-0" />
                            </button>
                        </template>
                    </div>
                </template>

                <!-- Search results -->
                <template v-else>
                    <div class="px-3 py-1.5">
                        <button
                            v-for="p in filteredProviders"
                            :key="(p as any).id"
                            @click="$emit('toggleProvider', p)"
                            class="flex items-center gap-2 w-full px-2 py-1.5 rounded-md cursor-pointer transition-colors text-left"
                            :class="isSelected(p) ? 'bg-primary/10' : 'hover:bg-accented'"
                        >
                            <img
                                :src="`https://models.dev/logos/${(p as any).id}.svg`"
                                class="w-5 h-5 rounded shrink-0"
                                @error="($event.target as HTMLImageElement).style.display = 'none'"
                            />
                            <span
                                class="flex-1 text-sm truncate"
                                :class="isSelected(p) ? 'text-primary font-medium' : 'text-toned'"
                            >
                                {{ (p as any).name }}
                            </span>
                            <span class="text-xs text-muted tabular-nums">{{ (p as any).modelCount }}</span>
                            <UIcon v-if="isSelected(p)" name="i-lucide-check" class="size-3.5 text-primary shrink-0" />
                        </button>
                        <div v-if="!filteredProviders.length" class="py-6 text-center text-sm text-muted">
                            {{ $t("catalog.noProviderMatch") }}
                        </div>
                    </div>
                </template>
            </div>
        </template>
    </UPopover>
</template>

<script setup lang="ts">
    defineProps<{
        selectedProviders: any[]
        topProviders: any[]
        groupedProviders: { letter: string; items?: any[] }[]
        filteredProviders: any[]
        providerTriggerLabel: string
        isSelected: (p: { id: string }) => boolean
    }>()

    const providerSearch = defineModel<string>('providerSearch', { required: true })

    defineEmits<{
        toggleProvider: [provider: any]
        removeProvider: [provider: any]
        clearProviders: []
    }>()
</script>
