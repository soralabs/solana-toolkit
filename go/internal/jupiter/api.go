package jupiter

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

func GetQuoteData(quoteRequest QuoteRequest) (interface{}, error) {
	client := resty.New()

	queryParams := map[string]string{
		"inputMint":   quoteRequest.InputMint,
		"outputMint":  quoteRequest.OutputMint,
		"amount":      quoteRequest.Amount,
		"slippageBps": quoteRequest.Slippage,
	}

	response, err := client.R().SetQueryParams(queryParams).Get("https://quote.jup.ag/v6/quote")
	if err != nil {
		return nil, err
	}

	if response.StatusCode() != 200 {
		return nil, fmt.Errorf("failed to get quote: %s", response.String())
	}

	return response.Body(), nil
}

func GetSwapTransaction(swapRequest SwapRequest) (interface{}, error) {
	client := resty.New()

	jsonBody, err := json.Marshal(swapRequest)
	if err != nil {
		return nil, err
	}

	response, err := client.R().SetBody(jsonBody).Post("https://swap.jup.ag/v6/swap")
	if err != nil {
		return nil, err
	}

	return response.Body(), nil
}
