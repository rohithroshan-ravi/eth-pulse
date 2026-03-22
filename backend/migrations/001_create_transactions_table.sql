-- Migration: 001_create_transactions_table
-- Description: Create transactions table to store whale transaction data

CREATE TABLE IF NOT EXISTS public.transactions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  tx_hash text NOT NULL UNIQUE,
  from_address text NOT NULL,
  to_address text,
  value_wei numeric(78, 0) NOT NULL,
  value_eth numeric(38, 18) NOT NULL,
  gas_price_wei numeric(78, 0),
  gas_price_gwei numeric(20, 2),
  inserted_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

-- Create indexes for common queries
CREATE INDEX IF NOT EXISTS idx_transactions_inserted_at ON public.transactions(inserted_at DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_value_eth ON public.transactions(value_eth DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_from_address ON public.transactions(from_address);
CREATE INDEX IF NOT EXISTS idx_transactions_to_address ON public.transactions(to_address);
