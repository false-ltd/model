<template>
    <div class="flex flex-col gap-1">
        <div class="flex justify-between items-center text-sm">
            <span class="text-toned">{{ label }}</span>
            <span class="font-semibold text-default">{{ formatted }}</span>
        </div>
        <div class="h-2 bg-elevated rounded overflow-hidden">
            <div class="h-full rounded transition-all" :style="{ width: pct + '%', background: color || 'var(--ui-color-primary-500)' }" />
        </div>
    </div>
</template>

<script setup lang="ts">
    const props = defineProps<{
        label: string;
        value: number;
        max: number;
        color?: string;
    }>();

    const pct = computed(() => Math.min(100, props.max > 0 ? (props.value / props.max) * 100 : 0));
    const formatted = computed(() => {
        if (props.value >= 1_000_000) return (props.value / 1_000_000).toFixed(0) + "M tokens";
        if (props.value >= 1_000) return (props.value / 1_000).toFixed(0) + "K tokens";
        return String(props.value) + " tokens";
    });
</script>
