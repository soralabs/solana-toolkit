package onchain_actions

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/ilkamo/jupiter-go/jupiter"

	"github.com/gagliardetto/solana-go"
	"github.com/soralabs/solana-toolkit/go/internal/pumpfun"
	toolkit "github.com/soralabs/toolkit/go"
)

type OnchainActionsTool struct {
	toolkit.Tool

	mu sync.Mutex

	rpcClient *rpc.Client

	jupClient *jupiter.ClientWithResponses
}

func NewOnchainActionsTool(rpcClient *rpc.Client) *OnchainActionsTool {
	jupClient, err := jupiter.NewClientWithResponses(jupiter.DefaultAPIURL)
	if err != nil {
		panic(fmt.Errorf("failed to create Jupiter client: %w", err))
	}

	return &OnchainActionsTool{
		rpcClient: rpcClient,
		jupClient: jupClient,
	}
}

func (t *OnchainActionsTool) GetName() string {
	return "onchain_actions"
}

func (t *OnchainActionsTool) GetDescription() string {
	return "Perform solana onchain actions: buy, sell, create, transfer."
}

func (t *OnchainActionsTool) GetSchema() toolkit.Schema {
	return toolkit.Schema{
		Parameters: json.RawMessage(`{
            "type": "object",
            "required": ["action", "params"],
            "properties": {
                "action": {
                    "type": "string",
                    "description": "The type of onchain action to perform",
                    "enum": [
                        "buy",
                        "sell",
                        "create",
						"transfer"
                    ]
                },
                "params": {
                    "type": "object",
                    "description": "Parameters specific to the action being performed",
                    "properties": {
                        "destination": {
                            "type": "string",
                            "description": "Destination wallet or account address for transfers"
                        },
                        "token_mint": {
                            "type": "string",
                            "description": "Token mint address for all actions except create"
                        },
                        "amount": {
                            "type": "number",
                            "description": "Amount of tokens/SOL (not in lamports) to transfer or interact with"
                        },
                        "token_name": {
                            "type": "string",
                            "description": "Name of the token to be created, when the action is create"
                        },
                        "token_symbol": {
                            "type": "string",
                            "description": "Symbol/ticker of the token to be created, when the action is create"
                        },
                    }
                }
            }
        }`),
	}
}

func (t *OnchainActionsTool) Execute(ctx context.Context, params json.RawMessage) (json.RawMessage, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	var input OnchainActionsInput
	if err := json.Unmarshal(params, &input); err != nil {
		return nil, fmt.Errorf("failed to parse parameters: %w", err)
	}

	// Validate private key for all actions
	if input.Params.Source == "" {
		return nil, fmt.Errorf("source/private key is required")
	}

	// Parse private key
	privateKey, err := solana.PrivateKeyFromBase58(input.Params.Source)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	var result interface{}

	switch input.Action {
	case ActionTransfer:
		if input.Params.Destination == "" || input.Params.TokenMint == "" || input.Params.Amount <= 0 {
			return nil, fmt.Errorf("invalid transfer parameters")
		}

		destPubKey, err := solana.PublicKeyFromBase58(input.Params.Destination)
		if err != nil {
			return nil, fmt.Errorf("invalid destination address: %w", err)
		}

		tokenMint, err := solana.PublicKeyFromBase58(input.Params.TokenMint)
		if err != nil {
			return nil, fmt.Errorf("invalid token mint address: %w", err)
		}

		// Convert SOL to lamports (1 SOL = 1e9 lamports)
		var lamports uint64
		if tokenMint.Equals(WSOL_MINT) {
			lamports = uint64(input.Params.Amount * 1e9)
		} else {
			lamports = uint64(input.Params.Amount)
		}

		sig, err := t.Transfer(ctx, privateKey, destPubKey, tokenMint, lamports)
		if err != nil {
			return nil, fmt.Errorf("failed to create transfer transaction: %w", err)
		}

		result = OnchainActionsOutput{
			Signature: sig.String(),
		}

	case ActionBuy, ActionSell:
		if input.Params.TokenMint == "" || input.Params.Amount <= 0 {
			return nil, fmt.Errorf("invalid swap parameters")
		}

		// Set up quote parameters
		quoteParams := jupiter.GetQuoteParams{
			Amount: jupiter.AmountParameter(input.Params.Amount),
		}

		if input.Action == ActionBuy {
			quoteParams.InputMint = "So11111111111111111111111111111111111111112" // SOL
			quoteParams.OutputMint = input.Params.TokenMint
		} else {
			quoteParams.InputMint = input.Params.TokenMint
			quoteParams.OutputMint = "So11111111111111111111111111111111111111112" // SOL
		}

		sig, err := t.Swap(ctx, quoteParams, privateKey)
		if err != nil {
			return nil, err
		}

		result = OnchainActionsOutput{
			Signature: sig.String(),
		}

	case ActionCreate:
		if input.Params.TokenName == "" || input.Params.TokenSymbol == "" {
			return nil, fmt.Errorf("invalid create parameters")
		}

		mintWallet := solana.NewWallet()
		sig, err := t.CreateToken(ctx, CreateTokenParams{
			TokenInfo: pumpfun.CreateTokenInformation{
				Name:   input.Params.TokenName,
				Symbol: input.Params.TokenSymbol,
			},
			Mint:            mintWallet,
			UserPrivateKey:  privateKey,
			BuyAmount:       input.Params.Amount,
			SlippagePercent: 10, // Default slippage
		})
		if err != nil {
			return nil, err
		}

		mintAddress := mintWallet.PublicKey().String()

		result = OnchainActionsOutput{
			Signature:   sig.String(),
			MintAddress: &mintAddress,
		}

	default:
		return nil, fmt.Errorf("unsupported action: %s", input.Action)
	}

	return json.Marshal(result)
}
