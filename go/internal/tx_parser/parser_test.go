package tx_parser

import (
	"context"
	"os"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
)

func setupTest(t *testing.T) *rpc.Client {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("failed to load .env file: %v", err)
	}
	return rpc.New(os.Getenv("RPC_URL"))
}

func testTransaction(t *testing.T, signature string) []*SwapInfo {
	rpcClient := setupTest(t)
	maxSupportedTxVersion := uint64(0)

	tx := solana.MustSignatureFromBase58(signature)
	txResult, err := rpcClient.GetTransaction(context.Background(), tx, &rpc.GetTransactionOpts{
		MaxSupportedTransactionVersion: &maxSupportedTxVersion,
	})
	if err != nil {
		t.Fatalf("failed to get transaction: %v", err)
	}

	parser, err := New(txResult)
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}

	swapInfo, err := parser.ParseTransaction()
	if err != nil {
		t.Fatalf("failed to parse transaction: %v", err)
	}

	// Log swap details
	for i, swap := range swapInfo {
		t.Logf("Swap %d:", i+1)
		t.Logf("  Protocol: %s", swap.Protocol)
		t.Logf("  Token In: %s (Amount: %d, Decimals: %d)",
			swap.TokenIn.Mint.String(),
			swap.TokenIn.Amount,
			swap.TokenIn.Decimals)
		t.Logf("  Token Out: %s (Amount: %d, Decimals: %d)",
			swap.TokenOut.Mint.String(),
			swap.TokenOut.Amount,
			swap.TokenOut.Decimals)
	}

	return swapInfo
}

func TestJupiterDCAParser(t *testing.T) {
	swaps := testTransaction(t, "3dKEeJANexZP1UkZbJG2gCJ8FQB67Rj4AGeRdKdMoJ2tU9RG1NoQpDcmfvQJqdBgyy2NBajvCYjN2rPMSNfyep3p")
	if len(swaps) == 0 {
		t.Fatal("expected at least one swap")
	}
	if swaps[0].Protocol != SwapTypeJupiterDCA {
		t.Errorf("expected Jupiter DCA swap, got %s", swaps[0].Protocol)
	}
}

func TestJupiterParser(t *testing.T) {
	// Direct swap
	t.Run("Direct Swap", func(t *testing.T) {
		swaps := testTransaction(t, "4675XfT8skvgXBCuzsixz2ZeEJ7ZQvunnoaihrbb9w7EixiSvVq8udW1jf1Bfhk6qpjZZfA8i5myP2ThtX2kwwWR")
		if len(swaps) == 0 {
			t.Fatal("expected at least one swap")
		}
		if swaps[0].Protocol != SwapTypeJupiter {
			t.Errorf("expected Jupiter swap, got %s", swaps[0].Protocol)
		}
	})

	// Multi-hop swap
	t.Run("Multi-hop Swap", func(t *testing.T) {
		swaps := testTransaction(t, "4675XfT8skvgXBCuzsixz2ZeEJ7ZQvunnoaihrbb9w7EixiSvVq8udW1jf1Bfhk6qpjZZfA8i5myP2ThtX2kwwWR")
		if len(swaps) == 0 {
			t.Fatal("expected at least one swap")
		}
		if swaps[0].Protocol != SwapTypeJupiter {
			t.Errorf("expected Jupiter swap, got %s", swaps[0].Protocol)
		}
	})
}

func TestMoonshotParser(t *testing.T) {
	// Buy transaction
	t.Run("Buy", func(t *testing.T) {
		swaps := testTransaction(t, "4t2XxKesjUvLSSEt1ioYuTQqQhvWnzxqmvxA3pdUxcaC3fWwXypSACj2Kd7pAGuFEsQxacMzf3KLXzaTrvKDKqb8")
		if len(swaps) == 0 {
			t.Fatal("expected at least one swap")
		}
		if swaps[0].Protocol != SwapTypeMoonshot {
			t.Errorf("expected Moonshot swap, got %s", swaps[0].Protocol)
		}
	})

	// Sell transaction
	t.Run("Sell", func(t *testing.T) {
		swaps := testTransaction(t, "2RtozpqJhH3rH5jd7oLuZ3yS5se9DvHAjQDx6Le7ZdExFfhq6c3M7L5xsu4NEzfYKHjAeLB7zvjtRcy6r5j9WtTt")
		if len(swaps) == 0 {
			t.Fatal("expected at least one swap")
		}
		if swaps[0].Protocol != SwapTypeMoonshot {
			t.Errorf("expected Moonshot swap, got %s", swaps[0].Protocol)
		}
	})
}

func TestRaydiumParser(t *testing.T) {
	// Standard swap
	t.Run("Standard Swap", func(t *testing.T) {
		swaps := testTransaction(t, "5W7Kqwr2cs7mGBaQCzw23NL2zYHRdznfvXMZtSCsN7h8aib6Ra2mCuchKr6ubvEcfMa3zF7uioR4AcTkBYG4sVzA")
		if len(swaps) == 0 {
			t.Fatal("expected at least one swap")
		}
		if swaps[0].Protocol != SwapTypeRaydium {
			t.Errorf("expected Raydium swap, got %s", swaps[0].Protocol)
		}
	})
}

func TestOrcaParser(t *testing.T) {
	// Whirlpool swap
	t.Run("Whirlpool Swap", func(t *testing.T) {
		swaps := testTransaction(t, "3bAv6WNXKtiRfUvqB1UjQ2twuDMpGvgnzgcaQwqi6KwEMY2TgHyYKJZ1oUWkrgjZ9xWFsCZfhsXdKFQJEh3CB8ET")
		if len(swaps) == 0 {
			t.Fatal("expected at least one swap")
		}
		if swaps[len(swaps)-1].Protocol != SwapTypeRaydium {
			t.Errorf("expected Orca swap, got %s", swaps[len(swaps)-1].Protocol)
		}
	})
}

func TestMeteoraParser(t *testing.T) {
	swaps := testTransaction(t, "V3Zt8dDKdpZ8MJ8CuvZYnnawu68fm6Edn3z48VBsgf2t2mkWhkpcpQWB2iorZx1bod6dZcbopg17QJFe61ye24k")
	if len(swaps) == 0 {
		t.Fatal("expected at least one swap")
	}
	if swaps[0].Protocol != SwapTypeMeteora {
		t.Errorf("expected Meteora swap, got %s", swaps[0].Protocol)
	}
}

func TestOKXParser(t *testing.T) {
	swaps := testTransaction(t, "27wChRDnfQwZuG1Q7YuM9VQ8yJmXZvjfgB25fBqex1WhE39GZ4ZdcABJ3u39mH2uJbUGb9nTYctQVZsKjJaUt4Jt")
	if len(swaps) == 0 {
		t.Fatal("expected at least one swap")
	}
	if swaps[0].Protocol != SwapTypeOKX {
		t.Errorf("expected OKX swap, got %s", swaps[0].Protocol)
	}
}
