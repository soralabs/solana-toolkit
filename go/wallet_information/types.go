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
}
