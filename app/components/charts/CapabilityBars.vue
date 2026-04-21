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
    type Chart as ChartType,
    BarElement,
    BarController,
    CategoryScale,
    LinearScale,
    Tooltip,
    type ChartOptions,
} from "chart.js";

Chart.register(BarElement, BarController, CategoryScale, LinearScale, Tooltip);

const props = defineProps<{
    title: string;
    subtitle: string;
    capabilities: { label: string; count: number; total: number; pct: number; color: string }[];
}>();

const { canvasRef } = useChart(
    () => {
        if (!props.capabilities.length) return null;

        const c = chartColors();

        return {
            type: "bar",
            data: {
                labels: props.capabilities.map((cap) => cap.label),
                datasets: [
                    {
                        label: "count",
                        data: props.capabilities.map((cap) => cap.pct),
                        backgroundColor: props.capabilities.map((cap) => cap.color + "cc"),
                        hoverBackgroundColor: props.capabilities.map((cap) => cap.color),
                        borderRadius: 4,
                        borderSkipped: false,
                        barThickness: 24,
                    },
                ],
            },
            options: {
                indexAxis: "y",
                responsive: true,
                maintainAspectRatio: false,
                animation: chartAnimation,
                scales: {
                    x: {
                        display: false,
                        max: 100,
                    },
                    y: {
                        grid: { display: false },
                        ticks: { color: c.text, font: chartTickFont() },
                        border: { display: false },
                    },
                },
                plugins: {
                    tooltip: {
                        ...chartTooltip,
                        callbacks: {
                            label: (ctx) => {
                                const cap = props.capabilities[ctx.dataIndex];
                                return cap ? ` ${cap.count} / ${cap.total} (${cap.pct}%)` : "";
                            },
                        },
                    },
                    legend: { display: false },
                },
            } satisfies ChartOptions<"bar">,
            plugins: [
                {
                    id: "pctLabels",
                    afterDatasetsDraw(chart: ChartType) {
                        const { ctx } = chart;
                        const meta = chart.getDatasetMeta(0);
                        for (let i = 0; i < meta.data.length; i++) {
                            const bar = meta.data[i]!;
                            const cap = props.capabilities[i];
                            if (!cap) continue;
                            ctx.save();
                            ctx.font = chartFont();
                            ctx.fillStyle = cap.pct > 85 ? "#fff" : cap.color;
                            ctx.textAlign = cap.pct > 85 ? "right" : "left";
                            ctx.textBaseline = "middle";
                            const x = cap.pct > 85 ? bar.x - 6 : bar.x + 6;
                            ctx.fillText(`${cap.pct}% (${cap.count})`, x, bar.y);
                            ctx.restore();
                        }
                    },
                },
            ],
        } as any;
    },
    () => props.capabilities,
);
</script>
