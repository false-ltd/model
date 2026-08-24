/**
 * Single owner of the atlas point cloud. Every consumer (homepage hero,
 * interactive atlas, providers enrichment) must go through this composable
 * so concurrent callers share one request instead of racing three fetches.
 */
export function useAtlasData() {
    const config = useRuntimeConfig();
    return useAsyncData("atlas-data", () =>
        $fetch<any>(`${config.public.apiBase}/api/v1/atlas`),
    { dedupe: "defer" });
}
