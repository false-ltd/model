import { serverSupabaseClient } from "#supabase/server";

export default defineEventHandler(async (event) => {
    const client = await serverSupabaseClient(event);

    const { data, error } = await client
        .from("providers")
        .select("synced_at")
        .not("synced_at", "is", null)
        .order("synced_at", { ascending: false })
        .limit(1)
        .single();

    if (error) {
        return { data: { synced_at: null } };
    }

    return { data: { synced_at: data?.synced_at ?? null } };
});
