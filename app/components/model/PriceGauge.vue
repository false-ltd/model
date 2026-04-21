<template>
    <div style="height: 100%; min-height: 200px; position: relative">
        <canvas ref="canvasRef" />
    </div>
</template>

<script setup lang="ts">
    import { Chart, BarController, CategoryScale, LinearScale, BarElement, Tooltip, type ChartMeta, type Chart as ChartType } from "chart.js";

    Chart.register(BarController, CategoryScale, LinearScale, BarElement, Tooltip);

    const props = defineProps<{
        fields: { key: string; label: string; value: number | null; display: string }[];
    }>();

    const COLOR_INDEX: Record<string, number> = {
        input: 0, output: 1, cache_read: 2, cache_write: 3,
        reasoning: 4, input_audio: 5, output_audio: 6,
    };

    const { canvasRef } = useChart(
        () => {
            const c = chartColors();
            const fields = props.fields;

            return {
                type: "bar",
                data: {
                    labels: fields.map((f) => f.label),
                    datasets: [{
                        data: fields.map((f) => f.value ?? 0),
                        backgroundColor: fields.map((f) => chartColor(COLOR_INDEX[f.key] ?? 0) + "cc"),
                        hoverBackgroundColor: fields.map((f) => chartColor(COLOR_INDEX[f.key] ?? 0)),
                        borderRadius: 4,
                        borderWidth: 0,
                        borderSkipped: false,
                    }],
                },
                plugins: [{
                    id: "priceLabels",
                    afterDraw(chart: ChartType) {
                        const ctx = chart.ctx;
                        const meta = chart.getDatasetMeta(0) as ChartMeta<"bar">;
                        const xScale = chart.scales.x!;
                        ctx.save();
                        ctx.font = `${chartFont("600", 11).replace(chartFontFamily, "monospace")}`;
                        ctx.textBaseline = "middle";
                        meta.data.forEach((bar, i) => {
                            const display = fields[i]?.display || "";
                            const x = Math.min(bar.x + 6, xScale.right - ctx.measureText(display).width - 4);
                            const y = bar.y;
                            ctx.fillStyle = fields[i]?.value === 0 ? getCSSVar("--color-success") : c.textDefault;
                            ctx.fillText(display, x, y);
                        });
                        ctx.restore();
                    },
                }],
                options: {
                    indexAxis: "y",
                    responsive: true,
                    maintainAspectRatio: false,
                    layout: { padding: { right: 60 } },
                    animation: chartAnimation,
                    plugins: {
                        legend: { display: false },
                        tooltip: {
                            ...chartTooltip,
                            bodyFont: { size: 11, family: "monospace" },
                            callbacks: {
                                label: (ctx: { dataIndex: number; parsed: { x: number } }) => ` ${fields[ctx.dataIndex]?.display || ctx.parsed.x}`,
                            },
                        },
                    },
                    scales: {
                        x: {
                            grid: { color: c.grid },
                            ticks: {
                                font: chartTickFont(10),
                                color: c.text,
                                callback: (val: any) => "$" + val,
                            },
                        },
                        y: {
                            grid: { display: false },
                            ticks: { font: chartTickFont(), color: c.text },
                        },
                    },
                },
            } as any;
        },
        () => props.fields,
    );
</script>
