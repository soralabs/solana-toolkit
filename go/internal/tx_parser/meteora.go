package tx_parser

import (
	"encoding/binary"
	"fmt"

	"github.com/gagliardetto/solana-go"
)

// MeteoraParser handles parsing Meteora protocol swaps
type MeteoraParser struct {
	seenInstructionPairs map[string]bool
}

// NewMeteoraParser creates a new Meteora parser instance
func NewMeteoraParser() *MeteoraParser {
	return &MeteoraParser{
		seenInstructionPairs: make(map[string]bool),
	}
}

var METEORA_SWAP_DISCRIMINATOR = []byte{0xf8, 0xc6, 0x9e, 0x91, 0xe1, 0x75, 0x87, 0xc8}

// CanHandle checks if this parser can handle the given instruction
func (p *MeteoraParser) CanHandle(instruction solana.CompiledInstruction, accountKeys []solana.PublicKey) bool {
	programID := accountKeys[instruction.ProgramIDIndex]
	return programID.Equals(METEORA_PROGRAM_ID) || programID.Equals(METEORA_POOLS_PROGRAM_ID)
}

// TransferCheckData represents a token transfer check instruction
type TransferCheckData struct {
	Amount    uint64
	Decimals  uint8
	Mint      solana.PublicKey
	Authority solana.PublicKey
}

// ParseInstruction processes the Meteora instruction and returns swap information
func (p *MeteoraParser) ParseInstruction(instruction solana.CompiledInstruction, instructionIndex int, ctx *TransactionContext) ([]*SwapInfo, error) {
	// Process transfers in each group of inner instructions
	var swaps []*SwapInfo

	for _, innerSet := range ctx.Meta.InnerInstructions {
		if innerSet.Index == uint16(instructionIndex) {
			var lastTransferIndex = -1
			var transfers []TransferCheckData

			for i, innerInstr := range innerSet.Instructions {
				// Orca only uses regular token transfers
				if isMeteoraTransferChecked(innerInstr, ctx.AccountKeys) {
					// If this is not consecutive with the last transfer, reset
					if lastTransferIndex != -1 && i != lastTransferIndex+1 {
						transfers = nil
					}
					lastTransferIndex = i

					transfer, err := p.parseTransferCheck(innerInstr, ctx)
					if err != nil {
						transfers = nil
						continue
					}
					transfers = append(transfers, transfer)

					// When we have a pair of consecutive transfers, build a swap
					if len(transfers) == 2 {
						if p.seenInstructionPairs[innerInstr.Data.String()+innerSet.Instructions[i-1].Data.String()] {
							transfers = nil
							continue
						}
						p.seenInstructionPairs[innerInstr.Data.String()+innerSet.Instructions[i-1].Data.String()] = true

						swap, err := p.buildSwapInfo(transfers[0], transfers[1], ctx)
						if err != nil {
							transfers = nil
							continue
						}
						swaps = append(swaps, swap)
						transfers = nil
					}
				} else {
					// Reset transfers if we see a non-transfer instruction
					transfers = nil
					lastTransferIndex = -1
				}
			}
		}
	}

	if len(swaps) == 0 {
		return nil, fmt.Errorf("no valid Meteora swaps found")
	}

	return swaps, nil
}

// isTransferChecked checks if the instruction is a token transfer check
func isMeteoraTransferChecked(instruction solana.CompiledInstruction, accountKeys []solana.PublicKey) bool {
	if len(instruction.Accounts) < 4 || len(instruction.Data) < 9 {
		return false
	}

	progID := accountKeys[instruction.ProgramIDIndex]
	if !progID.Equals(solana.TokenProgramID) && !progID.Equals(solana.Token2022ProgramID) {
		return false
	}

	return instruction.Data[0] == 12 // TransferChecked instruction
}

// parseTransferCheck extracts transfer data from the instruction
func (p *MeteoraParser) parseTransferCheck(instruction solana.CompiledInstruction, ctx *TransactionContext) (TransferCheckData, error) {
	amount := binary.LittleEndian.Uint64(instruction.Data[1:9])
	decimals := instruction.Data[9]

	return TransferCheckData{
		Amount:    amount,
		Decimals:  decimals,
		Mint:      ctx.AccountKeys[instruction.Accounts[1]], // Token mint
		Authority: ctx.AccountKeys[instruction.Accounts[3]], // Authority
	}, nil
}

// buildSwapInfo creates the final SwapInfo from the transfer data
func (p *MeteoraParser) buildSwapInfo(transfer1, transfer2 TransferCheckData, ctx *TransactionContext) (*SwapInfo, error) {
	swapInfo := &SwapInfo{
		Protocol: SwapTypeMeteora,
		TokenIn: TokenInfo{
			Mint:     transfer1.Mint,
			Amount:   transfer1.Amount,
			Decimals: transfer1.Decimals,
		},
		TokenOut: TokenInfo{
			Mint:     transfer2.Mint,
			Amount:   transfer2.Amount,
			Decimals: transfer2.Decimals,
		},
	}

	return swapInfo, nil
}
