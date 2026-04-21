export const formatTokens = (n: number | null): string => {
    if (n == null || !n) return "—";
    if (n >= 1_000_000) return (n / 1_000_000).toFixed(0) + "M";
    if (n >= 1_000) return (n / 1_000).toFixed(0) + "K";
    return String(n);
};

export const formatScore = (n: number): string => {
    if (n >= 1_000_000) return (n / 1_000_000).toFixed(1) + "M";
    if (n >= 1_000) return (n / 1_000).toFixed(0) + "K";
    return String(n);
};

export const modalityIconMap: Record<string, string> = {
    text: "i-lucide-type",
    image: "i-lucide-image",
    audio: "i-lucide-mic",
    video: "i-lucide-video",
    pdf: "i-lucide-file-text",
};

export const modalityIcon = (mod: string) => modalityIconMap[mod] || "i-lucide-circle-dot";

const modalityInputClasses: Record<string, string> = {
    image: "bg-emerald-500/10 text-emerald-600 dark:text-emerald-400",
    audio: "bg-violet-500/10 text-violet-600 dark:text-violet-400",
    video: "bg-amber-500/10 text-amber-600 dark:text-amber-400",
    pdf: "bg-rose-500/10 text-rose-500 dark:text-rose-400",
};

export const modalityInputClass = (mod: string) =>
    modalityInputClasses[mod] || "bg-accented text-toned";

export const modalityOutputClass = () => "bg-sky-500/10 text-sky-500 dark:text-sky-400";

/** Read a CSS custom property value at runtime */
export const getCSSVar = (name: string, fallback = ""): string => {
    if (!import.meta.client) return fallback;
    return getComputedStyle(document.documentElement).getPropertyValue(name).trim() || fallback;
};

/** Get a chart palette color by index (0–5) */
export const chartColor = (i: number): string => getCSSVar(`--chart-${i % 6}`);

/** Semantic color helpers */
export const successColor = () => getCSSVar("--color-success");
export const dangerColor = () => getCSSVar("--color-danger");
export const mutedColor = () => getCSSVar("--ui-text-muted");

export const formatMonth = (m: string): string => {
    const [y, mo] = m.split("-");
    if (!mo) return y ?? m;
    const months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
    return (months[Number(mo) - 1] || mo) + " " + y!.slice(2);
};
