<template>
    <div class="grid border-b border-default" :style="{ gridTemplateColumns: `160px repeat(${models.length}, 200px)`, minWidth: `${160 + models.length * 200}px` }">
        <div class="bg-elevated border-r border-b border-default p-2.5 px-4 text-xs text-muted flex items-center">{{ label }}</div>
        <div v-for="m in models" :key="m.id + fieldKey" class="border-r border-b border-default last:border-r-0 p-2.5 text-center">
            <template v-if="boolean">
                <span class="text-sm" :class="m[fieldKey] ? 'text-success font-semibold' : 'text-muted'">{{ m[fieldKey] ? "✓" : "—" }}</span>
            </template>
            <template v-else-if="nested">
                <span class="text-xs font-mono text-toned">{{ getNestedValue(m, fieldKey) || "—" }}</span>
            </template>
            <template v-else>
                <span class="text-sm font-semibold" :class="bestClass && isBest?.(m, fieldKey) ? bestClass : 'text-default'">
                    {{ format ? format(m[fieldKey]) : (m[fieldKey] || "—") }}{{ bestClass && isBest?.(m, fieldKey) ? " ★" : "" }}
                </span>
            </template>
        </div>
    </div>
</template>

<script setup lang="ts">
defineProps<{
    label: string;
    models: any[];
    fieldKey: string;
    format?: (v: any) => string;
    boolean?: boolean;
    nested?: boolean;
    monospace?: boolean;
    isBest?: (m: any, field: string) => boolean;
    bestClass?: string;
}>();

const getNestedValue = (obj: any, path: string) => {
    return path.split(".").reduce((o, k) => o?.[k], obj);
};
</script>
