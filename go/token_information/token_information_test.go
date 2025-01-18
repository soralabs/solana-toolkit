package token_information

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/gagliardetto/solana-go/rpc"

	"github.com/joho/godotenv"
)

func TestTokenInformation(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("failed to load .env file: %v", err)
	}

	soraTokenAddress := "89nnWMkWeF9LSJvAWcN2JFQfeWdDk6diKEckeToEU1hE"

	tokenInformationTool := NewTokenInformationTool(
		rpc.New(os.Getenv("RPC_URL")),
	)

	input := TokenInformationInput{
		TokenAddress: soraTokenAddress,
	}

	jsonInput, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("failed to marshal input: %v", err)
	}

	tokenInformation, err := tokenInformationTool.Execute(context.Background(), jsonInput)
	if err != nil {
		t.Fatalf("failed to execute token information tool: %v", err)
	}

	var tokenInformationOutput TokenInformationOutput
	if err := json.Unmarshal(tokenInformation, &tokenInformationOutput); err != nil {
		t.Fatalf("failed to unmarshal token information output: %v", err)
	}

	t.Logf("token information: %+v", tokenInformationOutput)
}
