import { serverSupabaseClient } from "#supabase/server";

export default defineEventHandler(async (event) => {
    const client = await serverSupabaseClient(event);

    const { data: providers, error: pErr } = await client
        .from("providers")
        .select("provider_id, name, npm, env, doc, api, models(count)")
        .order("name");

    if (pErr) throw createError({ statusCode: 500, statusMessage: pErr.message });

    const data = (providers || []).map((p: any) => ({
        id: p.provider_id,
        name: p.name,
        npm: p.npm,
        env: p.env || [],
        doc: p.doc,
        api: p.api,
        modelCount: p.models?.[0]?.count || 0,
    }));

    return { data, meta: { count: data.length } };
});
