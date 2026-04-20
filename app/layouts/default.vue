<template>
    <div class="min-h-screen flex flex-col bg-default transition-colors duration-200">
        <AppHeader />
        <main class="flex-1 w-full px-1 md:px-1 lg:px-4 py-2">
            <slot />
        </main>
        <AppFooter />
    </div>
</template>

<script setup lang="ts">
    const { initAutoSync, stopAutoSync } = useAutoSync();
    onMounted(() => initAutoSync());
    onUnmounted(() => stopAutoSync());
</script>
