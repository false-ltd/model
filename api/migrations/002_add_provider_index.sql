-- Incremental migration for deployments initialized with an older init.sql.
-- Applies the same schema change as init.sql (idx_models_provider), which
-- backs provider filtering, the stats provider JOIN/GROUP BY, and the
-- providers model counts query.
ALTER TABLE models ADD INDEX idx_models_provider (provider_id);
