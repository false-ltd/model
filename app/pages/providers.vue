<template>
    <div>
        <!-- Header -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-5">
            <div>
                <h1 class="text-xl font-bold text-default">{{ t("providers.title") }}</h1>
                <div class="text-sm text-muted">{{ filtered.length }} {{ t("providers.total") }}</div>
            </div>
            <div class="flex items-center gap-2">
                <UInput
                    v-model="search"
                    :placeholder="t('providers.search')"
                    icon="i-lucide-search"
                    class="w-full sm:w-56"
                    size="sm"
                />
                <!-- View toggle (desktop) -->
                <div class="hidden sm:flex bg-default border border-default rounded-lg p-0.5 shrink-0">
                    <button
                        @click="viewMode = 'grid'"
                        class="rounded-md px-2 py-1 cursor-pointer transition-colors"
                        :class="viewMode === 'grid' ? 'bg-primary text-white' : 'text-toned hover:text-default'"
                    >
                        <UIcon name="i-lucide-layout-grid" class="size-3.5" />
                    </button>
                    <button
                        @click="viewMode = 'table'"
                        class="rounded-md px-2 py-1 cursor-pointer transition-colors"
                        :class="viewMode === 'table' ? 'bg-primary text-white' : 'text-toned hover:text-default'"
                    >
                        <UIcon name="i-lucide-table" class="size-3.5" />
                    </button>
                </div>
            </div>
        </div>

        <!-- Grid view -->
        <div v-if="viewMode === 'grid'" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
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

        <!-- Table view -->
        <div v-else class="border border-default rounded-xl overflow-hidden overflow-x-auto">
            <table class="w-full text-sm">
                <thead>
                    <tr class="border-b border-default bg-elevated/50">
                        <th class="text-left px-4 py-2.5 text-xs font-medium text-muted uppercase tracking-wider">{{ t("providers.title") }}</th>
                        <th class="text-left px-4 py-2.5 text-xs font-medium text-muted uppercase tracking-wider">ID</th>
                        <th class="text-center px-4 py-2.5 text-xs font-medium text-muted uppercase tracking-wider">{{ t("providers.models") }}</th>
                        <th class="hidden md:table-cell text-left px-4 py-2.5 text-xs font-medium text-muted uppercase tracking-wider">NPM</th>
                        <th class="hidden lg:table-cell text-left px-4 py-2.5 text-xs font-medium text-muted uppercase tracking-wider">API</th>
                        <th class="text-center px-4 py-2.5 text-xs font-medium text-muted uppercase tracking-wider">{{ t("providers.docs") }}</th>
                    </tr>
                </thead>
                <tbody>
                    <NuxtLink
                        v-for="p in filtered"
                        :key="p.id"
                        :to="localePath(`/catalog?providers=${p.id}`)"
                        custom
                    >
                        <template #default="{ navigate }">
                            <tr
                                class="border-b border-default last:border-b-0 hover:bg-elevated/50 transition-colors cursor-pointer"
                                @click="navigate"
                            >
                                <td class="px-4 py-3">
                                    <div class="flex items-center gap-2">
                                        <img
                                            :src="`https://models.dev/logos/${p.id}.svg`"
                                            class="w-7 h-7 rounded-lg p-0.5 shrink-0"
                                            @error="($event.target as HTMLImageElement).style.display = 'none'"
                                        />
                                        <span class="font-medium text-default truncate">{{ p.name }}</span>
                                    </div>
                                </td>
                                <td class="px-4 py-3 text-toned font-mono text-xs">{{ p.id }}</td>
                                <td class="px-4 py-3 text-center">
                                    <span class="text-xs font-semibold text-toned bg-elevated px-2 py-0.5 rounded">
                                        {{ p.modelCount }}
                                    </span>
                                </td>
                                <td class="hidden md:table-cell px-4 py-3 text-toned font-mono text-xs truncate max-w-40">{{ p.npm || "—" }}</td>
                                <td class="hidden lg:table-cell px-4 py-3 text-toned font-mono text-xs truncate max-w-60">{{ p.api || "—" }}</td>
                                <td class="px-4 py-3 text-center">
                                    <a
                                        v-if="p.doc"
                                        :href="p.doc"
                                        target="_blank"
                                        class="text-primary hover:underline text-xs"
                                        @click.stop
                                    >
                                        {{ t("providers.docs") }} ↗
                                    </a>
                                    <span v-else class="text-muted">—</span>
                                </td>
                            </tr>
                        </template>
                    </NuxtLink>
                </tbody>
            </table>
        </div>

        <!-- Empty state -->
        <div v-if="filtered.length === 0 && search" class="text-center py-16 text-sm text-muted">
            {{ t("providers.noResults") }}
        </div>
    </div>
</template>

<script setup lang="ts">
    const localePath = useLocalePath();
    const { t } = useI18n();
    const { viewMode, search, filtered } = useProviders();
</script>
