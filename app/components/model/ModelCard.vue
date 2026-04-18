<template>
    <NuxtLink
        :to="localePath(`/model/${model.id}`)"
        class="block bg-default border border-default rounded-xl p-5 hover:border-primary/50 hover:shadow-md transition-all no-underline"
    >
        <div class="flex items-center gap-2.5 mb-3.5">
            <img
                :src="`https://models.dev/logos/${model.provider_id}.svg`"
                :alt="providerName"
                class="w-9 h-9 rounded-lg p-1"
                @error="($event.target as HTMLImageElement).style.display = 'none'"
            />
            <div class="min-w-0 flex-1">
                <div class="font-semibold text-sm text-default truncate">{{ model.name }}</div>
                <div class="text-xs text-muted">{{ model.family }}</div>
            </div>
            <ModelBadge v-if="isFree" type="free" />
        </div>

        <div class="flex gap-1 mb-3.5 flex-wrap">
            <ModelBadge v-if="model.reasoning" type="reasoning" />
            <ModelBadge v-if="model.tool_call" type="tool_call" />
            <ModelBadge v-if="hasVision" type="vision" />
            <ModelBadge v-if="model.attachment" type="attachment" />
        </div>

        <div class="border-t border-default pt-3 flex justify-between text-xs">
            <div class="flex-1">
                <div class="text-xs text-muted">{{ $t("catalog.input") }} / {{ $t("catalog.output") }}</div>
                <div class="font-semibold text-default" v-if="isFree">{{ $t("common.free") }}</div>
                <div class="font-semibold text-default" v-else>${{ model.cost_input }} / ${{ model.cost_output }}</div>
            </div>
            <div class="flex items-end gap-2">
                <div class="text-right">
                    <div class="text-xs text-muted">{{ $t("catalog.context") }}</div>
                    <div class="font-semibold text-default">{{ formatTokens(model.limit_context) }}</div>
                </div>
                <button
                    @click.prevent.stop="toggleCompare"
                    class="rounded-lg px-2 py-1 text-xs cursor-pointer transition-colors border"
                    :class="isInCompare ? 'bg-primary text-white border-primary' : 'bg-elevated text-toned border-default hover:border-primary hover:text-primary'"
                >
                    {{ isInCompare ? "✓" : "+" }}
                </button>
            </div>
        </div>
    </NuxtLink>
</template>

<script setup lang="ts">
    const localePath = useLocalePath();
    const { modelIds, addModel, removeModel } = useCompare();

    const props = defineProps<{
        model: any;
    }>();

    const providerName = computed(() => props.model.providers?.name || props.model.provider_id);
    const isFree = computed(() => (props.model.cost_input ?? 0) === 0 && (props.model.cost_output ?? 0) === 0);
    const hasVision = computed(() => props.model.modalities_input?.includes("image"));
    const isInCompare = computed(() => modelIds.value.includes(props.model.id));

    const toggleCompare = () => {
        if (isInCompare.value) {
            removeModel(props.model.id);
        } else {
            addModel(props.model.id);
        }
    };
</script>
