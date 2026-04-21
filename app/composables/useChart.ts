import { Chart, type ChartConfiguration } from "chart.js";

/**
 * Minimal chart lifecycle helper.
 * Provides canvas ref, render/destroy, and auto-watches a source.
 */
export function useChart(
    getChartConfig: () => ChartConfiguration | null,
    watchSource: () => unknown,
) {
    const canvasRef = ref<HTMLCanvasElement | null>(null);
    let chartInstance: Chart | null = null;

    const render = () => {
        chartInstance?.destroy();
        chartInstance = null;
        if (!canvasRef.value) return;
        const config = getChartConfig();
        if (!config) return;
        chartInstance = new Chart(canvasRef.value, config);
    };

    watch(watchSource, () => nextTick(render), { deep: true });
    onMounted(render);
    onUnmounted(() => {
        chartInstance?.destroy();
        chartInstance = null;
    });

    return { canvasRef };
}
