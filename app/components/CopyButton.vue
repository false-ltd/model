<template>
    <button @click="copy" class="cursor-pointer opacity-30 hover:opacity-100 transition-opacity">
        <UIcon :name="copied ? 'i-lucide-copy-check' : 'i-lucide-copy'" class="size-4 text-muted" />
    </button>
</template>

<script setup lang="ts">
    const props = defineProps<{ value: string }>();
    const copied = ref(false);

    const copy = async () => {
        try {
            await navigator.clipboard.writeText(props.value);
        } catch {}
        copied.value = true;
        setTimeout(() => {
            copied.value = false;
        }, 1500);
    };
</script>
