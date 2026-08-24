import Lenis from "lenis";

/**
 * Inertial smooth scrolling (desktop pointer only; native scrolling on
 * touch keeps platform gestures intact). Disabled entirely for users who
 * prefer reduced motion.
 */
export default defineNuxtPlugin(() => {
    if (window.matchMedia("(prefers-reduced-motion: reduce)").matches) return;
    if (window.matchMedia("(pointer: coarse)").matches) return;

    const lenis = new Lenis({
        lerp: 0.12,
        smoothWheel: true,
        // Native scrolling inside tagged containers (popovers, tables,
        // palettes) — belt-and-braces alongside the built-in
        // data-lenis-prevent attribute check.
        prevent: (node) => node instanceof HTMLElement && !!node.closest("[data-lenis-prevent]"),
    });

    let raf = 0;
    const loop = (time: number) => {
        lenis.raf(time);
        raf = requestAnimationFrame(loop);
    };
    raf = requestAnimationFrame(loop);

    // Route changes should always land at the top of the new page.
    const router = useRouter();
    router.afterEach(() => lenis.scrollTo(0, { immediate: true }));
});
