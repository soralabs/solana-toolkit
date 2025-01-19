package tx_parser

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
)

// MoonshotParser handles parsing Moonshot protocol swaps
type MoonshotParser struct{}

// NewMoonshotParser creates a new Moonshot parser instance
func NewMoonshotParser() *MoonshotParser {
	return &MoonshotParser{}
}

// TradeType represents buy/sell direction
type TradeType int

const (
	TradeTypeBuy TradeType = iota
	TradeTypeSell
)

// Trade instruction discriminators
var (
	MOONSHOT_BUY_INSTRUCTION  = [8]byte{102, 6, 61, 18, 1, 218, 235, 234}
	MOONSHOT_SELL_INSTRUCTION = [8]byte{51, 230, 133, 164, 1, 127, 131, 173}
)

// MoonshotTradeData represents a decoded trade instruction
type MoonshotTradeData struct {
	TradeType        TradeType
	TokenMint        solana.PublicKey
	TokenAmount      uint64
	CollateralAmount uint64
}

// CanHandle checks if this parser can handle the given instruction
func (p *MoonshotParser) CanHandle(instruction solana.CompiledInstruction, accountKeys []solana.PublicKey) bool {
	if !accountKeys[instruction.ProgramIDIndex].Equals(MOONSHOT_PROGRAM_ID) {
		return false
	}

	// Verify instruction format
	if len(instruction.Data) != 33 || len(instruction.Accounts) != 11 {
		return false
	}

	// Check instruction discriminator
	decodedBytes, err := base58.Decode(instruction.Data.String())
	if err != nil || len(decodedBytes) < 8 {
		return false
	}

	discriminator := decodedBytes[:8]
	return bytes.Equal(discriminator, MOONSHOT_BUY_INSTRUCTION[:]) ||
		bytes.Equal(discriminator, MOONSHOT_SELL_INSTRUCTION[:])
}

// ParseInstruction processes the Moonshot instruction and returns swap information
func (p *MoonshotParser) ParseInstruction(instruction solana.CompiledInstruction, instructionIndex int, ctx *TransactionContext) ([]*SwapInfo, error) {
	var swaps []*SwapInfo

	// Decode the instruction first
	tradeData, err := p.decodeMoonshotInstruction(instruction, ctx.AccountKeys)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Moonshot instruction: %w", err)
	}

	// Get token and SOL balance changes
	tokenAmount, solAmount, err := p.getBalanceChanges(tradeData.TokenMint, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance changes: %w", err)
	}

	// Build swap info
	swapInfo, err := p.buildSwapInfo(tradeData, tokenAmount, solAmount, ctx)
	if err == nil {
		swaps = append(swaps, swapInfo)
	}

	if len(swaps) == 0 {
		return nil, fmt.Errorf("no valid Moonshot swaps found")
	}

	return swaps, nil
}

// decodeMoonshotInstruction extracts trade information from the instruction
func (p *MoonshotParser) decodeMoonshotInstruction(instruction solana.CompiledInstruction, accountKeys []solana.PublicKey) (*MoonshotTradeData, error) {
	decodedBytes, err := base58.Decode(instruction.Data.String())
	if err != nil {
		return nil, fmt.Errorf("failed to decode instruction data: %w", err)
	}

	// Check instruction discriminator
	discriminator := decodedBytes[:8]
	var tradeType TradeType

	switch {
	case bytes.Equal(discriminator, MOONSHOT_BUY_INSTRUCTION[:]):
		tradeType = TradeTypeBuy
	case bytes.Equal(discriminator, MOONSHOT_SELL_INSTRUCTION[:]):
		tradeType = TradeTypeSell
	default:
		return nil, fmt.Errorf("invalid Moonshot instruction discriminator")
	}

	// Get token mint from instruction accounts
	tokenMint := accountKeys[instruction.Accounts[6]]

	return &MoonshotTradeData{
		TradeType: tradeType,
		TokenMint: tokenMint,
	}, nil
}

// getBalanceChanges calculates token and SOL balance changes
func (p *MoonshotParser) getBalanceChanges(tokenMint solana.PublicKey, ctx *TransactionContext) (tokenAmount, solAmount uint64, err error) {
	// Get signer's public key
	signer := ctx.Transaction.Message.AccountKeys[0]

	// Get token balance change
	tokenChange, err := p.getTokenBalanceChange(tokenMint, signer, ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get token balance change: %w", err)
	}

	// Get SOL balance change
	solChange, err := p.getNativeSolBalanceChange(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get SOL balance change: %w", err)
	}

	return uint64(abs(tokenChange)), uint64(abs(solChange)), nil
}

// getTokenBalanceChange calculates the token balance change for the signer
func (p *MoonshotParser) getTokenBalanceChange(mint, owner solana.PublicKey, ctx *TransactionContext) (int64, error) {
	var preAmount, postAmount int64

	// Find pre-balance
	for _, balance := range ctx.Meta.PreTokenBalances {
		if balance.Mint.Equals(mint) && balance.Owner.Equals(owner) {
			preAmount, _ = strconv.ParseInt(balance.UiTokenAmount.Amount, 10, 64)
			break
		}
	}

	// Find post-balance
	for _, balance := range ctx.Meta.PostTokenBalances {
		if balance.Mint.Equals(mint) && balance.Owner.Equals(owner) {
			postAmount, _ = strconv.ParseInt(balance.UiTokenAmount.Amount, 10, 64)
			break
		}
	}

	return postAmount - preAmount, nil
}

// getNativeSolBalanceChange calculates the SOL balance change for the signer
func (p *MoonshotParser) getNativeSolBalanceChange(ctx *TransactionContext) (int64, error) {
	if len(ctx.Meta.PostBalances) == 0 || len(ctx.Meta.PreBalances) == 0 {
		return 0, fmt.Errorf("insufficient balance information")
	}

	// Calculate SOL balance change for the first account (signer)
	return int64(ctx.Meta.PostBalances[0]) - int64(ctx.Meta.PreBalances[0]), nil
}

// buildSwapInfo creates the final SwapInfo
func (p *MoonshotParser) buildSwapInfo(tradeData *MoonshotTradeData, tokenAmount, solAmount uint64, ctx *TransactionContext) (*SwapInfo, error) {
	swapInfo := &SwapInfo{
		Protocol: SwapTypeMoonshot,
	}

	// Get decimals for both tokens
	tokenDecimals := ctx.GetMintDecimals(tradeData.TokenMint)
	solDecimals := ctx.GetMintDecimals(NATIVE_SOL_PROGRAM_ID)

	if tradeData.TradeType == TradeTypeBuy {
		// Buying tokens with SOL
		swapInfo.TokenIn = TokenInfo{
			Mint:     NATIVE_SOL_PROGRAM_ID,
			Amount:   solAmount,
			Decimals: solDecimals,
		}
		swapInfo.TokenOut = TokenInfo{
			Mint:     tradeData.TokenMint,
			Amount:   tokenAmount,
			Decimals: tokenDecimals,
		}
	} else {
		// Selling tokens for SOL
		swapInfo.TokenIn = TokenInfo{
			Mint:     tradeData.TokenMint,
			Amount:   tokenAmount,
			Decimals: tokenDecimals,
		}
		swapInfo.TokenOut = TokenInfo{
			Mint:     NATIVE_SOL_PROGRAM_ID,
			Amount:   solAmount,
			Decimals: solDecimals,
		}
	}

	return swapInfo, nil
}
