<template>
    <div style="height: 100%; min-height: 200px; position: relative">
        <canvas ref="canvasRef" />
    </div>
</template>

<script setup lang="ts">
    import { Chart, BarController, CategoryScale, LinearScale, BarElement, Tooltip, type ChartMeta, type Chart as ChartType } from "chart.js";

    Chart.register(BarController, CategoryScale, LinearScale, BarElement, Tooltip);

    const props = defineProps<{
        contextLabel: string;
        inputLabel: string;
        outputLabel: string;
        context: number;
        input: number;
        output: number;
    }>();

    const formatVal = (n: number) => {
        if (!n) return "\u2014";
        return formatTokens(n);
    };

    const { canvasRef } = useChart(
        () => {
            const c = chartColors();

            const values = [props.context || 0, props.input || 0, props.output || 0];
            const labels = [formatVal(props.context), formatVal(props.input), formatVal(props.output)];

            return {
                type: "bar",
                data: {
                    labels: [props.contextLabel, props.inputLabel, props.outputLabel],
                    datasets: [{
                        data: values,
                        backgroundColor: [0, 1, 2].map((i) => chartColor(i) + "cc"),
                        hoverBackgroundColor: [0, 1, 2].map((i) => chartColor(i)),
                        borderRadius: 4,
                        borderWidth: 0,
                        borderSkipped: false,
                    }],
                },
                plugins: [{
                    id: "limitLabels",
                    afterDraw(chart: ChartType) {
                        const ctx = chart.ctx;
                        const meta = chart.getDatasetMeta(0) as ChartMeta<"bar">;
                        const xScale = chart.scales.x!;
                        ctx.save();
                        ctx.font = "600 11px monospace";
                        ctx.textBaseline = "middle";
                        meta.data.forEach((bar, i) => {
                            const label = labels[i] ?? "";
                            const x = Math.min(bar.x + 6, xScale.right - ctx.measureText(label).width - 4);
                            ctx.fillStyle = c.textDefault;
                            ctx.fillText(label, x, bar.y);
                        });
                        ctx.restore();
                    },
                }],
                options: {
                    indexAxis: "y",
                    responsive: true,
                    maintainAspectRatio: false,
                    layout: { padding: { right: 50 } },
                    animation: chartAnimation,
                    plugins: {
                        legend: { display: false },
                        tooltip: {
                            ...chartTooltip,
                            bodyFont: { size: 11, family: "monospace" },
                            callbacks: {
                                label: (ctx: { parsed: { x: number } }) => ` ${formatVal(ctx.parsed.x)} tokens`,
                            },
                        },
                    },
                    scales: {
                        x: {
                            grid: { color: c.grid },
                            ticks: {
                                font: chartTickFont(10),
                                color: c.text,
                                callback: (val: any) => formatVal(val as number),
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
        () => [props.context, props.input, props.output],
    );
</script>
