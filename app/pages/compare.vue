<template>
    <div>
        <div class="flex items-center justify-between mb-5">
            <div>
                <h1 class="text-xl font-bold text-default">{{ t("compare.title") }}</h1>
                <div class="text-sm text-muted">{{ compareModels.length }} / 4 {{ t("compare.modelsSelected") }}</div>
            </div>
            <div class="flex items-center gap-2">
                <button
                    v-if="compareModels.length > 0"
                    :title="t('compare.copyMd')"
                    class="inline-flex items-center gap-1.5 rounded-lg px-3 py-2 text-sm bg-elevated text-toned border border-default hover:border-primary hover:text-primary cursor-pointer transition-colors"
                    @click="copyMarkdown"
                >
                    <UIcon name="i-lucide-clipboard-copy" class="size-3.5" />
                    <span class="hidden sm:inline">{{ t("compare.copyMd") }}</span>
                </button>
                <button
                    v-if="compareModels.length > 1"
                    @click="diffOnly = !diffOnly"
                    class="inline-flex items-center gap-1.5 rounded-lg px-3 py-2 text-sm cursor-pointer transition-colors border"
                    :class="
                        diffOnly
                            ? 'bg-primary text-white border-primary'
                            : 'bg-elevated text-toned border-default hover:border-primary hover:text-primary'
                    "
                >
                    <UIcon :name="diffOnly ? 'i-lucide-git-compare-arrows' : 'i-lucide-list'" class="size-3.5" />
                    {{ diffOnly ? t("compare.diffOnly") : t("compare.showAll") }}
                </button>
                <CompareModelPicker v-if="modelIds.length < 4" :model-ids="modelIds" />
                <button
                    v-if="compareModels.length > 0"
                    @click="clearAll"
                    class="bg-elevated rounded-lg px-3.5 py-2 text-sm text-error cursor-pointer hover:bg-accented transition-colors"
                >
                    {{ t("common.clearAll") }}
                </button>
            </div>
        </div>

        <div v-if="compareModels.length === 0" class="flex flex-col items-center justify-center py-24 text-center">
            <div class="size-20 rounded-full bg-default border border-default flex items-center justify-center mb-5">
                <UIcon name="i-lucide-git-compare" class="size-8 text-muted" />
            </div>
            <div class="text-base font-medium text-default mb-1">{{ t("compare.empty") }}</div>
            <NuxtLink :to="localePath('/catalog')" class="inline-flex items-center gap-1 text-sm text-primary hover:underline mt-1">
                {{ t("nav.catalog") }}
                <UIcon name="i-lucide-arrow-right" class="size-3.5" />
            </NuxtLink>
        </div>

        <template v-else>
            <div v-if="!isMobile" class="grid grid-cols-1 lg:grid-cols-2 gap-4 mb-5">
                <div class="bg-default border border-default rounded-xl overflow-hidden">
                    <div class="flex items-center justify-between px-4 py-3 border-b border-default">
                        <span class="flex items-center gap-2 text-sm font-semibold text-default">
                            <UIcon name="i-lucide-dollar-sign" class="size-4 text-primary" />
                            {{ t("compare.pricing") }}
                        </span>
                        <div v-if="hasPaidModels" class="flex gap-1">
                            <button
                                v-for="tab in pricingTabs"
                                :key="tab.key"
                                @click="pricingTab = tab.key"
                                class="rounded-md px-2 py-0.5 text-[10px] cursor-pointer transition-colors"
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
                    <div v-if="hasPaidModels" class="p-4">
                        <CompareChart
                            type="bar"
                            :labels="modelNames"
                            :datasets="pricingDatasets(pricingTab)"
                            :height="200"
                        />
                    </div>
                    <div v-else class="p-6 text-center">
                        <span
                            class="inline-flex items-center gap-1.5 bg-success/10 text-success text-sm font-medium px-4 py-2 rounded-lg"
                        >
                            <UIcon name="i-lucide-badge-dollar-sign" class="size-4" />
                            {{ t("common.free") }}
                        </span>
                    </div>
                </div>

                <div class="bg-default border border-default rounded-xl overflow-hidden">
                    <div class="flex items-center gap-2 px-4 py-3 border-b border-default">
                        <UIcon name="i-lucide-gauge" class="size-4 text-primary" />
                        <span class="text-sm font-semibold text-default">{{ t("compare.limits") }}</span>
                    </div>
                    <div class="p-4">
                        <CompareChart
                            type="horizontalBar"
                            :labels="modelNames"
                            :datasets="limitDatasets"
                            :height="200"
                        />
                    </div>
                </div>
            </div>

            <div v-if="isMobile" class="space-y-4">
                <div class="flex flex-wrap gap-2">
                    <div
                        v-for="(m, idx) in compareModels"
                        :key="m.id"
                        class="inline-flex items-center gap-1.5 rounded-full pl-1.5 pr-2 py-1 bg-default border border-default"
                    >
                        <span class="size-2 rounded-full shrink-0" :style="{ backgroundColor: modelColors[idx] }" />
                        <span class="text-xs font-medium text-default truncate max-w-25">{{ m.name }}</span>
                        <button
                            @click="removeModel(m.id)"
                            class="text-muted hover:text-error cursor-pointer transition-colors shrink-0"
                        >
                            <UIcon name="i-lucide-x" class="size-3" />
                        </button>
                    </div>
                </div>

                <CompareMobileFieldList
                    :title="t('compare.pricing')"
                    icon="i-lucide-dollar-sign"
                    :fields="visPricing"
                    :models="compareModels"
                    :colors="modelColors"
                    :is-best="isCheapest"
                    best-highlight="success"
                />
                <CompareMobileFieldList
                    :title="t('compare.limits')"
                    icon="i-lucide-gauge"
                    :fields="visLimits"
                    :models="compareModels"
                    :colors="modelColors"
                    :is-best="isBest"
                    best-highlight="primary"
                />
                <CompareMobileFieldList
                    :title="t('compare.capabilities')"
                    icon="i-lucide-zap"
                    :fields="visCapabilities"
                    :models="compareModels"
                    :colors="modelColors"
                    layout="compact"
                />
                <CompareMobileFieldList
                    :title="t('compare.timeline')"
                    icon="i-lucide-clock"
                    :fields="visTimeline"
                    :models="compareModels"
                    :colors="modelColors"
                />
            </div>

            <div v-else class="border border-default rounded-xl overflow-x-auto bg-default">
                <div class="grid border-b border-default" :style="gridStyle">
                    <div class="bg-elevated border-r border-default" />
                    <div
                        v-for="(m, idx) in compareModels"
                        :key="m.id"
                        draggable="true"
                        :title="t('compare.reorder')"
                        class="border-r border-default last:border-r-0 py-4 px-3 text-center relative group overflow-hidden cursor-grab active:cursor-grabbing transition-colors"
                        :class="[
                            dragOver === idx && dragFrom !== idx ? 'bg-primary/10' : '',
                            dragFrom === idx ? 'opacity-40' : '',
                        ]"
                        :style="{ borderLeftColor: modelColors[idx], borderLeftWidth: '3px' }"
                        @dragstart="onDragStart(idx)"
                        @dragover.prevent="onDragOver(idx)"
                        @drop.prevent="onDrop(idx)"
                        @dragend="onDragEnd"
                    >
                        <UIcon
                            name="i-lucide-grip-vertical"
                            class="absolute top-1.5 left-1.5 size-3 text-muted opacity-0 group-hover:opacity-60 transition-opacity"
                        />
                        <button
                            @click="removeModel(m.id)"
                            class="absolute top-2 right-2 size-5 flex items-center justify-center rounded text-muted hover:text-error hover:bg-error/10 opacity-0 group-hover:opacity-100 transition-all cursor-pointer"
                        >
                            <UIcon name="i-lucide-x" class="size-3" />
                        </button>
                        <ProviderLogo :provider-id="m.provider_id" cls="w-9 h-9 rounded-lg mx-auto mb-2" />
                        <NuxtLink
                            :to="localePath(`/model/${m.id}`)"
                            class="font-semibold text-sm text-default hover:text-primary transition-colors no-underline"
                            >{{ m.name }}</NuxtLink
                        >
                        <div class="text-xs text-muted mb-2">{{ m.family }}</div>
                        <div class="flex gap-0.5 justify-center flex-wrap mb-2">
                            <span v-if="m.cost_input === 0 && m.cost_output === 0" :class="inlineBadgeClass('free')">{{
                                t("common.free")
                            }}</span>
                            <span v-if="m.reasoning" :class="inlineBadgeClass('reasoning')">R</span>
                            <span v-if="m.tool_call" :class="inlineBadgeClass('tool_call')">T</span>
                            <span v-if="isVision(m)" :class="inlineBadgeClass('vision')">V</span>
                            <span v-if="m.open_weights" :class="inlineBadgeClass('open_weights')">O</span>
                        </div>
                        <div v-if="m.cost_input > 0 || m.cost_output > 0" class="text-[11px] text-toned">
                            ${{ m.cost_input }} / ${{ m.cost_output }}
                        </div>
                        <div v-else-if="m.cost_input === 0" class="text-[11px] text-success font-medium">{{ t("common.free") }}</div>
                    </div>
                </div>

                <SectionHeader v-if="visIdentity.length" :label="t('compare.identity')" icon="i-lucide-fingerprint" />
                <ComparisonRow
                    v-for="field in visIdentity"
                    :key="field.key"
                    :label="field.label"
                    :models="compareModels"
                    :field-key="field.key"
                    :format="field.format"
                    :column-style="gridStyle"
                />

                <SectionHeader v-if="visPricing.length" :label="t('compare.pricing')" icon="i-lucide-dollar-sign" />
                <ComparisonRow
                    v-for="field in visPricing"
                    :key="field.key"
                    :label="field.label"
                    :models="compareModels"
                    :field-key="field.key"
                    :format="field.format"
                    :is-best="isCheapest"
                    best-class="text-success"
                    :best-label="t('compare.best')"
                    :column-style="gridStyle"
                />

                <SectionHeader v-if="vis(limitsFields).length" :label="t('compare.limits')" icon="i-lucide-gauge" />
                <ComparisonRow
                    v-for="field in visLimits"
                    :key="field.key"
                    :label="field.label"
                    :models="compareModels"
                    :field-key="field.key"
                    :format="formatTokens"
                    :is-best="isBest"
                    best-class="text-primary"
                    :best-label="t('compare.best')"
                    :column-style="gridStyle"
                />

                <SectionHeader v-if="vis(capabilitiesFields).length" :label="t('compare.capabilities')" icon="i-lucide-zap" />
                <ComparisonRow
                    v-for="cap in visCapabilities"
                    :key="cap.key"
                    :label="cap.label"
                    :models="compareModels"
                    :field-key="cap.key"
                    boolean
                    :column-style="gridStyle"
                />

                <SectionHeader v-if="!(diffOnly && modalitiesUniform)" :label="t('compare.modalities')" icon="i-lucide-layers" />
                <div v-if="!(diffOnly && modalitiesUniform)" class="grid border-b border-default" :style="gridStyle">
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
                                :class="modalityClass(mod)"
                            >
                                <UIcon :name="modalityIcon(mod)" class="size-3.5" />
                            </span>
                        </div>
                    </div>
                </div>
                <div v-if="!(diffOnly && modalitiesUniform)" class="grid border-b border-default" :style="gridStyle">
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
                                :class="modalityClass(mod)"
                            >
                                <UIcon :name="modalityIcon(mod)" class="size-3.5" />
                            </span>
                        </div>
                    </div>
                </div>

                <template v-if="compareModels.some((m) => m.interleaved)">
                    <SectionHeader :label="t('compare.interleaved')" icon="i-lucide-shuffle" />
                    <div class="grid border-b border-default" :style="gridStyle">
                        <div
                            class="bg-elevated border-r border-default p-2.5 px-4 text-xs text-muted flex items-center"
                        >
                            {{ t("compare.interleaved") }}
                        </div>
                        <div
                            v-for="m in compareModels"
                            :key="m.id + 'interleaved'"
                            class="border-r border-default p-2.5 text-center text-sm"
                            :class="m.interleaved ? 'text-success font-semibold' : 'text-muted'"
                        >
                        <template v-if="m.interleaved">
                            <span
                                v-if="typeof m.interleaved === 'object' && m.interleaved !== null && 'field' in (m.interleaved as object)"
                                class="font-mono"
                            >{{ (m.interleaved as any).field }}</span>
                            <UIcon v-else name="i-lucide-check" class="size-4 mx-auto" />
                        </template>
                        <span v-else class="text-muted">—</span>
                        </div>
                    </div>
                </template>

                <SectionHeader v-if="visTimeline.length" :label="t('compare.timeline')" icon="i-lucide-clock" />
                <ComparisonRow
                    v-for="field in visTimeline"
                    :key="field.key"
                    :label="field.label"
                    :models="compareModels"
                    :field-key="field.key"
                    :column-style="gridStyle"
                />

                <SectionHeader v-if="visIntegration.length" :label="t('compare.integration')" icon="i-lucide-plug" />
                <ComparisonRow
                    v-for="field in visIntegration"
                    :key="field.key"
                    :label="field.label"
                    :models="compareModels"
                    :field-key="field.key"
                    nested
                />
            </div>
        </template>
    </div>
