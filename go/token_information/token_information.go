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

// NewTokenInformationTool creates a new instance of TokenInformationTool with the given RPC client.
// This tool is responsible for fetching information about Solana tokens.
func NewTokenInformationTool(rpcClient *rpc.Client) *TokenInformationTool {
	return &TokenInformationTool{
		rpcClient: rpcClient,
	}
}

// GetName returns the name of the tool.
// This is used to identify the tool in the broader toolkit.
func (t *TokenInformationTool) GetName() string {
	return "token_information"
}

// GetDescription returns a description of the tool.
// This is used to provide context on what the tool does, such as fetching token-related information.
func (t *TokenInformationTool) GetDescription() string {
	return "Fetch information like name, symbol, price, etc. of a token"
}

// GetSchema returns the schema that defines the expected input parameters for this tool.
// It specifies that a "token_address" is required and provides a description of this parameter.
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

// Execute is the main function that executes the token information retrieval process.
// It handles the following steps:
// 1. Parses the input parameters.
// 2. Validates and converts the token address.
// 3. Fetches metadata, holder count, and pair information for the token.
// 4. Checks whether the token is associated with the PumpFun platform.
// 5. Consolidates and formats the data into a JSON response.
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
