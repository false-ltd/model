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
    import { Chart, LineController, LineElement, PointElement, LinearScale, CategoryScale, Filler, Tooltip } from "chart.js";

    Chart.register(LineController, LineElement, PointElement, LinearScale, CategoryScale, Filler, Tooltip);

    const props = defineProps<{
        title: string;
        subtitle: string;
        items: { month: string; count: number; total: number }[];
        lineColor?: string;
        fillColor?: string;
    }>();

    const canvasRef = ref<HTMLCanvasElement | null>(null);
    let chart: Chart | null = null;

    const renderChart = () => {
        chart?.destroy();
        chart = null;
        if (!canvasRef.value || !props.items.length) return;

        const lineColor = props.lineColor || chartColor(0);
        const fillColor = props.fillColor || lineColor + "1a";

        chart = new Chart(canvasRef.value, {
            type: "line",
            data: {
                labels: props.items.map((i) => formatMonth(i.month)),
                datasets: [
                    {
                        label: "Models",
                        data: props.items.map((i) => i.total),
                        borderColor: lineColor,
                        backgroundColor: fillColor,
                        fill: true,
                        tension: 0.3,
                        pointRadius: 3,
                        pointHoverRadius: 5,
                        pointBackgroundColor: lineColor,
                        borderWidth: 2,
                    },
                ],
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: { display: false },
                    tooltip: {
                        callbacks: {
                            title: (ctx) => props.items[ctx[0].dataIndex]?.month || "",
                            label: (ctx) => {
                                const item = props.items[ctx.dataIndex];
                                return item ? `${item.total} total (+${item.count} new)` : "";
                            },
                        },
                    },
                },
                scales: {
                    x: { grid: { display: false }, ticks: { font: { size: 10 } } },
                    y: { grid: { color: "rgba(0,0,0,0.06)" }, ticks: { font: { size: 10 } }, beginAtZero: true },
                },
            },
        });
    };

    watch(() => props.items, () => nextTick(renderChart), { deep: true });
    onMounted(() => nextTick(renderChart));
    onUnmounted(() => chart?.destroy());
</script>
