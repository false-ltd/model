/** Standard animation config for all charts */
export const chartAnimation = { duration: 600, easing: "easeOutQuart" as const };

/** Standard font family for all chart text */
export const chartFontFamily = "Inter, -apple-system, BlinkMacSystemFont, sans-serif";

/** Plugin-drawn text font helper */
export const chartFont = (weight = "bold", size = 10) =>
    `${weight} ${size}px ${chartFontFamily}`;

/** Theme-aware colors from CSS variables */
export const chartColors = () => ({
    text: getCSSVar("--ui-text-muted"),
    textDefault: getCSSVar("--ui-text-default"),
    grid: getCSSVar("--ui-border-default"),
    bg: getCSSVar("--ui-bg-default"),
});

/** Shared Chart.js tooltip defaults */
export const chartTooltip = {
    backgroundColor: "rgba(0,0,0,0.75)",
    titleFont: { size: 12, family: chartFontFamily },
    bodyFont: { size: 11, family: chartFontFamily },
    titleColor: "#fff",
    bodyColor: "#e5e5e5",
    padding: 8,
    cornerRadius: 6,
};

/** Standard tick font config */
export const chartTickFont = (size = 11) => ({ size, family: chartFontFamily });

/** Standard legend label config */
export const chartLegendLabels = (color: string) => ({
    boxWidth: 10,
    boxHeight: 10,
    padding: 12,
    color,
    font: { size: 11, family: chartFontFamily },
});
