package onchain_actions

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/soralabs/solana-toolkit/go/internal/jupiter"
)

// Swap creates and signs a swap transaction with priority fees
func (t *OnchainActionsTool) Swap(swapRequest jupiter.SwapRequest, signer solana.PrivateKey) (*solana.Transaction, error) {
	swapTransactionResp, err := jupiter.GetSwapTransaction(swapRequest)
	if err != nil {
		return nil, err
	}

	// Extract transaction string from response
	swapTransactionBytes, ok := swapTransactionResp.([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected swap transaction response type")
	}

	swapTransactionStr := string(swapTransactionBytes)

	// Decode base64 transaction
	txBytes, err := base64.StdEncoding.DecodeString(swapTransactionStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 transaction: %w", err)
	}

	// Deserialize the transaction
	tx, err := solana.TransactionFromBytes(txBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize transaction: %w", err)
	}

	// Sign the transaction
	tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(signer.PublicKey()) {
			return &signer
		}
		return nil
	})

	return tx, nil
}

// SendSwapTransaction sends a signed swap transaction to the Solana network
func (t *OnchainActionsTool) SendSwapTransaction(ctx context.Context, signedTx *solana.Transaction) (solana.Signature, error) {
	sig, err := t.rpcClient.SendTransaction(ctx, signedTx)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to send transaction: %w", err)
	}

	// Wait for confirmation with retries
	deadline := time.Now().Add(30 * time.Second)

	for time.Now().Before(deadline) {
		status, err := t.rpcClient.GetSignatureStatuses(ctx, true, sig)
		if err != nil {
			return sig, fmt.Errorf("failed to get transaction status: %w", err)
		}

		if status.Value[0] != nil {
			if status.Value[0].Err != nil {
				return sig, fmt.Errorf("transaction failed: %v", status.Value[0].Err)
			}
			if status.Value[0].Confirmations != nil && *status.Value[0].Confirmations > 0 {
				return sig, nil
			}
		}

		time.Sleep(500 * time.Millisecond)
	}

	return sig, fmt.Errorf("transaction confirmation timeout")
}
