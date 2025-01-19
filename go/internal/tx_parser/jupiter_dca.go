package tx_parser

import (
	"bytes"
	"fmt"
	"time"

	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
)

// JupiterDCAParser handles parsing Jupiter DCA protocol swaps
type JupiterDCAParser struct{}

// NewJupiterDCAParser creates a new Jupiter DCA parser instance
func NewJupiterDCAParser() *JupiterDCAParser {
	return &JupiterDCAParser{}
}

// CanHandle checks if this parser can handle the given instruction
func (p *JupiterDCAParser) CanHandle(instruction solana.CompiledInstruction, accountKeys []solana.PublicKey) bool {
	return accountKeys[instruction.ProgramIDIndex].Equals(JUPITER_DCA_PROGRAM_ID)
}

// ParseInstruction processes the Jupiter DCA instruction and returns swap information
func (p *JupiterDCAParser) ParseInstruction(instruction solana.CompiledInstruction, instructionIndex int, ctx *TransactionContext) ([]*SwapInfo, error) {
	var swaps []*SwapInfo

	// Process events in inner instructions
	for _, innerSet := range ctx.Meta.InnerInstructions {
		if innerSet.Index == uint16(instructionIndex) {
			for _, innerInstr := range innerSet.Instructions {
				if event, err := p.parseJupiterEvent(innerInstr, ctx.AccountKeys); err == nil {
					swap, err := p.buildSwapInfo(event, ctx)
					if err != nil {
						continue
					}
					swaps = append(swaps, swap)
				}
			}
		}
	}

	if len(swaps) == 0 {
		return nil, fmt.Errorf("no valid DCA swaps found")
	}

	return swaps, nil
}

// JupiterDCAEvent represents a single swap event in the Jupiter DCA protocol
type JupiterDCAEvent struct {
	UserKey          solana.PublicKey
	DCAKey           solana.PublicKey
	InDeposited      uint64
	InputMint        solana.PublicKey
	OutputMint       solana.PublicKey
	CycleFrequency   int64
	InAmountPerCycle uint64
	CreatedAt        int64
}

// parseJupiterEvent decodes a single Jupiter swap event
func (p *JupiterDCAParser) parseJupiterEvent(instruction solana.CompiledInstruction, accountKeys []solana.PublicKey) (*JupiterDCAEvent, error) {
	decodedBytes, err := base58.Decode(instruction.Data.String())
	if err != nil {
		return nil, fmt.Errorf("failed to decode instruction data: %w", err)
	}

	// Check event discriminator (same as regular Jupiter)
	if len(decodedBytes) < 16 || !bytes.Equal(decodedBytes[:16], JUPITER_DCA_EVENT_DISCRIMINATOR[:]) {
		return nil, fmt.Errorf("invalid Jupiter event discriminator")
	}

	// Decode event data using binary decoder
	decoder := ag_binary.NewBorshDecoder(decodedBytes[16:])
	var event JupiterDCAEvent
	if err := decoder.Decode(&event); err != nil {
		return nil, fmt.Errorf("failed to decode Jupiter DCA event: %w", err)
	}

	return &event, nil
}

// buildSwapInfo creates the final SwapInfo from the Jupiter DCA event
func (p *JupiterDCAParser) buildSwapInfo(event *JupiterDCAEvent, ctx *TransactionContext) (*SwapInfo, error) {
	if event.InAmountPerCycle == 0 {
		return nil, fmt.Errorf("invalid DCA amount")
	}

	swapInfo := &SwapInfo{
		Protocol:  SwapTypeJupiterDCA,
		Timestamp: time.Unix(event.CreatedAt, 0),
		TokenIn: TokenInfo{
			Mint:     event.InputMint,
			Amount:   event.InDeposited,
			Decimals: ctx.GetMintDecimals(event.InputMint),
		},
		TokenOut: TokenInfo{
			Mint:     event.OutputMint,
			Decimals: ctx.GetMintDecimals(event.OutputMint),
		},
	}

	return swapInfo, nil
}
