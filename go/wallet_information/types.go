package wallet_information

type WalletInformationInput struct {
	Wallet string `json:"wallet"`
}

type Risk struct {
	TokenHoneypotRatio float64 `json:"token_honeypot_ratio"`
	NoBuyHoldRatio     float64 `json:"no_buy_hold_ratio"`
	SellPassBuyRatio   float64 `json:"sell_pass_buy_ratio"`
	FastTxRatio        float64 `json:"fast_tx_ratio"`
}

type Token struct {
	Address       string `json:"address"`
	TokenAddress  string `json:"token_address"`
	Symbol        string `json:"symbol"`
	Name          string `json:"name"`
	Decimals      int    `json:"decimals"`
	Logo          string `json:"logo"`
	PriceChange6h string `json:"price_change_6h"`
	IsShowAlert   bool   `json:"is_show_alert"`
	IsHoneypot    *bool  `json:"is_honeypot"`
}

type Holding struct {
	Token               Token    `json:"token"`
	Balance             string   `json:"balance"`
	USDValue            string   `json:"usd_value"`
	RealizedProfit30d   string   `json:"realized_profit_30d"`
	RealizedProfit      string   `json:"realized_profit"`
	RealizedPNL         string   `json:"realized_pnl"`
	RealizedPNL30d      string   `json:"realized_pnl_30d"`
	UnrealizedProfit    string   `json:"unrealized_profit"`
	UnrealizedPNL       string   `json:"unrealized_pnl"`
	TotalProfit         string   `json:"total_profit"`
	TotalProfitPNL      string   `json:"total_profit_pnl"`
	AvgCost             string   `json:"avg_cost"`
	AvgSold             string   `json:"avg_sold"`
	Buy30d              int      `json:"buy_30d"`
	Sell30d             int      `json:"sell_30d"`
	Sells               int      `json:"sells"`
	Price               string   `json:"price"`
	Cost                string   `json:"cost"`
	PositionPercent     string   `json:"position_percent"`
	LastActiveTimestamp int64    `json:"last_active_timestamp"`
	HistorySoldIncome   string   `json:"history_sold_income"`
	HistoryBoughtCost   string   `json:"history_bought_cost"`
	StartHoldingAt      *int64   `json:"start_holding_at"`
	EndHoldingAt        *int64   `json:"end_holding_at"`
	Liquidity           string   `json:"liquidity"`
	WalletTokenTags     []string `json:"wallet_token_tags"`
}

type WalletInformationOutput struct {
	Wallet string `json:"wallet"`
	// Social Information
	TwitterUsername *string `json:"twitter_username,omitempty"`
	TwitterName     *string `json:"twitter_name,omitempty"`
	TwitterFansNum  int     `json:"twitter_fans_num"`
	ENS             *string `json:"ens,omitempty"`
	Avatar          *string `json:"avatar,omitempty"`
	Name            *string `json:"name,omitempty"`

	// Balance Information
	SolBalance string  `json:"sol_balance"`
	TotalValue float64 `json:"total_value"`
	TokenNum   int     `json:"token_num"`

	// Trading Performance
	UnrealizedProfit float64 `json:"unrealized_profit"`
	RealizedProfit   float64 `json:"realized_profit"`
	PNL              float64 `json:"pnl"`
	PNL1d            float64 `json:"pnl_1d"`
	PNL7d            float64 `json:"pnl_7d"`
	PNL30d           float64 `json:"pnl_30d"`
	Winrate          float64 `json:"winrate"`

	// Trading Activity
	Buy1d   int `json:"buy_1d"`
	Sell1d  int `json:"sell_1d"`
	Buy7d   int `json:"buy_7d"`
	Sell7d  int `json:"sell_7d"`
	Buy30d  int `json:"buy_30d"`
	Sell30d int `json:"sell_30d"`
	Buy     int `json:"buy"`
	Sell    int `json:"sell"`

	// Risk Information
	Risk Risk `json:"risk"`

	// Additional Information
	Tags           []string       `json:"tags"`
	TagRank        map[string]int `json:"tag_rank"`
	FollowersCount int            `json:"followers_count"`
	IsContract     bool           `json:"is_contract"`
	UpdatedAt      int64          `json:"updated_at"`

	// Holdings
	Holdings []Holding `json:"holdings"`
}
