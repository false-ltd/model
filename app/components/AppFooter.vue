<template>
    <UFooter>
        <template #left>
            <div class="flex items-center gap-2 text-sm text-muted">
                <span class="w-1.5 h-1.5 bg-success rounded-full" />
            </div>
        </template>

        <template #right> </template>
    </UFooter>
</template>

<script setup lang="ts">
    const loading = ref(false);
    const lastSynced = ref("never");

    const refreshData = async () => {
        loading.value = true;
        try {
            const result = await $fetch<{ data: { synced_at: string } }>("/api/sync", { method: "POST" });
            lastSynced.value = new Date(result.data.synced_at).toLocaleString();
        } catch (e) {
            console.error("Sync failed", e);
        } finally {
            loading.value = false;
        }
    };
</script>
