package tx_parser

import (
	"fmt"

	"github.com/gagliardetto/solana-go/rpc"
)

// Parser is the main transaction parser
type Parser struct {
	ctx      *TransactionContext
	handlers map[SwapType]SwapParser
}

// New creates a new transaction parser
func New(txResult *rpc.GetTransactionResult) (*Parser, error) {
	tx, err := txResult.Transaction.GetTransaction()
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	// Combine all account keys
	allKeys := append(tx.Message.AccountKeys, txResult.Meta.LoadedAddresses.Writable...)
	allKeys = append(allKeys, txResult.Meta.LoadedAddresses.ReadOnly...)

	ctx := &TransactionContext{
		Transaction: tx,
		Meta:        txResult.Meta,
		AccountKeys: allKeys,
	}

	if err := ctx.ExtractMintDecimals(); err != nil {
		return nil, fmt.Errorf("failed to extract mint decimals: %w", err)
	}

	parser := &Parser{
		ctx:      ctx,
		handlers: make(map[SwapType]SwapParser),
	}

	// Register protocol parsers
	parser.registerHandlers()

	return parser, nil
}

// registerHandlers initializes all protocol-specific parsers
func (p *Parser) registerHandlers() {
	p.handlers[SwapTypeJupiter] = NewJupiterParser()
	p.handlers[SwapTypeJupiterDCA] = NewJupiterDCAParser()
	p.handlers[SwapTypePumpFun] = NewPumpFunParser()
	p.handlers[SwapTypeRaydium] = NewRaydiumParser()
	p.handlers[SwapTypeOrca] = NewOrcaParser()
	p.handlers[SwapTypeMeteora] = NewMeteoraParser()
	p.handlers[SwapTypeMoonshot] = NewMoonshotParser()
	p.handlers[SwapTypeOKX] = NewOKXParser()
}

// ParseTransaction parses the transaction and returns all swap information
func (p *Parser) ParseTransaction() ([]*SwapInfo, error) {
	var allSwaps []*SwapInfo

	// Process each outer instruction in the transaction
	for i, instruction := range p.ctx.Transaction.Message.Instructions {
		// Try each parser for outer instruction
		for _, handler := range p.handlers {
			if handler.CanHandle(instruction, p.ctx.AccountKeys) {
				swaps, err := handler.ParseInstruction(instruction, i, p.ctx)
				if err != nil {
					continue
				}
				for _, swap := range swaps {
					swap.Signers = p.ctx.Transaction.Message.Signers()
					swap.Signatures = p.ctx.Transaction.Signatures
				}
				allSwaps = append(allSwaps, swaps...)
				break // Found matching handler, no need to try others
			}
		}

		if len(allSwaps) > 0 {
			break
		}

		// Check inner instructions
		innerSwaps, err := p.parseInnerInstructions(i)
		if err != nil {
			continue
		}
		allSwaps = append(allSwaps, innerSwaps...)
	}

	if len(allSwaps) == 0 {
		return nil, fmt.Errorf("no valid swaps found in transaction")
	}

	// Remove duplicate swap sets
	allSwaps = p.removeDuplicateSwapSets(allSwaps)

	return allSwaps, nil
}

// parseInnerInstructions processes inner instructions for a given outer instruction index
func (p *Parser) parseInnerInstructions(index int) ([]*SwapInfo, error) {
	var swaps []*SwapInfo

	// Find inner instructions for this index
	for _, innerSet := range p.ctx.Meta.InnerInstructions {
		if innerSet.Index == uint16(index) {
			// Try each inner instruction with each parser
			for _, innerInstr := range innerSet.Instructions {
				for _, handler := range p.handlers {
					if handler.CanHandle(innerInstr, p.ctx.AccountKeys) {
						innerSwaps, err := handler.ParseInstruction(innerInstr, index, p.ctx)
						if err != nil {
							continue
						}
						for _, swap := range innerSwaps {
							swap.Signers = p.ctx.Transaction.Message.Signers()
							swap.Signatures = p.ctx.Transaction.Signatures
						}
						swaps = append(swaps, innerSwaps...)
						break // Found matching handler, no need to try others
					}
				}
			}
		}
	}

	return swaps, nil
}

// removeDuplicateSwapSets removes consecutive sets of swaps that have matching token pairs and amounts
func (p *Parser) removeDuplicateSwapSets(swaps []*SwapInfo) []*SwapInfo {
	if len(swaps) < 4 {
		return swaps
	}

	// Helper function to check if two swaps match
	swapsMatch := func(a, b *SwapInfo) bool {
		return a.TokenIn.Mint.Equals(b.TokenIn.Mint) &&
			a.TokenOut.Mint.Equals(b.TokenOut.Mint) &&
			a.TokenIn.Amount == b.TokenIn.Amount &&
			a.TokenOut.Amount == b.TokenOut.Amount
	}

	// Helper function to check if two sets of consecutive swaps match
	setsMatch := func(set1Start, set1End, set2Start, set2End int) bool {
		if set1End-set1Start != set2End-set2Start {
			return false
		}
		for i := 0; i <= set1End-set1Start; i++ {
			if !swapsMatch(swaps[set1Start+i], swaps[set2Start+i]) {
				return false
			}
		}
		return true
	}

	// Mark swaps to remove
	toRemove := make(map[int]bool)
	for setSize := 2; setSize <= len(swaps)/2; setSize++ {
		for i := 0; i <= len(swaps)-2*setSize; i++ {
			// Skip if this position is already marked for removal
			if toRemove[i] {
				continue
			}

			// Look for matching sets after this position
			for j := i + setSize; j <= len(swaps)-setSize; j++ {
				// Skip if this position is already marked for removal
				if toRemove[j] {
					continue
				}

				if setsMatch(i, i+setSize-1, j, j+setSize-1) {
					// Mark the second set for removal
					for k := j; k < j+setSize; k++ {
						toRemove[k] = true
					}
				}
			}
		}
	}

	// Create new slice without marked swaps
	result := make([]*SwapInfo, 0, len(swaps))
	for i, swap := range swaps {
		if !toRemove[i] {
			result = append(result, swap)
		}
	}

	return result
}
