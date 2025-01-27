package onchain_actions

import (
	"context"
	"os"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/ilkamo/jupiter-go/jupiter"
	"github.com/joho/godotenv"
	"github.com/soralabs/solana-toolkit/go/internal/pumpfun"
)

func TestNewOnchainActionsTool(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	rpcClient := rpc.New(os.Getenv("RPC_URL"))
	tool := NewOnchainActionsTool(rpcClient)

	if tool == nil {
		t.Fatal("Expected non-nil tool")
	}
	if tool.rpcClient != rpcClient {
		t.Error("Expected rpcClient to be set correctly")
	}
}

func TestSwap(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	rpcClient := rpc.New(os.Getenv("RPC_URL"))
	tool := NewOnchainActionsTool(rpcClient)

	wallet := solana.MustPrivateKeyFromBase58(os.Getenv("PRIVATE_KEY"))

	// First get a quote
	quoteReq := jupiter.GetQuoteParams{
		InputMint:  "So11111111111111111111111111111111111111112",  // SOL
		OutputMint: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", // USDC
		Amount:     1000000,                                        // 0.01 SOL
	}

	ctx := context.Background()

	// Test swap transaction creation
	tx, err := tool.Swap(ctx, quoteReq, wallet)
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
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	ctx := context.Background()
	rpcClient := rpc.New(os.Getenv("RPC_URL"))
	wallet := solana.MustPrivateKeyFromBase58(os.Getenv("PRIVATE_KEY"))
	tool := NewOnchainActionsTool(rpcClient)

	// Generate test wallets
	toWallet := solana.NewWallet()

	// Test transfer
	tx, err := tool.Transfer(ctx, wallet, toWallet.PublicKey(), 1000000)
	if err != nil {
		t.Errorf("Transfer error (expected during test): %v", err)
	}
	if tx != nil {
		t.Log("Successfully created transfer transaction")
	}
}

func TestCreateToken(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	ctx := context.Background()

	rpcClient := rpc.New(os.Getenv("RPC_URL"))
	wallet := solana.MustPrivateKeyFromBase58(os.Getenv("PRIVATE_KEY"))

	tool := NewOnchainActionsTool(rpcClient)

	mintWallet := solana.NewWallet()

	// Test token creation
	err = tool.CreateToken(ctx, CreateTokenParams{
		TokenInfo: pumpfun.CreateTokenInformation{
			Name:     "Test Token",
			Symbol:   "TEST",
			ImageURI: "https://example.com/image.png",
		},
		Mint:            mintWallet,
		UserPrivateKey:  wallet,
		BuyAmount:       0.1,
		SlippagePercent: 10,
	})
	if err != nil {
		t.Fatalf("Create token error (expected during test): %v", err)
	}

	t.Log("Successfully created token with mint", mintWallet.PublicKey().String())
}
