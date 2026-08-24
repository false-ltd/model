<script setup lang="ts">
import { en, zh_cn } from '@nuxt/ui/locale'

const { locale } = useI18n()
const { t } = useI18n()

const uiLocales: Record<string, typeof en> = { en, zh: zh_cn }

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
    <UApp :locale="uiLocales[locale] ?? en">
        <NuxtLoadingIndicator />
        <IntroOverlay />
        <NuxtLayout>
            <NuxtPage />
        </NuxtLayout>
    </UApp>
</template>
