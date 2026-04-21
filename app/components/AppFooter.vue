<template>
    <UFooter>
        <template #left>
            <div class="flex items-center gap-1.5 text-xs text-muted">
                <span>{{ t("footer.source") }}</span>
                <NuxtLink
                    to="https://models.dev"
                    target="_blank"
                    class="text-toned hover:text-primary transition-colors no-underline"
                >
                    models.dev
                    <UIcon name="i-lucide-external-link" class="size-2.5 align-[-1px] ml-0.5" />
                </NuxtLink>
            </div>
        </template>

        <template #right>
            <div class="flex items-center gap-2.5">
                <button
                    @click="triggerSync()"
                    :disabled="syncing"
                    class="inline-flex items-center gap-1.5 text-xs text-toned hover:text-primary cursor-pointer transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                    <UIcon
                        :name="syncing ? 'i-lucide-loader-2' : 'i-lucide-refresh-cw'"
                        class="size-3"
                        :class="syncing ? 'animate-spin' : ''"
                    />
                    {{ syncing ? t("footer.syncing") : t("footer.refresh") }}
                </button>

                <div class="flex items-center gap-1.5 text-xs text-muted">
                    <span class="size-1.5 rounded-full bg-success animate-pulse shrink-0" />
                    <span v-if="syncMessage" class="hidden sm:inline">{{ syncMessage }}</span>
                    <span v-else class="hidden sm:inline">
                        {{ formattedSyncTime ? `${t("footer.lastSync")} ${formattedSyncTime}` : t("footer.autoSync") }}
                    </span>
                </div>
            </div>
        </template>
    </UFooter>
</template>

<script setup lang="ts">
    const { t } = useI18n();
    const { syncing, formattedSyncTime, syncMessage, triggerSync } = useAutoSync();
</script>
