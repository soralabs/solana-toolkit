package transaction_information

import (
	"context"
	"encoding/json"
	"fmt"

	"sync"

	"github.com/gagliardetto/solana-go/rpc"

	toolkit "github.com/soralabs/toolkit/go"
)

type TransactionInformationTool struct {
	toolkit.Tool

	mu sync.Mutex

	rpcClient *rpc.Client
}

func NewTransactionInformationTool(rpcClient *rpc.Client) *TransactionInformationTool {
	return &TransactionInformationTool{
		rpcClient: rpcClient,
	}
}

func (t *TransactionInformationTool) GetName() string {
	return "transaction_information"
}

func (t *TransactionInformationTool) GetDescription() string {
	return "Fetch information about a solana transaction"
}

func (t *TransactionInformationTool) GetSchema() toolkit.Schema {
	return toolkit.Schema{
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["hash"],
			"properties": {
				"hash": {
					"type": "string",
					"description": "The hash of the transaction"
				}
			}
		}`),
	}
}

func (t *TransactionInformationTool) Execute(ctx context.Context, params json.RawMessage) (json.RawMessage, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	var input TransactionInformationInput
	if err := json.Unmarshal(params, &input); err != nil {
		return nil, fmt.Errorf("failed to parse parameters: %w", err)
	}

	return json.Marshal(TransactionInformationOutput{
		Hash: input.Hash,
	})
}
