<template>
    <Teleport to="body">
        <div
            v-if="modelIds.length > 0"
            class="fixed bottom-5 left-1/2 -translate-x-1/2 z-50 bg-default border border-default rounded-xl shadow-lg px-4 py-2.5 flex items-center gap-3 transition-all max-w-[90vw]"
        >
            <div class="text-sm text-muted whitespace-nowrap">
                {{ modelIds.length }} {{ $t("common.selected") }}
            </div>

            <!-- Selected models -->
            <div class="flex gap-1.5 overflow-x-auto">
                <div
                    v-for="m in selectedModels"
                    :key="m.id"
                    class="flex items-center gap-1.5 bg-elevated rounded-lg pl-1 pr-2 py-1 shrink-0 border border-default group/item"
                >
                    <img
                        :src="`https://models.dev/logos/${m.provider_id}.svg`"
                        class="w-5 h-5 rounded shrink-0"
                        @error="($event.target as HTMLImageElement).style.display = 'none'"
                    />
                    <span class="text-xs text-default font-medium whitespace-nowrap">{{ m.name }}</span>
                    <button
                        @click="removeModel(m.id)"
                        class="opacity-40 group-hover/item:opacity-100 cursor-pointer transition-opacity"
                    >
                        <UIcon name="i-lucide-x" class="size-3 text-muted" />
                    </button>
                </div>
            </div>

            <NuxtLink
                :to="localePath('/compare')"
                class="bg-primary text-white rounded-lg px-4 py-1.5 text-sm font-medium no-underline hover:bg-primary transition-colors whitespace-nowrap shrink-0"
            >
                {{ $t("common.compare") }}
            </NuxtLink>
            <button
                @click="clearAll"
                class="text-xs text-muted hover:text-error cursor-pointer transition-colors shrink-0"
            >
                {{ $t("common.clearAll") }}
            </button>
        </div>
    </Teleport>
</template>

<script setup lang="ts">
    const localePath = useLocalePath();
    const { modelIds, removeModel, clearAll } = useCompare();

    const { data: result } = await useAsyncData(
        "compare-bar-models",
        () =>
            $fetch<{ data: any[] }>("/api/compare", {
                params: { ids: modelIds.value.join(",") },
            }),
        { watch: [modelIds] },
    );

    const selectedModels = computed(() => result.value?.data || []);
</script>
