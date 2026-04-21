<template>
    <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium" :class="badgeClass">
        {{ label }}
    </span>
</template>

<script setup lang="ts">
import { getBadgeInfo } from "~/utils/badge";

const props = defineProps<{
    type: "family" | "reasoning" | "tool_call" | "attachment" | "vision" | "free" | "open_weights" | "temperature" | "structured_output" | "audio" | "alpha" | "beta" | "deprecated";
}>();

const info = computed(() => getBadgeInfo(props.type));
const label = computed(() => info.value.label);
const badgeClass = computed(() => info.value.cls);
</script>
