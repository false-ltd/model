import { serverSupabaseClient } from "#supabase/server";

export default defineEventHandler(async (event) => {
    const id = Number(getRouterParam(event, "id"));
    if (!id) throw createError({ statusCode: 400, statusMessage: "Invalid model ID" });

    const client = await serverSupabaseClient(event);

    const { data, error } = await client
        .from("models")
        .select("*, providers(provider_id, name, npm, env, doc, api)")
        .eq("id", id)
        .single();

    if (error || !data) throw createError({ statusCode: 404, statusMessage: "Model not found" });

    return { data };
});
