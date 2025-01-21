package tx_parser

import (
	"encoding/binary"
	"fmt"

	"github.com/gagliardetto/solana-go"
)

// OrcaParser handles parsing Orca protocol swaps
type OrcaParser struct {
	seenInstructionPairs map[string]bool
}

// NewOrcaParser creates a new Orca parser instance
func NewOrcaParser() *OrcaParser {
	return &OrcaParser{
		seenInstructionPairs: make(map[string]bool),
	}
}

// CanHandle checks if this parser can handle the given instruction
func (p *OrcaParser) CanHandle(instruction solana.CompiledInstruction, accountKeys []solana.PublicKey) bool {
	return accountKeys[instruction.ProgramIDIndex].Equals(ORCA_PROGRAM_ID)
}

// ParseInstruction processes the Orca instruction and returns swap information
func (p *OrcaParser) ParseInstruction(instruction solana.CompiledInstruction, instructionIndex int, ctx *TransactionContext) ([]*SwapInfo, error) {
	var swaps []*SwapInfo

	// Process transfers in each group of inner instructions
	for _, innerSet := range ctx.Meta.InnerInstructions {
		if innerSet.Index == uint16(instructionIndex) {
			var lastTransferIndex = -1
			var currentTransfers []TokenInfo
			for i, innerInstr := range innerSet.Instructions {
				// Orca only uses regular token transfers
				if isOrcaTransfer(innerInstr, ctx.AccountKeys) {
					// If this is not consecutive with the last transfer, reset
					if lastTransferIndex != -1 && i != lastTransferIndex+1 {
						currentTransfers = nil
					}
					lastTransferIndex = i

					transfer, err := p.processTransfer(innerInstr, ctx)
					if err != nil {
						currentTransfers = nil
						continue
					}
					currentTransfers = append(currentTransfers, *transfer)

					// When we have a pair of consecutive transfers, build a swap
					if len(currentTransfers) == 2 {
						if p.seenInstructionPairs[innerInstr.Data.String()+innerSet.Instructions[i-1].Data.String()] {
							currentTransfers = nil
							continue
						}
						p.seenInstructionPairs[innerInstr.Data.String()+innerSet.Instructions[i-1].Data.String()] = true

						swap, err := p.buildSwapInfo(currentTransfers[0], currentTransfers[1], ctx)
						if err != nil {
							currentTransfers = nil
							continue
						}
						swaps = append(swaps, swap)
						currentTransfers = nil
					}
				} else {
					// Reset transfers if we see a non-transfer instruction
					currentTransfers = nil
					lastTransferIndex = -1
				}
			}
		}
	}

	if len(swaps) == 0 {
		return nil, fmt.Errorf("no valid Orca swaps found")
	}

	return swaps, nil
}

// isOrcaTransfer checks if the instruction is a token transfer
func isOrcaTransfer(instr solana.CompiledInstruction, accountKeys []solana.PublicKey) bool {
	if len(instr.Accounts) < 3 || len(instr.Data) < 9 {
		return false
	}

	progID := accountKeys[instr.ProgramIDIndex]
	if !progID.Equals(solana.TokenProgramID) {
		return false
	}

	return instr.Data[0] == 3 // Transfer instruction
}

// processTransfer extracts transfer information from the instruction
func (p *OrcaParser) processTransfer(instr solana.CompiledInstruction, ctx *TransactionContext) (*TokenInfo, error) {
	if len(instr.Data) < 9 {
		return nil, fmt.Errorf("invalid transfer instruction data")
	}

	amount := binary.LittleEndian.Uint64(instr.Data[1:9])

	// Get source and destination accounts
	sourceAcc := ctx.AccountKeys[instr.Accounts[0]]
	destAcc := ctx.AccountKeys[instr.Accounts[1]]

	// Find token mint from either source or destination account
	mint := p.findTokenMint(sourceAcc, destAcc, ctx)
	if mint == (solana.PublicKey{}) {
		return nil, fmt.Errorf("could not determine token mint")
	}

	return &TokenInfo{
		Mint:     mint,
		Amount:   amount,
		Decimals: ctx.GetMintDecimals(mint),
	}, nil
}

// findTokenMint looks up the mint for token accounts
func (p *OrcaParser) findTokenMint(source, dest solana.PublicKey, ctx *TransactionContext) solana.PublicKey {
	// Check both pre and post token balances
	balances := append(ctx.Meta.PreTokenBalances, ctx.Meta.PostTokenBalances...)

	for _, balance := range balances {
		accKey := ctx.AccountKeys[balance.AccountIndex]
		if accKey.Equals(source) || accKey.Equals(dest) {
			return balance.Mint
		}
	}

	return solana.PublicKey{}
}

// buildSwapInfo creates a SwapInfo from a pair of transfers
func (p *OrcaParser) buildSwapInfo(transfer1, transfer2 TokenInfo, ctx *TransactionContext) (*SwapInfo, error) {
	if transfer1.Mint.Equals(transfer2.Mint) {
		return nil, fmt.Errorf("same token in both transfers")
	}

	swapInfo := &SwapInfo{
		Protocol: SwapTypeOrca,
	}

	// Find input token (transferred from signer)
	signers := ctx.Transaction.Message.Signers()
	found := false

	// Check first transfer
	for _, balance := range ctx.Meta.PreTokenBalances {
		if balance.Mint.Equals(transfer1.Mint) {
			owner := balance.Owner
			for _, signer := range signers {
				if owner.Equals(signer) {
					swapInfo.TokenIn = transfer1
					swapInfo.TokenOut = transfer2
					found = true
					break
				}
			}
			break
		}
	}

	// If first transfer wasn't from signer, use second transfer
	if !found {
		swapInfo.TokenIn = transfer2
		swapInfo.TokenOut = transfer1
	}

	return swapInfo, nil
}
