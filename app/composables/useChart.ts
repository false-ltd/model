import { Chart, type ChartConfiguration } from "chart.js";

/**
 * Minimal chart lifecycle helper.
 * Provides canvas ref, render/destroy, and auto-watches a source.
 * Re-renders on color-mode changes (colors are read from CSS variables
 * at render time) and skips the entry animation on re-renders.
 */
export function useChart(
    getChartConfig: () => ChartConfiguration | null,
    watchSource: () => unknown,
) {
    const canvasRef = ref<HTMLCanvasElement | null>(null);
    let chartInstance: Chart | null = null;

    const render = (animate = false) => {
        chartInstance?.destroy();
        chartInstance = null;
        if (!canvasRef.value) return;
        const config = getChartConfig();
        if (!config) return;
        if (!animate) config.options = { ...config.options, animation: false };
        chartInstance = new Chart(canvasRef.value, config);
    };

    watch(watchSource, () => nextTick(() => render(false)), { deep: true });

    const colorMode = useColorMode();
    watch(
        () => colorMode.value,
        () => nextTick(() => render(false)),
    );

    onMounted(() => render(true));
    onUnmounted(() => {
        chartInstance?.destroy();
        chartInstance = null;
    });

    return { canvasRef };
}
