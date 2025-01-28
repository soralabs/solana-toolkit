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
	TokenMint   string  `json:"token_mint"`
	Amount      float64 `json:"amount"`

	// Create params
	TokenName   string `json:"token_name"`
	TokenSymbol string `json:"token_symbol"`
}

type OnchainActionsInput struct {
	Action Action `json:"action"`
	Params Params `json:"params"`
}

type OnchainActionsOutput struct {
	Signature   string  `json:"signature"`
	MintAddress *string `json:"mint"`
}
