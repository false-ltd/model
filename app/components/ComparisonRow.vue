<template>
    <div class="grid border-b border-default" :style="columnStyle || fallbackStyle">
        <div class="bg-elevated border-r border-b border-default p-2.5 px-4 text-xs text-muted flex items-center">{{ label }}</div>
        <div
            v-for="m in models"
            :key="m.id + fieldKey"
            class="border-r border-b border-default last:border-r-0 p-2.5 text-center transition-colors"
            :class="cellClass(m)"
        >
            <template v-if="boolean">
                <span class="text-sm" :class="m[fieldKey] ? 'text-success font-semibold' : 'text-muted'">{{ m[fieldKey] ? "✓" : "—" }}</span>
            </template>
            <template v-else-if="nested">
                <span class="text-xs font-mono text-toned">{{ getNestedValue(m, fieldKey) || "—" }}</span>
            </template>
            <template v-else>
                <span class="text-sm font-semibold" :class="valueClass(m)">
                    {{ format ? format(m[fieldKey]) : (m[fieldKey] || "—") }}
                </span>
                <span v-if="bestClass && isBest?.(m, fieldKey)" class="text-[9px] font-medium text-success ml-1">{{ bestLabel }}</span>
            </template>
        </div>
    </div>
</template>

<script setup lang="ts">
const props = defineProps<{
    label: string;
    models: any[];
    fieldKey: string;
    format?: (v: any) => string;
    boolean?: boolean;
    nested?: boolean;
    monospace?: boolean;
    isBest?: (m: any, field: string) => boolean;
    bestClass?: string;
    bestLabel?: string;
    columnStyle?: Record<string, string>;
}>();

const fallbackStyle = computed(() => {
    const n = props.models.length;
    return {
        gridTemplateColumns: `160px repeat(${n}, 1fr)`,
    };
});

const hasWinner = computed(() =>
    props.bestClass && props.models.some((m) => props.isBest?.(m, props.fieldKey)),
);

const cellClass = (m: any) => {
    if (!props.bestClass) return "";
    const winner = props.isBest?.(m, props.fieldKey);
    if (winner) return "bg-success/10";
    if (hasWinner.value) return "text-toned";
    return "";
};

const valueClass = (m: any) => {
    if (!props.bestClass) return "text-default";
    const winner = props.isBest?.(m, props.fieldKey);
    if (winner) return props.bestClass;
    if (hasWinner.value) return "text-toned";
    return "text-default";
};

const getNestedValue = (obj: any, path: string) => {
    return path.split(".").reduce((o, k) => o?.[k], obj);
};
</script>
