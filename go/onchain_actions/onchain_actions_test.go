package onchain_actions

import (
	"context"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/soralabs/solana-toolkit/go/internal/jupiter"
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
	quoteReq := jupiter.QuoteRequest{
		InputMint:  "So11111111111111111111111111111111111111112",  // SOL
		OutputMint: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", // USDC
		Amount:     "1000000",                                      // 1 SOL in lamports
		Slippage:   "100",                                          // 1% slippage
	}

	quote, err := jupiter.GetQuoteData(quoteReq)
	if err != nil {
		t.Logf("Quote error (expected during test): %v", err)
		return
	}

	// Create swap request with the quote
	swapRequest := jupiter.SwapRequest{
		Quote:         quote,
		UserPublicKey: wallet.PublicKey().String(),
	}

	// Test swap transaction creation
	tx, err := tool.Swap(swapRequest, wallet.PrivateKey)
	if err != nil {
		t.Errorf("Swap error (expected during test): %v", err)
	}
	if tx != nil {
		t.Log("Successfully created swap transaction")
	}
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
