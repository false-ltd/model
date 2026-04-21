<template>
    <Teleport to="body">
        <Transition name="slide-up">
            <div v-if="modelIds.length > 0" class="fixed bottom-5 right-5 z-50 flex flex-col items-end gap-2">
                <!-- FAB button -->
                <button
                    @click="open = !open"
                    class="relative size-12 rounded-full bg-linear-to-br from-amber-500 to-amber-600 text-white flex items-center justify-center shadow-lg shadow-amber-500/30 hover:shadow-amber-500/50 transition-shadow cursor-pointer"
                >
                    <Transition name="fab-icon" mode="out-in">
                        <UIcon v-if="!open" name="i-lucide-git-compare" class="size-5" />
                        <UIcon v-else name="i-lucide-x" class="size-5" />
                    </Transition>
                    <!-- Badge -->
                    <span
                        class="absolute -top-1 -right-1 min-w-4.5 h-4.5 rounded-full bg-white text-amber-600 text-[10px] font-bold flex items-center justify-center px-1 shadow"
                    >
                        {{ modelIds.length }}
                    </span>
                </button>

                <!-- Popover panel -->
                <Transition name="popover">
                    <div
                        v-if="open"
                        class="absolute bottom-16 right-0 w-65 bg-default border border-default rounded-xl shadow-xl overflow-hidden"
                    >
                        <!-- Header -->
                        <div class="flex items-center px-3.5 py-2.5 border-b border-default">
                            <span class="text-xs font-semibold text-default">{{ t("compare.selectedModels") }}</span>
                            <span class="text-[10px] text-muted ml-auto mr-2.5">{{ modelIds.length }}/4</span>
                            <button
                                @click="
                                    clearAll();
                                    open = false;
                                "
                                class="text-[10px] text-amber-500 hover:text-amber-600 cursor-pointer transition-colors"
                            >
                                {{ t("common.clearAll") }}
                            </button>
                        </div>

                        <!-- Model list -->
                        <div class="py-1 max-h-60 overflow-y-auto">
                            <div
                                v-for="m in selectedModels"
                                :key="m.id"
                                class="group flex items-center gap-2 px-3.5 py-1.5 hover:bg-accented/50 transition-colors"
                            >
                                <div
                                    class="w-7 h-7 rounded-full bg-elevated shrink-0 flex items-center justify-center text-[10px] font-bold text-muted relative overflow-hidden"
                                >
                                    <span>{{ m.provider_id?.charAt(0).toUpperCase() }}</span>
                                    <ProviderLogo :provider-id="m.provider_id" cls="absolute inset-0 w-full h-full object-cover rounded-full" />
                                </div>
                                <div class="flex-1 min-w-0">
                                    <div class="text-xs font-medium text-default truncate">{{ m.name }}</div>
                                    <div class="text-[10px] text-muted">{{ m.provider_id }}</div>
                                </div>
                                <button
                                    @click="removeModel(m.id)"
                                    class="size-5 rounded-full bg-error/10 text-error flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer shrink-0"
                                >
                                    <UIcon name="i-lucide-x" class="size-2.5" />
                                </button>
                            </div>
                        </div>

                        <!-- Action -->
                        <div class="px-3.5 py-2.5 border-t border-default">
                            <NuxtLink
                                :to="localePath('/compare')"
                                class="block bg-linear-to-r from-amber-500 to-amber-600 text-white text-center rounded-lg py-2 text-xs font-semibold no-underline hover:from-amber-600 hover:to-amber-700 transition-all"
                            >
                                {{ t("common.compare") }}
                            </NuxtLink>
                        </div>
                    </div>
                </Transition>

                <!-- Backdrop -->
                <Transition name="fade">
                    <div v-if="open" class="fixed inset-0 -z-10" @click="open = false" />
                </Transition>
            </div>
        </Transition>
    </Teleport>
</template>

<script setup lang="ts">
    const localePath = useLocalePath();
    const config = useRuntimeConfig();
    const { t } = useI18n();
    const { modelIds, removeModel, clearAll } = useCompare();

    const open = ref(false);

    const modelIdsKey = computed(() => [...modelIds.value].sort().join(","));
    const { data: result } = await useAsyncData(
        "compare-bar-models",
        () => {
            if (!modelIds.value.length) return Promise.resolve({ data: [] });
            return $fetch<{ data: any[] }>(`${config.public.apiBase}/api/v1/compare`, {
                params: { ids: modelIds.value.join(",") },
            });
        },
        { watch: [modelIdsKey] },
    );

    const selectedModels = computed(() => result.value?.data || []);

    const route = useRoute();
    watch(
        () => route.fullPath,
        () => {
            open.value = false;
        },
    );

    onMounted(() => {
        const onKeydown = (e: KeyboardEvent) => {
            if (e.key === "Escape" && open.value) open.value = false;
        };
        document.addEventListener("keydown", onKeydown);
        onUnmounted(() => document.removeEventListener("keydown", onKeydown));
    });
</script>

<style scoped>
    .slide-up-enter-active,
    .slide-up-leave-active {
        transition: all 0.3s ease;
    }
    .slide-up-enter-from,
    .slide-up-leave-to {
        opacity: 0;
        transform: translateY(20px);
    }

    .popover-enter-active,
    .popover-leave-active {
        transition: all 0.2s ease;
    }
    .popover-enter-from,
    .popover-leave-to {
        opacity: 0;
        transform: translateY(8px) scale(0.95);
    }

    .fab-icon-enter-active,
    .fab-icon-leave-active {
        transition: all 0.15s ease;
    }
    .fab-icon-enter-from,
    .fab-icon-leave-to {
        opacity: 0;
        transform: rotate(90deg);
    }

    .fade-enter-active,
    .fade-leave-active {
        transition: opacity 0.2s ease;
    }
    .fade-enter-from,
    .fade-leave-to {
        opacity: 0;
    }
</style>
