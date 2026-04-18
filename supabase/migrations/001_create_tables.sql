-- Providers table: integer PK for internal joins, provider_id as the text slug for logos/URLs
CREATE TABLE IF NOT EXISTS providers (
    id BIGSERIAL PRIMARY KEY,
    provider_id TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    npm TEXT NOT NULL,
    env TEXT[] NOT NULL DEFAULT '{}',
    doc TEXT,
    api TEXT,
    synced_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE providers IS 'AI model providers (e.g. OpenAI, Anthropic)';
COMMENT ON COLUMN providers.id IS 'Auto-incrementing integer primary key';
COMMENT ON COLUMN providers.provider_id IS 'Provider identifier slug (e.g. "openai", "anthropic"), used in logo URLs and API lookups';
COMMENT ON COLUMN providers.name IS 'Display name of the provider';
COMMENT ON COLUMN providers.npm IS 'AI SDK NPM package name (e.g. "@ai-sdk/openai")';
COMMENT ON COLUMN providers.env IS 'Environment variable keys used for authentication (e.g. ["OPENAI_API_KEY"])';
COMMENT ON COLUMN providers.doc IS 'Link to the provider''s documentation';
COMMENT ON COLUMN providers.api IS 'OpenAI-compatible API endpoint, required when using @ai-sdk/openai-compatible as the npm package';
COMMENT ON COLUMN providers.synced_at IS 'Timestamp of last data sync from models.dev';

-- Models table: integer PK for clean URLs, model_id as the unique string identifier from models.dev
CREATE TABLE IF NOT EXISTS models (
    id BIGSERIAL PRIMARY KEY,
    model_id TEXT NOT NULL UNIQUE,
    provider_id TEXT NOT NULL REFERENCES providers(provider_id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    family TEXT,
    attachment BOOLEAN NOT NULL DEFAULT false,
    reasoning BOOLEAN NOT NULL DEFAULT false,
    tool_call BOOLEAN NOT NULL DEFAULT false,
    structured_output BOOLEAN,
    temperature BOOLEAN DEFAULT true,
    knowledge TEXT,
    release_date TEXT,
    last_updated TEXT,
    open_weights BOOLEAN NOT NULL DEFAULT false,
    status TEXT CHECK (status IN ('alpha', 'beta', 'deprecated')),
    interleaved JSONB,
    modalities_input TEXT[] NOT NULL DEFAULT '{}',
    modalities_output TEXT[] NOT NULL DEFAULT '{}',
    cost_input DECIMAL(10, 4),
    cost_output DECIMAL(10, 4),
    cost_reasoning DECIMAL(10, 4),
    cost_cache_read DECIMAL(10, 4),
    cost_cache_write DECIMAL(10, 4),
    cost_input_audio DECIMAL(10, 4),
    cost_output_audio DECIMAL(10, 4),
    limit_context INTEGER,
    limit_input INTEGER,
    limit_output INTEGER
);

COMMENT ON TABLE models IS 'AI model specifications, pricing, and capabilities';
COMMENT ON COLUMN models.id IS 'Auto-incrementing integer primary key for clean URLs (/model/1, /model/2)';
COMMENT ON COLUMN models.model_id IS 'Unique model identifier used by AI SDK (e.g. "gpt-4o", "claude-3.5-sonnet")';
COMMENT ON COLUMN models.provider_id IS 'Foreign key to providers table';
COMMENT ON COLUMN models.name IS 'Display name of the model';
COMMENT ON COLUMN models.family IS 'Model family grouping (e.g. "gpt-4o", "claude-3.5")';
COMMENT ON COLUMN models.attachment IS 'Whether the model supports file attachments';
COMMENT ON COLUMN models.reasoning IS 'Whether the model supports reasoning / chain-of-thought';
COMMENT ON COLUMN models.tool_call IS 'Whether the model supports tool / function calling';
COMMENT ON COLUMN models.structured_output IS 'Whether the model supports structured output feature';
COMMENT ON COLUMN models.temperature IS 'Whether the model supports temperature control';
COMMENT ON COLUMN models.knowledge IS 'Knowledge cutoff date in YYYY-MM or YYYY-MM-DD format';
COMMENT ON COLUMN models.release_date IS 'First public release date in YYYY-MM or YYYY-MM-DD format';
COMMENT ON COLUMN models.last_updated IS 'Most recent update date in YYYY-MM or YYYY-MM-DD format';
COMMENT ON COLUMN models.open_weights IS 'Whether the model''s trained weights are publicly available';
COMMENT ON COLUMN models.status IS 'Model lifecycle status: alpha (testing), beta, or deprecated (no longer served)';
COMMENT ON COLUMN models.interleaved IS 'Interleaved reasoning support: true for general, or {"field": "reasoning_content"} for specific format';
COMMENT ON COLUMN models.modalities_input IS 'Supported input modalities (e.g. ["text", "image", "audio", "video", "pdf"])';
COMMENT ON COLUMN models.modalities_output IS 'Supported output modalities (e.g. ["text"])';
COMMENT ON COLUMN models.cost_input IS 'Cost per million input tokens in USD';
COMMENT ON COLUMN models.cost_output IS 'Cost per million output tokens in USD';
COMMENT ON COLUMN models.cost_reasoning IS 'Cost per million reasoning tokens in USD';
COMMENT ON COLUMN models.cost_cache_read IS 'Cost per million cached read tokens in USD';
COMMENT ON COLUMN models.cost_cache_write IS 'Cost per million cached write tokens in USD';
COMMENT ON COLUMN models.cost_input_audio IS 'Cost per million audio input tokens in USD (if billed separately)';
COMMENT ON COLUMN models.cost_output_audio IS 'Cost per million audio output tokens in USD (if billed separately)';
COMMENT ON COLUMN models.limit_context IS 'Maximum context window in tokens';
COMMENT ON COLUMN models.limit_input IS 'Maximum input tokens';
COMMENT ON COLUMN models.limit_output IS 'Maximum output tokens';

-- Indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_models_provider_id ON models(provider_id);
CREATE INDEX IF NOT EXISTS idx_models_family ON models(family);
CREATE INDEX IF NOT EXISTS idx_models_cost_input ON models(cost_input);
CREATE INDEX IF NOT EXISTS idx_models_limit_context ON models(limit_context);
CREATE INDEX IF NOT EXISTS idx_models_reasoning ON models(reasoning) WHERE reasoning = true;
CREATE INDEX IF NOT EXISTS idx_models_open_weights ON models(open_weights) WHERE open_weights = true;
CREATE INDEX IF NOT EXISTS idx_models_status ON models(status) WHERE status IS NOT NULL;
