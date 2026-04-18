<template>
    <div class="text-xs text-muted uppercase tracking-wider font-medium p-3 px-4 bg-elevated border-b border-default">{{ label }}</div>
</template>

<script setup lang="ts">
defineProps<{ label: string }>();
</script>
