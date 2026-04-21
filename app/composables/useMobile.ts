export const useMobile = (breakpoint = 768) => {
    const isMobile = ref(false);

    if (import.meta.client) {
        const mq = window.matchMedia(`(max-width: ${breakpoint - 1}px)`);
        isMobile.value = mq.matches;
        const handler = (e: MediaQueryListEvent) => { isMobile.value = e.matches; };
        mq.addEventListener("change", handler);
        onUnmounted(() => mq.removeEventListener("change", handler));
    }

    return { isMobile };
};
