package wallet_information

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
)

func TestNewWalletInformationTool(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	rpcClient := rpc.New(os.Getenv("RPC_URL"))
	tool := NewWalletInformationTool(rpcClient)

	if tool == nil {
		t.Fatal("Expected non-nil tool")
	}
	if tool.rpcClient != rpcClient {
		t.Error("Expected rpcClient to be set correctly")
	}
}

func TestGetWalletInformation(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	ctx := context.Background()
	rpcClient := rpc.New(os.Getenv("RPC_URL"))
	tool := NewWalletInformationTool(rpcClient)

	// Use a testing wallet
	testWallet := os.Getenv("TEST_WALLET")

	input := WalletInformationInput{
		Wallet: testWallet,
	}

	// Convert input to JSON
	params, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal input: %v", err)
	}

	// Execute the tool
	result, err := tool.Execute(ctx, params)
	if err != nil {
		t.Fatalf("Failed to execute tool: %v", err)
	}

	// Parse the result
	var output WalletInformationOutput
	if err := json.Unmarshal(result, &output); err != nil {
		t.Fatalf("Failed to unmarshal output: %v", err)
	}

	// Verify basic output structure
	if output.Wallet != testWallet {
		t.Errorf("Expected wallet %s, got %s", testWallet, output.Wallet)
	}

	// Log some interesting information
	t.Logf("Wallet Information for %s:", testWallet)
	t.Logf("SOL Balance: %s", output.SolBalance)
	t.Logf("Total Value: %.2f", output.TotalValue)
	t.Logf("Token Count: %d", output.TokenNum)
	t.Logf("PNL (24h): %.2f", output.PNL1d*100)
	t.Logf("Win Rate: %.2f%%", output.Winrate*100)
}
