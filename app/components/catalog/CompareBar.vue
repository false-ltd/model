<template>
    <Teleport to="body">
        <Transition name="slide-up">
            <div
                v-if="modelIds.length > 0"
                class="fixed bottom-4 left-1/2 -translate-x-1/2 z-50 bg-default border border-default rounded-xl shadow-lg pl-3 pr-2 py-2 flex items-center gap-2.5 max-w-[90vw]"
            >
                <!-- Model logos (max 5 visible) -->
                <div class="flex items-center -space-x-1.5">
                    <UTooltip v-for="(m, i) in visibleModels" :key="m.id" :text="m.name">
                        <div class="relative group/item">
                            <img
                                :src="`https://models.dev/logos/${m.provider_id}.svg`"
                                class="w-7 h-7 rounded-full border-2 border-default bg-elevated shrink-0"
                                @error="($event.target as HTMLImageElement).style.display = 'none'"
                            />
                            <button
                                @click="removeModel(m.id)"
                                class="absolute -top-1 -right-1 size-3.5 rounded-full bg-error text-white flex items-center justify-center opacity-0 group-hover/item:opacity-100 transition-opacity cursor-pointer"
                            >
                                <UIcon name="i-lucide-x" class="size-2" />
                            </button>
                        </div>
                    </UTooltip>
                    <!-- Overflow indicator -->
                    <span
                        v-if="selectedModels.length > maxVisible"
                        class="w-7 h-7 rounded-full bg-elevated border-2 border-default flex items-center justify-center text-[10px] font-bold text-muted"
                    >
                        +{{ selectedModels.length - maxVisible }}
                    </span>
                </div>

                <!-- Count -->
                <span class="text-xs text-muted whitespace-nowrap pl-1">
                    {{ modelIds.length }} {{ $t("common.selected") }}
                </span>

                <!-- Actions -->
                <div class="flex items-center gap-1.5 ml-1">
                    <NuxtLink
                        :to="localePath('/compare')"
                        class="bg-primary text-white rounded-lg px-3.5 py-1.5 text-xs font-medium no-underline hover:bg-primary transition-colors whitespace-nowrap"
                    >
                        {{ $t("common.compare") }}
                    </NuxtLink>
                    <button
                        @click="clearAll"
                        class="text-[11px] text-muted hover:text-error cursor-pointer transition-colors whitespace-nowrap px-1.5"
                    >
                        {{ $t("common.clearAll") }}
                    </button>
                </div>
            </div>
        </Transition>
    </Teleport>
</template>

<script setup lang="ts">
    const localePath = useLocalePath();
    const { modelIds, removeModel, clearAll } = useCompare();

    const maxVisible = 5;

    const { data: result } = await useAsyncData(
        "compare-bar-models",
        () =>
            $fetch<{ data: any[] }>("/api/compare", {
                params: { ids: modelIds.value.join(",") },
            }),
        { watch: [modelIds] },
    );

    const selectedModels = computed(() => result.value?.data || []);
    const visibleModels = computed(() => selectedModels.value.slice(0, maxVisible));
</script>
