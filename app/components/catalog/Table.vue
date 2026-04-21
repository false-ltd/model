<template>
    <div class="border border-default rounded-xl overflow-hidden">
        <!-- Table -->
        <div class="overflow-x-auto">
        <UTable
            ref="table"
            sticky
            v-model:column-pinning="columnPinning"
            v-model:column-visibility="columnVisibility"
            :data="models"
            :columns="columns"
            :loading="loading"
            :ui="{
                tbody: '[&>tr]:hover:bg-accented/50 [&>tr]:transition-colors [&>tr]:cursor-pointer',
            }"
            class="flex-1 max-h-[calc(100vh)]"
            @row-click="(row: any) => handleRowClick(row)"
        >
            <template #loading>
                <div v-for="i in 8" :key="i" class="flex items-center gap-4 px-4 py-2.5">
                    <USkeleton class="size-5 rounded" />
                    <USkeleton class="h-4 w-40 rounded" />
                    <USkeleton class="h-4 w-16 rounded ml-auto" />
                    <USkeleton class="h-4 w-14 rounded" />
                    <USkeleton class="h-4 w-14 rounded" />
                </div>
            </template>
        </UTable>
        </div>

        <!-- Pagination -->
        <div v-if="totalPages > 1" class="flex items-center justify-between px-4 py-2.5 border-t border-default">
            <div class="flex items-center gap-2">
                <div class="text-xs text-muted">{{ rangeStart }}-{{ rangeEnd }} {{ t("common.of") }} {{ totalItems }}</div>
                <div class="flex items-center gap-0.5 bg-default border border-default rounded-md p-0.5">
                    <button
                        v-for="size in [20, 50, 100]"
                        :key="size"
                        @click="emit('changePageSize', size)"
                        class="rounded px-1.5 py-0.5 text-[11px] cursor-pointer transition-colors"
                        :class="pageSize === size ? 'bg-primary text-white' : 'text-muted hover:text-default'"
                    >
                        {{ size }}
                    </button>
                </div>
            </div>
            <UPagination
                :page="currentPage"
                :items-per-page="pageSize"
                :total="totalItems"
                @update:page="(p: number) => emit('goToPage', p)"
            />
        </div>
    </div>
</template>

