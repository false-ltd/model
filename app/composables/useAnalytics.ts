interface Model {
    id: string;
    provider_id: string;
    name: string;
    family: string;
    attachment: boolean;
    reasoning: boolean;
    tool_call: boolean;
    temperature: boolean;
    open_weights: boolean;
    modalities_input: string[];
    cost_input: number | null;
    cost_output: number | null;
    cost_cache_read: number | null;
    cost_cache_write: number | null;
    limit_context: number | null;
    limit_output: number | null;
    providers: { id: string; name: string };
}

export const useAnalytics = () => {
    const median = (nums: number[]): number => {
        const sorted = [...nums].sort((a, b) => a - b);
        const mid = Math.floor(sorted.length / 2);
        return sorted.length % 2 ? sorted[mid] : (sorted[mid - 1] + sorted[mid]) / 2;
    };

    const computeStats = (models: Model[]) => {
        const totalModels = models.length;
        const providers = new Set(models.map((m) => m.provider_id));
        const freeModels = models.filter((m) => (m.cost_input ?? 0) === 0 && (m.cost_output ?? 0) === 0);
        const reasoningModels = models.filter((m) => m.reasoning);
        const nonFreeModels = models.filter((m) => (m.cost_input ?? 0) > 0);
        const medianInputPrice = nonFreeModels.length > 0 ? median(nonFreeModels.map((m) => m.cost_input!)) : 0;
        const contextSizes = models.filter((m) => m.limit_context).map((m) => m.limit_context!);
        const medianContext = contextSizes.length > 0 ? median(contextSizes) : 0;

        return {
            totalProviders: providers.size,
            totalModels,
            freeModelsCount: freeModels.length,
            freeModelsPct: totalModels > 0 ? Math.round((freeModels.length / totalModels) * 100) : 0,
            medianInputPrice,
            reasoningCount: reasoningModels.length,
            reasoningPct: totalModels > 0 ? Math.round((reasoningModels.length / totalModels) * 100) : 0,
            medianContext,
        };
    };

    const topKByField = (models: Model[], field: "cost_input" | "cost_output" | "limit_context" | "limit_output", k: number = 10, order: "asc" | "desc" = "asc") => {
        return [...models]
            .filter((m) => m[field] != null)
            .sort((a, b) => order === "asc" ? (a[field]! - b[field]!) : (b[field]! - a[field]!))
            .slice(0, k);
    };

    const priceTierDistribution = (models: Model[]) => {
        const tiers = { free: 0, budget: 0, standard: 0, premium: 0, enterprise: 0 };
        for (const m of models) {
            const input = m.cost_input ?? 0;
            const output = m.cost_output ?? 0;
            const maxPrice = Math.max(input, output);
            if (maxPrice === 0) tiers.free++;
            else if (maxPrice < 1) tiers.budget++;
            else if (maxPrice < 5) tiers.standard++;
            else if (maxPrice < 15) tiers.premium++;
            else tiers.enterprise++;
        }
        return tiers;
    };

    const capabilityCoverage = (models: Model[]) => {
        const total = models.length;
        return {
            reasoning: { count: models.filter((m) => m.reasoning).length, total },
            tool_call: { count: models.filter((m) => m.tool_call).length, total },
            vision: { count: models.filter((m) => m.modalities_input.includes("image")).length, total },
            attachment: { count: models.filter((m) => m.attachment).length, total },
            temperature: { count: models.filter((m) => m.temperature).length, total },
        };
    };

    const topProviders = (models: Model[], k: number = 10) => {
        const counts: Record<string, { id: string; name: string; count: number }> = {};
        for (const m of models) {
            const pid = m.provider_id;
            if (!counts[pid]) counts[pid] = { id: pid, name: m.providers?.name || pid, count: 0 };
            counts[pid].count++;
        }
        return Object.values(counts).sort((a, b) => b.count - a.count).slice(0, k);
    };

    return { computeStats, topKByField, priceTierDistribution, capabilityCoverage, topProviders };
};
