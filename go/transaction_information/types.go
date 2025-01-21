package transaction_information

import "github.com/soralabs/solana-toolkit/go/internal/tx_parser"

type TransactionInformationInput struct {
	Hash string `json:"hash"`
}

type TransactionInformationOutput struct {
	Hash  string                `json:"hash"`
	Swaps []*tx_parser.SwapInfo `json:"swaps"`
}
