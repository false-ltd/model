# Model.dev some useful Schema

Models.dev is a comprehensive open-source database of AI model specifications, pricing, and features.

There's no single database with information about all the available AI models. We started Models.dev as a community-contributed project to address this. We also use it internally in opencode.

**API**
You can access this data through an API.

`curl https://models.dev/api.json`

Use the Model ID field to do a lookup on any model; it's the identifier used by AI SDK.

**Logos**
Provider logos are available at /logos/{provider}.svg where {provider} is the Provider ID.

`curl https://models.dev/logos/anthropic.svg`

If we don't have a provider's logo, a default logo is served instead.

**Provider Schema:**

- `name`: String - Display name of the provider
- `npm`: String - AI SDK Package name
- `env`: String[] - Environment variable keys used for auth
- `doc`: String - Link to the provider's documentation
- `api` _(optional)_: String - OpenAI-compatible API endpoint. Required only when using `@ai-sdk/openai-compatible` as the npm package

**Model Schema:**

- `name`: String — Display name of the model
- `attachment`: Boolean — Supports file attachments
- `reasoning`: Boolean — Supports reasoning / chain-of-thought
- `tool_call`: Boolean - Supports tool calling
- `structured_output` _(optional)_: Boolean — Supports structured output feature
- `temperature` _(optional)_: Boolean — Supports temperature control
- `knowledge` _(optional)_: String — Knowledge-cutoff date in `YYYY-MM` or `YYYY-MM-DD` format
- `release_date`: String — First public release date in `YYYY-MM` or `YYYY-MM-DD`
- `last_updated`: String — Most recent update date in `YYYY-MM` or `YYYY-MM-DD`
- `open_weights`: Boolean - Indicate the model's trained weights are publicly available
- `interleaved` _(optional)_: Boolean or Object — Supports interleaved reasoning. Use `true` for general support or an object with `field` to specify the format
- `interleaved.field`: String — Name of the interleaved field (`"reasoning_content"` or `"reasoning_details"`)
- `cost.input`: Number — Cost per million input tokens (USD)
- `cost.output`: Number — Cost per million output tokens (USD)
- `cost.reasoning` _(optional)_: Number — Cost per million reasoning tokens (USD)
- `cost.cache_read` _(optional)_: Number — Cost per million cached read tokens (USD)
- `cost.cache_write` _(optional)_: Number — Cost per million cached write tokens (USD)
- `cost.input_audio` _(optional)_: Number — Cost per million audio input tokens, if billed separately (USD)
- `cost.output_audio` _(optional)_: Number — Cost per million audio output tokens, if billed separately (USD)
- `limit.context`: Number — Maximum context window (tokens)
- `limit.input`: Number — Maximum input tokens
- `limit.output`: Number — Maximum output tokens
- `modalities.input`: Array of strings — Supported input modalities (e.g., ["text", "image", "audio", "video", "pdf"])
- `modalities.output`: Array of strings — Supported output modalities (e.g., ["text"])
- `status` _(optional)_: String — Supported status:
    - `alpha` - Indicate the model is in alpha testing
    - `beta` - Indicate the model is in beta testing
    - `deprecated` - Indicate the model is no longer served by the provider's public API

**Example:**

```json
{
    "minimax-coding-plan": {
        "id": "minimax-coding-plan",
        "env": ["MINIMAX_API_KEY"],
        "npm": "@ai-sdk/anthropic",
        "api": "https://api.minimax.io/anthropic/v1",
        "name": "MiniMax Coding Plan (minimax.io)",
        "doc": "https://platform.minimax.io/docs/coding-plan/intro",
        "models": {
            "MiniMax-M2.7": {
                "id": "MiniMax-M2.7",
                "name": "MiniMax-M2.7",
                "family": "minimax",
                "attachment": false,
                "reasoning": true,
                "tool_call": true,
                "temperature": true,
                "release_date": "2026-03-18",
                "last_updated": "2026-03-18",
                "modalities": { "input": ["text"], "output": ["text"] },
                "open_weights": true,
                "cost": { "input": 0, "output": 0, "cache_read": 0, "cache_write": 0 },
                "limit": { "context": 204800, "output": 131072 }
            },
            "MiniMax-M2.7-highspeed": {
                "id": "MiniMax-M2.7-highspeed",
                "name": "MiniMax-M2.7-highspeed",
                "family": "minimax",
                "attachment": false,
                "reasoning": true,
                "tool_call": true,
                "temperature": true,
                "release_date": "2026-03-18",
                "last_updated": "2026-03-18",
                "modalities": { "input": ["text"], "output": ["text"] },
                "open_weights": true,
                "cost": { "input": 0, "output": 0, "cache_read": 0, "cache_write": 0 },
                "limit": { "context": 204800, "output": 131072 }
            }
        }
    }
}
```