</template>

<script setup lang="ts">
    import { inlineBadgeClass } from "~/utils/badge";

    const { modelIds, removeModel, clearAll, reorderModels } = useCompare();
    const localePath = useLocalePath();
    const { t } = useI18n();
    const config = useRuntimeConfig();
    const toast = useToast();

    useSeoMeta({
        title: t("seo.compareTitle"),
        ogTitle: t("seo.compareTitle"),
        description: t("seo.compareDescription"),
        ogDescription: t("seo.compareDescription"),
        keywords: t("seo.compareKeywords"),
        twitterCard: "summary_large_image",
    });

    const pricingTab = ref("all");
    const { isMobile } = useMobile();

    // ---- differences-only mode ----
    const diffOnly = ref(false);

    const resolveVal = (m: any, key: string) => key.split(".").reduce((o: any, k: string) => o?.[k], m);

    const differs = (field: { key: string; format?: (v: any) => string }) => {
        const values = new Set(
            compareModels.value.map((m: any) => {
                const v = resolveVal(m, field.key);
                return field.format ? field.format(v) : (v ?? "—");
            }),
        );
        return values.size > 1;
    };

    function vis<T extends { key: string; format?: (v: any) => string }>(fields: T[] | undefined): T[] {
        const list = fields ?? [];
        return diffOnly.value ? list.filter(differs) : list;
    }

    const modalitiesUniform = computed(() => {
        const ins = compareModels.value.map((m: any) => (m.modalities_input || []).join(","));
        const outs = compareModels.value.map((m: any) => (m.modalities_output || []).join(","));
        return new Set(ins).size <= 1 && new Set(outs).size <= 1;
    });

    // ---- column drag reorder (desktop) ----
    const dragFrom = ref(-1);
    const dragOver = ref(-1);
    const onDragStart = (idx: number) => (dragFrom.value = idx);
    const onDragOver = (idx: number) => {
        if (dragFrom.value >= 0 && dragFrom.value !== idx) dragOver.value = idx;
    };
    const onDrop = (idx: number) => {
        if (dragFrom.value >= 0) reorderModels(dragFrom.value, idx);
        dragFrom.value = -1;
        dragOver.value = -1;
    };
    const onDragEnd = () => {
        dragFrom.value = -1;
        dragOver.value = -1;
    };

    const copyMarkdown = async () => {
        const names = compareModels.value.map((m: any) => m.name);
        const lines: string[] = [
            `| ${t("compare.title")} | ${names.join(" | ")} |`,
            `|---|${names.map(() => "---").join("|")}|`,
        ];
        const section = (label: string, fields: { key: string; label: string; format?: (v: any) => string }[]) => {
            const visFields = vis(fields);
            if (!visFields.length) return;
            lines.push(`| **${label}** | ${names.map(() => "").join(" | ")} |`);
            for (const f of visFields) {
                const vals = compareModels.value.map((m: any) => {
                    const v = resolveVal(m, f.key);
                    const s = f.format ? f.format(v) : (v ?? "—");
                    return String(s).replace(/\|/g, "\\|");
                });
                lines.push(`| ${f.label} | ${vals.join(" | ")} |`);
            }
        };
        section(t("compare.identity"), identityFields.value);
        section(t("compare.pricing"), pricingFields.value);
        section(t("compare.limits"), limitFields.value);
        section(t("compare.capabilities"), capabilityFields.value);
        section(t("compare.timeline"), timelineFields.value);
        section(t("compare.integration"), integrationFields.value);
        try {
            await navigator.clipboard.writeText(lines.join("\n"));
            toast.add({ title: t("compare.copiedMd"), color: "success", icon: "i-lucide-check-circle" });
        } catch {}
    };

    const { data: compareResult } = await useAsyncData(
        "compare-models",
        () =>
            $fetch<{ data: any[] }>(`${config.public.apiBase}/api/v1/compare`, {
                params: { ids: modelIds.value.join(",") },
            }),
        { watch: [modelIds] },
    );

    const compareModels = computed(() => compareResult.value?.data || []);
    const modelNames = computed(() => compareModels.value.map((m: any) => m.name));
    const isVision = isVisionModel;

    const gridStyle = computed(() => {
        const n = compareModels.value.length;
        return { gridTemplateColumns: `160px repeat(${n}, 1fr)` };
    });

    const {
        modelColors,
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

    // Pre-computed visible field lists — keeps templates off fragile
    // binding paths and guards against undefined during render.
    const visIdentity = computed(() => vis(identityFields.value));
    const visPricing = computed(() => vis(pricingFields.value));
    const visLimits = computed(() => vis(limitFields.value));
    const visCapabilities = computed(() => vis(capabilityFields.value));
    const visTimeline = computed(() => vis(timelineFields.value));
    const visIntegration = computed(() => vis(integrationFields.value));
</script>
