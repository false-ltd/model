<template>
    <div v-if="items.length" class="bg-default border border-default rounded-xl py-2.5 overflow-hidden">
        <div class="flex items-center gap-3 px-4">
            <span class="shrink-0 inline-flex items-center gap-1.5 text-[10px] font-semibold uppercase tracking-wider text-primary">
                <UIcon name="i-lucide-sparkle" class="size-3" />
                {{ title }}
            </span>
            <div class="ticker-viewport relative flex-1 overflow-hidden">
                <div class="ticker-track flex items-center gap-2 w-max">
                    <template v-for="(dup, d) in [0, 1]" :key="d">
                        <NuxtLink
                            v-for="m in items"
                            :key="`${d}-${m.id}`"
                            :to="localePath(`/model/${m.id}`)"
                            class="shrink-0 inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full border border-default bg-elevated text-xs text-toned no-underline hover:border-primary/50 hover:text-primary transition-colors"
                        >
                            <span class="size-1.5 rounded-full bg-primary shrink-0" />
                            <span class="max-w-40 truncate">{{ m.name }}</span>
                            <span class="text-muted tabular">{{ m.release_date?.slice(5) }}</span>
                        </NuxtLink>
                    </template>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
    /**
     * Recently indexed models marquee — the one homepage element that can
     * differ visit to visit. Loops seamlessly, pauses on hover, renders as
     * a static wrap under reduced motion.
     */
    const props = defineProps<{
        title: string;
        items: { id: number; name: string; release_date: string }[];
    }>();
    const localePath = useLocalePath();
    void props;
</script>

<style scoped>
    .ticker-track {
        animation: ticker-scroll 55s linear infinite;
    }
    .ticker-viewport:hover .ticker-track {
        animation-play-state: paused;
    }
    @keyframes ticker-scroll {
        from {
            transform: translateX(0);
        }
        to {
            transform: translateX(-50%);
        }
    }
    @media (prefers-reduced-motion: reduce) {
        .ticker-track {
            animation: none;
        }
        .ticker-viewport {
            overflow-x: auto;
        }
    }
</style>
