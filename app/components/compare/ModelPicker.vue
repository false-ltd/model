<template>
    <UPopover modal>
        <UButton color="primary" size="sm"> <span class="text-sm">+</span> {{ t("common.addModel") }} </UButton>

        <template #content>
            <div class="w-[320px] bg-default border border-default rounded-xl shadow-lg overflow-hidden">
                <div class="p-2.5 border-b border-default">
                    <input
                        v-model="query"
                        :placeholder="t('search.placeholder')"
                        class="w-full bg-elevated rounded-lg px-3 py-2 text-sm text-default border border-default focus:border-primary focus:outline-none"
                    />
                </div>
                <div class="max-h-65 overflow-y-auto">
                    <div v-if="loading" class="p-4 text-center text-sm text-muted">...</div>
                    <div v-else-if="results.length === 0" class="p-4 text-center text-sm text-muted">
                        {{ t("common.noResults") }}
                    </div>
                    <button
                        v-for="r in results"
                        :key="r.id"
                        @click="selectModel(r)"
                        :disabled="modelIds.includes(r.id)"
                        class="w-full flex items-center gap-2.5 px-3 py-2 text-left cursor-pointer hover:bg-accented transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
                    >
                        <ProviderLogo :provider-id="r.provider_id" cls="w-6 h-6 rounded shrink-0" />
                        <div class="min-w-0 flex-1">
                            <div class="text-sm text-default truncate">{{ r.name }}</div>
                            <div class="text-xs text-muted truncate">
                                {{ r.provider_id }} · {{ r.family }}
                            </div>
                        </div>
                        <span v-if="modelIds.includes(r.id)" class="text-xs text-success">✓</span>
                    </button>
                </div>
            </div>
        </template>
    </UPopover>
</template>

<script setup lang="ts">
    const props = defineProps<{
        modelIds: number[];
    }>();

    const emit = defineEmits<{
        select: [id: number];
    }>();

    const { t } = useI18n();
    const config = useRuntimeConfig();
    const { addModel } = useCompare();

    const query = ref("");
    const results = ref<any[]>([]);
    const loading = ref(false);
    let timer: ReturnType<typeof setTimeout> | null = null;

    const selectModel = (r: any) => {
        addModel(r.id);
    };

    watch(query, (q) => {
        if (timer) clearTimeout(timer);
        if (!q.trim()) {
            results.value = [];
            return;
        }
        timer = setTimeout(async () => {
            loading.value = true;
            try {
                const res = await $fetch<{ data: any[] }>(`${config.public.apiBase}/api/v1/models`, {
                    params: { q: q.trim(), page_size: 20 },
                });
                results.value = res.data || [];
            } catch {
                results.value = [];
            } finally {
                loading.value = false;
            }
        }, 300);
    });
</script>
