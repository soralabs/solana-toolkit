package tx_parser

import "github.com/gagliardetto/solana-go"

var (
	// Program IDs
	JUPITER_PROGRAM_ID     = solana.MustPublicKeyFromBase58("JUP6LkbZbjS1jKKwapdHNy74zcZ3tLUZoi5QNyVTaV4")
	JUPITER_DCA_PROGRAM_ID = solana.MustPublicKeyFromBase58("DCA265Vj8a9CEuX1eb1LWRnDT7uK6q1xMipnNyatn23M")

	PUMP_FUN_PROGRAM_ID = solana.MustPublicKeyFromBase58("6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P")

	RAYDIUM_V4_PROGRAM_ID                     = solana.MustPublicKeyFromBase58("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8")
	RAYDIUM_AMM_PROGRAM_ID                    = solana.MustPublicKeyFromBase58("routeUGWgWzqBWFcrCfv8tritsqukccJPu3q5GPP3xS")
	RAYDIUM_CPMM_PROGRAM_ID                   = solana.MustPublicKeyFromBase58("CPMMoo8L3F4NbTegBCKVNunggL7H1ZpdTHKxQB5qKP1C")
	RAYDIUM_CONCENTRATED_LIQUIDITY_PROGRAM_ID = solana.MustPublicKeyFromBase58("CAMMCzo5YL8w4VFF8KVHrK22GGUsp5VTaW7grrKgrWqK")

	METEORA_PROGRAM_ID       = solana.MustPublicKeyFromBase58("LBUZKhRxPF3XUpBCjp4YzTKgLccjZhTSDM9YuVaPwxo")
	METEORA_POOLS_PROGRAM_ID = solana.MustPublicKeyFromBase58("Eo7WjKq67rjJQSZxS6z3YkapzY3eMj6Xy8X5EQVn5UaB")

	MOONSHOT_PROGRAM_ID = solana.MustPublicKeyFromBase58("MoonCVVNZFSYkqNXP6bxHLPL6QQJiMagDL3qcqUQTrG")

	ORCA_PROGRAM_ID = solana.MustPublicKeyFromBase58("whirLbMiicVdio4qvUfM5KAg6Ct8VwpYzGff3uctyCc")

	OKX_PROGRAM_ID = solana.MustPublicKeyFromBase58("6m2CDdhRgxpH4WjvdzxAYbGxwdGUz5MziiL5jek2kBma")

	// Token Program IDs
	NATIVE_SOL_PROGRAM_ID = solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")
)

// Event Discriminators
var (
	JUPITER_ROUTE_EVENT_DISCRIMINATOR = [16]byte{228, 69, 165, 46, 81, 203, 154, 29, 64, 198, 205, 232, 38, 8, 113, 226}
	PUMPFUN_TRADE_EVENT_DISCRIMINATOR = [16]byte{228, 69, 165, 46, 81, 203, 154, 29, 189, 219, 127, 211, 78, 230, 97, 238}
	JUPITER_DCA_EVENT_DISCRIMINATOR   = [16]byte{
		0xe4, 0x45, 0xa5, 0x2e, 0x51, 0xcb, 0x9a, 0x1d,
		0xa6, 0xac, 0x61, 0x09, 0x4d, 0x4c, 0xbd, 0x6d,
	}
)

// Helper functions for program ID checks
func isRaydiumProgram(programID solana.PublicKey) bool {
	return programID.Equals(RAYDIUM_V4_PROGRAM_ID) ||
		programID.Equals(RAYDIUM_CPMM_PROGRAM_ID) ||
		programID.Equals(RAYDIUM_AMM_PROGRAM_ID) ||
		programID.Equals(RAYDIUM_CONCENTRATED_LIQUIDITY_PROGRAM_ID)
}

func isMeteoraProgramID(programID solana.PublicKey) bool {
	return programID.Equals(METEORA_PROGRAM_ID) ||
		programID.Equals(METEORA_POOLS_PROGRAM_ID)
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
