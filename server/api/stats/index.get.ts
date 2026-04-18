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
        cheapestInputRes, cheapestOutputRes,
        topContextRes, topOutputRes,
        allPricesRes,
        toolCallRes, visionRes, attachmentRes, temperatureRes,
        providerCountsRes,
        recentModelsRes,
        allModelDataRes,
        openCountRes,
        modalityDataRes,
    ] = await Promise.all([
        // Median input price
        client.from("models").select("cost_input").gt("cost_input", 0).order("cost_input", { ascending: true }),
        // Median context
        client.from("models").select("limit_context").not("limit_context", "is", null).order("limit_context", { ascending: true }),
        // Top 10 cheapest (input)
        client.from("models")
            .select("id, model_id, name, cost_input, cost_output, limit_context, providers(provider_id, name)")
            .order("cost_input", { ascending: true }).limit(10),
        // Top 10 cheapest (output)
        client.from("models")
            .select("id, model_id, name, cost_input, cost_output, limit_context, providers(provider_id, name)")
            .gt("cost_output", 0).order("cost_output", { ascending: true }).limit(10),
        // Top 10 context
        client.from("models")
            .select("id, model_id, name, limit_context, limit_output, providers(provider_id, name)")
            .not("limit_context", "is", null).order("limit_context", { ascending: false }).limit(10),
        // Top 10 output
        client.from("models")
            .select("id, model_id, name, limit_output, providers(provider_id, name)")
            .not("limit_output", "is", null).order("limit_output", { ascending: false }).limit(10),
        // Price tier distribution
        client.from("models").select("cost_input, cost_output"),
        // Capability counts
        client.from("models").select("id", { count: "exact", head: true }).eq("tool_call", true),
        client.from("models").select("id", { count: "exact", head: true }).contains("modalities_input", ["image"]),
        client.from("models").select("id", { count: "exact", head: true }).eq("attachment", true),
        client.from("models").select("id", { count: "exact", head: true }).eq("temperature", true),
        // Provider counts
        client.from("models").select("provider_id, providers(provider_id, name)"),
        // Release timeline — all models with release_date
        client.from("models")
            .select("name, release_date, provider_id, providers(provider_id, name)")
            .not("release_date", "is", null).order("release_date", { ascending: false }),
        // Value leaderboard + scatter data (single query for both)
        client.from("models")
            .select("id, model_id, name, cost_input, cost_output, limit_context, limit_output, reasoning, tool_call, provider_id, providers(provider_id, name)")
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

    const monthCounts: Record<string, number> = {};
    for (const m of recentModelsRes.data || []) {
        const month = (m.release_date as string).substring(0, 7);
        if (!month.includes("-")) continue;
        if (month < "2026-01") continue;
        monthCounts[month] = (monthCounts[month] || 0) + 1;
    }
    const releaseTimeline = Object.entries(monthCounts)
        .sort(([a], [b]) => a.localeCompare(b))
        .slice(-12)
        .map(([month, count]) => ({ month, count }));

    const allMonths = Object.entries(monthCounts).sort(([a], [b]) => a.localeCompare(b));
    const cumulativeReleases: { month: string; count: number; total: number }[] = [];
    let cumulative = 0;
    for (const [month, count] of allMonths) {
        cumulative += count;
        cumulativeReleases.push({ month, count, total: cumulative });
    }

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
            cheapestInput: cheapestInputRes.data || [],
            cheapestOutput: cheapestOutputRes.data || [],
            topContext: topContextRes.data || [],
            topOutput: topOutputRes.data || [],
            tiers, capabilities, topProviders: topProviderList,
            releaseTimeline, cumulativeReleases, modalities,
            weightsDistribution,
            contextDistribution: allModelData.map((m: any) => ({ limit_context: m.limit_context })),
        },
    };
});
