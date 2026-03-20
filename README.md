# eth-pulse

A high-concurrency, real-time Ethereum transaction visualizer built with Go and Next.js.

## Project Structure

```text
eth-pulse/
├── backend/
│   ├── .env.example
│   ├── go.mod
│   ├── main.go
│   └── internal/
│       ├── supastore/
│       │   └── client.go
│       ├── types/
│       │   └── transaction.go
│       └── worker/
│           └── alchemy.go
├── apps/
│   └── web/
│       ├── .env.example
│       ├── package.json
│       ├── next.config.mjs
│       ├── postcss.config.js
│       ├── tailwind.config.ts
│       ├── tsconfig.json
│       ├── app/
│       │   ├── globals.css
│       │   ├── layout.tsx
│       │   └── page.tsx
│       ├── components/
│       │   ├── live-feed.tsx
│       │   └── ui/
│       │       └── card.tsx
│       └── lib/
│           └── supabase/
│               └── client.ts
└── db/
		└── schema.sql
```

## Backend (Go + Echo + gorilla/websocket)

`backend/internal/worker/alchemy.go`:

- Connects to Alchemy Ethereum Mainnet WSS using `gorilla/websocket`.
- Subscribes to `alchemy_pendingTransactions`.
- Parses incoming transaction payloads.
- Filters whale transactions where `value > 5 ETH` (configurable via `WHALE_MIN_ETH`).
- Tracks latest gas price in gwei.
- Inserts filtered transactions into Supabase `transactions` table.

`backend/main.go`:

- Boots Echo API.
- Starts the worker with automatic reconnect.
- Exposes:
	- `GET /health`
	- `GET /metrics` (includes latest tracked gas value).

## Frontend (Next.js 14 + Supabase Realtime + Shadcn-style Cards)

`apps/web/components/live-feed.tsx`:

- Subscribes to `INSERT` events on `public.transactions` using Supabase Realtime.
- Prepends incoming rows in-memory for a live feed UX.
- Renders each transaction in Card components.
- Applies `animate-pulse` + ring highlight when a new event arrives.

## Database (Supabase/Postgres)

Run `db/schema.sql` in Supabase SQL editor to create `public.transactions` and enable realtime publication.

## Environment

Backend (`backend/.env`):

```env
PORT=8080
ALCHEMY_WSS_URL=wss://eth-mainnet.g.alchemy.com/v2/your-key
SUPABASE_URL=https://your-project-ref.supabase.co
SUPABASE_SERVICE_ROLE_KEY=your-service-role-key
WHALE_MIN_ETH=5
```

Frontend (`apps/web/.env.local`):

```env
NEXT_PUBLIC_SUPABASE_URL=https://your-project-ref.supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY=your-anon-key
```

## Local Run

1. Create Supabase table by running `db/schema.sql`.
2. Configure backend `.env` and frontend `.env.local`.
3. Start backend:

```bash
cd backend
go mod tidy
go run .
```

4. Start frontend:

```bash
cd apps/web
npm install
npm run dev
```

Open `http://localhost:3000` for the live whale transaction feed.
