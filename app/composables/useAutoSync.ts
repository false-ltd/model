const SYNC_INTERVAL_MS = 2 * 60 * 60 * 1000;
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
    const syncing = useState<boolean>("sync-syncing", () => false);
    const lastSyncedAt = useState<string | null>("sync-last-synced", () => null);
    const syncMessage = useState<string | null>("sync-message", () => null);

    let timerId: ReturnType<typeof setInterval> | null = null;
    let initialized = false;

    const formattedSyncTime = computed(() => {
        if (!lastSyncedAt.value) return null;
        return formatRelativeTime(lastSyncedAt.value, t);
    });

    const triggerSync = async () => {
        syncing.value = true;
        syncMessage.value = null;
        try {
            const result = await $fetch<{ data: { synced_at: string; skipped?: boolean } }>("/api/sync", { method: "POST" });
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

    const initAutoSync = () => {
        if (!import.meta.client || initialized) return;
        initialized = true;

        // Restore from localStorage to avoid flash
        const cached = localStorage.getItem(STORAGE_KEY);
        if (cached) lastSyncedAt.value = cached;

        // Fetch authoritative time from server
        $fetch<{ data: { synced_at: string | null } }>("/api/sync/status")
            .then((res) => {
                if (res.data.synced_at) {
                    lastSyncedAt.value = res.data.synced_at;
                    localStorage.setItem(STORAGE_KEY, res.data.synced_at);
                }
            })
            .catch(() => {});

        // Start periodic sync
        timerId = setInterval(() => {
            triggerSync();
        }, SYNC_INTERVAL_MS);
    };

    const stopAutoSync = () => {
        if (timerId) {
            clearInterval(timerId);
            timerId = null;
        }
    };

    return { syncing, lastSyncedAt, formattedSyncTime, syncMessage, triggerSync, initAutoSync, stopAutoSync };
}
