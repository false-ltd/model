/**
 * Shared star-map quality state. Low-spec devices (≤4GB reported RAM or
 * ≤4 logical cores) start in lite mode — smaller glows and dpr 1 — and the
 * atlas control lets anyone toggle it. Both canvas instances (hero
 * backdrop + interactive plot) react to the same state.
 */
export function useAtlasLite() {
    return useState<boolean>("atlas-lite", () => {
        if (typeof navigator === "undefined") return false;
        const nav = navigator as { deviceMemory?: number; hardwareConcurrency?: number };
        return (typeof nav.deviceMemory === "number" && nav.deviceMemory <= 4)
            || (typeof nav.hardwareConcurrency === "number" && nav.hardwareConcurrency <= 4);
    });
}
