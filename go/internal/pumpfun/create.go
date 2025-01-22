package pumpfun

import (
	"context"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	pump "github.com/soralabs/solana-toolkit/go/internal/pumpfun_anchor"
)

func CreateToken(ctx context.Context, request CreateTokenRequest) (solana.Signature, error) {
	// Derive bonding curve addresses
	bondingCurve, associatedBondingCurve, err := DeriveBondingCurveAddresses(request.Mint.PublicKey())
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to derive bonding curve addresses: %w", err)
	}

	// Get token metadata address
	metadata, _, err := solana.FindTokenMetadataAddress(request.Mint.PublicKey())
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to find token metadata address: %w", err)
	}

	// Build transaction instructions
	instructions := []solana.Instruction{
		// Set compute unit limit (250k - pump.fun default)
		computebudget.NewSetComputeUnitLimitInstruction(250000).Build(),
	}

	// Add compute unit price instruction if available
	if cupInst, err := getComputeUnitPriceInstruction(ctx, request.RpcClient, request.UserPrivateKey); err == nil {
		instructions = append(instructions, cupInst.Build())
	}

	// Create token instruction
	createInst := buildCreateTokenInstruction(
		request.TokenInfo,
		request.Mint.PublicKey(),
		bondingCurve,
		associatedBondingCurve,
		metadata,
		request.UserPrivateKey.PublicKey(),
	)
	instructions = append(instructions, createInst)

	// Add buy instructions if BuyAmount is specified
	if request.BuyAmount > 0 {
		// Get global account data
		global, err := GetGlobalAccount(ctx, request.RpcClient)
		if err != nil {
			return solana.Signature{}, fmt.Errorf("failed to get global account: %w", err)
		}

		buyInstructions, err := buildBuyInstructions(
			request.RpcClient,
			request.Mint.PublicKey(),
			request.UserPrivateKey.PublicKey(),
			request.BuyAmount,
			global,
		)
		if err != nil {
			return solana.Signature{}, fmt.Errorf("failed to build buy instructions: %w", err)
		}
		instructions = append(instructions, buyInstructions...)
	}

	// Get recent blockhash
	recent, err := request.RpcClient.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to get recent blockhash: %w", err)
	}

	// Build and sign transaction
	tx, err := solana.NewTransaction(
		instructions,
		recent.Value.Blockhash,
		solana.TransactionPayer(request.UserPrivateKey.PublicKey()),
	)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Sign transaction
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if request.UserPrivateKey.PublicKey().Equals(key) {
			return &request.UserPrivateKey
		}
		if request.Mint.PublicKey().Equals(key) {
			return &request.Mint.PrivateKey
		}
		return nil
	})
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send transaction
	sig, err := request.RpcClient.SendTransaction(ctx, tx)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to send transaction: %w", err)
	}

	// Poll for confirmation
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return solana.Signature{}, fmt.Errorf("transaction confirmation timed out")
		default:
			status, err := request.RpcClient.GetSignatureStatuses(ctx, true, sig)
			if err != nil {
				return solana.Signature{}, fmt.Errorf("failed to get signature status: %w", err)
			}

			if status.Value[0] != nil && status.Value[0].Confirmations != nil && *status.Value[0].Confirmations > 0 {
				return sig, nil
			}

			time.Sleep(500 * time.Millisecond)
		}
	}
}

func buildCreateTokenInstruction(
	info CreateTokenInformation,
	mint, bondingCurve, associatedBondingCurve, metadata, user solana.PublicKey,
) solana.Instruction {
	return pump.NewCreateInstruction(
		info.Name,
		info.Symbol,
		info.ImageURI,
		mint,
		MintAuthority,
		bondingCurve,
		associatedBondingCurve,
		GlobalPumpFunAddress,
		solana.TokenMetadataProgramID,
		metadata,
		user,
		system.ProgramID,
		token.ProgramID,
		associatedtokenaccount.ProgramID,
		solana.SysVarRentPubkey,
		EventAuthority,
		ProgramID,
	).Build()
}

func buildBuyInstructions(
	rpcClient *rpc.Client,
	mint solana.PublicKey,
	user solana.PublicKey,
	solAmount uint64,
	global *GlobalAccount,
) ([]solana.Instruction, error) {
	bondingCurve, associatedBondingCurve, err := DeriveBondingCurveAddresses(mint)
	if err != nil {
		return nil, fmt.Errorf("failed to get bonding curve data: %w", err)
	}

	var instructions []solana.Instruction

	// Get or create associated token account
	ata, _, err := solana.FindAssociatedTokenAddress(
		user,
		mint,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to derive associated token account: %w", err)
	}

	shouldCreateATA, err := shouldCreateAta(rpcClient, ata)
	if err != nil {
		return nil, fmt.Errorf("can't check if we should create ATA: %w", err)
	}
	if shouldCreateATA {
		ataInstr, err := associatedtokenaccount.NewCreateInstruction(user, user, mint).
			ValidateAndBuild()
		if err != nil {
			return nil, fmt.Errorf("can't create associated token account: %w", err)
		}
		instructions = append(instructions, ataInstr)
	}

	buyAmount, err := calculateInitialBuyAmount(solAmount, global)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate initial buy amount: %w", err)
	}

	buyInstr := pump.NewBuyInstruction(
		buyAmount,
		solAmount,
		GlobalPumpFunAddress,
		MintAuthority,
		mint,
		bondingCurve,
		associatedBondingCurve,
		ata,
		user,
		system.ProgramID,
		token.ProgramID,
		solana.SysVarRentPubkey,
		EventAuthority,
		ProgramID,
	)
	instructions = append(instructions, buyInstr.Build())
	return instructions, nil
}
