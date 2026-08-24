<template>
    <NuxtLink
        :to="localePath(`/catalog?providers=${provider.id}`)"
        class="group block bg-default border border-default rounded-xl no-underline transition-colors hover:border-primary/50"
        :class="featured ? 'p-4 hover:shadow-md' : 'px-3 py-2.5'"
    >
        <div class="flex items-center gap-2.5" :class="featured ? 'mb-3' : ''">
            <ProviderLogo
                :provider-id="provider.id"
                :cls="featured ? 'w-9 h-9 rounded-lg p-1 shrink-0' : 'w-6 h-6 rounded-md p-0.5 shrink-0'"
            />
            <div class="min-w-0 flex-1">
                <div
                    class="font-semibold text-default truncate group-hover:text-primary transition-colors"
                    :class="featured ? 'text-sm' : 'text-[13px]'"
                >
                    {{ provider.name }}
                </div>
                <div class="text-[11px] text-muted font-mono truncate">{{ provider.id }}</div>
            </div>
            <span
                class="text-xs font-semibold text-toned bg-elevated px-2 py-0.5 rounded tabular shrink-0"
                :class="featured ? 'text-sm' : ''"
            >
                {{ count }}
                <span v-if="featured" class="text-[10px] text-muted font-normal">{{ t("providers.models") }}</span>
            </span>
        </div>

        <!-- Featured: open-weights band + value line -->
        <template v-if="featured">
            <div class="flex items-center gap-2 mb-2.5">
                <div class="flex-1 h-1.5 rounded-full bg-elevated overflow-hidden">
                    <div class="h-full rounded-full bg-primary" :style="{ width: (e?.openPct ?? 0) + '%' }" />
                </div>
                <span class="text-[11px] tabular shrink-0" :class="openTone">{{ e?.openPct ?? 0 }}% {{ t("providers.openShare") }}</span>
            </div>
            <div class="grid grid-cols-2 gap-x-2 gap-y-1 text-[11px]">
                <div v-if="medianLabel" class="truncate">
                    <span class="text-muted">{{ t("providers.medianPrice") }}:</span>
                    <span class="text-toned font-mono ml-1 tabular">{{ medianLabel }}</span>
                </div>
                <div v-if="provider.npm" class="truncate">
                    <span class="text-muted">NPM:</span>
                    <span class="text-toned font-mono ml-1">{{ provider.npm }}</span>
                </div>
                <div v-if="provider.api_url" class="truncate col-span-2">
                    <span class="text-muted">API:</span>
                    <span class="text-toned font-mono ml-1">{{ provider.api_url }}</span>
                </div>
                <div v-if="provider.env?.length" class="truncate">
                    <span class="text-muted">ENV:</span>
                    <span class="text-toned font-mono ml-1">{{ provider.env[0] }}</span>
                </div>
                <div v-if="provider.doc_url" class="truncate">
                    <a
                        :href="provider.doc_url"
                        target="_blank"
                        class="text-primary hover:underline inline-flex items-center gap-0.5"
                        @click.stop
                    >
                        {{ t("providers.docs") }}
                        <UIcon name="i-lucide-arrow-up-right" class="size-3" />
                    </a>
                </div>
            </div>
        </template>

        <!-- Compact: median + open dot -->
        <template v-else-if="e">
            <div class="flex items-center gap-2 text-[11px] text-muted">
                <span v-if="e.median != null" class="tabular">${{ e.median < 10 ? e.median.toFixed(2) : e.median.toFixed(0) }} {{ t("providers.medianPrice") }}</span>
                <span class="tabular" :class="openTone">{{ e.openPct }}%</span>
            </div>
        </template>
    </NuxtLink>
</template>

<script setup lang="ts">
    /**
     * Provider card for the providers page. Featured cards (top tier) show
     * the open-weights band, median price, and integration details; compact
     * cards in the A-Z sections stay one-line dense.
     */
    const props = withDefaults(
        defineProps<{
            provider: any;
            enrich?: { count: number; openPct: number; median: number | null } | undefined;
            featured?: boolean;
            compact?: boolean;
        }>(),
        { featured: false, compact: false },
    );

    const localePath = useLocalePath();
    const { t } = useI18n();

    const count = computed(() => props.enrich?.count ?? props.provider.model_count ?? 0);

    const medianLabel = computed(() => {
        const m = props.enrich?.median;
        if (m == null) return null;
        return `$${m < 10 ? m.toFixed(2) : m.toFixed(0)}`;
    });

    const openTone = computed(() => {
        const pct = props.enrich?.openPct ?? 0;
        return pct >= 60 ? "text-success font-semibold" : pct >= 20 ? "text-toned" : "text-muted";
    });
</script>