<script setup lang="ts">
    import { h, resolveComponent } from "vue";
    import type { TableColumn } from "@nuxt/ui";
    import { inlineBadgeClass } from "~/utils/badge";
    import type { Model } from "~/types";

    const UCheckbox = resolveComponent("UCheckbox");
    const UButton = resolveComponent("UButton");
    const UIcon = resolveComponent("UIcon");
    const UTooltip = resolveComponent("UTooltip");
    const NuxtLinkComp = resolveComponent("NuxtLink");

    const props = defineProps<{
        models: Model[];
        loading: boolean;
        sortField: string;
        sortOrder: string;
        totalItems: number;
        totalPages: number;
        currentPage: number;
        pageSize: number;
    }>();

    const emit = defineEmits<{
        sort: [field: string];
        goToPage: [page: number];
        changePageSize: [size: number];
    }>();

    const localePath = useLocalePath();
    const { modelIds, addModel, removeModel } = useCompare();
    const { t } = useI18n();
    const toast = useToast();
    const table = useTemplateRef<any>("table");

    const copiedId = ref("");
    const copyModelId = async (id: string) => {
        try {
            await navigator.clipboard.writeText(id);
        } catch {}
        copiedId.value = id;
        setTimeout(() => {
            copiedId.value = "";
        }, 1500);
    };

    const handleRowClick = (row: any) => {
        if (row.original?.id) {
            navigateTo(localePath(`/model/${row.original.id}`));
        }
    };

    const showMaxToast = () => {
        toast.add({
            title: t('compare.maxReached'),
            description: t('compare.maxReachedHint'),
            color: 'warning' as const,
        });
    };

    const colLabels = computed<Record<string, string>>(() => ({
        name: t("catalog.colProviderModel"),
        family: t("catalog.colFamily"),
        provider_id: t("catalog.colProviderId"),
        model_id: t("catalog.colModelId"),
        cost_input: t("catalog.colInputCost"),
        cost_output: t("catalog.colOutputCost"),
        cost_reasoning: t("catalog.colReasoningCost"),
        cost_cache_read: t("catalog.colCacheRead"),
        cost_cache_write: t("catalog.colCacheWrite"),
        cost_input_audio: t("catalog.colInputAudio"),
        cost_output_audio: t("catalog.colOutputAudio"),
        limit_context: t("catalog.colContext"),
        limit_input: t("catalog.colInputLimit"),
        limit_output: t("catalog.colMaxOutput"),
        modalities_input: t("catalog.colInputMod"),
        modalities_output: t("catalog.colOutputMod"),
        knowledge: t("catalog.colKnowledge"),
        release_date: t("catalog.colReleased"),
        last_updated: t("catalog.colLastUpdated"),
    }));

    // Sortable header renderer
    function sortHeader(column: { id: string }, label: string) {
        const field = column.id;
        const active = props.sortField === field;
        const directionIcon =
            props.sortOrder === "asc" ? "i-lucide-arrow-up-narrow-wide" : "i-lucide-arrow-down-wide-narrow";
        return h(
            "button",
            {
                class:
                    "group/sort inline-flex items-center gap-1 whitespace-nowrap cursor-pointer transition-colors hover:text-default" +
                    (active ? " text-primary" : " text-muted"),
                onClick: () => emit("sort", field),
            },
            [
                h("span", label),
                h(
                    "span",
                    {
                        class: active
                            ? "inline-flex"
                            : "inline-flex opacity-0 group-hover/sort:opacity-100 transition-opacity",
                    },
                    h(UIcon, { name: active ? directionIcon : "i-lucide-arrow-up-down", class: "size-3.5" }),
                ),
            ],
        );
    }

    // Cost cell renderer
    const costCell =
        (key: string) =>
        ({ row }: { row: any }) => {
            const cost = row.original[key];
            if (cost == null) return "\u2014";
            return cost === 0 ? h("span", { class: "text-success font-semibold" }, t("common.free")) : `$${cost}`;
        };

    // Column definitions
    const columns = computed<TableColumn<any>[]>(() => [
        {
            id: "select",
            header: () => h("span", { class: "text-muted" }, "Select"),
            cell: ({ row }: { row: any }) =>
                h("div", { onClick: (e: Event) => e.stopPropagation() }, [
                    h(UCheckbox, {
                        modelValue: modelIds.value.includes(row.original.id),
                        "onUpdate:modelValue": () => {
                            modelIds.value.includes(row.original.id)
                                ? removeModel(row.original.id)
                                : (() => {
                                    const result = addModel(row.original.id);
                                    if (!result.added && result.reason === 'max') showMaxToast();
                                })();
                        },
                        "aria-label": "Select row",
                    }),
                ]),
            size: 50,
            enableHiding: false,
        },
        {
            id: "name",
            accessorKey: "name",
            header: ({ column }: { column: { id: string } }) => sortHeader(column, t("catalog.colProviderModel")),
            cell: ({ row }: { row: any }) => {
                const m = row.original;
                const rowChildren: any[] = [];
                if (!isMobile.value) {
                    rowChildren.push(
                        h("div", {
                            class: "w-6 h-6 rounded shrink-0 bg-elevated flex items-center justify-center text-[9px] font-bold text-muted relative overflow-hidden",
                        }, [
                            h("span", m.provider_id?.charAt(0).toUpperCase()),
                            h("img", {
                                src: `https://models.dev/logos/${m.provider_id}.svg`,
                                class: "absolute inset-0 w-full h-full object-cover rounded",
                                onError: (e: Event) => {
                                    (e.target as HTMLImageElement).style.display = "none";
                                },
                            }) as any,
                        ]),
                    );
                }
                const textCol: any[] = [
                    h(
                        NuxtLinkComp,
                        {
                            to: localePath(`/model/${m.id}`),
                            class: "font-medium text-default truncate no-underline hover:text-primary transition-colors",
                        },
                        () => m.name,
                    ),
                ];
                const badges: any[] = [];
                if (m.cost_input === 0) badges.push(h("span", { class: inlineBadgeClass("free") }, t("common.free")));
                if (m.reasoning) badges.push(h("span", { class: inlineBadgeClass("reasoning") }, t("catalog.reasoning")));
                if (m.tool_call) badges.push(h("span", { class: inlineBadgeClass("tool_call") }, t("catalog.colToolCall")));
                if (m.open_weights) badges.push(h("span", { class: inlineBadgeClass("open_weights") }, t("common.open")));
                if (m.structured_output) badges.push(h("span", { class: inlineBadgeClass("structured_output") }, t("catalog.colStructured")));
                if (m.attachment) badges.push(h("span", { class: inlineBadgeClass("attachment") }, t("catalog.colAttachment")));
                if (m.temperature) badges.push(h("span", { class: inlineBadgeClass("temperature") }, t("catalog.colTemperature")));
                if (badges.length) {
                    textCol.push(h("div", { class: "flex items-center gap-1 mt-0.5 flex-wrap" }, badges));
                }
                rowChildren.push(h("div", { class: "flex flex-col min-w-0" }, textCol));
                return h("div", { class: "flex items-center gap-2 min-w-0", onClick: (e: Event) => e.stopPropagation() }, rowChildren);
            },
            size: isMobile.value ? 130 : 260,
            enableHiding: false,
        },
        {
            id: "family",
            accessorKey: "family",
            header: ({ column }: { column: { id: string } }) => sortHeader(column, t("catalog.colFamily")),
            cell: ({ row }: { row: any }) => row.original.family || "\u2014",
            size: 120,
        },
        {
            id: "provider_id",
            accessorKey: "provider_id",
            header: ({ column }: { column: { id: string } }) => sortHeader(column, t("catalog.colProviderId")),
            cell: ({ row }: { row: any }) =>
                h("span", { class: "text-toned font-mono text-xs" }, row.original.provider_id || "\u2014"),
            size: 120,
        },
        {
            id: "model_id",
            accessorKey: "model_id",
            header: ({ column }: { column: { id: string } }) => sortHeader(column, t("catalog.colModelId")),
            cell: ({ row }: { row: any }) => {
                const mid = row.original.model_id;
                if (!mid) return "\u2014";
                return h(
                    "button",
                    {
                        class: "group/copy inline-flex items-center gap-1 text-toned font-mono text-xs cursor-pointer transition-colors hover:text-default",
                        onClick: (e: Event) => {
                            e.preventDefault();
                            e.stopPropagation();
                            copyModelId(mid);
                        },
                    },
                    [
                        h("span", { class: "truncate" }, mid),
                        h(UIcon, {
                            name: copiedId.value === mid ? "i-lucide-copy-check" : "i-lucide-copy",
                            class: "size-3 shrink-0 opacity-40 group-hover/copy:opacity-100 transition-opacity",
                        }),
                    ],
                );
            },
            size: 180,
        },
        {
            id: "cost_input",
            accessorKey: "cost_input",
            header: ({ column }: { column: { id: string } }) => sortHeader(column, t("catalog.colInputCost")),
            cell: costCell("cost_input"),
            size: 100,
        },
        {
            id: "cost_output",
            accessorKey: "cost_output",
            header: ({ column }: { column: { id: string } }) => sortHeader(column, t("catalog.colOutputCost")),
            cell: costCell("cost_output"),
            size: 100,
        },
        {
            id: "cost_reasoning",
            accessorKey: "cost_reasoning",
            header: ({ column }: { column: { id: string } }) => sortHeader(column, t("catalog.colReasoningCost")),
            cell: costCell("cost_reasoning"),
            size: 100,
        },
        {
            id: "cost_cache_read",
            accessorKey: "cost_cache_read",
            header: ({ column }: { column: { id: string } }) => sortHeader(column, t("catalog.colCacheRead")),
            cell: costCell("cost_cache_read"),
            size: 100,
        },
        {
            id: "cost_cache_write",
            accessorKey: "cost_cache_write",
            header: ({ column }: { column: { id: string } }) => sortHeader(column, t("catalog.colCacheWrite")),
            cell: costCell("cost_cache_write"),
            size: 100,
        },
        {
            id: "cost_input_audio",
            accessorKey: "cost_input_audio",
            header: ({ column }: { column: { id: string } }) => sortHeader(column, t("catalog.colInputAudio")),
            cell: costCell("cost_input_audio"),
            size: 100,
        },
        {
            id: "cost_output_audio",
            accessorKey: "cost_output_audio",
            header: ({ column }: { column: { id: string } }) => sortHeader(column, t("catalog.colOutputAudio")),
            cell: costCell("cost_output_audio"),
            size: 100,
        },
        {
            id: "limit_context",
            accessorKey: "limit_context",
            header: ({ column }: { column: { id: string } }) => sortHeader(column, t("catalog.colContext")),
            cell: ({ row }: { row: any }) => formatTokens(row.original.limit_context),
            size: 100,
        },
        {
            id: "limit_input",
            accessorKey: "limit_input",
            header: ({ column }: { column: { id: string } }) => sortHeader(column, t("catalog.colInputLimit")),
            cell: ({ row }: { row: any }) => formatTokens(row.original.limit_input),
            size: 100,
        },
        {
            id: "limit_output",
            accessorKey: "limit_output",
            header: ({ column }: { column: { id: string } }) => sortHeader(column, t("catalog.colMaxOutput")),
            cell: ({ row }: { row: any }) => formatTokens(row.original.limit_output),
            size: 100,
        },
        {
            id: "modalities_input",
            accessorKey: "modalities_input",
            header: () => h("span", { class: "text-muted" }, t("catalog.colInputMod")),
            cell: ({ row }: { row: any }) => {
                const mods: string[] = row.original.modalities_input || [];
                if (!mods.length) return h("span", { class: "text-muted" }, "\u2014");
                return h(
                    "div",
                    { class: "flex gap-1 justify-center" },
                    mods.map((mod: string) =>
                        h(UTooltip, { text: mod }, () =>
                            h(
                                "span",
                                {
                                    class: `inline-flex items-center justify-center size-6 rounded-md border ${modalityClass(mod)}`,
                                },
                                h(UIcon, { name: modalityIcon(mod), class: "size-3.5" }),
                            ),
                        ),
                    ),
                );
            },
            size: 120,
        },
        {
            id: "modalities_output",
            accessorKey: "modalities_output",
            header: () => h("span", { class: "text-muted" }, t("catalog.colOutputMod")),
            size: 140,
            cell: ({ row }: { row: any }) => {
                const mods: string[] = row.original.modalities_output || [];
                if (!mods.length) return h("span", { class: "text-muted" }, "\u2014");
                return h(
                    "div",
                    { class: "flex gap-1 justify-center " },
                    mods.map((mod: string) =>
                        h(UTooltip, { text: mod }, () =>
                            h(
                                "span",
                                {
                                    class: `inline-flex items-center justify-center size-6 rounded-md border ${modalityClass(mod)}`,
                                },
                                h(UIcon, { name: modalityIcon(mod), class: "size-3.5" }),
                            ),
                        ),
                    ),
                );
            },
        },
        {
            id: "knowledge",
            accessorKey: "knowledge",
            header: () => h("span", { class: "text-muted" }, t("catalog.colKnowledge")),
            cell: ({ row }: { row: any }) => h("span", { class: "text-toned" }, row.original.knowledge || "\u2014"),
            size: 100,
        },
        {
            id: "release_date",
            accessorKey: "release_date",
            header: ({ column }: { column: { id: string } }) => sortHeader(column, t("catalog.colReleased")),
            cell: ({ row }: { row: any }) => h("span", { class: "text-toned" }, row.original.release_date || "\u2014"),
            size: 110,
        },
        {
            id: "last_updated",
            accessorKey: "last_updated",
            header: ({ column }: { column: { id: string } }) => sortHeader(column, t("catalog.colLastUpdated")),
            cell: ({ row }: { row: any }) => h("span", { class: "text-toned" }, row.original.last_updated || "\u2014"),
            size: 110,
        },
    ]);

    const columnPinning = ref<{ left: string[] }>({ left: ["select", "name"] });

    const mobileHiddenColumns = [
        "select", "family", "provider_id", "cost_reasoning", "cost_cache_read", "cost_cache_write",
        "cost_input_audio", "cost_output_audio", "limit_input", "knowledge", "last_updated",
    ];
    const { isMobile } = useMobile();
    const columnVisibility = ref<Record<string, boolean>>({});
    onMounted(() => {
        if (isMobile.value) {
            columnVisibility.value = Object.fromEntries(mobileHiddenColumns.map((col) => [col, false]));
            columnPinning.value = { left: [] };
        }
        watch(isMobile, (mobile) => {
            if (mobile) {
                columnVisibility.value = Object.fromEntries(mobileHiddenColumns.map((col) => [col, false]));
                columnPinning.value = { left: [] };
            } else {
                columnPinning.value = { left: ["select", "name"] };
            }
        });
    });

    const columnMenuItems = computed(() => {
        const api: any = table.value?.tableApi;
        if (!api) return [];
        return (api.getAllColumns() as any[])
            .filter((col: any) => col.getCanHide())
            .map((col: any) => ({
                label: colLabels.value[col.id as string] || col.id,
                type: "checkbox" as const,
                checked: col.getIsVisible() as boolean,
                onUpdateChecked(checked: boolean) {
                    api.getColumn(col.id)?.toggleVisibility(!!checked);
                },
                onSelect(e: Event) {
                    e.preventDefault();
                },
            }));
    });

    const rangeStart = computed(() => (props.currentPage - 1) * props.pageSize + 1);
    const rangeEnd = computed(() => Math.min(props.currentPage * props.pageSize, props.totalItems));

    defineExpose({ columnMenuItems });
</script>
