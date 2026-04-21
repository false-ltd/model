<template>
    <div class="bg-default border border-default rounded-xl p-3 sm:p-5 hover:border-primary/40 transition-all">
        <div class="flex justify-between items-center mb-2.5">
            <div class="text-xs text-muted uppercase tracking-wider font-medium">{{ t("detail.quickStart") }}</div>
            <CopyButton :value="code" />
        </div>
        <div class="flex gap-1 mb-2.5">
            <button
                v-for="tab in tabs"
                :key="tab.key"
                @click="activeSdk = tab.key"
                class="px-2.5 py-1 text-xs rounded-md cursor-pointer transition-colors"
                :class="activeSdk === tab.key ? 'bg-primary text-white' : 'bg-elevated text-toned hover:text-default'"
            >
                {{ tab.label }}
            </button>
        </div>
        <div class="bg-elevated rounded-lg p-4 font-mono text-sm text-default leading-relaxed overflow-x-auto">
            <pre>{{ code }}</pre>
        </div>
    </div>
</template>

<script setup lang="ts">
    const { t } = useI18n();
    const props = defineProps<{ model: any }>();

    const activeSdk = ref("ai");
    const tabs = computed(() => [
        { key: "ai", label: t("detail.sdkAi") },
        { key: "curl", label: t("detail.sdkCurl") },
        { key: "python", label: t("detail.sdkPython") },
    ]);

    const code = computed(() => {
        const m = props.model;
        if (!m) return "";
        const npm = m.providers?.npm || "";
        const importName = npm.replace("@ai-sdk/", "");
        const modelId = m.model_id;
        const apiBase = m.providers?.api_url || "";
        switch (activeSdk.value) {
            case "curl":
                return `curl ${apiBase}/chat/completions \\
  -H "Content-Type: application/json" \\
  -H "Authorization: Bearer $API_KEY" \\
  -d '{
    "model": "${modelId}",
    "messages": [{"role": "user", "content": "Hello"}]
  }'`;
            case "python":
                return `from openai import OpenAI

client = OpenAI(
    base_url="${apiBase}/v1",
    api_key="$API_KEY",
)

response = client.chat.completions.create(
    model="${modelId}",
    messages=[{"role": "user", "content": "Hello"}],
)`;
            default:
                return `import { ${importName} } from '${npm}';\n\nconst model = ${importName}('${modelId}');`;
        }
    });
</script>
