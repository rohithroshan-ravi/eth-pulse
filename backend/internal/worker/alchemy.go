package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"

	"eth-pulse/backend/internal/supastore"
	"eth-pulse/backend/internal/types"
)

type alchemySubscribeMessage struct {
	Method string `json:"method"`
	Params struct {
		Result pendingTx `json:"result"`
	} `json:"params"`
}

type pendingTx struct {
	Hash     string `json:"hash"`
	From     string `json:"from"`
	To       string `json:"to"`
	Value    string `json:"value"`
	GasPrice string `json:"gasPrice"`
}

type AlchemyWorker struct {
	alchemyWSSURL string
	supabase      *supastore.Client
	minWei        *big.Int
	latestGasGwei atomic.Uint64
}

func NewAlchemyWorker(alchemyWSSURL string, supabaseClient *supastore.Client, minETH float64) *AlchemyWorker {
	ethWei := new(big.Float).Mul(big.NewFloat(minETH), big.NewFloat(1e18))
	minWei := new(big.Int)
	ethWei.Int(minWei)

	return &AlchemyWorker{
		alchemyWSSURL: alchemyWSSURL,
		supabase:      supabaseClient,
		minWei:        minWei,
	}
}

func (w *AlchemyWorker) LatestGasGwei() float64 {
	return float64(w.latestGasGwei.Load()) / 1_000
}

func (w *AlchemyWorker) Start(ctx context.Context) error {
	dialer := websocket.Dialer{HandshakeTimeout: 10 * time.Second}
	conn, _, err := dialer.DialContext(ctx, w.alchemyWSSURL, nil)
	if err != nil {
		return fmt.Errorf("dial alchemy websocket: %w", err)
	}
	defer conn.Close()

	subscribe := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "eth_subscribe",
		"params": []any{
			"alchemy_pendingTransactions",
			map[string]any{"hashesOnly": false},
		},
	}
	if err := conn.WriteJSON(subscribe); err != nil {
		return fmt.Errorf("send subscribe message: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return fmt.Errorf("read websocket message: %w", err)
			}

			if err := w.processMessage(msg); err != nil {
				log.Printf("worker: skip message: %v", err)
			}
		}
	}
}

func (w *AlchemyWorker) processMessage(raw []byte) error {
	var msg alchemySubscribeMessage
	if err := json.Unmarshal(raw, &msg); err != nil {
		return fmt.Errorf("decode message: %w", err)
	}

	if msg.Method != "eth_subscription" || msg.Params.Result.Hash == "" {
		return nil
	}

	tx := msg.Params.Result
	valueWei, ok := new(big.Int).SetString(trimHex(tx.Value), 16)
	if !ok {
		return fmt.Errorf("invalid value hex: %s", tx.Value)
	}
	if valueWei.Cmp(w.minWei) <= 0 {
		return nil
	}

	if tx.GasPrice != "" {
		if gasWei, ok := new(big.Int).SetString(trimHex(tx.GasPrice), 16); ok {
			gweiMillis := new(big.Int).Div(new(big.Int).Mul(gasWei, big.NewInt(1_000)), big.NewInt(1e9))
			w.latestGasGwei.Store(gweiMillis.Uint64())
		}
	}

	valueEth := new(big.Float).Quo(new(big.Float).SetInt(valueWei), big.NewFloat(1e18))
	valueEthStr := valueEth.Text('f', 6)

	whaleTx := types.WhaleTransaction{
		TxHash:      tx.Hash,
		FromAddress: tx.From,
		ToAddress:   tx.To,
		ValueWei:    valueWei.String(),
		ValueETH:    valueEthStr,
		GasPriceWei: hexToDecString(tx.GasPrice),
	}

	if err := w.supabase.InsertWhaleTransaction(whaleTx); err != nil {
		return fmt.Errorf("persist whale tx: %w", err)
	}

	log.Printf("whale tx inserted: hash=%s value_eth=%s", whaleTx.TxHash, whaleTx.ValueETH)
	return nil
}

func trimHex(h string) string {
	if len(h) >= 2 && h[:2] == "0x" {
		return h[2:]
	}
	return h
}

func hexToDecString(h string) string {
	if h == "" {
		return ""
	}
	bi, ok := new(big.Int).SetString(trimHex(h), 16)
	if !ok {
		return ""
	}
	return bi.String()
}
