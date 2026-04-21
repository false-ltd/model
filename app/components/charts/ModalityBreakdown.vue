<template>
    <div class="bg-default border border-default rounded-xl p-5 hover:border-accented transition-colors">
        <div class="text-sm font-semibold text-default mb-1">{{ title }}</div>
        <div class="text-xs text-muted mb-3">{{ subtitle }}</div>
        <div style="height: 200px; position: relative">
            <canvas ref="canvasRef" />
        </div>
    </div>
</template>

<script setup lang="ts">
import {
    Chart,
    BarElement,
    BarController,
    CategoryScale,
    LinearScale,
    Tooltip,
    Legend,
    type ChartOptions,
} from "chart.js";

Chart.register(BarElement, BarController, CategoryScale, LinearScale, Tooltip, Legend);

const props = defineProps<{
    title: string;
    subtitle: string;
    input: { type: string; count: number }[];
    output: { type: string; count: number }[];
    totalModels: number;
}>();

const canvasRef = ref<HTMLCanvasElement | null>(null);
let chart: Chart | null = null;

const colorMap: Record<string, string> = {
    text: chartColor(0),
    image: successColor(),
    audio: chartColor(5),
    video: chartColor(4),
    pdf: chartColor(3),
};

const renderChart = () => {
    chart?.destroy();
    if (!canvasRef.value) return;

    const allTypes = [...new Set([...props.input.map((m) => m.type), ...props.output.map((m) => m.type)])];
    const inputMap = Object.fromEntries(props.input.map((m) => [m.type, m.count]));
    const outputMap = Object.fromEntries(props.output.map((m) => [m.type, m.count]));

    const mutedColor = getCSSVar("--ui-text-muted");

    chart = new Chart(canvasRef.value, {
        type: "bar",
        data: {
            labels: allTypes.map((t) => t.charAt(0).toUpperCase() + t.slice(1)),
            datasets: [
                {
                    label: "Input",
                    data: allTypes.map((t) => inputMap[t] || 0),
                    backgroundColor: allTypes.map((t) => (colorMap[t] || getCSSVar("--ui-text-muted")) + "cc"),
                    borderRadius: 4,
                    borderSkipped: false,
                },
                {
                    label: "Output",
                    data: allTypes.map((t) => outputMap[t] || 0),
                    backgroundColor: allTypes.map((t) => (colorMap[t] || getCSSVar("--ui-text-muted")) + "55"),
                    borderRadius: 4,
                    borderSkipped: false,
                },
            ],
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            scales: {
                x: {
                    grid: { display: false },
                    ticks: { color: mutedColor, font: { size: 11 } },
                    border: { display: false },
                },
                y: {
                    beginAtZero: true,
                    grid: { color: mutedColor + "20" },
                    ticks: { color: mutedColor, font: { size: 10 } },
                    border: { display: false },
                },
            },
            plugins: {
                tooltip: {
                    backgroundColor: "rgba(0,0,0,0.8)",
                    titleFont: { size: 12 },
                    bodyFont: { size: 11 },
                    padding: 8,
                    cornerRadius: 6,
                },
                legend: {
                    position: "bottom",
                    labels: { boxWidth: 10, boxHeight: 10, padding: 12, color: mutedColor, font: { size: 11 } },
                },
            },
        } satisfies ChartOptions<"bar">,
    });
};

watch([() => props.input, () => props.output], () => nextTick(renderChart), { deep: true });
onMounted(renderChart);
onUnmounted(() => chart?.destroy());
</script>
