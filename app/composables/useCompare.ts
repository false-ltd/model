export const useCompare = () => {
    const route = useRoute();
    const router = useRouter();

    const parseIds = (raw: string | undefined): number[] =>
        raw
            ? raw
                  .split(",")
                  .filter(Boolean)
                  .map(Number)
                  .filter((n) => !isNaN(n))
            : [];

    const readStorage = (): number[] => {
        if (!import.meta.client) return [];
        try {
            const raw = localStorage.getItem("compare-models");
            if (raw) return JSON.parse(raw);
        } catch {}
        return [];
    };

    const writeStorage = (ids: number[]) => {
        if (!import.meta.client) return;
        localStorage.setItem("compare-models", JSON.stringify(ids));
    };

    const syncToUrl = (ids: number[]) => {
        const query = { ...route.query };
        if (ids.length) query.models = ids.join(",");
        else delete query.models;
        router.replace({ query });
    };

    // Init: URL params take priority (for sharing), then fall back to localStorage
    const storage = useState<number[]>("compare-models", () => {
        const urlIds = parseIds(route.query.models as string);
        if (urlIds.length) {
            writeStorage(urlIds);
            return urlIds;
        }
        return readStorage();
    });

    // On mount, sync URL with current state (e.g. restored from localStorage)
    onMounted(() => {
        if (storage.value.length && !route.query.models) {
            syncToUrl(storage.value);
        }
    });

    const modelIds = computed(() => storage.value);

    const addModel = (id: number): { added: boolean; reason?: 'max' | 'duplicate' } => {
        if (storage.value.includes(id)) return { added: false, reason: 'duplicate' };
        if (storage.value.length >= 4) return { added: false, reason: 'max' };
        storage.value = [...storage.value, id];
        writeStorage(storage.value);
        syncToUrl(storage.value);
        return { added: true };
    };

    const removeModel = (id: number) => {
        storage.value = storage.value.filter((m: number) => m !== id);
        writeStorage(storage.value);
        syncToUrl(storage.value);
    };

    const clearAll = () => {
        storage.value = [];
        writeStorage([]);
        syncToUrl([]);
    };

    return { modelIds, addModel, removeModel, clearAll };
};
