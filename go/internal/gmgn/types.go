package gmgn

import (
	tls_client "github.com/bogdanfinn/tls-client"
)

type GMGN struct {
	tlsClient tls_client.HttpClient
}

type WalletInfoResponse struct {
	Code int            `json:"code"`
	Msg  string         `json:"msg"`
	Data WalletInfoData `json:"data"`
}

type WalletInfoData struct {
	TwitterBind        bool           `json:"twitter_bind"`
	TwitterFansNum     int            `json:"twitter_fans_num"`
	TwitterUsername    *string        `json:"twitter_username"`
	TwitterName        *string        `json:"twitter_name"`
	ENS                *string        `json:"ens"`
	Avatar             *string        `json:"avatar"`
	Name               *string        `json:"name"`
	EthBalance         string         `json:"eth_balance"`
	SolBalance         string         `json:"sol_balance"`
	TrxBalance         string         `json:"trx_balance"`
	Balance            string         `json:"balance"`
	TotalValue         float64        `json:"total_value"`
	UnrealizedProfit   float64        `json:"unrealized_profit"`
	UnrealizedPNL      float64        `json:"unrealized_pnl"`
	RealizedProfit     float64        `json:"realized_profit"`
	PNL                float64        `json:"pnl"`
	PNL1d              float64        `json:"pnl_1d"`
	PNL7d              float64        `json:"pnl_7d"`
	PNL30d             float64        `json:"pnl_30d"`
	RealizedProfit1d   float64        `json:"realized_profit_1d"`
	RealizedProfit7d   float64        `json:"realized_profit_7d"`
	RealizedProfit30d  float64        `json:"realized_profit_30d"`
	Winrate            float64        `json:"winrate"`
	AllPNL             float64        `json:"all_pnl"`
	TotalProfit        float64        `json:"total_profit"`
	TotalProfitPNL     float64        `json:"total_profit_pnl"`
	Buy1d              int            `json:"buy_1d"`
	Sell1d             int            `json:"sell_1d"`
	Buy30d             int            `json:"buy_30d"`
	Sell30d            int            `json:"sell_30d"`
	Buy7d              int            `json:"buy_7d"`
	Sell7d             int            `json:"sell_7d"`
	Buy                int            `json:"buy"`
	Sell               int            `json:"sell"`
	HistoryBoughtCost  float64        `json:"history_bought_cost"`
	TokenAvgCost       float64        `json:"token_avg_cost"`
	TokenSoldAvgProfit float64        `json:"token_sold_avg_profit"`
	TokenNum           int            `json:"token_num"`
	ProfitNum          int            `json:"profit_num"`
	Tags               []string       `json:"tags"`
	TagRank            map[string]int `json:"tag_rank"`
	FollowersCount     int            `json:"followers_count"`
	IsContract         bool           `json:"is_contract"`
	UpdatedAt          int64          `json:"updated_at"`
	RefreshRequestedAt *int64         `json:"refresh_requested_at"`
	AvgHoldingPeriod   float64        `json:"avg_holding_peroid"`
	Risk               RiskInfo       `json:"risk"`
}

type RiskInfo struct {
	TokenActive        string  `json:"token_active"`
	TokenHoneypot      string  `json:"token_honeypot"`
	TokenHoneypotRatio float64 `json:"token_honeypot_ratio"`
	NoBuyHold          string  `json:"no_buy_hold"`
	NoBuyHoldRatio     float64 `json:"no_buy_hold_ratio"`
	SellPassBuy        string  `json:"sell_pass_buy"`
	SellPassBuyRatio   float64 `json:"sell_pass_buy_ratio"`
	FastTx             string  `json:"fast_tx"`
	FastTxRatio        float64 `json:"fast_tx_ratio"`
}
