<template>
    <div aria-hidden="true" class="absolute inset-0 overflow-hidden pointer-events-none -z-10">
        <div class="absolute inset-0 ambient-grid opacity-60 dark:opacity-40" />
        <div
            class="ambient-blob ambient-blob-a -top-32 -left-24 size-[480px] bg-primary/15 dark:bg-primary/10"
        />
        <div
            class="ambient-blob ambient-blob-b top-10 right-[-140px] size-[420px] bg-chart-blob opacity-[0.12] dark:opacity-[0.14]"
        />
        <!-- The atlas starfield drifting slowly behind the hero -->
        <div class="absolute inset-0 opacity-45 dark:opacity-55 atlas-mask">
            <AtlasCanvas ambient />
        </div>
    </div>
</template>

<style scoped>
.bg-chart-blob {
    background-color: var(--chart-2);
}
.atlas-mask {
    -webkit-mask-image: linear-gradient(to bottom, black 40%, transparent 95%);
    mask-image: linear-gradient(to bottom, black 40%, transparent 95%);
}
</style>
