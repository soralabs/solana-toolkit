package gmgn

import (
	"encoding/json"
	"fmt"
	"io"

	http "github.com/bogdanfinn/fhttp"
)

func (g *GMGN) GetWalletInformation(wallet string) (*WalletInfoResponse, error) {
	url := fmt.Sprintf("https://gmgn.ai/defi/quotation/v1/smartmoney/sol/walletNew/%s", wallet)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header = http.Header{
		"user-agent": {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"},
		"referer":    {"https://gmgn.ai/"},
	}

	resp, err := g.tlsClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var response WalletInfoResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("API error: %s", response.Msg)
	}

	return &response, nil
}
