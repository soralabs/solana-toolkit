package token_information

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

const BONDING_CURVE_SEED = "bonding-curve"
const PUMP_FUN_PROGRAM_ID = "6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P"

// GetBondingCurvePDA returns the PDA (Program Derived Address) for a token's bonding curve
func GetBondingCurvePDA(mint solana.PublicKey) (solana.PublicKey, error) {
	seeds := [][]byte{
		[]byte(BONDING_CURVE_SEED),
		mint.Bytes(),
	}

	addr, _, err := solana.FindProgramAddress(seeds, solana.MustPublicKeyFromBase58(PUMP_FUN_PROGRAM_ID))
	if err != nil {
		return solana.PublicKey{}, err
	}

	return addr, nil
}

// IsPumpFunToken checks if a given mint address corresponds to a pump fun token
// by verifying if its bonding curve account exists
func (t *TokenInformationTool) IsPumpFunToken(mint solana.PublicKey) (bool, error) {
	bondingCurveAddr, err := GetBondingCurvePDA(mint)
	if err != nil {
		return false, fmt.Errorf("failed to get bonding curve PDA: %w", err)
	}

	account, err := t.rpcClient.GetAccountInfo(context.Background(), bondingCurveAddr)
	if err != nil {
		if err == rpc.ErrNotFound {
			return false, nil
		}

		return false, fmt.Errorf("failed to get account info: %w", err)
	}

	// If account data exists, it's a pump fun token
	return account != nil && len(account.Value.Data.GetBinary()) > 0, nil
}
