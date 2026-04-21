<template>
    <div class="text-xs text-muted uppercase tracking-wider font-medium p-3 px-4 bg-elevated border-b border-default flex items-center gap-1.5">
        <UIcon v-if="icon" :name="icon" class="size-3.5 text-primary" />
        {{ label }}
    </div>
</template>

<script setup lang="ts">
defineProps<{ label: string; icon?: string }>();
</script>
