package wallet_information

import (
	"context"
	"encoding/json"
	"fmt"

	"sync"

	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	"github.com/soralabs/solana-toolkit/go/internal/gmgn"
	toolkit "github.com/soralabs/toolkit/go"

	tls_client "github.com/bogdanfinn/tls-client"
)

type WalletInformationTool struct {
	toolkit.Tool

	mu sync.Mutex

	rpcClient *rpc.Client

	gmgnClient *gmgn.GMGN
}

func NewWalletInformationTool(rpcClient *rpc.Client) (*WalletInformationTool, error) {
	jar := tls_client.NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Chrome_131),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return nil, err
	}

	return &WalletInformationTool{
		rpcClient:  rpcClient,
		gmgnClient: gmgn.New(client),
	}, nil
}

func (t *WalletInformationTool) GetName() string {
	return "wallet_information"
}

func (t *WalletInformationTool) GetDescription() string {
	return "Fetch information about a solana wallet"
}

func (t *WalletInformationTool) GetSchema() toolkit.Schema {
	return toolkit.Schema{
		Parameters: json.RawMessage(`{
            "type": "object",
            "required": ["wallet"],
            "properties": {
                "wallet": {
                    "type": "string",
                    "description": "The wallet address"
                }
            }
        }`),
	}
}

func (t *WalletInformationTool) Execute(ctx context.Context, params json.RawMessage) (json.RawMessage, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	var input WalletInformationInput
	if err := json.Unmarshal(params, &input); err != nil {
		return nil, fmt.Errorf("failed to parse parameters: %w", err)
	}

	wallet, err := solana.PublicKeyFromBase58(input.Wallet)
	if err != nil {
		return nil, fmt.Errorf("failed to parse wallet: %w", err)
	}

	walletInfo, err := t.gmgnClient.GetWalletInformation(wallet.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet information: %w", err)
	}

	// Fetch holdings information
	holdings, err := t.gmgnClient.GetAllWalletHoldings(wallet.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet holdings: %w", err)
	}

	// Convert holdings
	convertedHoldings := make([]Holding, 0)
	for _, h := range holdings {
		for _, holding := range h.Holdings {
			convertedHoldings = append(convertedHoldings, convertHolding(holding))
		}
	}

	output := WalletInformationOutput{
		// Wallet Address
		Wallet: wallet.String(),

		// Social Information
		TwitterUsername: walletInfo.TwitterUsername,
		TwitterName:     walletInfo.TwitterName,
		TwitterFansNum:  walletInfo.TwitterFansNum,
		ENS:             walletInfo.ENS,
		Avatar:          walletInfo.Avatar,
		Name:            walletInfo.Name,

		// Balance Information
		SolBalance: walletInfo.SolBalance,
		TotalValue: walletInfo.TotalValue,
		TokenNum:   walletInfo.TokenNum,

		// Trading Performance
		UnrealizedProfit: walletInfo.UnrealizedProfit,
		RealizedProfit:   walletInfo.RealizedProfit,
		PNL:              walletInfo.PNL,
		PNL1d:            walletInfo.PNL1d,
		PNL7d:            walletInfo.PNL7d,
		PNL30d:           walletInfo.PNL30d,
		Winrate:          walletInfo.Winrate,

		// Trading Activity
		Buy1d:   walletInfo.Buy1d,
		Sell1d:  walletInfo.Sell1d,
		Buy7d:   walletInfo.Buy7d,
		Sell7d:  walletInfo.Sell7d,
		Buy30d:  walletInfo.Buy30d,
		Sell30d: walletInfo.Sell30d,
		Buy:     walletInfo.Buy,
		Sell:    walletInfo.Sell,

		// Risk Information
		Risk: Risk{
			TokenHoneypotRatio: walletInfo.Risk.TokenHoneypotRatio,
			NoBuyHoldRatio:     walletInfo.Risk.NoBuyHoldRatio,
			SellPassBuyRatio:   walletInfo.Risk.SellPassBuyRatio,
			FastTxRatio:        walletInfo.Risk.FastTxRatio,
		},

		// Additional Information
		Tags:           walletInfo.Tags,
		TagRank:        walletInfo.TagRank,
		FollowersCount: walletInfo.FollowersCount,
		IsContract:     walletInfo.IsContract,
		UpdatedAt:      walletInfo.UpdatedAt,

		// Add holdings information
		Holdings: convertedHoldings,
	}

	return json.Marshal(output)
}
