<template>
    <div>
        <!-- Header -->
        <div class="flex items-center justify-between mb-5">
            <div>
                <h1 class="text-xl font-bold text-default">{{ t("compare.title") }}</h1>
                <div class="text-sm text-muted">{{ compareModels.length }} / 10 {{ t("compare.modelsSelected") }}</div>
            </div>
            <div class="flex items-center gap-2">
                <CompareModelPicker v-if="modelIds.length < 10" :model-ids="modelIds" />
                <button
                    v-if="compareModels.length > 0"
                    @click="clearAll"
                    class="bg-elevated rounded-lg px-3.5 py-2 text-sm text-error cursor-pointer hover:bg-accented transition-colors"
                >
                    {{ t("common.clearAll") }}
                </button>
            </div>
        </div>

        <!-- Empty state -->
        <div v-if="compareModels.length === 0" class="bg-default border border-default rounded-xl p-16 text-center">
            <div class="text-sm text-muted mb-1">{{ t("compare.empty") }}</div>
            <div class="text-sm text-muted">
                <NuxtLink :to="localePath('/catalog')" class="text-primary hover:underline">{{
                    t("nav.catalog")
                }}</NuxtLink>
            </div>
        </div>

        <!-- Mobile: card view -->
        <div v-else-if="isMobile" class="space-y-3">
            <div
                v-for="m in compareModels"
                :key="m.id"
                class="bg-default border border-default rounded-xl p-4"
            >
                <div class="flex items-center gap-3 mb-3">
                    <img
                        :src="`https://models.dev/logos/${m.provider_id}.svg`"
                        class="w-9 h-9 rounded-lg shrink-0"
                        @error="($event.target as HTMLImageElement).style.display = 'none'"
                    />
                    <div class="min-w-0">
                        <NuxtLink
                            :to="localePath(`/model/${m.id}`)"
                            class="font-semibold text-sm text-default hover:text-primary transition-colors no-underline"
                        >{{ m.name }}</NuxtLink>
                        <div class="text-xs text-muted">{{ m.family }}</div>
                    </div>
                    <button @click="removeModel(m.id)" class="ml-auto text-xs text-error cursor-pointer hover:underline shrink-0">
                        ✕ {{ t("common.remove") }}
                    </button>
                </div>
                <div class="grid grid-cols-2 gap-x-4 gap-y-1 text-xs">
                    <template v-for="field in identityFields" :key="'i'+field.key">
                        <span class="text-muted">{{ field.label }}</span>
                        <span class="text-toned">{{ field.format ? field.format(m[field.key]) : (m[field.key] || '—') }}</span>
                    </template>
                    <template v-if="hasPaidModels" v-for="field in pricingFields" :key="'p'+field.key">
                        <span class="text-muted">{{ field.label }}</span>
                        <span class="text-toned">{{ field.format ? field.format(m[field.key]) : (m[field.key] ?? '—') }}</span>
                    </template>
                    <template v-for="field in limitFields" :key="'l'+field.key">
                        <span class="text-muted">{{ field.label }}</span>
                        <span class="text-toned">{{ formatTokens(m[field.key]) }}</span>
                    </template>
                    <template v-for="cap in capabilityFields" :key="'c'+cap.key">
                        <span class="text-muted">{{ cap.label }}</span>
                        <span :class="m[cap.key] ? 'text-success' : 'text-muted'">{{ m[cap.key] ? t('common.yes') : '—' }}</span>
                    </template>
                    <template v-for="field in timelineFields" :key="'t'+field.key">
                        <span class="text-muted">{{ field.label }}</span>
                        <span class="text-toned">{{ m[field.key] || '—' }}</span>
                    </template>
                </div>
            </div>
        </div>

        <!-- Desktop: Comparison grid -->
        <div v-else class="border border-default rounded-xl overflow-x-auto bg-default">
            <!-- Model headers -->
            <div
                class="grid border-b border-default"
                :style="{
                    gridTemplateColumns: `160px repeat(${compareModels.length}, 200px)`,
                    minWidth: `${160 + compareModels.length * 200}px`,
                }"
            >
                <div class="bg-elevated border-r border-default" />
                <div
                    v-for="m in compareModels"
                    :key="m.id"
                    class="border-r border-default last:border-r-0 p-4 text-center"
                >
                    <img
                        :src="`https://models.dev/logos/${m.provider_id}.svg`"
                        class="w-9 h-9 rounded-lg mx-auto mb-1.5"
                        @error="($event.target as HTMLImageElement).style.display = 'none'"
                    />
                    <NuxtLink
                        :to="localePath(`/model/${m.id}`)"
                        class="font-semibold text-sm text-default hover:text-primary transition-colors no-underline"
                        >{{ m.name }}</NuxtLink
                    >
                    <div class="text-xs text-muted">{{ m.family }}</div>
                    <button @click="removeModel(m.id)" class="text-xs text-error mt-1 cursor-pointer hover:underline">
                        ✕ {{ t("common.remove") }}
                    </button>
                </div>
            </div>

            <!-- Identity -->
            <SectionHeader :label="t('compare.identity')" />
            <ComparisonRow
                v-for="field in identityFields"
                :key="field.key"
                :label="field.label"
                :models="compareModels"
                :field-key="field.key"
                :format="field.format"
            />

            <!-- Pricing -->
            <SectionHeader :label="t('compare.pricing')" />
            <div
                v-if="hasPaidModels"
                class="cursor-pointer hover:bg-elevated/50 transition-colors"
                :style="{ minWidth: `${160 + compareModels.length * 200}px` }"
                @click="pricingModalOpen = true"
            >
                <CompareChart
                    type="bar"
                    :labels="compareModels.map((m) => m.name)"
                    :datasets="pricingDatasets(pricingTab)"
                    :tabs="pricingTabs"
                    :active-tab="pricingTab"
                    :height="160"
                    @tab-change="(k: string) => (pricingTab = k)"
                />
            </div>
            <div v-else class="p-5 border-b border-default text-center">
                <span class="inline-block bg-success/10 text-success text-sm font-medium px-4 py-2 rounded-lg">{{
                    t("common.free")
                }}</span>
            </div>
            <ComparisonRow
                v-for="field in pricingFields"
                :key="field.key"
                :label="field.label"
                :models="compareModels"
                :field-key="field.key"
                :format="field.format"
                :is-best="isCheapest"
                best-class="text-success"
            />

            <!-- Limits -->
            <SectionHeader :label="t('compare.limits')" />
            <div
                class="cursor-pointer hover:bg-elevated/50 transition-colors"
                :style="{ minWidth: `${160 + compareModels.length * 200}px` }"
                @click="limitsModalOpen = true"
            >
                <CompareChart
                    type="horizontalBar"
                    :labels="compareModels.map((m) => m.name)"
                    :datasets="limitDatasets"
                    :height="300"
                />
            </div>
            <ComparisonRow
                v-for="field in limitFields"
                :key="field.key"
                :label="field.label"
                :models="compareModels"
                :field-key="field.key"
                :format="formatTokens"
                :is-best="isBest"
                best-class="text-primary"
            />

            <!-- Capabilities -->
            <SectionHeader :label="t('compare.capabilities')" />
            <ComparisonRow
                v-for="cap in capabilityFields"
                :key="cap.key"
                :label="cap.label"
                :models="compareModels"
                :field-key="cap.key"
                boolean
            />

            <!-- Modalities -->
            <SectionHeader :label="t('compare.modalities')" />
            <div
                class="grid border-b border-default"
                :style="{
                    gridTemplateColumns: `160px repeat(${compareModels.length}, 200px)`,
                    minWidth: `${160 + compareModels.length * 200}px`,
                }"
            >
                <div
                    class="bg-elevated border-r border-b border-default p-2.5 px-4 text-xs text-muted flex items-center"
                >
                    {{ t("compare.input") }}
                </div>
                <div
                    v-for="m in compareModels"
                    :key="m.id + 'mod-in'"
                    class="border-r border-b border-default p-2.5 text-center"
                >
                    <div class="flex gap-1.5 justify-center flex-wrap">
                        <span
                            v-for="mod in m.modalities_input"
                            :key="mod"
                            :title="mod"
                            class="inline-flex items-center justify-center w-6 h-6 rounded-md"
                            :class="modalityInputClass(mod)"
                        >
                            <UIcon :name="modalityIcon(mod)" class="size-3.5" />
                        </span>
                    </div>
                </div>
            </div>
            <div
                class="grid border-b border-default"
                :style="{
                    gridTemplateColumns: `160px repeat(${compareModels.length}, 200px)`,
                    minWidth: `${160 + compareModels.length * 200}px`,
                }"
            >
                <div
                    class="bg-elevated border-r border-b border-default p-2.5 px-4 text-xs text-muted flex items-center"
                >
                    {{ t("compare.output") }}
                </div>
                <div
                    v-for="m in compareModels"
                    :key="m.id + 'mod-out'"
                    class="border-r border-b border-default p-2.5 text-center"
                >
                    <div class="flex gap-1.5 justify-center flex-wrap">
                        <span
                            v-for="mod in m.modalities_output"
                            :key="mod"
                            :title="mod"
                            class="inline-flex items-center justify-center w-6 h-6 rounded-md"
                            :class="modalityOutputClass()"
                        >
                            <UIcon :name="modalityIcon(mod)" class="size-3.5" />
                        </span>
                    </div>
                </div>
            </div>

            <!-- Interleaved -->
            <template v-if="compareModels.some((m) => m.interleaved)">
                <SectionHeader :label="t('compare.interleaved')" />
                <div
                    class="grid border-b border-default"
                    :style="{
                        gridTemplateColumns: `160px repeat(${compareModels.length}, 200px)`,
                        minWidth: `${160 + compareModels.length * 200}px`,
                    }"
                >
                    <div class="bg-elevated border-r border-default p-2.5 px-4 text-xs text-muted flex items-center">
                        {{ t("compare.interleaved") }}
                    </div>
                    <div
                        v-for="m in compareModels"
                        :key="m.id + 'interleaved'"
                        class="border-r border-default p-2.5 text-center text-sm"
                        :class="m.interleaved ? 'text-success font-semibold' : 'text-muted'"
                    >
                        {{ m.interleaved ? (typeof m.interleaved === "object" ? m.interleaved.field : "✓") : "—" }}
                    </div>
                </div>
            </template>

            <!-- Timeline -->
            <SectionHeader :label="t('compare.timeline')" />
            <ComparisonRow
                v-for="field in timelineFields"
                :key="field.key"
                :label="field.label"
                :models="compareModels"
                :field-key="field.key"
            />

            <!-- Integration -->
            <SectionHeader :label="t('compare.integration')" />
            <ComparisonRow
                v-for="field in integrationFields"
                :key="field.key"
                :label="field.label"
                :models="compareModels"
                :field-key="field.key"
                nested
                monospace
            />
        </div>

        <!-- Pricing chart modal -->
        <UModal v-model:open="pricingModalOpen">
            <template #header>
                <div class="flex items-center gap-3">
                    <span class="text-sm font-semibold text-default">{{ t("compare.pricing") }}</span>
                    <div class="flex gap-1">
                        <button
                            v-for="tab in pricingTabs"
                            :key="tab.key"
                            @click="pricingTab = tab.key"
                            class="rounded-md px-2.5 py-1 text-xs cursor-pointer transition-colors"
                            :class="
                                pricingTab === tab.key
                                    ? 'bg-primary text-white'
                                    : 'bg-elevated text-toned hover:bg-accented'
                            "
                        >
                            {{ tab.label }}
                        </button>
                    </div>
                </div>
            </template>
            <template #body>
                <div class="p-6">
                    <CompareChart
                        type="bar"
                        :labels="compareModels.map((m) => m.name)"
                        :datasets="pricingDatasets(pricingTab)"
                        :height="350"
                    />
                </div>
            </template>
        </UModal>

        <!-- Limits chart modal -->
        <UModal v-model:open="limitsModalOpen">
            <template #header>
                <span class="text-sm font-semibold text-default">{{ t("compare.limits") }}</span>
            </template>
            <template #body>
                <div class="p-6">
                    <CompareChart
                        type="horizontalBar"
                        :labels="compareModels.map((m) => m.name)"
                        :datasets="limitDatasets"
                        :height="300"
                    />
                </div>
            </template>
        </UModal>
    </div>
