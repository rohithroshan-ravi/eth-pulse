create extension if not exists "pgcrypto";

create table if not exists public.transactions (
  id uuid primary key default gen_random_uuid(),
  tx_hash text not null unique,
  from_address text not null,
  to_address text,
  value_wei numeric(78, 0) not null,
  value_eth numeric(38, 18) not null,
  gas_price_wei numeric(78, 0),
  inserted_at timestamptz not null default now()
);

alter publication supabase_realtime add table public.transactions;
