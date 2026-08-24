// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    compatibilityDate: "2025-07-15",
    devtools: { enabled: false },
    modules: ["@nuxt/ui", "@nuxtjs/i18n"],
    ssr: false,
    colorMode: {
        preference: "system",
        fallback: "light",
        classSuffix: "",
    },
    css: ["~/assets/css/main.css"],
    app: {
        buildAssetsDir: "assets",
        pageTransition: { name: "page", mode: "out-in" },
        head: {
            link: [{ rel: "icon", type: "image/svg+xml", href: "/favicon.svg" }],
            meta: [
                { name: "applicable-device", content: "pc,mobile" },
                { name: "description", content: "Browse, compare, and analyze AI/LLM model pricing, capabilities, and performance across all major providers." },
                { name: "baidu-site-verification", content: "CODE" },
                { name: "sogou_site_verification", content: "CODE" },
                { name: "360-site-verification", content: "CODE" },
            ],
        },
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
    nitro: {
        compressPublicAssets: true,
    },
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
        // Fonts are bundled locally via @fontsource-variable/* — remote
        // providers are disabled so builds never depend on network access.
        providers: {
            google: false,
            googleicons: false,
            bunny: false,
            fontshare: false,
        },
        families: [
            { name: "Inter Variable", provider: "none" },
            { name: "Space Grotesk Variable", provider: "none" },
        ],
    },
});
