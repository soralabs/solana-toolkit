package onchain_actions

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
)

// Transfer creates and signs a SOL transfer transaction
func (t *OnchainActionsTool) Transfer(
	ctx context.Context,
	from solana.PrivateKey,
	to solana.PublicKey,
	amount uint64,
) (*solana.Transaction, error) {
	// Create transfer instruction
	transferIx := system.NewTransferInstruction(
		amount,
		from.PublicKey(),
		to,
	).Build()

	// Get recent blockhash
	recent, err := t.rpcClient.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent blockhash: %w", err)
	}

	// Create transaction
	tx, err := solana.NewTransaction(
		[]solana.Instruction{transferIx},
		recent.Value.Blockhash,
		solana.TransactionPayer(from.PublicKey()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Sign transaction
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(from.PublicKey()) {
			return &from
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	return tx, nil
}
