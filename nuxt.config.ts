// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    compatibilityDate: "2025-07-15",
    devtools: { enabled: false },
    modules: ["@nuxt/ui", "@nuxtjs/i18n", "@nuxtjs/seo"],
    ssr: false,
    colorMode: {
        preference: "system",
        fallback: "light",
        classSuffix: "",
    },
    css: ["~/assets/css/main.css"],
    app: {
        buildAssetsDir: "assets",
        head: {
            link: [{ rel: "icon", type: "image/svg+xml", href: "/favicon.svg" }],
        },
    },
    site: {
        url: process.env.NUXT_PUBLIC_SITE_URL || "https://model.false.ltd",
        name: "AI Model Catalog",
        description:
            "Browse, compare, and analyze AI/LLM model pricing, capabilities, and performance across all major providers.",
        defaultLocale: "en",
    },
    sitemap: {
        zeroRuntime: true,
    },
    linkChecker: {
        failOnError: false,
    },
    i18n: {
        locales: [
            { code: "en", name: "English", file: "en.json" },
            { code: "zh", name: "中文", file: "zh.json" },
        ],
        defaultLocale: "en",
        langDir: "locales",
        strategy: "prefix_except_default",
        detectBrowserLanguage: {
            useCookie: true,
            cookieKey: "i18n_redirected",
            redirectOn: "all",
        },
    },
    runtimeConfig: {
        public: {
            apiBase: process.env.NUXT_PUBLIC_API_BASE || "",
        },
    },
    nitro: {},
    vite: {
        optimizeDeps: {
            include: ["chart.js"],
        },
    },
    ui: {
        theme: {
            defaultVariants: {
                color: "neutral",
                size: "md",
            },
        },
    },
    fonts: {
        experimental: {
            processCSSVariables: false,
        },
        defaults: {
            weights: [400, 500, 600, 700],
        },
        families: [{ name: "Inter", provider: "none" }],
    },
});
