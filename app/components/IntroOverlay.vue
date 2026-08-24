<template>
    <Transition name="intro">
        <div
            v-if="visible"
            class="fixed inset-0 z-[100] bg-default flex flex-col items-center justify-center pointer-events-auto cursor-pointer"
            @click="dismiss"
        >
            <div class="flex items-center gap-3 mb-5">
                <UIcon name="i-lucide-square-m" class="size-8 text-primary" />
                <span class="font-display font-bold text-2xl text-default tracking-tight">AI Model Atlas</span>
            </div>
            <div class="w-48 h-px bg-default border border-default overflow-hidden rounded-full">
                <div class="h-full bg-primary intro-bar" />
            </div>
            <div class="mt-4 text-[10px] text-muted uppercase tracking-[0.25em] tabular">
                {{ t("intro.indexing") }} {{ count.toLocaleString() }}+
            </div>
        </div>
    </Transition>
</template>

<script setup lang="ts">
    /**
     * First-visit-only intro: wordmark, a hairline sweep and a quick count
     * up. One session, once; click to skip; never shown with reduced motion.
     */
    const { t } = useI18n();
    const visible = ref(false);
    const count = ref(0);
    const SESSION_KEY = "model-intro-seen";

    let raf = 0;
    const dismiss = () => {
        visible.value = false;
        try {
            sessionStorage.setItem(SESSION_KEY, "1");
        } catch {}
    };

    onMounted(() => {
        if (window.matchMedia("(prefers-reduced-motion: reduce)").matches) return;
        try {
            if (sessionStorage.getItem(SESSION_KEY)) return;
        } catch {}
        visible.value = true;

        const to = 4300;
        const start = performance.now();
        const duration = 1000;
        const tick = (now: number) => {
            const p = Math.min((now - start) / duration, 1);
            count.value = Math.round((1 - Math.pow(1 - p, 3)) * to);
            if (p < 1) raf = requestAnimationFrame(tick);
        };
        raf = requestAnimationFrame(tick);
        setTimeout(dismiss, 1500);
    });

    onUnmounted(() => cancelAnimationFrame(raf));
</script>

<style scoped>
    .intro-leave-active {
        transition:
            opacity 0.5s ease,
            transform 0.5s var(--ease-out-quint);
    }
    .intro-leave-to {
        opacity: 0;
        transform: scale(1.04);
    }
    .intro-bar {
        width: 100%;
        transform-origin: left;
        animation: intro-sweep 1.2s var(--ease-out-quint) forwards;
    }
    @keyframes intro-sweep {
        from {
            transform: scaleX(0);
        }
        to {
            transform: scaleX(1);
        }
    }
</style>
