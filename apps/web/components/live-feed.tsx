"use client";

import { useEffect, useMemo, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { supabase } from "@/lib/supabase/client";
import type { RealtimePostgresInsertPayload } from "@supabase/supabase-js";

type WhaleTx = {
  id: string;
  tx_hash: string;
  from_address: string;
  to_address: string | null;
  value_eth: string;
  gas_price_wei: string | null;
  inserted_at: string;
};

export default function LiveFeed() {
  const [txs, setTxs] = useState<WhaleTx[]>([]);
  const [freshIds, setFreshIds] = useState<Set<string>>(new Set());

  useEffect(() => {
    const channel = supabase
      .channel("transactions-feed")
      .on(
        "postgres_changes",
        {
          event: "INSERT",
          schema: "public",
          table: "transactions",
        },
        (payload: RealtimePostgresInsertPayload<WhaleTx>) => {
          const row = payload.new as WhaleTx;
          setTxs((prev: WhaleTx[]) => [row, ...prev].slice(0, 30));
          setFreshIds((prev: Set<string>) => {
            const next = new Set(prev);
            next.add(row.id);
            return next;
          });

          window.setTimeout(() => {
            setFreshIds((prev: Set<string>) => {
              const next = new Set(prev);
              next.delete(row.id);
              return next;
            });
          }, 1800);
        }
      )
      .subscribe();

    return () => {
      supabase.removeChannel(channel);
    };
  }, []);

  const cards = useMemo(
    () =>
      txs.map((tx: WhaleTx) => {
        const gasGwei = tx.gas_price_wei
          ? (Number(tx.gas_price_wei) / 1e9).toFixed(2)
          : "n/a";

        return (
          <Card
            key={tx.id}
            className={`transition-all duration-300 ${
              freshIds.has(tx.id) ? "animate-pulse ring-2 ring-emerald-300" : ""
            }`}
          >
            <CardHeader>
              <CardTitle className="font-mono text-sm text-emerald-100">
                {tx.tx_hash.slice(0, 12)}...{tx.tx_hash.slice(-8)}
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-2 text-sm text-slate-100">
              <p>
                <span className="text-slate-300">From:</span> {tx.from_address}
              </p>
              <p>
                <span className="text-slate-300">To:</span> {tx.to_address ?? "Contract Creation"}
              </p>
              <p>
                <span className="text-slate-300">Value:</span>{" "}
                <span className="font-semibold text-emerald-300">{tx.value_eth} ETH</span>
              </p>
              <p>
                <span className="text-slate-300">Gas:</span> {gasGwei} GWEI
              </p>
            </CardContent>
          </Card>
        );
      }),
    [freshIds, txs]
  );

  return (
    <section className="mx-auto w-full max-w-5xl space-y-6">
      <div>
        <h2 className="text-3xl font-bold tracking-tight text-white">Live Whale Feed</h2>
        <p className="text-slate-300">Incoming Ethereum transactions above 5 ETH.</p>
      </div>
      <div className="grid gap-4 sm:grid-cols-2">{cards}</div>
    </section>
  );
}
