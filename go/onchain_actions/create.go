package onchain_actions

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/soralabs/solana-toolkit/go/internal/pumpfun"
)

type CreateTokenParams struct {
	TokenInfo       pumpfun.CreateTokenInformation
	Mint            *solana.Wallet
	UserPrivateKey  solana.PrivateKey
	BuyAmount       float64
	SlippagePercent float64
}

func (o *OnchainActionsTool) CreateToken(ctx context.Context, params CreateTokenParams) error {
	_, err := pumpfun.CreateToken(ctx, pumpfun.CreateTokenRequest{
		RpcClient:       o.rpcClient,
		TokenInfo:       params.TokenInfo,
		Mint:            params.Mint,
		UserPrivateKey:  params.UserPrivateKey,
		BuyAmount:       params.BuyAmount,
		SlippagePercent: params.SlippagePercent,
	})

	if err != nil {
		return fmt.Errorf("failed to create token: %w", err)
	}

	return nil
}
