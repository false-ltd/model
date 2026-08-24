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
                >{{ model.provider?.name }}</NuxtLink
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
                        <span class="text-xs font-mono text-muted">{{ model.provider?.name }}</span>
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
                    :title="t('detail.share')"
                    class="size-9 flex items-center justify-center rounded-lg border border-default bg-elevated text-toned hover:text-primary hover:border-primary/50 transition-colors cursor-pointer"
                    @click="copyLink"
                >
                    <UIcon :name="linkCopied ? 'i-lucide-check' : 'i-lucide-share-2'" class="size-4" :class="linkCopied ? 'text-success' : ''" />
                </button>
                <div class="flex items-center rounded-lg border border-default overflow-hidden">
                    <button
                        v-if="prevModel"
                        :title="`${t('detail.prevModel')}: ${prevModel.name}`"
                        class="size-9 flex items-center justify-center bg-elevated text-toned hover:text-primary hover:bg-accented transition-colors cursor-pointer"
                        @click="goSibling(prevModel.id)"
                    >
                        <UIcon name="i-lucide-arrow-left" class="size-4" />
                    </button>
                    <button
                        v-if="nextModel"
                        :title="`${t('detail.nextModel')}: ${nextModel.name}`"
                        class="size-9 flex items-center justify-center bg-elevated text-toned border-l border-default hover:text-primary hover:bg-accented transition-colors cursor-pointer"
                        @click="goSibling(nextModel.id)"
                    >
                        <UIcon name="i-lucide-arrow-right" class="size-4" />
                    </button>
                </div>
                <button
                    @click="
                        isInCompare
                            ? removeModel(model.id)
                            : addModel({ id: model.id, name: model.name, provider_id: model.provider_id })
                    "
                    class="inline-flex items-center gap-1.5 px-3.5 py-2 rounded-lg border text-sm font-medium cursor-pointer transition-all active:scale-95"
                    :class="isInCompare ? 'bg-primary text-white border-primary' : 'bg-elevated text-toned border-default hover:border-primary hover:text-primary'"
                >
                    <UIcon name="i-lucide-git-compare" class="size-3.5" />
                    <UIcon :name="isInCompare ? 'i-lucide-check' : 'i-lucide-plus'" class="size-3.5" />
                    {{ t("common.compare") }}
                </button>
                <a
                    v-if="model.provider?.doc_url"
                    :href="model.provider.doc_url"
                    target="_blank"
                    class="inline-flex items-center gap-1.5 px-3.5 py-2 rounded-lg bg-elevated text-toned text-sm font-medium no-underline hover:bg-accented transition-all active:scale-95"
                >
                    <UIcon name="i-lucide-book-open" class="size-3.5" />
                    {{ t("detail.docsLink") }}
                    <UIcon name="i-lucide-external-link" class="size-3" />
                </a>
                <a
                    v-if="model.provider?.api_url"
                    :href="model.provider.api_url"
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
                    <div v-if="priceBadges.length" class="ml-auto flex flex-wrap gap-1.5 normal-case tracking-normal">
                        <span
                            v-for="badge in priceBadges"
                            :key="badge.label"
                            class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-medium border"
                            :class="badge.cls"
                        >
                            <UIcon :name="badge.icon" class="size-3" />
                            {{ badge.label }}
                        </span>
                    </div>
                </div>
                <div class="flex-1 min-h-0">
                    <ModelPriceGauge :fields="pricingFields" />
                </div>
            </div>

            <div class="flex flex-col bg-default border border-default rounded-xl p-5">
                <div class="flex items-center gap-2 text-xs font-semibold uppercase tracking-widest text-muted mb-5">
                    <UIcon name="i-lucide-gauge" class="size-3.5 text-primary" />
                    {{ t("detail.limits") }}
                    <span
                        v-if="ctxPercentile != null"
                        class="ml-auto inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-medium border border-default bg-elevated text-toned normal-case tracking-normal"
                    >
                        <UIcon name="i-lucide-trending-up" class="size-3 text-primary" />
                        {{ t("detail.ctxBeats", { p: ctxPercentile }) }}
                    </span>
                </div>
                <div class="flex-1 min-h-0">
                    <ModelLimitGauge
                        :context-label="t('detail.contextWindow')"
                        :input-label="t('detail.maxInput')"
                        :output-label="t('detail.maxOutput')"
                        :context="model.limit_context ?? 0"
                        :input="model.limit_input ?? 0"
                        :output="model.limit_output ?? 0"
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
                                <NuxtTime v-if="item.date" :datetime="item.date" month="short" day="numeric" year="numeric" class="text-xs font-mono text-default" />
                                <span v-else class="text-xs font-mono text-default">{{ item.text }}</span>
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

                    <div v-if="model.provider?.npm">
                        <div class="text-xs text-muted mb-2 uppercase tracking-wider">{{ t("detail.npmPackage") }}</div>
                        <div class="flex items-center gap-2 bg-elevated border border-default rounded-lg px-3 py-2 group/npm hover:border-primary/30 transition-colors">
                            <UIcon name="i-lucide-package" class="size-3.5 text-muted group-hover/npm:text-success transition-colors shrink-0" />
                            <code class="text-xs font-mono text-success flex-1 truncate">{{ model.provider.npm }}</code>
                            <CopyButton :value="model.provider.npm" />
                        </div>
                    </div>

                    <div v-if="model.provider?.env?.length" :class="model.provider?.npm ? 'mt-4' : ''">
                        <div class="text-xs text-muted mb-2 uppercase tracking-wider">{{ t("detail.envVariable") }}</div>
                        <div class="space-y-1">
                            <div
                                v-for="envKey in model.provider.env"
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
                    {{ t("detail.similarModelsSub", { provider: model.provider?.name }) }}
                </span>
            </div>
            <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-4">
                <TransitionGroup name="card-list">
                    <ModelCard v-for="m in similarModels" :key="m.id" :model="m" />
                </TransitionGroup>
            </div>
        </section>
    </div>
