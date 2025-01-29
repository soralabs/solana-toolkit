package gmgn

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"time"

	http "github.com/bogdanfinn/fhttp"
)

func (g *GMGN) GetWalletInformation(wallet string) (*WalletInfoData, error) {
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

	return &response.Data, nil
}

func (g *GMGN) GetWalletHoldings(wallet string) (*WalletHoldingsData, error) {
	url := fmt.Sprintf("https://gmgn.ai/api/v1/wallet_holdings/sol/%s", wallet)

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

	var response WalletHoldingsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return &response.Data, nil
}

func (g *GMGN) GetWalletHoldingsNext(wallet string, cursor string, limit int) (*WalletHoldingsData, error) {
	baseURL := fmt.Sprintf("https://gmgn.ai/api/v1/wallet_holdings/sol/%s", wallet)

	values := url.Values{}
	values.Set("limit", strconv.Itoa(limit))
	if cursor != "" {
		values.Set("cursor", cursor)
	}

	url := baseURL + "?" + values.Encode()

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

	var response WalletHoldingsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("API error: %s", response.Message)
	}

	return &response.Data, nil
}

func (g *GMGN) GetAllWalletHoldings(wallet string) ([]*WalletHoldingsData, error) {
	var allHoldings []*WalletHoldingsData

	// Get first page
	firstPage, err := g.GetWalletHoldings(wallet)
	if err != nil {
		return nil, err
	}
	allHoldings = append(allHoldings, firstPage)

	// Keep getting next pages while there's a cursor
	cursor := firstPage.Next
	for cursor != "" {
		time.Sleep(500 * time.Millisecond)

		nextPage, err := g.GetWalletHoldingsNext(wallet, cursor, 50)
		if err != nil {
			return nil, err
		}
		allHoldings = append(allHoldings, nextPage)
		cursor = nextPage.Next
	}

	return allHoldings, nil
}
