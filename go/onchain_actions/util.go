package onchain_actions

import (
	"context"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"
)

// sendTransacton sends a signed transaction to the Solana network
func (t *OnchainActionsTool) sendTransacton(ctx context.Context, signedTx *solana.Transaction) (*solana.Signature, error) {
	sig, err := t.rpcClient.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	// Wait for confirmation with retries
	deadline := time.Now().Add(30 * time.Second)

	for time.Now().Before(deadline) {
		status, err := t.rpcClient.GetSignatureStatuses(ctx, true, sig)
		if err != nil {
			return &sig, fmt.Errorf("failed to get transaction status: %w", err)
		}

		if status.Value[0] != nil {
			if status.Value[0].Err != nil {
				return &sig, fmt.Errorf("transaction failed: %v", status.Value[0].Err)
			}
			if status.Value[0].Confirmations != nil && *status.Value[0].Confirmations > 0 {
				return &sig, nil
			}
		}

		time.Sleep(500 * time.Millisecond)
	}

	return &sig, fmt.Errorf("transaction confirmation timeout")
}
