package onchain_actions

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/ilkamo/jupiter-go/jupiter"

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
                        "amount": {
                            "type": "number",
                            "description": "Amount of tokens/SOL in LAMPORTS to transfer or interact with"
                        },
                        "token_mint": {
                            "type": "string",
                            "description": "Token mint address for all actions"
                        },
                        "slippage": {
                            "type": "number",
                            "description": "Maximum acceptable slippage percentage for swaps"
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

	return json.Marshal(OnchainActionsOutput{})
}