</template>

<script setup lang="ts">
    // Retrigger the page transition when navigating between model pages.
    definePageMeta({ key: (route) => route.fullPath });

    const localePath = useLocalePath();
    const { modelIds, addModel, removeModel } = useCompare();
    const { t } = useI18n();

    const { model, pricingFields } = await useModelDetail();

    useSeoMeta({
        title: () => t("seo.detailTitle", { name: model.value?.name || "", provider: model.value?.provider?.name || "" }),
        ogTitle: () => t("seo.detailTitle", { name: model.value?.name || "", provider: model.value?.provider?.name || "" }),
        description: () => t("seo.detailDescription", { name: model.value?.name || "", provider: model.value?.provider?.name || "" }),
        ogDescription: () => t("seo.detailDescription", { name: model.value?.name || "", provider: model.value?.provider?.name || "" }),
        keywords: () => t("seo.detailKeywords", { name: model.value?.name || "", provider: model.value?.provider?.name || "" }),
        twitterCard: "summary",
    });

    const isInCompare = computed(() => model.value ? modelIds.value.includes(model.value.id) : false);

    const isVision = computed(() => model.value ? isVisionModel(model.value) : false);
    const hasAudio = computed(() =>
        model.value?.modalities_input?.includes("audio") || model.value?.modalities_output?.includes("audio"),
    );

    const hasIntegration = computed(() =>
        !!(model.value?.provider?.npm || model.value?.provider?.env?.length),
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

    // ---- reference context vs the whole catalog (shared stats cache) ----
    const { data: statsResult } = await useAsyncData("overview-stats", () =>
        $fetch<any>(`${config.public.apiBase}/api/v1/stats`),
    );

    const priceBadges = computed<{ label: string; cls: string; icon: string }[]>(() => {
        const stats = statsResult.value?.data?.stats;
        if (!stats?.medianInputPrice) return [];
        const badges: { label: string; cls: string; icon: string }[] = [];
        const ratio = (cost: number | null, median: number) =>
            cost == null || cost === 0 || !median ? null : cost / median;
        const tone = (r: number | null) =>
            r == null
                ? null
                : r <= 0.5
                  ? { cls: "bg-success/10 border-success/30 text-success", icon: "i-lucide-trending-down" }
                  : r <= 2
                    ? { cls: "bg-elevated border-default text-toned", icon: "i-lucide-equal" }
                    : { cls: "bg-amber-500/10 border-amber-500/30 text-amber-600 dark:text-amber-400", icon: "i-lucide-trending-up" };

        const rIn = ratio(model.value?.cost_input ?? null, stats.medianInputPrice);
        const tIn = tone(rIn);
        if (tIn) {
            badges.push({
                label:
                    rIn! <= 2
                        ? `${t("detail.inputPrice")} ${rIn! <= 0.5 ? t("detail.belowMedian") : t("detail.nearMedian")}`
                        : t("detail.aboveMedian", { field: t("detail.inputPrice"), ratio: rIn!.toFixed(1) }),
                ...tIn,
            });
        }
        // Output median isn't in stats; derive an approximate tone from input ratio only.
        return badges;
    });

    const ctxPercentile = computed<number | null>(() => {
        const dist: { limit_context: number }[] = statsResult.value?.data?.contextDistribution || [];
        const ctx = model.value?.limit_context;
        if (!dist.length || !ctx) return null;
        const below = dist.filter((d) => d.limit_context < ctx).length;
        return Math.round((below / dist.length) * 100);
    });

    // ---- same-provider prev/next navigation ----
    const { data: similarResult } = await useAsyncData(
        `similar-${model.value?.provider_id}`,
        () =>
            $fetch<any>(`${config.public.apiBase}/api/v1/models`, {
                params: { providers: model.value?.provider_id, page_size: 24, sort: "name", order: "asc" },
            }),
        { watch: [() => model.value?.provider_id] },
    );
    const siblingCycle = computed(() => {
        const list = (similarResult.value?.data ?? []).filter((m: any) => m.id !== model.value?.id);
        return list;
    });
    const prevModel = computed(() => {
        if (!siblingCycle.value.length) return null;
        const cur = model.value;
        const before = siblingCycle.value.filter((m: any) => m.name < (cur?.name ?? ""));
        return before.length ? before[before.length - 1] : siblingCycle.value[siblingCycle.value.length - 1];
    });
    const nextModel = computed(() => {
        if (!siblingCycle.value.length) return null;
        const cur = model.value;
        const after = siblingCycle.value.filter((m: any) => m.name > (cur?.name ?? ""));
        return after.length ? after[0] : siblingCycle.value[0];
    });
    const goSibling = (id: number) => navigateTo(localePath(`/model/${id}`));
    const linkCopied = ref(false);
    const copyLink = async () => {
        try {
            await navigator.clipboard.writeText(window.location.href);
            linkCopied.value = true;
            setTimeout(() => (linkCopied.value = false), 2000);
        } catch {}
    };
    onMounted(() => {
        const onKeydown = (e: KeyboardEvent) => {
            const el = e.target as HTMLElement;
            if (el && (el.tagName === "INPUT" || el.tagName === "TEXTAREA" || el.isContentEditable)) return;
            if (e.key === "ArrowLeft" && prevModel.value) goSibling(prevModel.value.id);
            if (e.key === "ArrowRight" && nextModel.value) goSibling(nextModel.value.id);
        };
        window.addEventListener("keydown", onKeydown);
        onUnmounted(() => window.removeEventListener("keydown", onKeydown));
    });
    const similarModels = computed(() => siblingCycle.value.slice(0, 4));

    const capabilityItems = computed(() => [
        { key: "reasoning", label: t("detail.reasoning"), value: model.value?.reasoning },
        { key: "tool_call", label: t("detail.toolCall"), value: model.value?.tool_call },
        { key: "structured_output", label: t("detail.structuredOutput"), value: model.value?.structured_output },
        { key: "attachment", label: t("detail.attachment"), value: model.value?.attachment },
        { key: "temperature", label: t("detail.temperature"), value: model.value?.temperature },
        { key: "open_weights", label: t("detail.openWeights"), value: model.value?.open_weights },
    ]);

    const timelineItems = computed(() => [
        { label: t("detail.released"), date: model.value?.release_date ?? null, text: "\u2014" },
        { label: t("detail.updated"), date: model.value?.last_updated ?? null, text: "\u2014" },
        { label: t("detail.knowledge"), date: null, text: model.value?.knowledge || "\u2014" },
        { label: t("detail.status"), date: null, text: model.value?.status || "\u2014" },
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
