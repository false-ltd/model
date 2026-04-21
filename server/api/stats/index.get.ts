import { serverSupabaseClient } from "#supabase/server";

const median = (arr: number[]): number => {
    if (arr.length === 0) return 0;
    return arr.length % 2
        ? arr[Math.floor(arr.length / 2)]
        : (arr[arr.length / 2 - 1] + arr[arr.length / 2]) / 2;
};

export default defineEventHandler(async (event) => {
    const client = await serverSupabaseClient(event);

    // Batch 1: KPI counts
    const [providersRes, totalRes, freeRes, reasoningRes] = await Promise.all([
        client.from("providers").select("id", { count: "exact", head: true }),
        client.from("models").select("id", { count: "exact", head: true }),
        client.from("models").select("id", { count: "exact", head: true }).eq("cost_input", 0).eq("cost_output", 0),
        client.from("models").select("id", { count: "exact", head: true }).eq("reasoning", true),
    ]);

    const totalProviders = providersRes.count || 0;
    const totalModels = totalRes.count || 0;
    const freeModelsCount = freeRes.count || 0;
    const reasoningCount = reasoningRes.count || 0;

    // Batch 2: All independent queries in parallel
    const [
        pricesRes, contextsRes,
        allPricesRes,
        toolCallRes, visionRes, attachmentRes, temperatureRes,
        providerCountsRes,
        allModelDataRes,
        openCountRes,
        modalityDataRes,
    ] = await Promise.all([
        // Median input price
        client.from("models").select("cost_input").gt("cost_input", 0).order("cost_input", { ascending: true }),
        // Median context
        client.from("models").select("limit_context").not("limit_context", "is", null).order("limit_context", { ascending: true }),
        // Price tier distribution
        client.from("models").select("cost_input, cost_output"),
        // Capability counts
        client.from("models").select("id", { count: "exact", head: true }).eq("tool_call", true),
        client.from("models").select("id", { count: "exact", head: true }).contains("modalities_input", ["image"]),
        client.from("models").select("id", { count: "exact", head: true }).eq("attachment", true),
        client.from("models").select("id", { count: "exact", head: true }).eq("temperature", true),
        // Provider counts
        client.from("models").select("provider_id, providers(provider_id, name)"),
        // Context distribution data
        client.from("models")
            .select("limit_context")
            .gt("limit_context", 0),
        // Open weights count
        client.from("models").select("id", { count: "exact", head: true }).eq("open_weights", true),
        // Modality breakdown
        client.from("models").select("modalities_input, modalities_output"),
    ]);

    // Derive results from batch 2
    const medianInputPrice = median((pricesRes.data || []).map((p: any) => p.cost_input));
    const medianContext = median((contextsRes.data || []).map((c: any) => c.limit_context));

    const tiers = { free: 0, budget: 0, standard: 0, premium: 0, enterprise: 0 };
    for (const m of allPricesRes.data || []) {
        const maxPrice = Math.max(m.cost_input || 0, m.cost_output || 0);
        if (maxPrice === 0) tiers.free++;
        else if (maxPrice < 1) tiers.budget++;
        else if (maxPrice < 5) tiers.standard++;
        else if (maxPrice < 15) tiers.premium++;
        else tiers.enterprise++;
    }

    const capabilities = {
        reasoning: { count: reasoningCount, total: totalModels },
        tool_call: { count: toolCallRes.count || 0, total: totalModels },
        vision: { count: visionRes.count || 0, total: totalModels },
        attachment: { count: attachmentRes.count || 0, total: totalModels },
        temperature: { count: temperatureRes.count || 0, total: totalModels },
    };

    const providerMap: Record<string, { id: string; name: string; count: number }> = {};
    for (const m of providerCountsRes.data || []) {
        const pid = m.provider_id;
        if (!providerMap[pid]) {
            providerMap[pid] = { id: pid, name: (m.providers as any)?.name || pid, count: 0 };
        }
        providerMap[pid].count++;
    }
    const topProviderList = Object.values(providerMap).sort((a, b) => b.count - a.count).slice(0, 5);

    const inputModCounts: Record<string, number> = {};
    const outputModCounts: Record<string, number> = {};
    for (const m of modalityDataRes.data || []) {
        for (const mod of m.modalities_input || []) inputModCounts[mod] = (inputModCounts[mod] || 0) + 1;
        for (const mod of m.modalities_output || []) outputModCounts[mod] = (outputModCounts[mod] || 0) + 1;
    }
    const modalities = {
        input: Object.entries(inputModCounts).map(([type, count]) => ({ type, count })).sort((a, b) => b.count - a.count),
        output: Object.entries(outputModCounts).map(([type, count]) => ({ type, count })).sort((a, b) => b.count - a.count),
    };

    const allModelData = allModelDataRes.data || [];

    const contextDistribution = allModelData.map((m: any) => ({ limit_context: m.limit_context }));

    const openCount = openCountRes.count || 0;
    const weightsDistribution = {
        open: openCount,
        closed: totalModels - openCount,
        openPct: totalModels > 0 ? Math.round((openCount / totalModels) * 100) : 0,
    };

    return {
        data: {
            stats: {
                totalProviders, totalModels, freeModelsCount,
                freeModelsPct: totalModels > 0 ? Math.round((freeModelsCount / totalModels) * 100) : 0,
                medianInputPrice, reasoningCount,
                reasoningPct: totalModels > 0 ? Math.round((reasoningCount / totalModels) * 100) : 0,
                medianContext,
            },
            tiers, capabilities, topProviders: topProviderList,
            modalities, weightsDistribution, contextDistribution,
        },
    };
});
