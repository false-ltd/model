<template>
    <div class="min-h-screen flex flex-col bg-default transition-colors duration-200">
        <AppHeader />
        <main class="flex-1 w-full px-3 sm:px-4 md:px-6 py-2">
            <slot />
        </main>
        <AppFooter />
    </div>
</template>

<script setup lang="ts">
    const { initSync } = useAutoSync();
    onMounted(() => initSync());
</script>
