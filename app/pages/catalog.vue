<template>
    <div>
        <CatalogFilters
            v-model="viewMode"
            v-model:search="searchQuery"
            v-model:providerSearch="providerSearch"
            v-model:selectedInputTypes="selectedInputTypes"
            v-model:selectedOutputTypes="selectedOutputTypes"
            v-model:priceRange="priceRange"
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

        <!-- Grid view -->
        <div v-if="viewMode === 'grid'" class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-4">
            <ModelCard v-for="model in models" :key="model.id" :model="model" />
        </div>

        <!-- Table view -->
        <CatalogTable
            v-else
            ref="catalogTable"
            :models="models"
            :sort-field="sortField"
            :sort-order="sortOrder"
            :total-items="catalogData.meta.total"
            :total-pages="catalogData.meta.totalPages"
            :current-page="currentPage"
            :page-size="catalogData.meta.pageSize"
            @sort="toggleSort"
            @go-to-page="goToPage"
        />

        <CatalogCompareBar />
    </div>
</template>

<script setup lang="ts">
    const {
        viewMode,
        searchQuery,
        sortField,
        sortOrder,
        selectedProviders,
        selectedInputTypes,
        selectedOutputTypes,
        priceRange,
        filters,
        filterToggles,
        toggleFilter,
        catalogData,
        models,
        fetchCatalog,
        toggleSort,
        goToPage,
        currentPage,
        inputTypeItems,
        outputTypeItems,
    } = useCatalog();

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
    const columnMenuItems = computed<any>(() => catalogTable.value?.columnMenuItems || []);
</script>
