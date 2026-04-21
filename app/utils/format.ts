import type { Model } from "~/types";

export const formatTokens = (n: number | null): string => {
    if (n == null || !n) return "—";
    if (n >= 1_000_000) return (n / 1_000_000).toFixed(0) + "M";
    if (n >= 1_000) return (n / 1_000).toFixed(0) + "K";
    return String(n);
};

const modalityIconMap: Record<string, string> = {
    text: "i-lucide-type",
    image: "i-lucide-image",
    audio: "i-lucide-mic",
    video: "i-lucide-video",
    pdf: "i-lucide-file-text",
};

export const modalityIcon = (mod: string) => modalityIconMap[mod] || "i-lucide-circle-dot";

const modalityClassMap: Record<string, string> = {
    text: "bg-sky-500/10 border-sky-500/30 text-sky-600 dark:text-sky-400",
    image: "bg-green-500/10 border-green-500/30 text-green-600 dark:text-green-400",
    audio: "bg-sky-500/10 border-sky-500/30 text-sky-600 dark:text-sky-400",
    video: "bg-rose-500/10 border-rose-500/30 text-rose-600 dark:text-rose-400",
    pdf: "bg-emerald-500/10 border-emerald-500/30 text-emerald-600 dark:text-emerald-400",
};

export const modalityClass = (mod: string) =>
    modalityClassMap[mod] || "bg-accented border-default text-toned";

/** Read a CSS custom property value at runtime */
export const getCSSVar = (name: string, fallback = ""): string => {
    if (!import.meta.client) return fallback;
    return getComputedStyle(document.documentElement).getPropertyValue(name).trim() || fallback;
};

/** Get a chart palette color by index (0–5) */
export const chartColor = (i: number): string => getCSSVar(`--chart-${i % 6}`);

export const successColor = () => getCSSVar("--color-success");
export const dangerColor = () => getCSSVar("--color-danger");

export const isVisionModel = (m: Model): boolean => m?.modalities_input?.includes("image") ?? false;

export const isFreeModel = (m: Model): boolean => (m?.cost_input ?? 0) === 0 && (m?.cost_output ?? 0) === 0;
