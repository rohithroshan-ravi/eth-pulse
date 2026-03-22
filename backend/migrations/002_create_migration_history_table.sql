-- Migration: 002_create_migration_history_table
-- Description: Track applied migrations

CREATE TABLE IF NOT EXISTS public.migration_history (
  id SERIAL PRIMARY KEY,
  migration_name text NOT NULL UNIQUE,
  executed_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_migration_history_name ON public.migration_history(migration_name);
