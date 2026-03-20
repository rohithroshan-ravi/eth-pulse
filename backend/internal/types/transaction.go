package types

import "time"

type WhaleTransaction struct {
	TxHash      string    `json:"tx_hash"`
	FromAddress string    `json:"from_address"`
	ToAddress   string    `json:"to_address,omitempty"`
	ValueWei    string    `json:"value_wei"`
	ValueETH    string    `json:"value_eth"`
	GasPriceWei string    `json:"gas_price_wei,omitempty"`
	InsertedAt  time.Time `json:"inserted_at,omitempty"`
}
