<template>
    <div :style="{ height: (height || 160) + 'px', position: 'relative' }">
        <canvas ref="canvasRef" />
    </div>
</template>

<script setup lang="ts">
    import { Chart, BarController, CategoryScale, LinearScale, BarElement, Tooltip, Legend } from "chart.js";

    Chart.register(BarController, CategoryScale, LinearScale, BarElement, Tooltip, Legend);

    const props = defineProps<{
        type: "bar" | "horizontalBar";
        labels: string[];
        datasets: { label: string; data: number[]; backgroundColor: string | string[]; borderColor?: string | string[]; borderRadius?: number }[];
        height?: number;
    }>();

    const { canvasRef } = useChart(
        () => {
            if (!props.labels.length || !props.datasets.length) return null;

            const c = chartColors();

            return {
                type: "bar",
                data: {
                    labels: props.labels,
                    datasets: props.datasets.map((ds) => ({ ...ds, borderWidth: ds.borderColor ? 1 : 0 })),
                },
                options: {
                    indexAxis: (props.type === "horizontalBar" ? "y" : "x") as "x" | "y",
                    responsive: true,
                    maintainAspectRatio: false,
                    animation: chartAnimation,
                    plugins: {
                        legend: {
                            display: props.datasets.length > 1,
                            position: "bottom" as const,
                            labels: chartLegendLabels(c.text),
                        },
                        tooltip: {
                            ...chartTooltip,
                            titleFont: { size: 12, weight: "bold" as const, family: chartFontFamily },
                            padding: 10,
                            cornerRadius: 8,
                            displayColors: true,
                            boxPadding: 4,
                        },
                    },
                    scales: {
                        x: { grid: { display: false }, ticks: { font: chartTickFont(), color: c.text } },
                        y: { grid: { color: c.grid }, ticks: { font: chartTickFont(), color: c.text } },
                    },
                },
            } as any;
        },
        () => [props.labels, props.datasets],
    );
</script>
