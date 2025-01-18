package token_information

import (
	"context"
	"encoding/json"
	"fmt"

	"sync"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	"github.com/soralabs/solana-toolkit/go/internal/dexscreener"
	"github.com/soralabs/solana-toolkit/go/internal/pumpfun"
	toolkit "github.com/soralabs/toolkit/go"
)

type TokenInformationTool struct {
	toolkit.Tool

	mu sync.Mutex

	rpcClient *rpc.Client
}

func NewTokenInformationTool(rpcClient *rpc.Client) *TokenInformationTool {
	return &TokenInformationTool{
		rpcClient: rpcClient,
	}
}

func (t *TokenInformationTool) GetName() string {
	return "token_information"
}

func (t *TokenInformationTool) GetDescription() string {
	return "Fetch information like name, symbol, price, etc. of a token"
}

func (t *TokenInformationTool) GetSchema() toolkit.Schema {
	return toolkit.Schema{
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["token_address"],
			"properties": {
				"token_address": {
					"type": "string",
					"description": "The address of the token"
				}
			}
		}`),
	}
}

func (t *TokenInformationTool) Execute(ctx context.Context, params json.RawMessage) (json.RawMessage, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	var input TokenInformationInput
	if err := json.Unmarshal(params, &input); err != nil {
		return nil, fmt.Errorf("failed to parse parameters: %w", err)
	}

	tokenAddress, err := solana.PublicKeyFromBase58(input.TokenAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token address: %w", err)
	}

	metadata, err := t.getMetadata(ctx, tokenAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata: %w", err)
	}

	isPumpFunToken, err := t.IsPumpFunToken(tokenAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to check if token is pump fun: %w", err)
	}

	pairs, err := dexscreener.GetPairInformation(ctx, input.TokenAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get pair information: %w", err)
	}

	if !isPumpFunToken && len(pairs) == 0 {
		return nil, fmt.Errorf("not a valid tradeable token")
	}

	holderCount, err := t.getHolderCount(ctx, tokenAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get holder count: %w", err)
	}

	if isPumpFunToken && len(pairs) == 0 {
		pfInfo, err := pumpfun.GetTokenInformation(ctx, input.TokenAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to get PumpFun token information: %w", err)
		}

		return json.Marshal(TokenInformationOutput{
			Metadata:       metadata,
			IsPumpFunToken: true,
			USDMarketCap:   fmt.Sprintf("%f", pfInfo.UsdMarketCap),
			Socials:        make([]Social, 0),
			HolderCount:    holderCount,
		})
	}

	mainPair := pairs[0]

	socials := make([]Social, len(mainPair.Info.Socials))
	for i, social := range mainPair.Info.Socials {
		socials[i] = Social{
			Type: social.Type,
			URL:  social.URL,
		}
	}

	return json.Marshal(TokenInformationOutput{
		Metadata:       metadata,
		IsPumpFunToken: isPumpFunToken,
		USDMarketCap:   fmt.Sprintf("%f", mainPair.MarketCap),
		Socials:        socials,
		HolderCount:    holderCount,
		PriceChange: &PriceChange{
			H24: mainPair.PriceChange.H24,
			H6:  mainPair.PriceChange.H6,
			H1:  mainPair.PriceChange.H1,
			M5:  mainPair.PriceChange.M5,
		},
	})
}
