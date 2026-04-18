<template>
    <div>
        <div class="flex items-center justify-between mb-5">
            <div>
                <h1 class="text-xl font-bold text-default">{{ t("providers.title") }}</h1>
                <div class="text-sm text-muted">{{ providersList.length }} {{ t("providers.total") }}</div>
            </div>
            <div class="flex items-center gap-2">
                <UInput
                    v-model="search"
                    :placeholder="t('providers.search')"
                    icon="i-lucide-search"
                    class="w-56"
                    size="sm"
                />
            </div>
        </div>

        <!-- Provider grid -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
            <NuxtLink
                v-for="p in filtered"
                :key="p.id"
                :to="localePath(`/catalog?providers=${p.id}`)"
                class="bg-default border border-default rounded-xl p-4 hover:border-accented transition-colors no-underline group"
            >
                <div class="flex items-center gap-3 mb-3">
                    <img
                        :src="`https://models.dev/logos/${p.id}.svg`"
                        class="w-9 h-9 rounded-lg p-1"
                        @error="($event.target as HTMLImageElement).style.display = 'none'"
                    />
                    <div class="min-w-0">
                        <div class="text-sm font-semibold text-default group-hover:text-primary transition-colors truncate">
                            {{ p.name }}
                        </div>
                        <div class="text-xs text-muted font-mono">{{ p.id }}</div>
                    </div>
                    <div class="ml-auto shrink-0">
                        <span class="text-xs font-semibold text-toned bg-elevated px-2 py-0.5 rounded">
                            {{ p.modelCount }} {{ t("providers.models") }}
                        </span>
                    </div>
                </div>
                <div class="grid grid-cols-2 gap-2 text-xs">
                    <div v-if="p.npm" class="truncate">
                        <span class="text-muted">NPM:</span>
                        <span class="text-toned font-mono ml-1">{{ p.npm }}</span>
                    </div>
                    <div v-if="p.api" class="truncate col-span-2">
                        <span class="text-muted">API:</span>
                        <span class="text-toned font-mono ml-1">{{ p.api }}</span>
                    </div>
                    <div v-if="p.env?.length" class="truncate">
                        <span class="text-muted">ENV:</span>
                        <span class="text-toned font-mono ml-1">{{ p.env[0] }}</span>
                    </div>
                    <div v-if="p.doc" class="truncate">
                        <a
                            :href="p.doc"
                            target="_blank"
                            class="text-primary hover:underline"
                            @click.stop
                        >
                            {{ t("providers.docs") }} ↗
                        </a>
                    </div>
                </div>
            </NuxtLink>
        </div>

        <!-- Empty state -->
        <div
            v-if="filtered.length === 0 && search"
            class="text-center py-16 text-sm text-muted"
        >
            {{ t("providers.noResults") }}
        </div>
    </div>
</template>

<script setup lang="ts">
    const localePath = useLocalePath();
    const { t } = useI18n();

    const search = ref("");

    const { data: result } = await useAsyncData("providers-list", () =>
        $fetch<{ data: any[] }>("/api/providers"),
    );

    const providersList = computed(() => result.value?.data || []);

    const filtered = computed(() => {
        if (!search.value.trim()) return providersList.value;
        const q = search.value.toLowerCase();
        return providersList.value.filter(
            (p: any) =>
                p.name.toLowerCase().includes(q) ||
                p.id.toLowerCase().includes(q) ||
                p.npm?.toLowerCase().includes(q),
        );
    });
</script>
