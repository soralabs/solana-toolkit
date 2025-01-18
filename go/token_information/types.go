package token_information

import "github.com/gagliardetto/solana-go"

type TokenInformationInput struct {
	TokenAddress string `json:"token_address"`
}

type TokenInformationOutput struct {
	Metadata       Metadata `json:"metadata"`
	USDMarketCap   string   `json:"usd_market_cap"`
	IsPumpFunToken bool     `json:"is_pump_fun_token"`
	Socials        []Social `json:"socials"`
	HolderCount    int      `json:"holder_count"`

	PriceChange *PriceChange `json:"price_change"`
}

type Metadata struct {
	Key                 uint8            `json:"key"`
	UpdateAuthority     solana.PublicKey `json:"update_authority"`
	Mint                solana.PublicKey `json:"mint"`
	Data                MetadataData     `json:"data"`
	PrimarySaleHappened bool             `json:"primary_sale_happened"`
	IsMutable           bool             `json:"is_mutable"`
}

type MetadataData struct {
	Name                 string            `json:"name"`
	Symbol               string            `json:"symbol"`
	Uri                  string            `json:"uri"`
	SellerFeeBasisPoints uint16            `json:"seller_fee_basis_points"`
	CreatorCount         uint8             `json:"creator_count"`
	Creators             []MetadataCreator `json:"creators"`
}

type MetadataCreator struct {
	Address  solana.PublicKey `json:"address"`
	Verified bool             `json:"verified"`
	Share    uint8            `json:"share"`
}

type Social struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type PriceChange struct {
	H24 float64 `json:"h24"`
	H6  float64 `json:"h6"`
	H1  float64 `json:"h1"`
	M5  float64 `json:"m5"`
}
