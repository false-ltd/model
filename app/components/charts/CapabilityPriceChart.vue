<template>
    <div class="bg-default border border-default rounded-xl p-5 hover:border-accented transition-colors h-full">
        <div class="text-sm font-semibold text-default mb-1">{{ title }}</div>
        <div class="text-xs text-muted mb-3">{{ subtitle }}</div>
        <div style="height: 220px; position: relative">
            <canvas ref="canvasRef" />
        </div>
    </div>
</template>

<script setup lang="ts">
    import {
        Chart,
        LineController,
        LineElement,
        PointElement,
        CategoryScale,
        LinearScale,
        Tooltip,
        type ChartOptions,
    } from "chart.js";

    Chart.register(LineController, LineElement, PointElement, CategoryScale, LinearScale, Tooltip);

    /**
     * Median input cost per capability count — "what does capability cost".
     */
    const props = defineProps<{
        title: string;
        subtitle: string;
        points: { k: number; median: number }[];
        xLabel: string;
    }>();

    const { canvasRef } = useChart(
        () => {
            if (!props.points.length) return null;
            const c = chartColors();

            return {
                type: "line",
                data: {
                    labels: props.points.map((p) => `${p.k}`),
                    datasets: [
                        {
                            data: props.points.map((p) => p.median),
                            borderColor: chartColor(0),
                            backgroundColor: chartColor(0),
                            tension: 0.35,
                            fill: false,
                            pointRadius: 4,
                            pointHoverRadius: 6,
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
                            title: { display: true, text: props.xLabel, color: c.text, font: chartTickFont(10) },
                            border: { display: false },
                        },
                        y: {
                            beginAtZero: true,
                            grid: { color: c.grid },
                            ticks: {
                                color: c.text,
                                font: chartTickFont(10),
                                callback: (v: any) => `$${v}`,
                            },
                            border: { display: false },
                        },
                    },
                    plugins: {
                        tooltip: {
                            ...chartTooltip,
                            callbacks: {
                                label: (ctx: any) => ` $${ctx.parsed.y.toFixed(2)} / 1M tokens`,
                            },
                        },
                        legend: { display: false },
                    },
                } satisfies ChartOptions<"line">,
                plugins: [
                    {
                        id: "valueLabels",
                        afterDatasetsDraw(chart: Chart) {
                            const { ctx } = chart;
                            const meta = chart.getDatasetMeta(0);
                            for (let i = 0; i < meta.data.length; i++) {
                                const pt = meta.data[i]!;
                                const val = props.points[i]?.median;
                                if (val == null) continue;
                                ctx.save();
                                ctx.font = chartFont("600", 10);
                                ctx.fillStyle = chartColor(0);
                                ctx.textAlign = "center";
                                ctx.textBaseline = "bottom";
                                ctx.fillText(`$${val < 10 ? val.toFixed(1) : val.toFixed(0)}`, pt.x, pt.y - 7);
                                ctx.restore();
                            }
                        },
                    },
                ],
            } as any;
        },
        () => props.points,
    );
</script>
