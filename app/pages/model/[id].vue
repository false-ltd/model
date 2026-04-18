<template>
    <div v-if="model">
        <!-- Breadcrumb -->
        <nav class="text-sm text-muted mb-5">
            <NuxtLink :to="localePath('/catalog')" class="text-toned hover:text-primary transition-colors">{{
                t("nav.catalog")
            }}</NuxtLink>
            <span class="mx-1.5">/</span>
            <NuxtLink
                :to="localePath(`/catalog?providers=${model.provider_id}`)"
                class="text-toned hover:text-primary transition-colors"
                >{{ model.providers?.name }}</NuxtLink
            >
            <span class="mx-1.5">/</span>
            <span class="text-default font-medium">{{ model.name }}</span>
        </nav>

        <!-- Hero -->
        <div class="bg-default border border-default rounded-xl p-6 mb-5 hover:border-primary/40 transition-all">
            <div class="flex flex-col sm:flex-row items-start gap-4">
                <img
                    :src="`https://models.dev/logos/${model.provider_id}.svg`"
                    class="w-12 h-12 rounded-xl p-1.5"
                    @error="($event.target as HTMLImageElement).style.display = 'none'"
                />
                <div class="flex-1 min-w-0">
                    <div class="flex items-center gap-2.5 mb-1 flex-wrap">
                        <h1 class="text-2xl font-bold text-default">{{ model.name }}</h1>
                        <ModelBadge v-if="model.reasoning" type="reasoning" />
                        <ModelBadge v-if="model.open_weights" type="open_weights" />
                    </div>
                    <div class="text-sm text-muted mb-3">
                        {{ model.model_id }} · {{ t("detail.family") }}: {{ model.family }} · {{ t("detail.released")
                        }}
                        {{ model.release_date }} · {{ t("detail.updated") }} {{ model.last_updated }}
                    </div>
                    <div class="flex gap-1.5 flex-wrap">
                        <ModelBadge v-if="model.tool_call" type="tool_call" />
                        <ModelBadge v-if="model.attachment" type="attachment" />
                        <ModelBadge v-if="model.temperature" type="temperature" />
                        <ModelBadge
                            v-for="mod in model.modalities_input"
                            :key="mod"
                            :type="mod === 'image' ? 'vision' : 'family'"
                        />
                    </div>
                </div>
                <div class="flex gap-2 shrink-0">
                    <button
                        @click="isInCompare ? removeModel(model.id) : addModel(model.id)"
                        class="rounded-lg px-4 py-2 text-sm font-medium cursor-pointer transition-colors"
                        :class="isInCompare ? 'bg-primary text-white' : 'bg-elevated text-toned border border-default hover:border-primary hover:text-primary'"
                    >
                        {{ isInCompare ? "✓" : "+" }} {{ t("common.compare") }}
                    </button>
                    <a
                        v-if="model.providers?.doc"
                        :href="model.providers.doc"
                        target="_blank"
                        class="bg-elevated rounded-lg px-4 py-2 text-sm text-toned no-underline hover:bg-accented transition-colors"
                        >{{ t("detail.docsLink") }} ↗</a
                    >
                </div>
            </div>
        </div>

        <!-- Pricing -->
        <div class="bg-default border border-default rounded-xl p-5 mb-4 hover:border-primary/40 transition-all">
            <div class="text-xs text-muted uppercase tracking-wider font-medium mb-3.5">
                {{ t("detail.pricing") }} <span class="normal-case tracking-normal">{{ t("detail.perMTokens") }}</span>
            </div>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
                <div v-for="field in pricingFields" :key="field.key" class="bg-elevated rounded-lg p-3">
                    <div class="text-xs text-muted mb-0.5">{{ field.label }}</div>
                    <div class="flex items-center gap-1.5">
                        <span class="text-xl font-bold text-default">{{ field.display }}</span>
                        <CopyButton :value="field.display" />
                    </div>
                </div>
            </div>
        </div>

        <!-- Limits -->
        <div class="bg-default border border-default rounded-xl p-5 mb-4 hover:border-primary/40 transition-all">
            <div class="text-xs text-muted uppercase tracking-wider font-medium mb-3.5">{{ t("detail.limits") }}</div>
            <div class="flex flex-col gap-3">
                <ModelLimitBar
                    :label="t('detail.contextWindow')"
                    :value="model.limit_context"
                    :max="2_000_000"
                    color="#e11d48"
                />
                <ModelLimitBar
                    :label="t('detail.maxInput')"
                    :value="model.limit_input"
                    :max="2_000_000"
                    :color="chartColors[1]!"
                />
                <ModelLimitBar
                    :label="t('detail.maxOutput')"
                    :value="model.limit_output"
                    :max="128_000"
                    :color="chartColors[2]!"
                />
            </div>
        </div>

        <!-- Capabilities -->
        <div class="bg-default border border-default rounded-xl p-5 mb-4 hover:border-primary/40 transition-all">
            <div class="text-xs text-muted uppercase tracking-wider font-medium mb-3.5">
                {{ t("detail.capabilities") }}
            </div>
            <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-3">
                <div
                    v-for="cap in capabilityItems"
                    :key="cap.key"
                    class="bg-elevated rounded-lg p-3 text-center"
                    :class="cap.value ? 'ring-1 ring-success/30' : ''"
                >
                    <div class="text-lg mb-0.5" :class="cap.value ? 'text-success' : 'text-muted'">
                        {{ cap.value ? "✓" : "—" }}
                    </div>
                    <div class="text-xs" :class="cap.value ? 'text-default font-medium' : 'text-muted'">
                        {{ cap.label }}
                    </div>
                </div>
            </div>
        </div>

        <!-- Modalities -->
        <div class="bg-default border border-default rounded-xl p-5 mb-4 hover:border-primary/40 transition-all">
            <div class="text-xs text-muted uppercase tracking-wider font-medium mb-3.5">
                {{ t("detail.modalities") }}
            </div>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                    <div class="text-xs text-muted mb-2">{{ t("detail.inputModalities") }}</div>
                    <div class="flex gap-2 flex-wrap">
                        <span
                            v-for="mod in model.modalities_input || []"
                            :key="mod"
                            class="inline-flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg text-sm font-medium"
                            :class="modalityInputClass(mod)"
                        >
                            <UIcon :name="modalityIcon(mod)" class="size-4" />
                            {{ mod }}
                        </span>
                        <span v-if="!model.modalities_input?.length" class="text-sm text-muted">—</span>
                    </div>
                </div>
                <div>
                    <div class="text-xs text-muted mb-2">{{ t("detail.outputModalities") }}</div>
                    <div class="flex gap-2 flex-wrap">
                        <span
                            v-for="mod in model.modalities_output || []"
                            :key="mod"
                            class="inline-flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg text-sm font-medium"
                            :class="modalityOutputClass()"
                        >
                            <UIcon :name="modalityIcon(mod)" class="size-4" />
                            {{ mod }}
                        </span>
                        <span v-if="!model.modalities_output?.length" class="text-sm text-muted">—</span>
                    </div>
                </div>
            </div>
        </div>

        <!-- Timeline -->
        <div class="bg-default border border-default rounded-xl p-5 mb-4 hover:border-primary/40 transition-all">
            <div class="text-xs text-muted uppercase tracking-wider font-medium mb-3.5">{{ t("detail.timeline") }}</div>
            <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
                <div class="bg-elevated rounded-lg p-3">
                    <div class="text-xs text-muted mb-0.5">{{ t("detail.released") }}</div>
                    <div class="text-sm font-semibold text-default">{{ model.release_date || "—" }}</div>
                </div>
                <div class="bg-elevated rounded-lg p-3">
                    <div class="text-xs text-muted mb-0.5">{{ t("detail.updated") }}</div>
                    <div class="text-sm font-semibold text-default">{{ model.last_updated || "—" }}</div>
                </div>
                <div class="bg-elevated rounded-lg p-3">
                    <div class="text-xs text-muted mb-0.5">{{ t("detail.knowledge") }}</div>
                    <div class="text-sm font-semibold text-default">{{ model.knowledge || "—" }}</div>
                </div>
                <div class="bg-elevated rounded-lg p-3">
                    <div class="text-xs text-muted mb-0.5">{{ t("detail.status") }}</div>
                    <div class="text-sm font-semibold text-default">{{ model.status || "—" }}</div>
                </div>
            </div>
        </div>

        <!-- Integration -->
        <div class="bg-default border border-default rounded-xl p-5 mb-4 hover:border-primary/40 transition-all">
            <div class="text-xs text-muted uppercase tracking-wider font-medium mb-3.5">
                {{ t("detail.integration") }}
            </div>
            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                <div>
                    <div class="text-xs text-muted mb-1">{{ t("detail.provider") }}</div>
                    <div class="flex items-center gap-1.5">
                        <img
                            :src="`https://models.dev/logos/${model.provider_id}.svg`"
                            class="w-5 h-5 rounded"
                            @error="($event.target as HTMLImageElement).style.display = 'none'"
                        />
                        <span class="text-sm font-medium text-default">{{ model.providers?.name }}</span>
                    </div>
                </div>
                <div>
                    <div class="text-xs text-muted mb-1">{{ t("detail.npmPackage") }}</div>
                    <div class="flex items-center gap-1">
                        <span class="text-sm font-medium text-default bg-elevated px-2 py-0.5 rounded font-mono">{{
                            model.providers?.npm
                        }}</span>
                        <CopyButton :value="model.providers?.npm || ''" />
                    </div>
                </div>
                <div>
                    <div class="text-xs text-muted mb-1">{{ t("detail.apiEndpoint") }}</div>
                    <div class="flex items-center gap-1">
                        <span class="text-xs text-toned font-mono bg-elevated px-2 py-0.5 rounded">{{
                            model.providers?.api
                        }}</span>
                        <CopyButton :value="model.providers?.api || ''" />
                    </div>
                </div>
                <div>
                    <div class="text-xs text-muted mb-1">{{ t("detail.envVariable") }}</div>
                    <div class="flex items-center gap-1">
                        <span class="text-sm font-mono text-default bg-elevated px-2 py-0.5 rounded">{{
                            (model.providers?.env || [])[0]
                        }}</span>
                        <CopyButton :value="(model.providers?.env || [])[0] || ''" />
                    </div>
                </div>
                <div>
                    <div class="text-xs text-muted mb-1">{{ t("detail.modelId") }}</div>
                    <div class="flex items-center gap-1">
                        <span class="text-sm font-mono text-default bg-elevated px-2 py-0.5 rounded">{{
                            model.model_id
                        }}</span>
                        <CopyButton :value="model.model_id" />
                    </div>
                </div>
                <div>
                    <div class="text-xs text-muted mb-1">{{ t("detail.documentation") }}</div>
                    <a
                        v-if="model.providers?.doc"
                        :href="model.providers.doc"
                        target="_blank"
                        class="text-sm text-primary underline hover:no-underline"
                        >{{ t("detail.docsLink") }} ↗</a
                    >
                </div>
            </div>
        </div>

        <!-- Quick Start -->
        <div class="bg-default border border-default rounded-xl p-5 hover:border-primary/40 transition-all">
            <div class="flex justify-between items-center mb-2.5">
                <div class="text-xs text-muted uppercase tracking-wider font-medium">{{ t("detail.quickStart") }}</div>
                <CopyButton :value="quickStartCode" />
            </div>
            <div class="bg-[#1a1a1a] rounded-lg p-4 font-mono text-sm text-gray-200 leading-relaxed overflow-x-auto">
                <pre>{{ quickStartCode }}</pre>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
    const localePath = useLocalePath();
    const { modelIds, addModel, removeModel } = useCompare();
    const { t } = useI18n();

    const { model, pricingFields, quickStartCode } = await useModelDetail();

    const isInCompare = computed(() => modelIds.value.includes(model.value?.id));

    const capabilityItems = computed(() => [
        { key: "reasoning", label: t("detail.reasoning"), value: model.value?.reasoning },
        { key: "tool_call", label: t("detail.toolCall"), value: model.value?.tool_call },
        { key: "attachment", label: t("detail.attachment"), value: model.value?.attachment },
        { key: "structured_output", label: t("detail.structuredOutput"), value: model.value?.structured_output },
        { key: "temperature", label: t("detail.temperature"), value: model.value?.temperature },
    ]);
</script>
