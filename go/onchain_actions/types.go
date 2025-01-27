package onchain_actions

type Action string

const (
	ActionBuy      Action = "buy"
	ActionSell     Action = "sell"
	ActionTransfer Action = "transfer"
	ActionCreate   Action = "create"
)

type Params struct {
	Source      string  `json:"source"`
	Destination string  `json:"destination"`
	Amount      float64 `json:"amount"`
	TokenMint   string  `json:"token_mint"`
	Slippage    float64 `json:"slippage"`
	TokenName   string  `json:"token_name"`
	TokenSymbol string  `json:"token_symbol"`
}

type OnchainActionsInput struct {
	Action Action `json:"action"`
	Params Params `json:"params"`
}

type OnchainActionsOutput struct {
}
