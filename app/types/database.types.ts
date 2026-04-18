export type Json =
  | string
  | number
  | boolean
  | null
  | { [key: string]: Json | undefined }
  | Json[]

export interface Database {
  public: {
    Tables: {
      providers: {
        Row: {
          id: number
          provider_id: string
          name: string
          npm: string
          env: string[]
          doc: string | null
          api: string | null
          synced_at: string | null
        }
        Insert: {
          id?: number
          provider_id: string
          name: string
          npm: string
          env?: string[]
          doc?: string | null
          api?: string | null
          synced_at?: string | null
        }
        Update: {
          id?: number
          provider_id?: string
          name?: string
          npm?: string
          env?: string[]
          doc?: string | null
          api?: string | null
          synced_at?: string | null
        }
      }
      models: {
        Row: {
          id: number
          model_id: string
          provider_id: string
          name: string
          family: string | null
          attachment: boolean
          reasoning: boolean
          tool_call: boolean
          structured_output: boolean | null
          temperature: boolean | null
          knowledge: string | null
          release_date: string | null
          last_updated: string | null
          open_weights: boolean
          status: string | null
          interleaved: Json | null
          modalities_input: string[]
          modalities_output: string[]
          cost_input: number | null
          cost_output: number | null
          cost_reasoning: number | null
          cost_cache_read: number | null
          cost_cache_write: number | null
          cost_input_audio: number | null
          cost_output_audio: number | null
          limit_context: number | null
          limit_input: number | null
          limit_output: number | null
        }
        Insert: {
          id?: number
          model_id: string
          provider_id: string
          name: string
          family?: string | null
          attachment?: boolean
          reasoning?: boolean
          tool_call?: boolean
          structured_output?: boolean | null
          temperature?: boolean | null
          knowledge?: string | null
          release_date?: string | null
          last_updated?: string | null
          open_weights?: boolean
          status?: string | null
          interleaved?: Json | null
          modalities_input?: string[]
          modalities_output?: string[]
          cost_input?: number | null
          cost_output?: number | null
          cost_reasoning?: number | null
          cost_cache_read?: number | null
          cost_cache_write?: number | null
          cost_input_audio?: number | null
          cost_output_audio?: number | null
          limit_context?: number | null
          limit_input?: number | null
          limit_output?: number | null
        }
        Update: {
          id?: number
          model_id?: string
          provider_id?: string
          name?: string
          family?: string | null
          attachment?: boolean
          reasoning?: boolean
          tool_call?: boolean
          structured_output?: boolean | null
          temperature?: boolean | null
          knowledge?: string | null
          release_date?: string | null
          last_updated?: string | null
          open_weights?: boolean
          status?: string | null
          interleaved?: Json | null
          modalities_input?: string[]
          modalities_output?: string[]
          cost_input?: number | null
          cost_output?: number | null
          cost_reasoning?: number | null
          cost_cache_read?: number | null
          cost_cache_write?: number | null
          cost_input_audio?: number | null
          cost_output_audio?: number | null
          limit_context?: number | null
          limit_input?: number | null
          limit_output?: number | null
        }
      }
    }
    Views: {}
    Functions: {}
    Enums: {}
  }
}
