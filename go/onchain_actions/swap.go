package onchain_actions

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/ilkamo/jupiter-go/jupiter"
)

// Swap creates and signs a swap transaction with priority fees
func (t *OnchainActionsTool) Swap(ctx context.Context, quoteRequest jupiter.GetQuoteParams, signer solana.PrivateKey) (*solana.Signature, error) {
	// Get quote using Jupiter client
	quoteResponse, err := t.jupClient.GetQuoteWithResponse(ctx, &quoteRequest)
	if err != nil {
		return nil, err
	}

	if quoteResponse.JSON200 == nil {
		return nil, fmt.Errorf("invalid GetQuoteWithResponse response")
	}

	// Setup swap request parameters
	prioritizationFeeLamports := jupiter.SwapRequest_PrioritizationFeeLamports{}
	if err = prioritizationFeeLamports.UnmarshalJSON([]byte(`"auto"`)); err != nil {
		return nil, err
	}

	dynamicComputeUnitLimit := true

	// Get swap transaction
	swapResponse, err := t.jupClient.PostSwapWithResponse(ctx, jupiter.PostSwapJSONRequestBody{
		PrioritizationFeeLamports: &prioritizationFeeLamports,
		QuoteResponse:             *quoteResponse.JSON200,
		UserPublicKey:             signer.PublicKey().String(),
		DynamicComputeUnitLimit:   &dynamicComputeUnitLimit,
	})
	if err != nil {
		return nil, err
	}

	if swapResponse.JSON200 == nil {
		return nil, fmt.Errorf("invalid PostSwapWithResponse response")
	}

	// Decode base64 transaction
	txBytes, err := base64.StdEncoding.DecodeString(swapResponse.JSON200.SwapTransaction)
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

	return t.sendTransacton(ctx, tx)
}
