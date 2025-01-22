package token_information

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/soralabs/solana-toolkit/go/internal/pumpfun"
)

// IsPumpFunToken checks if a given mint address corresponds to a pump fun token
// by verifying if its bonding curve account exists
func (t *TokenInformationTool) IsPumpFunToken(mint solana.PublicKey) (bool, error) {
	bondingCurveAddr, _, err := pumpfun.DeriveBondingCurveAddresses(mint)
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
