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
    type ChartOptions,
} from "chart.js";

Chart.register(BarElement, BarController, CategoryScale, LinearScale, Tooltip);

const props = defineProps<{
    title: string;
    subtitle: string;
    capabilities: { label: string; count: number; total: number; pct: number; color: string }[];
}>();

const canvasRef = ref<HTMLCanvasElement | null>(null);
let chart: Chart | null = null;

const renderChart = () => {
    chart?.destroy();
    if (!canvasRef.value || !props.capabilities.length) return;

    const mutedColor =
        getComputedStyle(document.documentElement).getPropertyValue("--ui-text-muted").trim() || "#a8a29e";
    const bgColor =
        getComputedStyle(document.documentElement).getPropertyValue("--ui-bg-elevated").trim() || "#f5f5f4";

    chart = new Chart(canvasRef.value, {
        type: "bar",
        data: {
            labels: props.capabilities.map((c) => c.label),
            datasets: [
                {
                    label: "count",
                    data: props.capabilities.map((c) => c.pct),
                    backgroundColor: props.capabilities.map((c) => c.color + "cc"),
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
            scales: {
                x: {
                    display: false,
                    max: 100,
                },
                y: {
                    grid: { display: false },
                    ticks: { color: mutedColor, font: { size: 11 } },
                    border: { display: false },
                },
            },
            plugins: {
                tooltip: {
                    backgroundColor: "rgba(0,0,0,0.8)",
                    titleFont: { size: 12 },
                    bodyFont: { size: 11 },
                    padding: 8,
                    cornerRadius: 6,
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
                afterDatasetsDraw(chart) {
                    const { ctx } = chart;
                    const meta = chart.getDatasetMeta(0);
                    for (let i = 0; i < meta.data.length; i++) {
                        const bar = meta.data[i]!;
                        const cap = props.capabilities[i];
                        if (!cap) continue;
                        ctx.save();
                        ctx.font = "bold 10px Inter, sans-serif";
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
    });
};

watch(() => props.capabilities, () => nextTick(renderChart), { deep: true });
onMounted(renderChart);
onUnmounted(() => chart?.destroy());
</script>
