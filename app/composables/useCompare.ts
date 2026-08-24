export interface CompareItem {
    id: number;
    name: string;
    provider_id: string;
}

const MAX_COMPARE = 4;

/**
 * Compare selection state. Stores a {id, name, provider_id} snapshot when a
 * model is added so the FAB list can render without an extra /compare
 * request; ids persist in localStorage and ?models= for sharing.
 */
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

    const readStorage = (): CompareItem[] => {
        if (!import.meta.client) return [];
        try {
            const raw = localStorage.getItem("compare-models");
            if (raw) return JSON.parse(raw);
        } catch {}
        return [];
    };

    const writeStorage = (items: CompareItem[]) => {
        if (!import.meta.client) return;
        localStorage.setItem("compare-models", JSON.stringify(items));
    };

    const syncToUrl = (items: CompareItem[]) => {
        const query = { ...route.query };
        if (items.length) query.models = items.map((m) => m.id).join(",");
        else delete query.models;
        router.replace({ query });
    };

    const idSet = (items: CompareItem[]) => new Set(items.map((m) => m.id));

    // Init: URL params take priority (for sharing), then fall back to localStorage.
    // Models restored from a bare id list carry no name yet; the FAB hydrates
    // them lazily from the /compare endpoint.
    const storage = useState<CompareItem[]>("compare-models", () => {
        const urlIds = parseIds(route.query.models as string);
        if (urlIds.length) {
            const stored = readStorage();
            const storedById = new Map(stored.map((m) => [m.id, m]));
            const items = urlIds.map((id) => storedById.get(id) ?? { id, name: "", provider_id: "" });
            writeStorage(items);
            return items;
        }
        return readStorage();
    });

    // On mount, sync URL with current state (e.g. restored from localStorage)
    onMounted(() => {
        if (storage.value.length && !route.query.models) {
            syncToUrl(storage.value);
        }
    });

    const modelIds = computed(() => storage.value.map((m) => m.id));
    const selectedModels = computed(() => storage.value);

    const addModel = (item: CompareItem): { added: boolean; reason?: "max" | "duplicate" } => {
        if (idSet(storage.value).has(item.id)) return { added: false, reason: "duplicate" };
        if (storage.value.length >= MAX_COMPARE) return { added: false, reason: "max" };
        // Keep the richest snapshot: preserve a known name when the caller
        // only knows the id (e.g. compare page remove/re-add).
        const existing = storage.value.find((m) => m.id === item.id);
        storage.value = [...storage.value, existing ?? item];
        writeStorage(storage.value);
        syncToUrl(storage.value);
        return { added: true };
    };

    const removeModel = (id: number) => {
        storage.value = storage.value.filter((m: CompareItem) => m.id !== id);
        writeStorage(storage.value);
        syncToUrl(storage.value);
    };

    const clearAll = () => {
        storage.value = [];
        writeStorage([]);
        syncToUrl([]);
    };

    // Drag-reorder columns on the compare page; persists through the same
    // channels as add/remove (localStorage + URL).
    const reorderModels = (from: number, to: number) => {
        if (from === to || from < 0 || to < 0 || from >= storage.value.length || to >= storage.value.length) return;
        const arr = [...storage.value];
        const [item] = arr.splice(from, 1);
        arr.splice(to, 0, item);
        storage.value = arr;
        writeStorage(arr);
        syncToUrl(arr);
    };

    const patchModel = (item: CompareItem) => {
        const idx = storage.value.findIndex((m) => m.id === item.id);
        if (idx === -1 || storage.value[idx].name) return;
        storage.value = storage.value.map((m, i) => (i === idx ? item : m));
        writeStorage(storage.value);
    };

    return { modelIds, selectedModels, addModel, removeModel, clearAll, reorderModels, patchModel };
};
