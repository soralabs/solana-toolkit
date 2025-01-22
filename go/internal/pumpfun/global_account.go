package pumpfun

import (
	"context"
	"encoding/binary"
	"fmt"

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
