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
     * Efficiency frontier: cheapest paid model per context tier, with the
     * winning model annotated at each step. Tiers that also have a free
     * option get a green marker.
     */
    const props = defineProps<{
        title: string;
        subtitle: string;
        items: { tier: string; price: number | null; model: string | null; free: boolean }[];
        freeLabel: string;
    }>();

    const { canvasRef } = useChart(
        () => {
            if (!props.items.length) return null;
            const c = chartColors();
            const freeColor = () => getCSSVar("--color-success");

            return {
                type: "line",
                data: {
                    labels: props.items.map((i) => i.tier),
                    datasets: [
                        {
                            data: props.items.map((i) => i.price),
                            borderColor: chartColor(4),
                            backgroundColor: chartColor(4),
                            tension: 0.15,
                            stepped: true,
                            fill: false,
                            pointRadius: props.items.map((i) => (i.free ? 5 : 4)),
                            pointBackgroundColor: props.items.map((i) =>
                                i.free ? freeColor() : chartColor(4),
                            ),
                            spanGaps: true,
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
                            ticks: { color: c.text, font: chartTickFont(10), callback: (v: any) => `$${v}` },
                            border: { display: false },
                        },
                    },
                    plugins: {
                        tooltip: {
                            ...chartTooltip,
                            callbacks: {
                                label: (ctx: any) => {
                                    const it = props.items[ctx.dataIndex];
                                    if (!it?.price) return "";
                                    const free = it.free ? ` · ${props.freeLabel}` : "";
                                    return ` $${it.price < 10 ? it.price.toFixed(2) : it.price.toFixed(0)}/M · ${it.model}${free}`;
                                },
                            },
                        },
                        legend: { display: false },
                    },
                } satisfies ChartOptions<"line">,
                plugins: [
                    {
                        id: "modelLabels",
                        afterDatasetsDraw(chart: Chart) {
                            const { ctx } = chart;
                            const meta = chart.getDatasetMeta(0);
                            for (let i = 0; i < meta.data.length; i++) {
                                const pt = meta.data[i]!;
                                const it = props.items[i];
                                if (!it?.model) continue;
                                const name = it.model.length > 16 ? it.model.slice(0, 14) + "…" : it.model;
                                ctx.save();
                                ctx.font = chartFont("600", 9);
                                ctx.fillStyle = it.free ? freeColor() : c.textDefault;
                                ctx.textAlign = "center";
                                ctx.textBaseline = "bottom";
                                ctx.fillText(name, pt.x, pt.y - 9);
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
