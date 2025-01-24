package onchain_actions

import (
	"context"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/ilkamo/jupiter-go/jupiter"
	"github.com/soralabs/solana-toolkit/go/internal/pumpfun"
)

func TestNewOnchainActionsTool(t *testing.T) {
	rpcClient := rpc.New("https://api.mainnet-beta.solana.com")
	tool := NewOnchainActionsTool(rpcClient)

	if tool == nil {
		t.Fatal("Expected non-nil tool")
	}
	if tool.rpcClient != rpcClient {
		t.Error("Expected rpcClient to be set correctly")
	}
}

func TestSwap(t *testing.T) {
	rpcClient := rpc.New("https://api.mainnet-beta.solana.com")
	tool := NewOnchainActionsTool(rpcClient)

	// Generate a test wallet
	wallet := solana.NewWallet()

	// First get a quote
	quoteReq := jupiter.GetQuoteParams{
		InputMint:  "So11111111111111111111111111111111111111112",  // SOL
		OutputMint: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", // USDC
		Amount:     1000000,                                        // 0.01 SOL
	}

	ctx := context.Background()

	// Test swap transaction creation
	tx, err := tool.Swap(ctx, quoteReq, wallet.PrivateKey)
	if err != nil {
		t.Fatalf("Swap error (expected during test): %v", err)
	}
	if tx != nil {
		t.Log("Successfully created swap transaction")
	}

	sentSwapTx, err := tool.SendSwapTransaction(ctx, tx)
	if err != nil {
		t.Fatalf("Send swap transaction error (expected during test): %v", err)
	}

	t.Logf("Successfully sent swap transaction: %s", sentSwapTx.String())
}

func TestTransfer(t *testing.T) {
	ctx := context.Background()
	rpcClient := rpc.New("https://api.mainnet-beta.solana.com")
	tool := NewOnchainActionsTool(rpcClient)

	// Generate test wallets
	fromWallet := solana.NewWallet()
	toWallet := solana.NewWallet()

	// Test transfer
	tx, err := tool.Transfer(ctx, fromWallet.PrivateKey, toWallet.PublicKey(), 1000000)
	if err != nil {
		t.Errorf("Transfer error (expected during test): %v", err)
	}
	if tx != nil {
		t.Log("Successfully created transfer transaction")
	}
}

func TestCreateToken(t *testing.T) {
	ctx := context.Background()
	rpcClient := rpc.New("https://api.mainnet-beta.solana.com")
	tool := NewOnchainActionsTool(rpcClient)

	// Generate test wallet
	wallet := solana.NewWallet()
	mintWallet := solana.NewWallet()

	// Test token creation
	err := tool.CreateToken(ctx, CreateTokenParams{
		TokenInfo: pumpfun.CreateTokenInformation{
			Name:     "Test Token",
			Symbol:   "TEST",
			ImageURI: "https://example.com/image.png",
		},
		Mint:           mintWallet,
		UserPrivateKey: wallet.PrivateKey,
		BuyAmount:      1000000000, // 1 SOL
	})
	if err != nil {
		t.Errorf("Create token error (expected during test): %v", err)
	}

	t.Log("Successfully created token")
}
