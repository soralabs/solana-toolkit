package onchain_actions

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
)

// Transfer creates and signs a transfer transaction for either SOL or SPL tokens
func (t *OnchainActionsTool) Transfer(
	ctx context.Context,
	from solana.PrivateKey,
	tokenMint solana.PublicKey,
	to solana.PublicKey,
	amount uint64,
) (*solana.Signature, error) {
	var instructions []solana.Instruction

	if tokenMint.Equals(WSOL_MINT) {
		// Native SOL transfer
		transferIx := system.NewTransferInstruction(
			amount,
			from.PublicKey(),
			to,
		).Build()
		instructions = append(instructions, transferIx)
	} else {
		// SPL token transfer
		fromATA, _, err := solana.FindAssociatedTokenAddress(from.PublicKey(), tokenMint)
		if err != nil {
			return nil, fmt.Errorf("failed to find source associated token account: %w", err)
		}

		toATA, _, err := solana.FindAssociatedTokenAddress(to, tokenMint)
		if err != nil {
			return nil, fmt.Errorf("failed to find destination associated token account: %w", err)
		}

		// Check if destination ATA exists
		_, err = t.rpcClient.GetAccountInfo(ctx, toATA)
		if err != nil || err == rpc.ErrNotFound {
			// Create destination ATA if it doesn't exist
			createATAIx, err := associatedtokenaccount.NewCreateInstruction(
				from.PublicKey(),
				to,
				tokenMint,
			).ValidateAndBuild()
			if err != nil {
				return nil, fmt.Errorf("failed to create ATA instruction: %w", err)
			}
			instructions = append(instructions, createATAIx)
		}

		// Transfer SPL token
		transferIx, err := token.NewTransferInstruction(
			amount,
			fromATA,
			toATA,
			from.PublicKey(),
			[]solana.PublicKey{},
		).ValidateAndBuild()
		if err != nil {
			return nil, fmt.Errorf("failed to create transfer instruction: %w", err)
		}
		instructions = append(instructions, transferIx)
	}

	// Get recent blockhash
	recent, err := t.rpcClient.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent blockhash: %w", err)
	}

	// Create transaction
	tx, err := solana.NewTransaction(
		instructions,
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

	return t.sendTransacton(ctx, tx)
}
