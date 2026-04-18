import { serverSupabaseClient } from "#supabase/server";

export default defineEventHandler(async (event) => {
    const query = getQuery(event);
    const client = await serverSupabaseClient(event);

    const page = Number(query.page) || 1;
    const pageSize = Number(query.pageSize) || 100;
    const search = query.q as string || "";
    const sort = (query.sort as string) || "name";
    const order = query.order === "desc" ? false : true;
    const providers = query.providers ? (query.providers as string).split(",") : [];
    const inputTypes = query.inputTypes ? (query.inputTypes as string).split(",") : [];
    const outputTypes = query.outputTypes ? (query.outputTypes as string).split(",") : [];
    const reasoning = query.reasoning === "true";
    const toolCall = query.toolCall === "true";
    const vision = query.vision === "true";
    const attachment = query.attachment === "true";
    const freeOnly = query.freeOnly === "true";
    const under1 = query.under1 === "true";
    const priceMin = Number(query.priceMin) || 0;
    const priceMax = Number(query.priceMax) ?? 100;
    const hasPriceFilter = priceMin > 0 || priceMax < 100;

    let qb = client
        .from("models")
        .select("*, providers(provider_id, name, npm, env, doc, api)", { count: "exact" });

    if (search) {
        qb = qb.or(`name.ilike.%${search}%,family.ilike.%${search}%,model_id.ilike.%${search}%`);
    }
    if (providers.length) qb = qb.in("provider_id", providers);
    if (inputTypes.length) qb = qb.contains("modalities_input", inputTypes);
    if (outputTypes.length) qb = qb.contains("modalities_output", outputTypes);
    if (reasoning) qb = qb.eq("reasoning", true);
    if (toolCall) qb = qb.eq("tool_call", true);
    if (vision) qb = qb.contains("modalities_input", ["image"]);
    if (attachment) qb = qb.eq("attachment", true);
    if (freeOnly) qb = qb.eq("cost_input", 0);
    if (under1) qb = qb.gt("cost_input", 0).lt("cost_input", 1);
    // Price range: models with null price always pass through
    if (hasPriceFilter) {
        qb = qb.or(`cost_input.is.null,cost_input.gte.${priceMin},cost_input.lte.${priceMax}`);
    }

    const from = (page - 1) * pageSize;
    const to = from + pageSize - 1;

    const { data, count, error } = await qb
        .order(sort, { ascending: order })
        .range(from, to);

    if (error) throw createError({ statusCode: 500, statusMessage: error.message });

    return {
        data: data || [],
        meta: {
            total: count || 0,
            page,
            pageSize,
            totalPages: Math.ceil((count || 0) / pageSize),
        },
    };
});
