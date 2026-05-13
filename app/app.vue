<script setup lang="ts">
import * as locales from '@nuxt/ui/locale'

const { locale } = useI18n()
const { t } = useI18n()

useHead({
    htmlAttrs: { lang: locale },
    script: [
        {
            type: 'application/ld+json',
            innerHTML: JSON.stringify({
                '@context': 'https://schema.org',
                '@type': 'WebSite',
                name: 'AI Model Catalog',
                alternateName: 'AI 模型目录',
                url: 'https://model.false.ltd',
                description: t('seo.indexDescription'),
                inLanguage: ['en', 'zh'],
                potentialAction: {
                    '@type': 'SearchAction',
                    target: 'https://model.false.ltd/catalog?q={search_term_string}',
                    'query-input': 'required name=search_term_string',
                },
            }),
        },
    ],
})
</script>

<template>
    <UApp :locale="locales[locale]">
        <NuxtLayout>
            <NuxtPage />
        </NuxtLayout>
    </UApp>
</template>
