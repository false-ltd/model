import { serverSupabaseClient } from "#supabase/server";

export default defineEventHandler(async (event) => {
    const query = getQuery(event);
    const client = await serverSupabaseClient(event);

    const ids = (query.ids as string || "")
        .split(",")
        .filter(Boolean)
        .map(Number)
        .filter((n) => !isNaN(n));

    if (!ids.length) return { data: [], meta: { count: 0 } };

    const { data, error } = await client
        .from("models")
        .select("*, providers(provider_id, name, npm, env, doc, api)")
        .in("id", ids);

    if (error) throw createError({ statusCode: 500, statusMessage: error.message });

    // Preserve the order from the query
    const orderMap = new Map(ids.map((id, i) => [id, i]));
    const sorted = (data || []).sort((a: any, b: any) =>
        (orderMap.get(a.id) ?? 0) - (orderMap.get(b.id) ?? 0)
    );

    return { data: sorted, meta: { count: sorted.length } };
});
