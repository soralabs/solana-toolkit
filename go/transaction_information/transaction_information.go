package transaction_information

import (
	"context"
	"encoding/json"
	"fmt"

	"sync"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	"github.com/soralabs/solana-toolkit/go/internal/tx_parser"
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

	txHash, err := solana.SignatureFromBase58(input.Hash)
	if err != nil {
		return nil, fmt.Errorf("failed to parse transaction hash: %w", err)
	}

	maxSupportedTxVersion := uint64(0)

	tx, err := t.rpcClient.GetTransaction(ctx, txHash, &rpc.GetTransactionOpts{
		MaxSupportedTransactionVersion: &maxSupportedTxVersion,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	parser, err := tx_parser.New(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to create parser: %w", err)
	}

	swaps, err := parser.ParseTransaction()
	if err != nil {
		return nil, fmt.Errorf("failed to parse transaction: %w", err)
	}

	return json.Marshal(TransactionInformationOutput{
		Hash:  input.Hash,
		Swaps: swaps,
	})
}
