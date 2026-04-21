<template>
    <div>
        <CatalogFilters
            v-model:search="searchQuery"
            v-model:providerSearch="providerSearch"
            v-model:selectedInputTypes="selectedInputTypes"
            v-model:selectedOutputTypes="selectedOutputTypes"
            v-model:priceRange="priceRange"
            v-model:outputPriceRange="outputPriceRange"
            :total="catalogData.meta.total"
            :column-menu-items="columnMenuItems"
            :filter-toggles="filterToggles"
            :filters="filters"
            :selected-providers="selectedProviders"
            :top-providers="topProviders"
            :grouped-providers="groupedProviders"
            :filtered-providers="filteredProviders"
            :provider-trigger-label="providerTriggerLabel"
            :is-selected="isSelected"
            :input-type-items="inputTypeItems"
            :output-type-items="outputTypeItems"
            @toggle-filter="toggleFilter"
            @toggle-provider="toggleProvider"
            @remove-provider="removeProvider"
            @clear-providers="clearProviders"
        />

        <!-- Empty state -->
        <div v-if="models.length === 0 && !loading" class="flex flex-col items-center justify-center py-20 text-center">
            <div class="size-16 rounded-full bg-default border border-default flex items-center justify-center mb-4">
                <UIcon name="i-lucide-search-x" class="size-6 text-muted" />
            </div>
            <div class="text-base font-medium text-default mb-1">{{ t("common.noResults") }}</div>
            <div class="text-sm text-muted">{{ t("catalog.noResultsHint") }}</div>
        </div>

        <Transition name="catalog-fade">
            <CatalogTable
                v-if="models.length > 0"
                ref="catalogTable"
                :models="models"
                :loading="loading"
                :sort-field="sortField"
                :sort-order="sortOrder"
                :total-items="catalogData.meta.total"
                :total-pages="catalogData.meta.total_pages"
                :current-page="currentPage"
                :page-size="catalogData.meta.page_size"
                @sort="toggleSort"
                @go-to-page="goToPage"
                @change-page-size="changePageSize"
            />
        </Transition>

        <CompareFab />
    </div>
</template>

<script setup lang="ts">
    const {
        searchQuery,
        sortField,
        sortOrder,
        selectedProviders,
        selectedInputTypes,
        selectedOutputTypes,
        priceRange,
        outputPriceRange,
        filters,
        filterToggles,
        toggleFilter,
        catalogData,
        models,
        loading,
        fetchCatalog,
        toggleSort,
        goToPage,
        changePageSize,
        currentPage,
        inputTypeItems,
        outputTypeItems,
    } = useCatalog();

    const { t } = useI18n();

    const {
        providerSearch,
        topProviders,
        groupedProviders,
        filteredProviders,
        providerTriggerLabel,
        isSelected,
        toggleProvider,
        removeProvider,
        clearProviders,
    } = await useProviderFilter(selectedProviders);

    await fetchCatalog();

    const catalogTable = useTemplateRef<any>("catalogTable");
    const columnMenuItems = computed<any[]>(() => catalogTable.value?.columnMenuItems || []);
</script>

<style scoped>
.catalog-fade-enter-active,
.catalog-fade-leave-active {
    transition: opacity 0.2s ease;
}
.catalog-fade-enter-from,
.catalog-fade-leave-to {
    opacity: 0;
}
</style>
