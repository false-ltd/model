type BadgeType =
    | "family" | "reasoning" | "tool_call" | "attachment" | "vision"
    | "free" | "open_weights" | "temperature" | "structured_output"
    | "audio" | "alpha" | "beta" | "deprecated";

const badgeMap: Record<BadgeType, { label: string; cls: string }> = {
    family: { label: "family", cls: "bg-primary/10 text-primary" },
    reasoning: { label: "reasoning", cls: "bg-amber-500/10 text-amber-600 dark:text-amber-400" },
    tool_call: { label: "tool_call", cls: "bg-blue-500/10 text-blue-600 dark:text-blue-400" },
    attachment: { label: "attachment", cls: "bg-rose-500/10 text-rose-500 dark:text-rose-400" },
    vision: { label: "vision", cls: "bg-green-500/10 text-green-600 dark:text-green-400" },
    audio: { label: "audio", cls: "bg-sky-500/10 text-sky-600 dark:text-sky-400" },
    free: { label: "FREE", cls: "bg-emerald-500/10 text-emerald-500" },
    open_weights: { label: "open weights", cls: "bg-teal-500/10 text-teal-500" },
    temperature: { label: "temperature", cls: "hover:bg-accented text-toned" },
    structured_output: { label: "structured output", cls: "bg-purple-500/10 text-purple-600 dark:text-purple-400" },
    alpha: { label: "alpha", cls: "bg-orange-500/10 text-orange-600 dark:text-orange-400" },
    beta: { label: "beta", cls: "bg-violet-500/10 text-violet-600 dark:text-violet-400" },
    deprecated: { label: "deprecated", cls: "bg-rose-500/10 text-rose-600 dark:text-rose-400" },
};

export const getBadgeInfo = (type: BadgeType) =>
    badgeMap[type] ?? { label: type, cls: "hover:bg-accented text-toned" };

/** Inline badge class for compact table cells (e.g. "bg-amber-500/10 text-amber-500") */
export const inlineBadgeClass = (type: BadgeType): string => {
    const info = getBadgeInfo(type);
    return `inline-flex items-center rounded px-1 py-px text-[10px] font-medium ${info.cls}`;
};
