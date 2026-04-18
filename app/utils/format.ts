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
    image: "bg-green-500/10 text-green-600 dark:text-green-400",
    audio: "bg-purple-500/10 text-purple-600 dark:text-purple-400",
    video: "bg-orange-500/10 text-orange-600 dark:text-orange-400",
    pdf: "bg-red-500/10 text-red-500 dark:text-red-400",
};

export const modalityInputClass = (mod: string) =>
    modalityInputClasses[mod] || "bg-accented text-toned";

export const modalityOutputClass = () => "bg-blue-500/10 text-blue-500 dark:text-blue-400";

export const chartColors = [
    "#e11d48", // rose
    "#f97316", // orange
    "#0d9488", // teal
    "#8b5cf6", // violet
    "#3b82f6", // blue
    "#eab308", // yellow
];

export const chartColor = (i: number): string => chartColors[i % chartColors.length]!;

export const formatMonth = (m: string): string => {
    const [y, mo] = m.split("-");
    if (!mo) return y ?? m;
    const months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
    return (months[Number(mo) - 1] || mo) + " " + y!.slice(2);
};
