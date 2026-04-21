<template>
    <header class="border-b border-default bg-default/80 backdrop-blur-md">
        <div class="w-full px-3 sm:px-4 md:px-6 h-14 flex items-center">
            <NuxtLink :to="localePath('/')" class="flex items-center">
                <div class="w-7 h-7 rounded-lg flex items-center justify-center text-xs font-bold">
                    <UIcon name="i-lucide-square-m" class="w-8 h-8" />
                </div>
                <span class="font-bold text-base text-default">odel</span>
            </NuxtLink>

            <!-- Desktop nav -->
            <div class="hidden md:flex flex-1 justify-between pl-8">
                <UNavigationMenu color="neutral" variant="link" :items="navItems" highlight />
            </div>

            <!-- Desktop actions -->
            <div class="hidden md:flex items-center gap-1">
                <UButton
                    icon="i-lucide-github"
                    color="neutral"
                    variant="ghost"
                    to="https://github.com/false-ltd/model"
                    target="_blank"
                    aria-label="GitHub"
                />
                <UColorModeButton />
                <UDropdownMenu :items="localeItems">
                    <button class="p-1.5 rounded-lg hover:bg-accented transition-colors flex items-center gap-1.5">
                        <UIcon name="i-lucide-languages" class="w-4 h-4 text-toned" />
                    </button>
                </UDropdownMenu>
            </div>

            <!-- Mobile spacer -->
            <div class="flex-1 md:hidden" />

            <!-- Mobile actions -->
            <div class="flex md:hidden items-center gap-1">
                <UColorModeButton />
                <UDropdownMenu :items="mobileMenuItems" :content="{ align: 'end' }">
                    <button class="p-1.5 rounded-lg hover:bg-accented transition-colors">
                        <UIcon name="i-lucide-menu" class="w-5 h-5 text-toned" />
                    </button>
                </UDropdownMenu>
            </div>
        </div>
    </header>
</template>

<script setup lang="ts">
    const { locale, locales, t } = useI18n();
    const switchLocalePath = useSwitchLocalePath();
    const localePath = useLocalePath();

    const navItems = computed(() => [
        { label: t("nav.overview"), to: localePath("/") },
        { label: t("nav.catalog"), to: localePath("/catalog") },
        { label: t("nav.providers"), to: localePath("/providers") },
        { label: t("nav.compare"), to: localePath("/compare") },
        { label: "False", icon: "i-lucide-terminal", to: "https://false.ltd", target: "_blank" },
    ]);

    const localeItems = computed(() =>
        (locales.value as { code: string; name: string }[]).map((l) => ({
            label: l.name,
            icon: locale.value === l.code ? "i-lucide-check" : undefined,
            to: switchLocalePath(l.code as "en" | "zh"),
        })),
    );

    const mobileMenuItems = computed(() => [
        ...navItems.value.map((item) => ({
            label: item.label,
            icon: (item as any).icon,
            to: item.to,
            target: (item as any).target,
        })),
        { type: "separator" as const },
        ...(locales.value as { code: string; name: string }[]).map((l) => ({
            label: l.name,
            icon: locale.value === l.code ? "i-lucide-check" : "i-lucide-languages",
            to: switchLocalePath(l.code as "en" | "zh"),
        })),
        { type: "separator" as const },
        {
            label: t("nav.github"),
            icon: "i-lucide-github",
            to: "https://github.com/false-ltd/model",
            target: "_blank",
        },
    ]);
</script>
