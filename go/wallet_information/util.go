package wallet_information

import "github.com/soralabs/solana-toolkit/go/internal/gmgn"

// Convert from gmgn.Holding to wallet_information.Holding
func convertHolding(h gmgn.Holding) Holding {
	return Holding{
		Token: Token{
			Address:       h.Token.Address,
			TokenAddress:  h.Token.TokenAddress,
			Symbol:        h.Token.Symbol,
			Name:          h.Token.Name,
			Decimals:      h.Token.Decimals,
			Logo:          h.Token.Logo,
			PriceChange6h: h.Token.PriceChange6h,
			IsShowAlert:   h.Token.IsShowAlert,
			IsHoneypot:    h.Token.IsHoneypot,
		},
		Balance:             h.Balance,
		USDValue:            h.USDValue,
		RealizedProfit30d:   h.RealizedProfit30d,
		RealizedProfit:      h.RealizedProfit,
		RealizedPNL:         h.RealizedPNL,
		RealizedPNL30d:      h.RealizedPNL30d,
		UnrealizedProfit:    h.UnrealizedProfit,
		UnrealizedPNL:       h.UnrealizedPNL,
		TotalProfit:         h.TotalProfit,
		TotalProfitPNL:      h.TotalProfitPNL,
		AvgCost:             h.AvgCost,
		AvgSold:             h.AvgSold,
		Buy30d:              h.Buy30d,
		Sell30d:             h.Sell30d,
		Sells:               h.Sells,
		Price:               h.Price,
		Cost:                h.Cost,
		PositionPercent:     h.PositionPercent,
		LastActiveTimestamp: h.LastActiveTimestamp,
		HistorySoldIncome:   h.HistorySoldIncome,
		HistoryBoughtCost:   h.HistoryBoughtCost,
		StartHoldingAt:      h.StartHoldingAt,
		EndHoldingAt:        h.EndHoldingAt,
		Liquidity:           h.Liquidity,
		WalletTokenTags:     h.WalletTokenTags,
	}
}
