export interface Provider {
    id: string;
    name: string;
    npm?: string;
    env?: string[];
    doc_url?: string;
    api_url?: string;
    api?: string;
    modelCount: number;
}

export interface Model {
    id: number;
    name: string;
    model_id: string;
    family: string | null;
    provider_id: string;
    providers: Provider | null;
    reasoning: boolean;
    tool_call: boolean;
    attachment: boolean;
    open_weights: boolean;
    structured_output: boolean;
    temperature: boolean;
    knowledge: string | null;
    release_date: string | null;
    last_updated: string | null;
    status: "alpha" | "beta" | "deprecated" | null;
    interleaved: unknown;
    cost_input: number | null;
    cost_output: number | null;
    cost_reasoning: number | null;
    cost_cache_read: number | null;
    cost_cache_write: number | null;
    cost_input_audio: number | null;
    cost_output_audio: number | null;
    limit_context: number | null;
    limit_input: number | null;
    limit_output: number | null;
    modalities_input: string[];
    modalities_output: string[];
}

export interface CatalogMeta {
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
}

export interface ApiResponse<T> {
    code: number;
    data: T;
    message: string;
    meta?: CatalogMeta;
}

export type SelectItem = { label: string; value: string };

export interface PricingField {
    key: string;
    label: string;
    value: number | null;
    display: string;
}
