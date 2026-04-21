<template>
    <div>
        <div class="flex items-center gap-2 px-1 mb-2">
            <UIcon :name="icon" class="size-3.5 text-primary" />
            <span class="text-xs font-semibold text-muted uppercase tracking-wider">{{ title }}</span>
        </div>
        <div :class="layout === 'compact' ? 'bg-default border border-default rounded-lg divide-y divide-default' : 'space-y-2'">
            <div
                v-for="field in fields"
                :key="field.key"
                :class="layout === 'compact' ? 'flex items-center gap-2 px-3 py-2' : 'bg-default border border-default rounded-lg px-3 py-2'"
            >
                <div v-if="layout !== 'compact'" class="text-[11px] text-muted mb-1.5">{{ field.label }}</div>
                <span v-else class="text-xs text-muted w-24 shrink-0">{{ field.label }}</span>
                <div :class="layout === 'compact' ? 'flex gap-2 flex-1 flex-wrap' : 'space-y-1'">
                    <div
                        v-for="(m, idx) in models"
                        :key="m.id"
                        class="flex items-center gap-2 text-xs"
                        :class="rowClass(m, field)"
                    >
                        <span class="rounded-full shrink-0" :class="layout === 'compact' ? 'size-1.5' : 'size-2'" :style="{ backgroundColor: colors[idx] }" />
                        <span v-if="layout !== 'compact'" class="text-muted min-w-0 truncate flex-1">{{ m.name }}</span>
                        <span class="shrink-0" :class="valueClass(m, field)">{{ formatValue(m, field) }}</span>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
    const props = defineProps<{
        title: string;
        icon: string;
        fields: { key: string; label: string; format?: (v: any) => string }[];
        models: any[];
        colors: string[];
        layout?: "default" | "compact";
        isBest?: (m: any, key: string) => boolean;
        bestHighlight?: string;
    }>();

    const rowClass = (m: any, field: { key: string }) => {
        if (!props.isBest?.(m, field.key)) return "";
        return props.bestHighlight === "success" ? "bg-success/5 -mx-1 px-1 rounded" : "bg-primary/5 -mx-1 px-1 rounded";
    };

    const valueClass = (m: any, field: { key: string }) => {
        if (props.layout === "compact") {
            return m[field.key] ? "text-success font-semibold" : "text-muted";
        }
        if (props.isBest?.(m, field.key)) {
            return props.bestHighlight === "success" ? "font-mono font-semibold text-success" : "font-mono font-semibold text-primary";
        }
        return "font-mono font-semibold text-default";
    };

    const formatValue = (m: any, field: { key: string; format?: (v: any) => string }) => {
        if (props.layout === "compact") return m[field.key] ? "✓" : "—";
        return field.format ? field.format(m[field.key]) : (m[field.key] ?? "—");
    };
</script>
