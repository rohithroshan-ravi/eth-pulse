package supastore

import (
	"fmt"

	supabase "github.com/supabase-community/supabase-go"

	"eth-pulse/backend/internal/types"
)

type Client struct {
	inner *supabase.Client
}

func New(url, serviceRoleKey string) (*Client, error) {
	client, err := supabase.NewClient(url, serviceRoleKey, nil)
	if err != nil {
		return nil, fmt.Errorf("init supabase client: %w", err)
	}
	return &Client{inner: client}, nil
}

func (c *Client) InsertWhaleTransaction(tx types.WhaleTransaction) error {
	payload := map[string]any{
		"tx_hash":       tx.TxHash,
		"from_address":  tx.FromAddress,
		"to_address":    tx.ToAddress,
		"value_wei":     tx.ValueWei,
		"value_eth":     tx.ValueETH,
		"gas_price_wei": tx.GasPriceWei,
	}

	_, _, err := c.inner.
		From("transactions").
		Insert(payload, false, "", "", "").
		Execute()
	if err != nil {
		return fmt.Errorf("insert whale transaction: %w", err)
	}

	return nil
}
