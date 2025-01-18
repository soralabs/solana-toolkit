package dexscreener

type PairInformation struct {
	ChainID     string `json:"chainId"`
	DexID       string `json:"dexId"`
	URL         string `json:"url"`
	PairAddress string `json:"pairAddress"`
	BaseToken   struct {
		Address string `json:"address"`
		Name    string `json:"name"`
		Symbol  string `json:"symbol"`
	} `json:"baseToken"`
	QuoteToken struct {
		Address string `json:"address"`
		Name    string `json:"name"`
		Symbol  string `json:"symbol"`
	} `json:"quoteToken"`
	PriceNative string `json:"priceNative"`
	PriceUsd    string `json:"priceUsd"`
	Txns        struct {
		M5 struct {
			Buys  int `json:"buys"`
			Sells int `json:"sells"`
		} `json:"m5"`
		H1 struct {
			Buys  int `json:"buys"`
			Sells int `json:"sells"`
		} `json:"h1"`
		H6 struct {
			Buys  int `json:"buys"`
			Sells int `json:"sells"`
		} `json:"h6"`
		H24 struct {
			Buys  int `json:"buys"`
			Sells int `json:"sells"`
		} `json:"h24"`
	} `json:"txns"`
	Volume struct {
		H24 float64 `json:"h24"`
		H6  float64 `json:"h6"`
		H1  float64 `json:"h1"`
		M5  float64 `json:"m5"`
	} `json:"volume"`
	PriceChange struct {
		M5  float64 `json:"m5"`
		H1  float64 `json:"h1"`
		H6  float64 `json:"h6"`
		H24 float64 `json:"h24"`
	} `json:"priceChange,omitempty"`
	Liquidity struct {
		Usd   float64 `json:"usd"`
		Base  float64 `json:"base"`
		Quote float64 `json:"quote"`
	} `json:"liquidity"`
	Fdv           float64 `json:"fdv"`
	MarketCap     float64 `json:"marketCap"`
	PairCreatedAt int64   `json:"pairCreatedAt"`
	Info          struct {
		ImageURL  string `json:"imageUrl"`
		Header    string `json:"header"`
		OpenGraph string `json:"openGraph"`
		Websites  []struct {
			Label string `json:"label"`
			URL   string `json:"url"`
		} `json:"websites"`
		Socials []struct {
			Type string `json:"type"`
			URL  string `json:"url"`
		} `json:"socials"`
	} `json:"info"`
	Labels       []string `json:"labels,omitempty"`
	PriceChange0 struct {
		H6  float64 `json:"h6"`
		H24 float64 `json:"h24"`
	} `json:"priceChange,omitempty"`
	PriceChange1 struct {
		H1  float64 `json:"h1"`
		H6  float64 `json:"h6"`
		H24 float64 `json:"h24"`
	} `json:"priceChange,omitempty"`
	PriceChange2 struct {
		H6  float64 `json:"h6"`
		H24 float64 `json:"h24"`
	} `json:"priceChange,omitempty"`
	PriceChange3 struct {
		H6  float64 `json:"h6"`
		H24 float64 `json:"h24"`
	} `json:"priceChange,omitempty"`
	PriceChange4 struct {
		H24 float64 `json:"h24"`
	} `json:"priceChange,omitempty"`
	PriceChange5 struct {
		H1  float64 `json:"h1"`
		H6  float64 `json:"h6"`
		H24 float64 `json:"h24"`
	} `json:"priceChange,omitempty"`
	PriceChange6 struct {
		H24 float64 `json:"h24"`
	} `json:"priceChange,omitempty"`
	PriceChange7 struct {
		H6  float64 `json:"h6"`
		H24 float64 `json:"h24"`
	} `json:"priceChange,omitempty"`
	PriceChange8 struct {
		H1  float64 `json:"h1"`
		H6  float64 `json:"h6"`
		H24 float64 `json:"h24"`
	} `json:"priceChange,omitempty"`
	PriceChange9 struct {
		H6  float64 `json:"h6"`
		H24 float64 `json:"h24"`
	} `json:"priceChange,omitempty"`
}

type PairInformationResponse struct {
	SchemaVersion string            `json:"schemaVersion"`
	Pairs         []PairInformation `json:"pairs"`
}
