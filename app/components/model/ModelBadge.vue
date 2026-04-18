<template>
    <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium" :class="badgeClass">
        {{ label }}
    </span>
</template>

<script setup lang="ts">
const props = defineProps<{
    type: "family" | "reasoning" | "tool_call" | "attachment" | "vision" | "free" | "open_weights" | "temperature" | "structured_output";
}>();

const badgeMap: Record<string, { label: string; class: string }> = {
    family: { label: "family", class: "bg-primary/10 text-primary" },
    reasoning: { label: "reasoning", class: "bg-amber-500/10 text-amber-600 dark:text-amber-400" },
    tool_call: { label: "tool_call", class: "bg-blue-500/10 text-blue-600 dark:text-blue-400" },
    attachment: { label: "attachment", class: "bg-red-500/10 text-red-500 dark:text-red-400" },
    vision: { label: "vision", class: "bg-green-500/10 text-green-600 dark:text-green-400" },
    free: { label: "FREE", class: "bg-success/10 text-success font-semibold" },
    open_weights: { label: "open weights", class: "bg-green-500/10 text-green-600 dark:text-green-400" },
    temperature: { label: "temperature", class: "hover:bg-accented text-toned" },
    structured_output: { label: "structured output", class: "bg-purple-500/10 text-purple-600 dark:text-purple-400" },
};

const computed_badge = computed(() => badgeMap[props.type] || { label: props.type, class: "hover:bg-accented text-toned" });
const label = computed(() => computed_badge.value.label);
const badgeClass = computed(() => computed_badge.value.class);
</script>
