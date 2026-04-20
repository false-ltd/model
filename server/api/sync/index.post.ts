import { serverSupabaseClient } from "#supabase/server";

export default defineEventHandler(async (event) => {
    const client = await serverSupabaseClient(event);

    // Check if synced within the last 10 minutes
    const { data: recent } = await client
        .from("providers")
        .select("synced_at")
        .not("synced_at", "is", null)
        .order("synced_at", { ascending: false })
        .limit(1)
        .single();

    if (recent?.synced_at) {
        const elapsed = Date.now() - new Date(recent.synced_at).getTime();
        if (elapsed < 10 * 60 * 1000) {
            return { data: { synced_at: recent.synced_at, skipped: true } };
        }
    }

    const result = await performSync();
    return { data: result };
});
