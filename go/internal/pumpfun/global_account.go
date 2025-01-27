package pumpfun

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// GlobalAccount represents the global state of the pump.fun program
type GlobalAccount struct {
	Discriminator               uint64
	Initialized                 bool
	Authority                   solana.PublicKey
	FeeRecipient                solana.PublicKey
	InitialVirtualTokenReserves uint64
	InitialVirtualSolReserves   uint64
	InitialRealTokenReserves    uint64
	TokenTotalSupply            uint64
	FeeBasisPoints              uint64
}

func (g *GlobalAccount) FromBuffer(data []byte) error {
	if len(data) < 8+1+32+32+8+8+8+8+8 {
		return fmt.Errorf("buffer too short")
	}

	g.Discriminator = binary.LittleEndian.Uint64(data[0:8])
	g.Initialized = data[8] != 0
	copy(g.Authority[:], data[9:41])
	copy(g.FeeRecipient[:], data[41:73])
	g.InitialVirtualTokenReserves = binary.LittleEndian.Uint64(data[73:81])
	g.InitialVirtualSolReserves = binary.LittleEndian.Uint64(data[81:89])
	g.InitialRealTokenReserves = binary.LittleEndian.Uint64(data[89:97])
	g.TokenTotalSupply = binary.LittleEndian.Uint64(data[97:105])
	g.FeeBasisPoints = binary.LittleEndian.Uint64(data[105:113])

	return nil
}

func (g *GlobalAccount) GetInitialBuyPrice(solAmount uint64) (uint64, error) {
	if solAmount <= 0 {
		return 0, nil
	}

	// Use big.Int for calculations to prevent overflow
	vSol := new(big.Int).SetUint64(g.InitialVirtualSolReserves)
	vToken := new(big.Int).SetUint64(g.InitialVirtualTokenReserves)

	// Add 5% buffer to solAmount for slippage
	amount := new(big.Int).SetUint64(solAmount)
	buffer := new(big.Int).Div(amount, big.NewInt(20)) // 5% = 1/20
	amount = new(big.Int).Add(amount, buffer)

	// Calculate k = x * y
	k := new(big.Int).Mul(vSol, vToken)

	// Calculate new sol reserves: i = x + amount
	newSolReserves := new(big.Int).Add(vSol, amount)

	// Calculate r = k/i (rounded up)
	r := new(big.Int).Div(k, newSolReserves)
	r.Add(r, big.NewInt(1)) // Add 1 to handle division rounding

	// Calculate s = vToken - r
	s := new(big.Int).Sub(vToken, r)

	// Check if s is negative
	if s.Sign() < 0 {
		return 0, fmt.Errorf("negative token amount calculated")
	}

	// Convert back to uint64, checking for overflow
	if !s.IsUint64() {
		return 0, fmt.Errorf("token amount overflow")
	}

	result := s.Uint64()
	if result < g.InitialRealTokenReserves {
		return result, nil
	}
	return g.InitialRealTokenReserves, nil
}

func GetGlobalAccount(ctx context.Context, rpcClient *rpc.Client) (*GlobalAccount, error) {
	accountInfo, err := rpcClient.GetAccountInfo(ctx, GlobalPumpFunAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get global account info: %w", err)
	}

	var global GlobalAccount
	if err := global.FromBuffer(accountInfo.Value.Data.GetBinary()); err != nil {
		return nil, fmt.Errorf("failed to parse global account data: %w", err)
	}

	return &global, nil
}
