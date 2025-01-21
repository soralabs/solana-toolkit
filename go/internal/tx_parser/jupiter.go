package tx_parser

import (
	"bytes"
	"fmt"

	ag_binary "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
)

// JupiterParser handles parsing Jupiter protocol swaps
type JupiterParser struct{}

// NewJupiterParser creates a new Jupiter parser instance
func NewJupiterParser() *JupiterParser {
	return &JupiterParser{}
}

// JupiterSwapEvent represents a single swap event in the Jupiter protocol
type JupiterSwapEvent struct {
	Amm          solana.PublicKey
	InputMint    solana.PublicKey
	InputAmount  uint64
	OutputMint   solana.PublicKey
	OutputAmount uint64
}

// CanHandle checks if this parser can handle the given instruction
func (p *JupiterParser) CanHandle(instruction solana.CompiledInstruction, accountKeys []solana.PublicKey) bool {
	return accountKeys[instruction.ProgramIDIndex].Equals(JUPITER_PROGRAM_ID)
}

// ParseInstruction processes the Jupiter instruction and returns swap information
func (p *JupiterParser) ParseInstruction(instruction solana.CompiledInstruction, instructionIndex int, ctx *TransactionContext) ([]*SwapInfo, error) {
	events, err := p.parseJupiterEvents(instruction, instructionIndex, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Jupiter events: %w", err)
	}

	// Group events by route (based on sequential mints)
	routes := p.groupEventsIntoRoutes(events)
	if len(routes) == 0 {
		return nil, fmt.Errorf("no valid swap routes found")
	}

	var swaps []*SwapInfo
	for _, routeEvents := range routes {
		swapInfo, err := p.processRoute(routeEvents, ctx)
		if err != nil {
			continue // Skip invalid routes but continue processing others
		}
		swaps = append(swaps, swapInfo)
	}

	if len(swaps) == 0 {
		return nil, fmt.Errorf("no valid swaps found")
	}

	return swaps, nil
}

// parseJupiterEvents extracts all swap events from the instruction
func (p *JupiterParser) parseJupiterEvents(instruction solana.CompiledInstruction, instructionIndex int, ctx *TransactionContext) ([]*JupiterSwapEvent, error) {
	var events []*JupiterSwapEvent

	for _, innerSet := range ctx.Meta.InnerInstructions {
		if innerSet.Index == uint16(instructionIndex) {
			for _, innerInstr := range innerSet.Instructions {
				event, err := p.parseJupiterEvent(innerInstr, ctx.AccountKeys)
				if err != nil {
					continue // Skip invalid events but continue processing
				}
				events = append(events, event)
			}
		}
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("no Jupiter swap events found")
	}

	return events, nil
}

// parseJupiterEvent decodes a single Jupiter swap event
func (p *JupiterParser) parseJupiterEvent(instruction solana.CompiledInstruction, accountKeys []solana.PublicKey) (*JupiterSwapEvent, error) {
	decodedBytes, err := base58.Decode(instruction.Data.String())
	if err != nil {
		return nil, fmt.Errorf("failed to decode instruction data: %w", err)
	}

	// Check event discriminator
	if len(decodedBytes) < 16 || !bytes.Equal(decodedBytes[:16], JUPITER_ROUTE_EVENT_DISCRIMINATOR[:]) {
		return nil, fmt.Errorf("invalid Jupiter event discriminator")
	}

	// Decode event data
	decoder := ag_binary.NewBorshDecoder(decodedBytes[16:])
	var event JupiterSwapEvent
	if err := decoder.Decode(&event); err != nil {
		return nil, fmt.Errorf("failed to decode Jupiter event: %w", err)
	}

	return &event, nil
}

// groupEventsIntoRoutes groups events into separate routes based on connected mints
func (p *JupiterParser) groupEventsIntoRoutes(events []*JupiterSwapEvent) [][]*JupiterSwapEvent {
	var routes [][]*JupiterSwapEvent
	var currentRoute []*JupiterSwapEvent

	for _, event := range events {
		if len(currentRoute) == 0 {
			currentRoute = append(currentRoute, event)
			continue
		}

		lastEvent := currentRoute[len(currentRoute)-1]
		if lastEvent.OutputMint.Equals(event.InputMint) {
			// Events are connected, add to current route
			currentRoute = append(currentRoute, event)
		} else {
			// Events are not connected, start new route
			if len(currentRoute) > 0 {
				routes = append(routes, currentRoute)
			}
			currentRoute = []*JupiterSwapEvent{event}
		}
	}

	// Add last route if exists
	if len(currentRoute) > 0 {
		routes = append(routes, currentRoute)
	}

	return routes
}

// processRoute converts a route of events into a SwapInfo
func (p *JupiterParser) processRoute(events []*JupiterSwapEvent, ctx *TransactionContext) (*SwapInfo, error) {
	if len(events) == 0 {
		return nil, fmt.Errorf("empty route")
	}

	// Start with first event's input and last event's output
	firstEvent := events[0]
	lastEvent := events[len(events)-1]

	swapInfo := &SwapInfo{
		Protocol: SwapTypeJupiter,
		TokenIn: TokenInfo{
			Mint:     firstEvent.InputMint,
			Amount:   firstEvent.InputAmount,
			Decimals: ctx.GetMintDecimals(firstEvent.InputMint),
		},
		TokenOut: TokenInfo{
			Mint:     lastEvent.OutputMint,
			Amount:   lastEvent.OutputAmount,
			Decimals: ctx.GetMintDecimals(lastEvent.OutputMint),
		},
	}

	// Validate no cyclic routes (input != output)
	if swapInfo.TokenIn.Mint.Equals(swapInfo.TokenOut.Mint) {
		return nil, fmt.Errorf("invalid route: input and output tokens are the same")
	}

	return swapInfo, nil
}
