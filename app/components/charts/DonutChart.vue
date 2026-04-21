<template>
    <div class="bg-default border border-default rounded-xl p-5 hover:border-accented transition-colors">
        <div class="text-sm font-semibold text-default mb-1">{{ title }}</div>
        <div class="text-xs text-muted mb-3">{{ subtitle }}</div>
        <div class="flex items-center gap-4">
            <div class="shrink-0" style="width: 140px; height: 140px; position: relative">
                <canvas ref="canvasRef" />
            </div>
            <div class="flex flex-col gap-2 flex-1 min-w-0">
                <div v-for="seg in segments" :key="seg.label" class="flex items-center gap-2 text-xs">
                    <div class="w-2.5 h-2.5 rounded-sm shrink-0" :style="{ background: seg.color }" />
                    <span class="text-toned truncate">{{ seg.label }}</span>
                    <span class="text-default font-semibold ml-auto">{{ seg.count }}</span>
                    <span class="text-muted w-10 text-right">{{ pct(seg.count) }}%</span>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { Chart, ArcElement, DoughnutController, Tooltip, type ChartOptions, type Chart as ChartType } from "chart.js";

Chart.register(ArcElement, DoughnutController, Tooltip);

interface Segment {
    label: string;
    count: number;
    color: string;
}

const props = defineProps<{
    title: string;
    subtitle: string;
    segments: Segment[];
    total: number;
}>();

const pct = (count: number) => (props.total > 0 ? Math.round((count / props.total) * 100) : 0);

const { canvasRef } = useChart(
    () => {
        if (!props.segments.length) return null;

        const c = chartColors();

        return {
            type: "doughnut",
            data: {
                labels: props.segments.map((s) => s.label),
                datasets: [
                    {
                        data: props.segments.map((s) => s.count),
                        backgroundColor: props.segments.map((s) => s.color),
                        borderWidth: 2,
                        borderColor: c.bg,
                        hoverOffset: 6,
                    },
                ],
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                animation: chartAnimation,
                cutout: "68%",
                plugins: {
                    tooltip: {
                        ...chartTooltip,
                        callbacks: {
                            label: (ctx) => ` ${ctx.label}: ${ctx.parsed} (${pct(ctx.parsed as number)}%)`,
                        },
                    },
                    legend: { display: false },
                },
            } satisfies ChartOptions<"doughnut">,
            plugins: [
                {
                    id: "centerText",
                    afterDraw(chart: ChartType) {
                        const { ctx, width, height } = chart;
                        ctx.save();
                        ctx.font = chartFont("bold", 20);
                        ctx.fillStyle = c.textDefault;
                        ctx.textAlign = "center";
                        ctx.textBaseline = "middle";
                        ctx.fillText(String(props.total), width / 2, height / 2 - 8);
                        ctx.font = chartFont("normal", 10);
                        ctx.fillStyle = c.text;
                        ctx.fillText("models", width / 2, height / 2 + 10);
                        ctx.restore();
                    },
                },
            ],
        } as any;
    },
    () => props.segments,
);
</script>
