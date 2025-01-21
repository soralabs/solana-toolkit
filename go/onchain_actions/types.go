package onchain_actions

type Action string

const (
	ActionBuy    Action = "buy"
	ActionSell   Action = "sell"
	ActionCreate Action = "create"
)

type Params struct {
	Source      string  `json:"source"`
	Destination string  `json:"destination"`
	Amount      float64 `json:"amount"`
	TokenMint   string  `json:"token_mint"`
	Slippage    float64 `json:"slippage"`
}

type OnchainActionsInput struct {
	Action Action `json:"action"`
	Params Params `json:"params"`
}

type OnchainActionsOutput struct {
}
