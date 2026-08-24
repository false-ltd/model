/**
 * Animated numeric count-up. Starts when the target becomes non-zero —
 * async data arrives after mount, so the target must be watched, not read
 * once. Snaps instantly under reduced motion; re-animates smoothly if the
 * target changes again (e.g. after a data sync).
 */
export function useCountUp(target: Ref<number>, duration = 1400) {
    const display = ref(0);
    let raf = 0;
    let lastTarget: number | null = null;

    const start = (to: number) => {
        lastTarget = to;
        const reduced = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
        if (reduced || to <= 0) {
            display.value = to;
            return;
        }
        cancelAnimationFrame(raf);
        const from = display.value;
        const t0 = performance.now();
        const easeOutExpo = (t: number) => (t === 1 ? 1 : 1 - Math.pow(2, -10 * t));
        const tick = (now: number) => {
            const progress = Math.min((now - t0) / duration, 1);
            display.value = Math.round(from + (to - from) * easeOutExpo(progress));
            if (progress < 1) raf = requestAnimationFrame(tick);
        };
        raf = requestAnimationFrame(tick);
    };

    onMounted(() => {
        watch(target, (to) => {
            if (to !== lastTarget) start(to);
        }, { immediate: true });
    });

    onUnmounted(() => cancelAnimationFrame(raf));

    return display;
}
