package token_information

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func (t *TokenInformationTool) getHolderCount(ctx context.Context, tokenAddress solana.PublicKey) (int, error) {
	filters := []rpc.RPCFilter{
		{
			DataSize: 165, // Token account size
		},
		{
			Memcmp: &rpc.RPCFilterMemcmp{
				Offset: 0,
				Bytes:  tokenAddress.Bytes(),
			},
		},
	}

	offset := uint64(64)
	length := uint64(8)

	opts := &rpc.GetProgramAccountsOpts{
		Filters:   filters,
		Encoding:  solana.EncodingBase64,
		DataSlice: &rpc.DataSlice{Offset: &offset, Length: &length},
	}

	accounts, err := t.rpcClient.GetProgramAccountsWithOpts(ctx, solana.TokenProgramID, opts)
	if err != nil {
		return 0, fmt.Errorf("failed to get program accounts: %w", err)
	}

	isZeroBytes := func(data []byte) bool {
		for _, b := range data {
			if b != 0 {
				return false
			}
		}
		return true
	}

	activeHolders := 0
	for _, account := range accounts {
		if !isZeroBytes(account.Account.Data.GetBinary()) {
			activeHolders++
		}
	}

	return activeHolders, nil
}
