import type { Model, PricingField, ApiResponse } from "~/types";

export async function useModelDetail() {
    const route = useRoute();
    const { t } = useI18n();
    const config = useRuntimeConfig();

    const modelId = Number(route.params.modelId);

    const { data: result, error } = await useAsyncData(`model-${modelId}`, () =>
        $fetch<ApiResponse<Model>>(`${config.public.apiBase}/api/v1/models/${modelId}`),
    );

    // Distinguish "model does not exist" from transport/server failures —
    // the latter must not be reported as 404.
    if (error.value) {
        const statusCode = (error.value as any)?.statusCode ?? (error.value as any)?.status;
        if (statusCode === 404) {
            throw createError({ statusCode: 404, statusMessage: "Model not found" });
        }
        throw createError({ statusCode: 500, statusMessage: t("common.loadError") });
    }

    const model = computed(() => result.value?.data ?? null);

    const pricingFields = computed((): PricingField[] => {
        const m = model.value!;
        const fields = [
            { key: "input", label: t("detail.inputPrice"), value: m.cost_input, display: m.cost_input != null ? `$${m.cost_input}` : "—" },
            { key: "output", label: t("detail.outputPrice"), value: m.cost_output, display: m.cost_output != null ? `$${m.cost_output}` : "—" },
        ];
        if (m.cost_cache_read != null)
            fields.push({ key: "cache_read", label: t("detail.cacheRead"), value: m.cost_cache_read, display: `$${m.cost_cache_read}` });
        if (m.cost_cache_write != null)
            fields.push({ key: "cache_write", label: t("detail.cacheWrite"), value: m.cost_cache_write, display: `$${m.cost_cache_write}` });
        if (m.cost_reasoning != null)
            fields.push({ key: "reasoning", label: t("detail.costReasoning"), value: m.cost_reasoning, display: `$${m.cost_reasoning}` });
        if (m.cost_input_audio != null)
            fields.push({ key: "input_audio", label: t("detail.costInputAudio"), value: m.cost_input_audio, display: `$${m.cost_input_audio}` });
        if (m.cost_output_audio != null)
            fields.push({ key: "output_audio", label: t("detail.costOutputAudio"), value: m.cost_output_audio, display: `$${m.cost_output_audio}` });
        return fields;
    });

    return { model, pricingFields };
}
