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

        <!-- Warnings -->
        <Transition name="slide-down">
            <div
                v-if="model.status === 'deprecated'"
                class="flex items-center gap-2 px-4 py-2.5 rounded-lg mb-4 border bg-rose-500/10 border-rose-500/30 text-rose-600 dark:text-rose-400"
            >
                <UIcon name="i-lucide-alert-triangle" class="size-4 shrink-0" />
                <span class="text-sm font-medium">{{ t("detail.deprecatedWarning") }}</span>
            </div>
        </Transition>
        <Transition name="slide-down">
            <div
                v-if="model.status === 'beta'"
                class="flex items-center gap-2 px-4 py-2.5 rounded-lg mb-4 border bg-violet-500/10 border-violet-500/30 text-violet-600 dark:text-violet-400"
            >
                <UIcon name="i-lucide-flask-conical" class="size-4 shrink-0" />
                <span class="text-sm font-medium">{{ t("detail.betaWarning") }}</span>
            </div>
        </Transition>

        <!-- Hero -->
        <div class="flex items-start justify-between gap-6 flex-wrap mb-4">
            <div class="space-y-3">
                <h1 class="text-4xl font-black tracking-tight text-default leading-none">{{ model.name }}</h1>

                <div class="flex items-center gap-2 flex-wrap">
                    <div class="inline-flex items-center gap-2 px-3 py-1 bg-elevated border border-default rounded-full">
                        <ProviderLogo :provider-id="model.provider_id" cls="size-4 rounded" />
                        <span class="text-xs font-mono text-muted">{{ model.providers?.name }}</span>
                    </div>
                    <span v-if="model.family" class="text-xs text-muted">/</span>
                    <span v-if="model.family" class="text-xs font-mono text-muted">{{ model.family }}</span>
                    <div
                        class="inline-flex items-center gap-1.5 font-mono text-xs text-muted bg-elevated border border-default px-2 py-0.5 rounded-md cursor-pointer hover:border-primary/50 hover:text-primary transition-all group"
                        @click="copyModelId"
                    >
                        {{ model.model_id }}
                        <UIcon
                            :name="modelIdCopied ? 'i-lucide-check' : 'i-lucide-copy'"
                            class="size-3 transition-all"
                            :class="modelIdCopied ? 'text-success' : 'text-muted group-hover:text-primary'"
                        />
                    </div>
                </div>

                <div class="flex flex-wrap gap-1.5 pt-1">
                    <ModelBadge v-if="model.reasoning" type="reasoning" />
                    <ModelBadge v-if="model.tool_call" type="tool_call" />
                    <ModelBadge v-if="isVision" type="vision" />
                    <ModelBadge v-if="hasAudio" type="audio" />
                    <ModelBadge v-if="model.attachment" type="attachment" />
                    <ModelBadge v-if="model.open_weights" type="open_weights" />
                    <ModelBadge v-if="model.status === 'alpha'" type="alpha" />
                    <ModelBadge v-if="model.status === 'beta'" type="beta" />
                    <ModelBadge v-if="model.status === 'deprecated'" type="deprecated" />
                </div>
            </div>

            <div class="flex items-center gap-2 shrink-0">
                <button
                    @click="isInCompare ? removeModel(model.id) : addModel(model.id)"
                    class="inline-flex items-center gap-1.5 px-3.5 py-2 rounded-lg border text-sm font-medium cursor-pointer transition-all active:scale-95"
                    :class="isInCompare ? 'bg-primary text-white border-primary' : 'bg-elevated text-toned border-default hover:border-primary hover:text-primary'"
                >
                    <UIcon name="i-lucide-git-compare" class="size-3.5" />
                    {{ isInCompare ? "✓" : "+" }} {{ t("common.compare") }}
                </button>
                <a
                    v-if="model.providers?.doc_url"
                    :href="model.providers.doc_url"
                    target="_blank"
                    class="inline-flex items-center gap-1.5 px-3.5 py-2 rounded-lg bg-elevated text-toned text-sm font-medium no-underline hover:bg-accented transition-all active:scale-95"
                >
                    <UIcon name="i-lucide-book-open" class="size-3.5" />
                    {{ t("detail.docsLink") }}
                    <UIcon name="i-lucide-external-link" class="size-3" />
                </a>
                <a
                    v-if="model.providers?.api_url"
                    :href="model.providers.api_url"
                    target="_blank"
                    class="inline-flex items-center gap-1.5 px-3.5 py-2 rounded-lg bg-elevated text-toned text-sm font-medium no-underline hover:bg-accented transition-all active:scale-95"
                >
                    <UIcon name="i-lucide-terminal" class="size-3.5" />
                    API
                </a>
            </div>
        </div>

        <!-- Sticky Section Nav -->
        <div class="sticky top-14 z-10 bg-default/90 backdrop-blur-xl border-b border-default -mx-4 px-4 mb-6">
            <nav class="flex gap-0.5 overflow-x-auto scrollbar-none">
                <button
                    v-for="sec in sections"
                    :key="sec.id"
                    @click="scrollTo(sec.id)"
                    class="relative px-3 py-2.5 text-sm font-medium whitespace-nowrap transition-colors cursor-pointer"
                    :class="activeSection === sec.id ? 'text-primary' : 'text-muted hover:text-default'"
                >
                    <UIcon :name="sec.icon" class="size-3.5 mr-1.5 align-[-2px]" />
                    {{ sec.label }}
                    <span
                        v-if="activeSection === sec.id"
                        class="absolute bottom-0 left-2 right-2 h-0.5 bg-primary rounded-full transition-all duration-300"
                    />
                </button>
            </nav>
        </div>

        <!-- Section: Pricing & Limits Charts -->
        <section id="pricing" class="grid grid-cols-1 lg:grid-cols-2 gap-4 mb-6 scroll-mt-28">
            <div class="flex flex-col bg-default border border-default rounded-xl p-5">
                <div class="flex items-center gap-2 text-xs font-semibold uppercase tracking-widest text-muted mb-5">
                    <UIcon name="i-lucide-dollar-sign" class="size-3.5 text-primary" />
                    {{ t("detail.pricing") }}
                    <span class="normal-case tracking-normal text-muted/60">{{ t("detail.perMTokens") }}</span>
                </div>
                <div class="flex-1 min-h-0">
                    <ModelPriceGauge :fields="pricingFields" />
                </div>
            </div>

            <div class="flex flex-col bg-default border border-default rounded-xl p-5">
                <div class="flex items-center gap-2 text-xs font-semibold uppercase tracking-widest text-muted mb-5">
                    <UIcon name="i-lucide-gauge" class="size-3.5 text-primary" />
                    {{ t("detail.limits") }}
                </div>
                <div class="flex-1 min-h-0">
                    <ModelLimitGauge
                        :context-label="t('detail.contextWindow')"
                        :input-label="t('detail.maxInput')"
                        :output-label="t('detail.maxOutput')"
                        :context="model.limit_context"
                        :input="model.limit_input"
                        :output="model.limit_output"
                    />
                </div>
            </div>
        </section>

        <!-- Section: Overview -->
        <section id="overview" class="mb-6 scroll-mt-28">
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
                <!-- Capabilities -->
                <div class="bg-default border border-default rounded-xl p-5 hover:shadow-sm transition-shadow">
                    <div class="flex items-center gap-2 text-xs font-semibold uppercase tracking-widest text-muted mb-4">
                        <UIcon name="i-lucide-zap" class="size-3.5 text-primary" />
                        {{ t("detail.capabilities") }}
                    </div>
                    <div
                        v-for="cap in capabilityItems"
                        :key="cap.key"
                        class="flex items-center justify-between py-2.5 border-b border-default last:border-0 group/cap"
                    >
                        <span class="text-sm text-muted group-hover/cap:text-default transition-colors">{{ cap.label }}</span>
                        <UIcon v-if="cap.value" name="i-lucide-check-circle-2" class="size-4 text-success" />
                        <span v-else class="text-xs text-muted">&mdash;</span>
                    </div>
                </div>

                <!-- Modalities + Timeline -->
                <div class="bg-default border border-default rounded-xl p-5 hover:shadow-sm transition-shadow">
                    <div class="flex items-center gap-2 text-xs font-semibold uppercase tracking-widest text-muted mb-4">
                        <UIcon name="i-lucide-layers" class="size-3.5 text-primary" />
                        {{ t("detail.modalities") }}
                    </div>
                    <div class="mb-4">
                        <div class="text-xs text-muted mb-2 uppercase tracking-wider">{{ t("detail.inputModalities") }}</div>
                        <div class="flex flex-wrap gap-1.5">
                            <span
                                v-for="mod in model.modalities_input || []"
                                :key="mod"
                                class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-md border border-default text-xs font-mono transition-transform hover:scale-105"
                                :class="modalityClass(mod)"
                            ><UIcon :name="modalityIcon(mod)" class="size-3" />{{ mod }}</span>
                            <span v-if="!model.modalities_input?.length" class="text-xs text-muted">&mdash;</span>
                        </div>
                    </div>
                    <div>
                        <div class="text-xs text-muted mb-2 uppercase tracking-wider">{{ t("detail.outputModalities") }}</div>
                        <div class="flex flex-wrap gap-1.5">
                            <span
                                v-for="mod in model.modalities_output || []"
                                :key="mod"
                                class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-md border border-default text-xs font-mono transition-transform hover:scale-105"
                                :class="modalityClass(mod)"
                            ><UIcon :name="modalityIcon(mod)" class="size-3" />{{ mod }}</span>
                            <span v-if="!model.modalities_output?.length" class="text-xs text-muted">&mdash;</span>
                        </div>
                    </div>

                    <div class="mt-5 pt-4 border-t border-default">
                        <div class="flex items-center gap-2 text-xs font-semibold uppercase tracking-widest text-muted mb-3">
                            <UIcon name="i-lucide-clock" class="size-3.5 text-primary" />
                            {{ t("detail.timeline") }}
                        </div>
                        <div class="space-y-2.5">
                            <div v-for="item in timelineItems" :key="item.label" class="flex justify-between items-center">
                                <span class="text-xs text-muted">{{ item.label }}</span>
                                <span class="text-xs font-mono text-default">{{ item.value }}</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <!-- Section: Quick Start + Integration -->
        <section id="quickstart" class="mb-6 scroll-mt-28">
            <div class="grid grid-cols-1 lg:grid-cols-3 gap-4">
                <div class="lg:col-span-2">
                    <ModelQuickStartPanel :model="model" />
                </div>
                <div v-if="hasIntegration" class="bg-default border border-default rounded-xl p-5 hover:shadow-sm transition-shadow">
                    <div class="flex items-center gap-2 text-xs font-semibold uppercase tracking-widest text-muted mb-4">
                        <UIcon name="i-lucide-plug" class="size-3.5 text-primary" />
                        {{ t("detail.integration") }}
                    </div>

                    <div v-if="model.providers?.npm">
                        <div class="text-xs text-muted mb-2 uppercase tracking-wider">{{ t("detail.npmPackage") }}</div>
                        <div class="flex items-center gap-2 bg-elevated border border-default rounded-lg px-3 py-2 group/npm hover:border-primary/30 transition-colors">
                            <UIcon name="i-lucide-package" class="size-3.5 text-muted group-hover/npm:text-success transition-colors shrink-0" />
                            <code class="text-xs font-mono text-success flex-1 truncate">{{ model.providers.npm }}</code>
                            <CopyButton :value="model.providers.npm" />
                        </div>
                    </div>

                    <div v-if="model.providers?.env?.length" :class="model.providers?.npm ? 'mt-4' : ''">
                        <div class="text-xs text-muted mb-2 uppercase tracking-wider">{{ t("detail.envVariable") }}</div>
                        <div class="space-y-1">
                            <div
                                v-for="envKey in model.providers.env"
                                :key="envKey"
                                class="font-mono text-xs text-primary dark:text-primary bg-primary/5 border border-primary/20 rounded px-2 py-1 cursor-pointer hover:bg-primary/10 transition-colors"
                                @click="copyEnvVar(envKey)"
                            >
                                {{ envKey }}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <!-- Interleaved Config -->
        <section v-if="model.interleaved" class="bg-default border border-default rounded-xl p-5 mb-6">
            <div class="flex items-center gap-2 text-xs font-semibold uppercase tracking-widest text-muted mb-3">
                <UIcon name="i-lucide-shuffle" class="size-3.5 text-primary" />
                Interleaved Reasoning
            </div>
            <pre class="text-xs font-mono text-default bg-elevated border border-default rounded-lg p-4 overflow-x-auto">{{ JSON.stringify(model.interleaved, null, 2) }}</pre>
        </section>

        <!-- Section: Similar Models -->
        <section v-if="similarModels.length" class="mb-6">
            <div class="flex items-center gap-2 text-xs font-semibold uppercase tracking-widest text-muted mb-4">
                <UIcon name="i-lucide-layers" class="size-3.5 text-primary" />
                {{ t("detail.similarModels") }}
                <span class="normal-case tracking-normal text-muted/60">
                    {{ t("detail.similarModelsSub", { provider: model.providers?.name }) }}
                </span>
            </div>
            <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-4">
                <TransitionGroup name="card-list">
                    <ModelCard v-for="m in similarModels" :key="m.id" :model="m" />
                </TransitionGroup>
            </div>
        </section>

        <CompareFab />
    </div>
