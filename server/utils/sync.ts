import { createClient } from "@supabase/supabase-js";

interface ApiProvider {
    id: string;
    name: string;
    npm: string;
    env: string[];
    doc: string;
    api?: string;
    models: Record<string, ApiModel>;
}

interface ApiModel {
    id: string;
    name: string;
    family: string;
    attachment: boolean;
    reasoning: boolean;
    tool_call: boolean;
    structured_output?: boolean;
    temperature?: boolean;
    knowledge?: string;
    release_date: string;
    last_updated: string;
    open_weights: boolean;
    interleaved?: boolean | { field: string };
    modalities: { input: string[]; output: string[] };
    cost: {
        input: number;
        output: number;
        cache_read?: number;
        cache_write?: number;
        reasoning?: number;
        input_audio?: number;
        output_audio?: number;
    };
    limit: { context: number; input?: number; output: number };
    status?: string;
}

export async function performSync() {
    const config = useRuntimeConfig();
    const client = createClient(config.public.supabase.url, config.public.supabase.key);

    const response = await $fetch<Record<string, ApiProvider>>("https://models.dev/api.json");

    const providerMap = new Map<string, any>();
    const modelMap = new Map<string, any>();

    for (const [providerId, provider] of Object.entries(response)) {
        providerMap.set(providerId, {
            provider_id: providerId,
            name: provider.name,
            npm: provider.npm,
            env: provider.env,
            doc: provider.doc,
            api: provider.api || null,
        });

        for (const [modelId, model] of Object.entries(provider.models || {})) {
            modelMap.set(modelId, {
                model_id: modelId,
                provider_id: providerId,
                name: model.name,
                family: model.family,
                attachment: model.attachment,
                reasoning: model.reasoning,
                tool_call: model.tool_call,
                structured_output: model.structured_output ?? null,
                temperature: model.temperature ?? true,
                knowledge: model.knowledge || null,
                release_date: model.release_date || null,
                last_updated: model.last_updated || null,
                open_weights: model.open_weights,
                status: model.status || null,
                interleaved: model.interleaved ? JSON.stringify(model.interleaved) : null,
                modalities_input: model.modalities?.input || [],
                modalities_output: model.modalities?.output || [],
                cost_input: model.cost?.input ?? null,
                cost_output: model.cost?.output ?? null,
                cost_reasoning: model.cost?.reasoning ?? null,
                cost_cache_read: model.cost?.cache_read ?? null,
                cost_cache_write: model.cost?.cache_write ?? null,
                cost_input_audio: model.cost?.input_audio ?? null,
                cost_output_audio: model.cost?.output_audio ?? null,
                limit_context: model.limit?.context ?? null,
                limit_input: model.limit?.input ?? null,
                limit_output: model.limit?.output ?? null,
            });
        }
    }

    const providers = [...providerMap.values()];
    const models = [...modelMap.values()];

    const { error: providerError } = await client.from("providers").upsert(providers, { onConflict: "provider_id" });

    if (providerError) {
        throw createError({ statusCode: 500, statusMessage: providerError.message });
    }

    const { error: modelError } = await client.from("models").upsert(models, { onConflict: "model_id" });

    if (modelError) {
        throw createError({ statusCode: 500, statusMessage: modelError.message });
    }

    const syncedAt = new Date().toISOString();
    await client.from("providers").update({ synced_at: syncedAt }).neq("provider_id", "");

    return {
        synced_at: syncedAt,
        providers: providers.length,
        models: models.length,
    };
}
