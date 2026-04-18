<template>
    <div class="p-5 border-b border-default">
        <div v-if="tabs?.length" class="flex gap-1 mb-3">
            <button
                v-for="tab in tabs"
                :key="tab.key"
                @click="$emit('tabChange', tab.key)"
                class="rounded-md px-2 py-0.5 text-[10px] cursor-pointer transition-colors"
                :class="activeTab === tab.key ? 'bg-primary text-white' : 'bg-elevated text-toned hover:bg-accented'"
            >
                {{ tab.label }}
            </button>
        </div>
        <div :style="{ height: (height || 160) + 'px', position: 'relative' }">
            <canvas ref="canvasRef" />
        </div>
    </div>
</template>

<script setup lang="ts">
    import { Chart, BarController, CategoryScale, LinearScale, BarElement, Tooltip, Legend } from "chart.js";

    Chart.register(BarController, CategoryScale, LinearScale, BarElement, Tooltip, Legend);

    const props = defineProps<{
        type: "bar" | "horizontalBar";
        labels: string[];
        datasets: { label: string; data: number[]; backgroundColor: string | string[] }[];
        tabs?: { key: string; label: string }[];
        activeTab?: string;
        height?: number;
    }>();

    defineEmits<{ tabChange: [key: string] }>();

    const canvasRef = ref<HTMLCanvasElement | null>(null);
    let chart: Chart | null = null;

    const getCSSColor = (v: string) => getComputedStyle(document.documentElement).getPropertyValue(v).trim();

    const renderChart = () => {
        chart?.destroy();
        chart = null;
        if (!canvasRef.value) return;
        if (!props.labels.length || !props.datasets.length) return;

        const textColor = getCSSColor("--ui-text");
        const gridColor = "rgba(0,0,0,0.06)";
        const isHorizontal = props.type === "horizontalBar";

        chart = new Chart(canvasRef.value, {
            type: "bar",
            data: { labels: props.labels, datasets: props.datasets },
            options: {
                indexAxis: isHorizontal ? "y" : "x",
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        display: props.datasets.length > 1,
                        position: "bottom",
                        labels: { boxWidth: 10, font: { size: 11 }, color: textColor },
                    },
                    tooltip: { enabled: true },
                },
                scales: {
                    x: { grid: { display: false }, ticks: { font: { size: 11 }, color: textColor } },
                    y: { grid: { color: gridColor }, ticks: { font: { size: 11 }, color: textColor } },
                },
            },
        });
    };

    watch(() => [props.labels, props.datasets], () => nextTick(renderChart), { deep: true });
    onMounted(() => nextTick(renderChart));
    onUnmounted(() => chart?.destroy());
</script>
