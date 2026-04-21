const STORAGE_KEY = "sync-last-synced-at";

function formatRelativeTime(iso: string, t: (key: string, params?: Record<string, string>) => string): string {
    const diffMs = Date.now() - new Date(iso).getTime();
    const minutes = Math.floor(diffMs / 60000);
    if (minutes < 1) return t("footer.syncAgo", { time: "1m" });
    if (minutes < 60) return t("footer.syncAgo", { time: `${minutes}m` });
    const hours = Math.floor(minutes / 60);
    if (hours < 24) return t("footer.syncAgo", { time: `${hours}h` });
    const days = Math.floor(hours / 24);
    return t("footer.syncAgo", { time: `${days}d` });
}

export function useAutoSync() {
    const { t } = useI18n();
    const config = useRuntimeConfig();
    const syncing = useState<boolean>("sync-syncing", () => false);
    const lastSyncedAt = useState<string | null>("sync-last-synced", () => null);
    const syncMessage = useState<string | null>("sync-message", () => null);

    const formattedSyncTime = computed(() => {
        if (!lastSyncedAt.value) return null;
        return formatRelativeTime(lastSyncedAt.value, t);
    });

    const triggerSync = async () => {
        syncing.value = true;
        syncMessage.value = null;
        try {
            const result = await $fetch<{ data: { synced_at: string; skipped?: boolean } }>(`${config.public.apiBase}/api/v1/sync`, { method: "POST" });
            lastSyncedAt.value = result.data.synced_at;
            if (import.meta.client) {
                localStorage.setItem(STORAGE_KEY, result.data.synced_at);
            }
            if (result.data.skipped) {
                syncMessage.value = t("footer.syncSkipped");
            }
        } catch (e) {
            console.error("Sync failed", e);
        } finally {
            syncing.value = false;
        }
    };

    const initSync = () => {
        if (!import.meta.client) return;

        // Restore from localStorage
        const cached = localStorage.getItem(STORAGE_KEY);
        if (cached) lastSyncedAt.value = cached;

        // Fetch authoritative time from server
        $fetch<{ data: { synced_at: string | null } }>(`${config.public.apiBase}/api/v1/sync/status`)
            .then((res) => {
                if (res.data.synced_at) {
                    lastSyncedAt.value = res.data.synced_at;
                    localStorage.setItem(STORAGE_KEY, res.data.synced_at);
                }
            })
            .catch(() => {});
    };

    return { syncing, lastSyncedAt, formattedSyncTime, syncMessage, triggerSync, initSync };
}
