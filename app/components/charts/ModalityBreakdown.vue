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

const colorMap: Record<string, string> = {
    text: chartColor(0),
    image: successColor(),
    audio: chartColor(5),
    video: chartColor(4),
    pdf: chartColor(3),
};

const { canvasRef } = useChart(
    () => {
        const allTypes = [...new Set([...props.input.map((m) => m.type), ...props.output.map((m) => m.type)])];
        const inputMap = Object.fromEntries(props.input.map((m) => [m.type, m.count]));
        const outputMap = Object.fromEntries(props.output.map((m) => [m.type, m.count]));

        const c = chartColors();

        return {
            type: "bar",
            data: {
                labels: allTypes.map((t) => t.charAt(0).toUpperCase() + t.slice(1)),
                datasets: [
                    {
                        label: "Input",
                        data: allTypes.map((t) => inputMap[t] || 0),
                        backgroundColor: allTypes.map((t) => (colorMap[t] || c.text) + "cc"),
                        hoverBackgroundColor: allTypes.map((t) => colorMap[t] || c.text),
                        borderRadius: 4,
                        borderSkipped: false,
                    },
                    {
                        label: "Output",
                        data: allTypes.map((t) => outputMap[t] || 0),
                        backgroundColor: allTypes.map((t) => (colorMap[t] || c.text) + "55"),
                        hoverBackgroundColor: allTypes.map((t) => (colorMap[t] || c.text) + "88"),
                        borderRadius: 4,
                        borderSkipped: false,
                    },
                ],
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                animation: chartAnimation,
                scales: {
                    x: {
                        grid: { display: false },
                        ticks: { color: c.text, font: chartTickFont() },
                        border: { display: false },
                    },
                    y: {
                        beginAtZero: true,
                        grid: { color: c.grid },
                        ticks: { color: c.text, font: chartTickFont(10) },
                        border: { display: false },
                    },
                },
                plugins: {
                    tooltip: {
                        ...chartTooltip,
                    },
                    legend: {
                        position: "bottom",
                        labels: chartLegendLabels(c.text),
                    },
                },
            } satisfies ChartOptions<"bar">,
        } as any;
    },
    () => [props.input, props.output],
);
</script>
