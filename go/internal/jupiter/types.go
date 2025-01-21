package jupiter

type QuoteRequest struct {
	InputMint  string `json:"inputMint"`
	OutputMint string `json:"outputMint"`
	Amount     string `json:"amount"`
	Slippage   string `json:"slippage"`
}

type SwapRequest struct {
	Quote         interface{} `json:"quoteResponse"`
	UserPublicKey string      `json:"userPublicKey"`
}
