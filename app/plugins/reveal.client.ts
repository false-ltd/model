import type { DirectiveBinding } from "vue";

/**
 * v-reveal — fade/translate an element in when it scrolls into view.
 * Optional modifier value sets the stagger delay in ms:
 *   <div v-reveal> / <div v-reveal="120">
 */
export default defineNuxtPlugin((nuxtApp) => {
    const reduced =
        typeof window !== "undefined" &&
        window.matchMedia("(prefers-reduced-motion: reduce)").matches;

    let observer: IntersectionObserver | null = null;
    if (!reduced) {
        observer = new IntersectionObserver(
            (entries) => {
                for (const entry of entries) {
                    if (entry.isIntersecting) {
                        entry.target.classList.add("reveal-in");
                        observer!.unobserve(entry.target);
                    }
                }
            },
            { threshold: 0.06, rootMargin: "0px 0px -32px 0px" },
        );
    }

    nuxtApp.vueApp.directive("reveal", {
        mounted(el: HTMLElement, binding: DirectiveBinding<number | undefined>) {
            if (!observer) {
                el.classList.add("reveal-in");
                return;
            }
            el.classList.add("reveal");
            if (binding.value) {
                el.style.setProperty("--reveal-delay", `${binding.value}ms`);
            }
            observer.observe(el);
        },
        unmounted(el: HTMLElement) {
            observer?.unobserve(el);
        },
    });
});
