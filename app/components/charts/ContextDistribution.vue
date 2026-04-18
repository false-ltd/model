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
    import { Chart, BarElement, BarController, CategoryScale, LinearScale, Tooltip, type ChartOptions } from "chart.js";

    Chart.register(BarElement, BarController, CategoryScale, LinearScale, Tooltip);

    const props = defineProps<{
        title: string;
        subtitle: string;
        items: { limit_context: number }[];
    }>();

    const canvasRef = ref<HTMLCanvasElement | null>(null);
    let chart: Chart | null = null;

    const buckets = computed(() => {
        const ranges = [
            { label: "< 8K", min: 0, max: 8192 },
            { label: "8K–32K", min: 8192, max: 32768 },
            { label: "32K–128K", min: 32768, max: 131072 },
            { label: "128K–256K", min: 131072, max: 262144 },
            { label: "256K–1M", min: 262144, max: 1048576 },
            { label: "> 1M", min: 1048576, max: Infinity },
        ];
        return ranges.map((r, i) => ({
            label: r.label,
            count: props.items.filter((m) => m.limit_context >= r.min && m.limit_context < r.max).length,
            color: chartColor(i),
        }));
    });

    const renderChart = () => {
        chart?.destroy();
        if (!canvasRef.value || !buckets.value.length) return;

        const mutedColor =
            getComputedStyle(document.documentElement).getPropertyValue("--ui-text-muted").trim() || "#a8a29e";

        chart = new Chart(canvasRef.value, {
            type: "bar",
            data: {
                labels: buckets.value.map((b) => b.label),
                datasets: [
                    {
                        data: buckets.value.map((b) => b.count),
                        backgroundColor: buckets.value.map((b) => b.color + "cc"),
                        borderRadius: 6,
                        borderSkipped: false,
                        barPercentage: 0.7,
                    },
                ],
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    x: {
                        grid: { display: false },
                        ticks: { color: mutedColor, font: { size: 11 } },
                        border: { display: false },
                    },
                    y: {
                        beginAtZero: true,
                        grid: { color: mutedColor + "20" },
                        ticks: { color: mutedColor, font: { size: 10 } },
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
                            label: (ctx) => ` ${ctx.parsed.y} models`,
                        },
                    },
                    legend: { display: false },
                },
            } satisfies ChartOptions<"bar">,
            plugins: [
                {
                    id: "barCountLabels",
                    afterDatasetsDraw(chart) {
                        const { ctx } = chart;
                        const meta = chart.getDatasetMeta(0);
                        for (let i = 0; i < meta.data.length; i++) {
                            const bar = meta.data[i]!;
                            const val = buckets.value[i]?.count;
                            if (val == null || val === 0) continue;
                            ctx.save();
                            ctx.font = "bold 10px Inter, sans-serif";
                            ctx.fillStyle = mutedColor;
                            ctx.textAlign = "center";
                            ctx.textBaseline = "bottom";
                            ctx.fillText(String(val), bar.x, bar.y - 4);
                            ctx.restore();
                        }
                    },
                },
            ],
        });
    };

    watch(buckets, () => nextTick(renderChart), { deep: true });
    onMounted(renderChart);
    onUnmounted(() => chart?.destroy());
</script>
