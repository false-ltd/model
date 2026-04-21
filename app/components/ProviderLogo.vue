<template>
    <img
        :src="`https://models.dev/logos/${providerId}.svg`"
        :alt="providerId"
        :class="cls"
        @error="($event.target as HTMLImageElement).style.display = 'none'"
    />
</template>

<script setup lang="ts">
    const props = withDefaults(defineProps<{
        providerId: string;
        cls?: string;
    }>(), {
        cls: "w-5 h-5 rounded shrink-0",
    });
</script>