</template>

<script setup lang="ts">
    const { modelIds, removeModel, clearAll } = useCompare();
    const localePath = useLocalePath();
    const { t } = useI18n();

    const pricingTab = ref("all");
    const pricingModalOpen = ref(false);
    const limitsModalOpen = ref(false);

    // Mobile detection
    const isMobile = ref(false);
    onMounted(() => {
        isMobile.value = window.innerWidth < 768;
        const mq = window.matchMedia("(min-width: 768px)");
        const handler = (e: MediaQueryListEvent) => { isMobile.value = !e.matches; };
        mq.addEventListener("change", handler);
        onUnmounted(() => mq.removeEventListener("change", handler));
    });

    // Fetch compare data
    const { data: compareResult } = await useAsyncData(
        "compare-models",
        () => $fetch<{ data: any[] }>("/api/compare", { params: { ids: modelIds.value.join(",") } }),
        { watch: [modelIds] },
    );

    const compareModels = computed(() => compareResult.value?.data || []);

    const {
        identityFields,
        pricingFields,
        limitFields,
        capabilityFields,
        timelineFields,
        integrationFields,
        hasPaidModels,
        isCheapest,
        isBest,
        pricingTabs,
        pricingDatasets,
        limitDatasets,
    } = useCompareData(compareModels);
</script>