</template>

<script setup lang="ts">
    const localePath = useLocalePath();
    const { modelIds, addModel, removeModel } = useCompare();
    const { t } = useI18n();

    const { model, pricingFields } = await useModelDetail();

    const isInCompare = computed(() => modelIds.value.includes(model.value?.id));

    const isVision = computed(() => isVisionModel(model.value));
    const hasAudio = computed(() =>
        model.value?.modalities_input?.includes("audio") || model.value?.modalities_output?.includes("audio"),
    );

    const hasIntegration = computed(() =>
        !!(model.value?.providers?.npm || model.value?.providers?.env?.length),
    );

    const modelIdCopied = ref(false);
    const copyModelId = async () => {
        try {
            await navigator.clipboard.writeText(model.value?.model_id || "");
        } catch {}
        modelIdCopied.value = true;
        setTimeout(() => { modelIdCopied.value = false; }, 2000);
    };

    const copyEnvVar = async (key: string) => {
        try {
            await navigator.clipboard.writeText(key);
        } catch {}
    };

    const activeSection = ref("pricing");
    const sections = computed(() => [
        { id: "pricing", label: t("detail.pricing"), icon: "i-lucide-dollar-sign" },
        { id: "overview", label: t("detail.overview"), icon: "i-lucide-layout-grid" },
        { id: "quickstart", label: t("detail.quickStart"), icon: "i-lucide-terminal" },
    ]);

    const scrollTo = (id: string) => {
        const el = document.getElementById(id);
        if (el) el.scrollIntoView({ behavior: "smooth" });
    };

    onMounted(() => {
        const observer = new IntersectionObserver(
            (entries) => {
                for (const entry of entries) {
                    if (entry.isIntersecting) {
                        activeSection.value = entry.target.id;
                    }
                }
            },
            { rootMargin: "-120px 0px -60% 0px", threshold: 0 },
        );
        sections.value.forEach((sec) => {
            const el = document.getElementById(sec.id);
            if (el) observer.observe(el);
        });
        onUnmounted(() => observer.disconnect());
    });

    const config = useRuntimeConfig();
    const { data: similarResult } = await useAsyncData(
        `similar-${model.value?.provider_id}`,
        () =>
            $fetch<any>(`${config.public.apiBase}/api/v1/models`, {
                params: { providers: model.value?.provider_id, page_size: 5 },
            }),
        { watch: [() => model.value?.provider_id] },
    );
    const similarModels = computed(() => {
        const list = similarResult.value?.data ?? [];
        return list.filter((m: any) => m.id !== model.value?.id).slice(0, 4);
    });

    const capabilityItems = computed(() => [
        { key: "reasoning", label: t("detail.reasoning"), value: model.value?.reasoning },
        { key: "tool_call", label: t("detail.toolCall"), value: model.value?.tool_call },
        { key: "structured_output", label: t("detail.structuredOutput"), value: model.value?.structured_output },
        { key: "attachment", label: t("detail.attachment"), value: model.value?.attachment },
        { key: "temperature", label: t("detail.temperature"), value: model.value?.temperature },
        { key: "open_weights", label: t("detail.openWeights"), value: model.value?.open_weights },
    ]);

    const timelineItems = computed(() => [
        { label: t("detail.released"), value: model.value?.release_date || "\u2014" },
        { label: t("detail.updated"), value: model.value?.last_updated || "\u2014" },
        { label: t("detail.knowledge"), value: model.value?.knowledge || "\u2014" },
        { label: t("detail.status"), value: model.value?.status || "\u2014" },
    ]);
</script>

<style scoped>
.slide-down-enter-active,
.slide-down-leave-active {
    transition: all 0.3s ease;
}
.slide-down-enter-from,
.slide-down-leave-to {
    opacity: 0;
    transform: translateY(-8px);
}

.card-list-enter-active {
    transition: all 0.4s ease;
}
.card-list-enter-from {
    opacity: 0;
    transform: translateY(12px);
}

.scrollbar-none::-webkit-scrollbar {
    display: none;
}
.scrollbar-none {
    -ms-overflow-style: none;
    scrollbar-width: none;
}
</style>
