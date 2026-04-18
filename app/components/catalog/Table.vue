<template>
    <div class="border border-default rounded-xl">
        <!-- Table -->
        <UTable
            ref="table"
            sticky
            v-model:column-pinning="columnPinning"
            v-model:column-visibility="columnVisibility"
            :data="models"
            :columns="columns"
            class="flex-1 max-h-[calc(100vh)]"
        />

        <!-- Pagination -->
        <div v-if="totalPages > 1" class="flex items-center justify-between px-4 py-2.5 border-t border-default">
            <div class="text-xs text-muted"></div>
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

    const UCheckbox = resolveComponent("UCheckbox");
    const UButton = resolveComponent("UButton");
    const UIcon = resolveComponent("UIcon");
    const UTooltip = resolveComponent("UTooltip");
    const NuxtLinkComp = resolveComponent("NuxtLink");

    const props = defineProps<{
        models: any[];
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
    }>();

    const localePath = useLocalePath();
    const { modelIds, addModel, removeModel } = useCompare();
    const { t } = useI18n();
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

    // Column visibility labels
    const colLabels: Record<string, string> = {
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
        tool_call: t("catalog.colToolCall"),
        reasoning: t("catalog.colReasoning"),
        modalities_input: t("catalog.colInputMod"),
        modalities_output: t("catalog.colOutputMod"),
        open_weights: t("catalog.colWeights"),
        structured_output: t("catalog.colStructured"),
        attachment: t("catalog.colAttachment"),
        temperature: t("catalog.colTemperature"),
        knowledge: t("catalog.colKnowledge"),
        release_date: t("catalog.colReleased"),
        last_updated: t("catalog.colLastUpdated"),
    };

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
                h(UCheckbox, {
                    modelValue: modelIds.value.includes(row.original.id),
                    "onUpdate:modelValue": () => {
                        modelIds.value.includes(row.original.id)
                            ? removeModel(row.original.id)
                            : addModel(row.original.id);
                    },
                    "aria-label": "Select row",
                }),
            size: 50,
            enableHiding: false,
        },
        {
            id: "name",
            accessorKey: "name",
            header: ({ column }: { column: { id: string } }) => sortHeader(column, t("catalog.colProviderModel")),
            cell: ({ row }: { row: any }) => {
                const m = row.original;
                return h("div", { class: "flex items-center gap-2 min-w-0" }, [
                    h("img", {
                        src: `https://models.dev/logos/${m.provider_id}.svg`,
                        class: "w-6 h-6 rounded shrink-0",
                        onError: (e: Event) => {
                            (e.target as HTMLImageElement).style.display = "none";
                        },
                    }),
                    h(
                        NuxtLinkComp,
                        {
                            to: localePath(`/model/${m.id}`),
                            class: "font-medium text-default truncate no-underline hover:text-primary transition-colors",
                        },
                        () => m.name,
                    ),
                ]);
            },
            size: 240,
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
            id: "tool_call",
            accessorKey: "tool_call",
            header: ({ column }: { column: { id: string } }) => sortHeader(column, t("catalog.colToolCall")),
            cell: ({ row }: { row: any }) => {
                const val = row.original.tool_call;
                return h("span", { class: val ? "text-success" : "text-muted" }, val ? t("common.yes") : "\u2014");
            },
            size: 80,
        },
        {
            id: "reasoning",
            accessorKey: "reasoning",
            header: ({ column }: { column: { id: string } }) => sortHeader(column, t("catalog.colReasoning")),
            cell: ({ row }: { row: any }) => {
                const val = row.original.reasoning;
                return h("span", { class: val ? "text-warning" : "text-muted" }, val ? t("common.yes") : "\u2014");
            },
            size: 80,
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
                                    class: `inline-flex items-center justify-center size-6 rounded-md ${modalityInputClass(mod)}`,
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
                                    class: `inline-flex items-center justify-center size-6 rounded-md ${modalityOutputClass()}`,
                                },
                                h(UIcon, { name: modalityIcon(mod), class: "size-3.5" }),
                            ),
                        ),
                    ),
                );
            },
        },
        {
            id: "open_weights",
            accessorKey: "open_weights",
            header: ({ column }: { column: { id: string } }) => sortHeader(column, t("catalog.colWeights")),
            cell: ({ row }: { row: any }) => {
                const val = row.original.open_weights;
                return h(
                    "span",
                    { class: val ? "text-success" : "text-muted" },
                    val ? t("common.open") : t("common.closed"),
                );
            },
            size: 80,
        },
        {
            id: "structured_output",
            accessorKey: "structured_output",
            header: ({ column }: { column: { id: string } }) => sortHeader(column, t("catalog.colStructured")),
            cell: ({ row }: { row: any }) => {
                const val = row.original.structured_output;
                return h("span", { class: val ? "text-success" : "text-muted" }, val ? t("common.yes") : "\u2014");
            },
            size: 80,
        },
        {
            id: "attachment",
            accessorKey: "attachment",
            header: ({ column }: { column: { id: string } }) => sortHeader(column, t("catalog.colAttachment")),
            cell: ({ row }: { row: any }) => {
                const val = row.original.attachment;
                return h("span", { class: val ? "text-primary" : "text-muted" }, val ? t("common.yes") : "\u2014");
            },
            size: 80,
        },
        {
            id: "temperature",
            accessorKey: "temperature",
            header: ({ column }: { column: { id: string } }) => sortHeader(column, t("catalog.colTemperature")),
            cell: ({ row }: { row: any }) => {
                const val = row.original.temperature;
                return h("span", { class: val ? "text-success" : "text-muted" }, val ? t("common.yes") : "\u2014");
            },
            size: 80,
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

    // Column pinning
    const columnPinning = ref({
        left: ["select", "name"],
    });

    // Column visibility
    const columnVisibility = ref<Record<string, boolean>>({});

    // Column visibility dropdown items
    const columnMenuItems = computed(() => {
        const api: any = table.value?.tableApi;
        if (!api) return [];
        return (api.getAllColumns() as any[])
            .filter((col: any) => col.getCanHide())
            .map((col: any) => ({
                label: colLabels[col.id as string] || col.id,
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

    // Pagination range
    const rangeStart = computed(() => (props.currentPage - 1) * props.pageSize + 1);
    const rangeEnd = computed(() => Math.min(props.currentPage * props.pageSize, props.totalItems));

    defineExpose({ columnMenuItems });
</script>
