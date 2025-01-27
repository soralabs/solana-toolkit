package pumpfun

import (
	"context"

	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
)

func DeriveBondingCurveAddresses(mint solana.PublicKey) (bondingCurve, associatedBondingCurve solana.PublicKey, err error) {
	seeds := [][]byte{
		[]byte("bonding-curve"),
		mint.Bytes(),
	}
	bondingCurve, _, err = solana.FindProgramAddress(seeds, ProgramID)
	if err != nil {
		return solana.PublicKey{}, solana.PublicKey{}, err
	}

	associatedBondingCurve, _, err = solana.FindAssociatedTokenAddress(
		bondingCurve,
		mint,
	)
	if err != nil {
		return solana.PublicKey{}, solana.PublicKey{}, err
	}

	return bondingCurve, associatedBondingCurve, nil
}

func shouldCreateAta(rpcClient *rpc.Client, ata solana.PublicKey) (bool, error) {
	account, err := rpcClient.GetAccountInfo(context.Background(), ata)
	if err != nil {
		// If we get an error, it likely means the account doesn't exist
		return true, nil
	}
	return account == nil, nil
}

func getComputeUnitPriceInstruction(ctx context.Context, rpcClient *rpc.Client, user solana.PrivateKey) (*computebudget.SetComputeUnitPrice, error) {
	out, err := rpcClient.GetRecentPrioritizationFees(
		ctx,
		solana.PublicKeySlice{
			user.PublicKey(),
			ProgramID,
			MintAuthority,
			GlobalPumpFunAddress,
			solana.TokenMetadataProgramID,
			system.ProgramID,
			token.ProgramID,
			associatedtokenaccount.ProgramID,
			solana.SysVarRentPubkey,
			EventAuthority,
		},
	)
	if err != nil {
		return nil, err
	}

	var median uint64
	for _, fee := range out {
		median += fee.PrioritizationFee
	}
	median /= uint64(len(out))

	return computebudget.NewSetComputeUnitPriceInstruction(median), nil
}
