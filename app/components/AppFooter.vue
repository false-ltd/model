<template>
    <UFooter>
        <template #left>
            <div class="flex items-center gap-2 text-sm text-muted"></div>
        </template>

        <template #right>
            <UButton
                icon="i-lucide-cloud-sync"
                :label="syncing ? $t('footer.syncing') : $t('footer.refresh')"
                size="md"
                :loading="syncing"
                @click="triggerSync()"
            />
            <span class="size-2 animate-pulse bg-green-500 w-1.5 h-1.5 rounded-full" />

            <span v-if="syncMessage" class="text-xs text-muted hidden sm:inline">
                {{ syncMessage }}
            </span>
            <span v-else class="text-xs text-muted hidden sm:inline">
                {{ formattedSyncTime ? `${$t("footer.lastSync")}: ${formattedSyncTime}` : $t("footer.autoSync") }}
            </span>
        </template>
    </UFooter>
</template>

<script setup lang="ts">
    const { syncing, formattedSyncTime, syncMessage, triggerSync } = useAutoSync();
</script>
