import type { DirectiveBinding } from "vue";

/**
 * Motion foundation:
 *  - Route transitions via the View Transitions API. Only real path
 *    changes are wrapped — query-only updates (catalog filtering) must not
 *    flash a full-page crossfade. Graceful no-op without browser support
 *    or with reduced motion; the CSS page transition remains as fallback
 *    and is disabled while view transitions are active (html.vt).
 *  - v-magnetic directive — element leans toward the cursor within a
 *    radius and springs back on leave (transform-only, GPU friendly).
 */
export default defineNuxtPlugin((nuxtApp) => {
    const reduced = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
    const supportsVT = typeof (document as any).startViewTransition === "function";

    if (supportsVT && !reduced) {
        // Install after the app has fully mounted: the initial navigation
        // (including i18n locale redirects) must never run inside a view
        // transition — a hydration-time transition can leave the snapshot
        // overlay stuck and freeze the page on the old pixels.
        nuxtApp.hook("app:mounted", () => {
            document.documentElement.classList.add("vt");
            const router = nuxtApp.$router as any;

            const wrap = (original: (...args: any[]) => any) => {
                return (to: any, ...rest: any[]) => {
                    let targetPath: string | undefined;
                    try {
                        targetPath = router.resolve(to)?.path;
                    } catch {
                        targetPath = undefined;
                    }
                    const samePath = targetPath === router.currentRoute.value?.path;
                    if (!targetPath || samePath) return original(to, ...rest);

                    const result = original(to, ...rest);
                    if (typeof result?.then === "function") {
                        const transition = (document as any).startViewTransition(() => result);
                        transition.finished.catch(() => {});
                        // Safety valve: if anything keeps the transition from
                        // settling, snap out of the frozen snapshot overlay.
                        setTimeout(() => transition.skipTransition?.(), 700);
                        return Promise.resolve();
                    }
                    return result;
                };
            };
            router.push = wrap(router.push.bind(router));
            router.replace = wrap(router.replace.bind(router));
        });
    }

    nuxtApp.vueApp.directive("magnetic", {
        mounted(el: HTMLElement, binding: DirectiveBinding<number | undefined>) {
            if (reduced || window.matchMedia("(pointer: coarse)").matches) return;
            if (binding.value) el.dataset.magneticStrength = String(binding.value);

            const onMove = (e: MouseEvent) => {
                const rect = el.getBoundingClientRect();
                const dx = e.clientX - (rect.left + rect.width / 2);
                const dy = e.clientY - (rect.top + rect.height / 2);
                const dist = Math.hypot(dx, dy);
                const max = Math.max(rect.width, rect.height) * 0.9;
                if (dist > max) return;
                const s = Number(el.dataset.magneticStrength ?? 14);
                el.style.transition = "transform 0.18s var(--ease-out-quint)";
                el.style.transform = `translate(${(dx / max) * s}px, ${(dy / max) * s}px)`;
            };
            const onLeave = () => {
                el.style.transition = "transform 0.45s var(--ease-out-quint)";
                el.style.transform = "";
            };

            el.addEventListener("mousemove", onMove);
            el.addEventListener("mouseleave", onLeave);
            (el as any).__magneticCleanup = () => {
                el.removeEventListener("mousemove", onMove);
                el.removeEventListener("mouseleave", onLeave);
            };
        },
        unmounted(el: HTMLElement) {
            (el as any).__magneticCleanup?.();
        },
    });
});
