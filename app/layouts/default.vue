<template>
    <div class="min-h-screen flex flex-col bg-default transition-colors duration-200">
        <AppHeader />
        <main class="relative flex-1 w-full px-3 sm:px-4 md:px-6 py-2">
            <slot />
        </main>
        <AppFooter />
        <AppCommandPalette v-model:open="paletteOpen" />
        <CompareFab />
    </div>
</template>

<script setup lang="ts">
    // Shared open state so the header button and the ⌘K shortcut inside
    // AppCommandPalette control the same palette.
    const paletteOpen = useState<boolean>("command-palette-open", () => false);

    const { initSync } = useAutoSync();
    onMounted(() => initSync());
</script>
