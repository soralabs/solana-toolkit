package tx_parser

import (
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// SwapType represents different DEX protocols
type SwapType string

const (
	SwapTypeJupiter    SwapType = "Jupiter"
	SwapTypeJupiterDCA SwapType = "JupiterDCA"
	SwapTypePumpFun    SwapType = "PumpFun"
	SwapTypeRaydium    SwapType = "Raydium"
	SwapTypeOrca       SwapType = "Orca"
	SwapTypeMeteora    SwapType = "Meteora"
	SwapTypeMoonshot   SwapType = "Moonshot"
	SwapTypeOKX        SwapType = "OKX"
	SwapTypeUnknown    SwapType = "Unknown"
)

// TokenInfo represents detailed information about a token
type TokenInfo struct {
	Mint     solana.PublicKey
	Amount   uint64
	Decimals uint8
}

// SwapInfo represents the parsed swap transaction data
type SwapInfo struct {
	Protocol   SwapType
	Signers    []solana.PublicKey
	Signatures []solana.Signature
	Timestamp  time.Time
	TokenIn    TokenInfo
	TokenOut   TokenInfo
}

// TransactionContext holds all the necessary context for parsing a transaction
type TransactionContext struct {
	Transaction  *solana.Transaction
	Meta         *rpc.TransactionMeta
	AccountKeys  []solana.PublicKey
	MintDecimals map[string]uint8 // map[mint_address]decimals
}

// SwapParser defines the interface for protocol-specific parsers
type SwapParser interface {
	// CanHandle checks if this parser can handle the given instruction
	CanHandle(instruction solana.CompiledInstruction, accountKeys []solana.PublicKey) bool

	// ParseInstruction processes a single instruction and returns swap information
	ParseInstruction(instruction solana.CompiledInstruction, instructionIndex int, ctx *TransactionContext) ([]*SwapInfo, error)
}

// GetMintDecimals returns the decimals for a given mint address
func (ctx *TransactionContext) GetMintDecimals(mint solana.PublicKey) uint8 {
	if decimals, exists := ctx.MintDecimals[mint.String()]; exists {
		return decimals
	}
	// Default to 9 decimals for unknown mints (common in Solana ecosystem)
	return 9
}

// ExtractMintDecimals processes the transaction to extract token decimal information
func (ctx *TransactionContext) ExtractMintDecimals() error {
	ctx.MintDecimals = make(map[string]uint8)

	// Process token balances from transaction metadata
	for _, balance := range ctx.Meta.PreTokenBalances {
		if !balance.Mint.IsZero() {
			ctx.MintDecimals[balance.Mint.String()] = uint8(balance.UiTokenAmount.Decimals)
		}
	}
	for _, balance := range ctx.Meta.PostTokenBalances {
		if !balance.Mint.IsZero() {
			ctx.MintDecimals[balance.Mint.String()] = uint8(balance.UiTokenAmount.Decimals)
		}
	}

	// Add Native SOL if not present
	if _, exists := ctx.MintDecimals[NATIVE_SOL_PROGRAM_ID.String()]; !exists {
		ctx.MintDecimals[NATIVE_SOL_PROGRAM_ID.String()] = 9
	}

	return nil
}
