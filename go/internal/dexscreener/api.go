package dexscreener

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

func GetPairInformation(ctx context.Context, tokenAddress string) ([]PairInformation, error) {
	client := resty.New()

	requestUrl := fmt.Sprintf("https://api.dexscreener.com/latest/dex/search?q=%s", tokenAddress)

	headers := map[string]string{
		"user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
	}

	response, err := client.R().SetContext(ctx).SetHeaders(headers).Get(requestUrl)
	if err != nil {
		return nil, err
	}

	var pairInformationResponse PairInformationResponse
	if err := json.Unmarshal(response.Body(), &pairInformationResponse); err != nil {
		return nil, err
	}

	return pairInformationResponse.Pairs, nil
}
