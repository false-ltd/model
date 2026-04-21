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
import { Chart, ArcElement, DoughnutController, Tooltip, type ChartOptions } from "chart.js";

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

const canvasRef = ref<HTMLCanvasElement | null>(null);
let chart: Chart | null = null;

const pct = (count: number) => (props.total > 0 ? Math.round((count / props.total) * 100) : 0);



const renderChart = () => {
    chart?.destroy();
    if (!canvasRef.value || !props.segments.length) return;

    const textColor = getCSSVar("--ui-text-default") || getCSSVar("--ui-text");
    const mutedColor = getCSSVar("--ui-text-muted");

    chart = new Chart(canvasRef.value, {
        type: "doughnut",
        data: {
            labels: props.segments.map((s) => s.label),
            datasets: [
                {
                    data: props.segments.map((s) => s.count),
                    backgroundColor: props.segments.map((s) => s.color),
                    borderWidth: 2,
                    borderColor: getCSSVar("--ui-bg"),
                    hoverOffset: 6,
                },
            ],
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            cutout: "68%",
            plugins: {
                tooltip: {
                    backgroundColor: "rgba(0,0,0,0.8)",
                    titleFont: { size: 12 },
                    bodyFont: { size: 11 },
                    padding: 8,
                    cornerRadius: 6,
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
                afterDraw(chart) {
                    const { ctx, width, height } = chart;
                    ctx.save();
                    ctx.font = "bold 20px Inter, sans-serif";
                    ctx.fillStyle = textColor;
                    ctx.textAlign = "center";
                    ctx.textBaseline = "middle";
                    ctx.fillText(String(props.total), width / 2, height / 2 - 8);
                    ctx.font = "10px Inter, sans-serif";
                    ctx.fillStyle = mutedColor;
                    ctx.fillText("models", width / 2, height / 2 + 10);
                    ctx.restore();
                },
            },
        ],
    });
};

watch(() => props.segments, () => nextTick(renderChart), { deep: true });
onMounted(renderChart);
onUnmounted(() => chart?.destroy());
</script>
