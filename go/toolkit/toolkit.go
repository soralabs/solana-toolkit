package toolkit

import (
	"github.com/gagliardetto/solana-go/rpc"

	"github.com/soralabs/solana-toolkit/go/onchain_actions"
	"github.com/soralabs/solana-toolkit/go/token_information"
	"github.com/soralabs/solana-toolkit/go/transaction_information"
	"github.com/soralabs/solana-toolkit/go/wallet_information"
	toolkit "github.com/soralabs/toolkit/go"
)

func New(rpcUrl string) (*toolkit.Toolkit, error) {
	rpcClient := rpc.New(rpcUrl)

	tk := toolkit.NewToolkit(
		toolkit.WithToolkitName("solana-toolkit"),
		toolkit.WithToolkitDescription("A toolkit for Solana operations"),
	)

	// Initialize and add all tools
	transactionTool, err := transaction_information.NewTransactionInformationTool(rpcClient)
	if err != nil {
		return nil, err
	}
	walletTool, err := wallet_information.NewWalletInformationTool(rpcClient)
	if err != nil {
		return nil, err
	}
	onchainTool, err := onchain_actions.NewOnchainActionsTool(rpcClient)
	if err != nil {
		return nil, err
	}
	tokenTool, err := token_information.NewTokenInformationTool(rpcClient)

	// Register tools
	if err := tk.RegisterTool(transactionTool); err != nil {
		return nil, err
	}
	if err := tk.RegisterTool(walletTool); err != nil {
		return nil, err
	}
	if err := tk.RegisterTool(onchainTool); err != nil {
		return nil, err
	}
	if err := tk.RegisterTool(tokenTool); err != nil {
		return nil, err
	}

	return tk, nil
}
