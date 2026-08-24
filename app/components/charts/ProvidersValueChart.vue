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
        BarController,
        BarElement,
        CategoryScale,
        LinearScale,
        Tooltip,
        type ChartOptions,
    } from "chart.js";

    Chart.register(BarController, BarElement, CategoryScale, LinearScale, Tooltip);

    /**
     * Provider value view: bar length = model count, inner band = open
     * weights share, tooltip adds median input price. Stacked segments
     * render the band without extra drawing code.
     */
    const props = defineProps<{
        title: string;
        subtitle: string;
        items: { name: string; open: number; closed: number; median: number | null }[];
        medianLabel: string;
    }>();

    const { canvasRef } = useChart(
        () => {
            if (!props.items.length) return null;
            const c = chartColors();
            const primary = () => getCSSVar("--ui-primary");

            return {
                type: "bar",
                data: {
                    labels: props.items.map((i) =>
                        i.name.length > 16 ? i.name.slice(0, 14) + "…" : i.name,
                    ),
                    datasets: [
                        {
                            label: "open",
                            data: props.items.map((i) => i.open),
                            backgroundColor: chartColor(0) + "cc",
                            borderRadius: 4,
                            borderSkipped: false,
                        },
                        {
                            label: "closed",
                            data: props.items.map((i) => i.closed),
                            backgroundColor: (getCSSVar("--ui-border-accented") || "#ccc") + "aa",
                            borderRadius: 4,
                            borderSkipped: false,
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
                            stacked: true,
                            display: false,
                            max: Math.max(...props.items.map((i) => i.open + i.closed)) * 1.05,
                        },
                        y: {
                            stacked: true,
                            grid: { display: false },
                            ticks: { color: c.text, font: chartTickFont() },
                            border: { display: false },
                        },
                    },
                    plugins: {
                        tooltip: {
                            ...chartTooltip,
                            callbacks: {
                                label: (ctx: any) => {
                                    const it = props.items[ctx.dataIndex];
                                    if (!it) return "";
                                    const median =
                                        it.median != null ? ` · ${props.medianLabel} $${it.median < 10 ? it.median.toFixed(2) : it.median.toFixed(0)}` : "";
                                    return ` ${it.open + it.closed} models · ${(it.open / (it.open + it.closed) * 100).toFixed(0)}% open${median}`;
                                },
                            },
                        },
                        legend: { display: false },
                    },
                } satisfies ChartOptions<"bar">,
                plugins: [
                    {
                        id: "countLabels",
                        afterDatasetsDraw(chart: Chart) {
                            const { ctx } = chart;
                            const meta = chart.getDatasetMeta(1);
                            for (let i = 0; i < meta.data.length; i++) {
                                const bar = meta.data[i]!;
                                const it = props.items[i];
                                if (!it) continue;
                                ctx.save();
                                ctx.font = chartFont();
                                ctx.fillStyle = c.text;
                                ctx.textAlign = "left";
                                ctx.textBaseline = "middle";
                                ctx.fillText(String(it.open + it.closed), bar.x + bar.width + 6, bar.y);
                                ctx.restore();
                            }
                        },
                    },
                ],
            } as any;
        },
        () => props.items,
    );
</script>
