<template>
    <div class="bg-default border border-default rounded-xl p-5 hover:border-accented transition-colors">
        <div class="flex justify-between items-center mb-3">
            <div>
                <div class="text-sm font-semibold text-default">{{ title }}</div>
                <div class="text-xs text-muted">{{ subtitle }}</div>
            </div>
            <div v-if="tabs?.length" class="flex gap-1">
                <button
                    v-for="tab in tabs"
                    :key="tab.key"
                    @click="$emit('tabChange', tab.key)"
                    class="rounded-md px-2 py-0.5 text-[10px] cursor-pointer transition-colors"
                    :class="
                        activeTab === tab.key
                            ? 'bg-primary text-white'
                            : 'bg-elevated text-toned hover:bg-accented'
                    "
                >
                    {{ tab.label }}
                </button>
            </div>
        </div>
        <div style="height: 180px; position: relative">
            <canvas ref="canvasRef" />
        </div>
        <div
            v-if="linkText"
            class="mt-2.5 text-center text-xs text-primary cursor-pointer hover:underline"
            @click="$emit('navigate')"
        >
            {{ linkText }}
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

interface ChartItem {
    label: string;
    value: number;
    display: string;
    color: string;
}

const props = defineProps<{
    title: string;
    subtitle: string;
    items: ChartItem[];
    tabs?: { key: string; label: string }[];
    activeTab?: string;
    linkText?: string;
}>();

defineEmits<{
    tabChange: [key: string];
    navigate: [];
}>();

const canvasRef = ref<HTMLCanvasElement | null>(null);
let chart: Chart | null = null;

const renderChart = () => {
    chart?.destroy();
    if (!canvasRef.value || !props.items.length) return;

    const mutedColor =
        getComputedStyle(document.documentElement).getPropertyValue("--ui-text-muted").trim() || "#a8a29e";

    chart = new Chart(canvasRef.value, {
        type: "bar",
        data: {
            labels: props.items.map((i) => i.label),
            datasets: [
                {
                    data: props.items.map((i) => i.value),
                    backgroundColor: props.items.map((i) => i.color + "cc"),
                    borderRadius: 4,
                    borderSkipped: false,
                    barThickness: 18,
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
                    ticks: {
                        color: mutedColor,
                        font: { size: 11 },
                        callback: function (this: any, value: number) {
                            const label = this.getLabelForValue(value) as string;
                            return label.length > 18 ? label.slice(0, 16) + "…" : label;
                        },
                    },
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
                        label: (ctx) => ` ${props.items[ctx.dataIndex]?.display || ctx.parsed.x}`,
                    },
                },
                legend: { display: false },
            },
        } satisfies ChartOptions<"bar">,
        plugins: [
            {
                id: "barLabels",
                afterDatasetsDraw(chart) {
                    const { ctx } = chart;
                    const meta = chart.getDatasetMeta(0);
                    for (let i = 0; i < meta.data.length; i++) {
                        const bar = meta.data[i]!;
                        const item = props.items[i];
                        if (!item) continue;
                        ctx.save();
                        ctx.font = "bold 10px Inter, sans-serif";
                        ctx.fillStyle = item.color;
                        ctx.textAlign = "left";
                        ctx.textBaseline = "middle";
                        ctx.fillText(item.display, bar.x + 6, bar.y);
                        ctx.restore();
                    }
                },
            },
        ],
    });
};

watch(() => props.items, () => nextTick(renderChart), { deep: true });
onMounted(renderChart);
onUnmounted(() => chart?.destroy());
</script>
