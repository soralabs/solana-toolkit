package tx_parser

import (
	"bytes"
	"fmt"
	"time"

	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
)

// PumpFunParser handles parsing PumpFun protocol swaps
type PumpFunParser struct{}

// NewPumpFunParser creates a new PumpFun parser instance
func NewPumpFunParser() *PumpFunParser {
	return &PumpFunParser{}
}

// PumpFunTradeEvent represents a single trade event in the PumpFun protocol
type PumpFunTradeEvent struct {
	Mint                 solana.PublicKey
	SolAmount            uint64
	TokenAmount          uint64
	IsBuy                bool
	User                 solana.PublicKey
	Timestamp            int64
	VirtualSolReserves   uint64
	VirtualTokenReserves uint64
}

// CanHandle checks if this parser can handle the given instruction
func (p *PumpFunParser) CanHandle(instruction solana.CompiledInstruction, accountKeys []solana.PublicKey) bool {
	programID := accountKeys[instruction.ProgramIDIndex]
	return programID.Equals(PUMP_FUN_PROGRAM_ID)
}

// ParseInstruction processes the PumpFun instruction and returns swap information
func (p *PumpFunParser) ParseInstruction(instruction solana.CompiledInstruction, instructionIndex int, ctx *TransactionContext) ([]*SwapInfo, error) {
	var swaps []*SwapInfo

	// Process all trade events in inner instructions
	for _, innerSet := range ctx.Meta.InnerInstructions {
		if innerSet.Index == uint16(instructionIndex) {
			for _, innerInstr := range innerSet.Instructions {
				event, err := p.parsePumpFunEvent(innerInstr, ctx.AccountKeys)
				if err != nil {
					continue // Skip invalid events but continue processing
				}

				swap, err := p.buildSwapInfo(event, ctx)
				if err != nil {
					continue
				}
				swaps = append(swaps, swap)
			}
		}
	}

	if len(swaps) == 0 {
		return nil, fmt.Errorf("no valid PumpFun swaps found")
	}

	return swaps, nil
}

// parsePumpFunEvent decodes a single PumpFun trade event
func (p *PumpFunParser) parsePumpFunEvent(instruction solana.CompiledInstruction, accountKeys []solana.PublicKey) (*PumpFunTradeEvent, error) {
	decodedBytes, err := base58.Decode(instruction.Data.String())
	if err != nil {
		return nil, fmt.Errorf("failed to decode instruction data: %w", err)
	}

	// Check event discriminator
	if len(decodedBytes) < 16 || !bytes.Equal(decodedBytes[:16], PUMPFUN_TRADE_EVENT_DISCRIMINATOR[:]) {
		return nil, fmt.Errorf("invalid PumpFun event discriminator")
	}

	// Decode event data
	decoder := ag_binary.NewBorshDecoder(decodedBytes[16:])
	var event PumpFunTradeEvent
	if err := decoder.Decode(&event); err != nil {
		return nil, fmt.Errorf("failed to decode PumpFun event: %w", err)
	}

	return &event, nil
}

// buildSwapInfo creates the final SwapInfo from the PumpFun event
func (p *PumpFunParser) buildSwapInfo(event *PumpFunTradeEvent, ctx *TransactionContext) (*SwapInfo, error) {
	if event.SolAmount == 0 || event.TokenAmount == 0 {
		return nil, fmt.Errorf("invalid amounts in PumpFun event")
	}

	swapInfo := &SwapInfo{
		Protocol:  SwapTypePumpFun,
		Timestamp: time.Unix(event.Timestamp, 0),
	}

	// Get token decimals
	tokenDecimals := ctx.GetMintDecimals(event.Mint)
	solDecimals := ctx.GetMintDecimals(NATIVE_SOL_PROGRAM_ID)

	if event.IsBuy {
		// Buying tokens with SOL
		swapInfo.TokenIn = TokenInfo{
			Mint:     NATIVE_SOL_PROGRAM_ID,
			Amount:   event.SolAmount,
			Decimals: solDecimals,
		}
		swapInfo.TokenOut = TokenInfo{
			Mint:     event.Mint,
			Amount:   event.TokenAmount,
			Decimals: tokenDecimals,
		}
	} else {
		// Selling tokens for SOL
		swapInfo.TokenIn = TokenInfo{
			Mint:     event.Mint,
			Amount:   event.TokenAmount,
			Decimals: tokenDecimals,
		}
		swapInfo.TokenOut = TokenInfo{
			Mint:     NATIVE_SOL_PROGRAM_ID,
			Amount:   event.SolAmount,
			Decimals: solDecimals,
		}
	}

	return swapInfo, nil
}
