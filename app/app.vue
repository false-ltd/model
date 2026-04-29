<script setup lang="ts">
import * as locales from '@nuxt/ui/locale'

const { locale } = useI18n()

useHead({
    htmlAttrs: { lang: locale },
})
</script>

<template>
    <UApp :locale="locales[locale]">
        <NuxtLayout>
            <NuxtPage />
        </NuxtLayout>
    </UApp>
</template>
