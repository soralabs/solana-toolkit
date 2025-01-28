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

func NewWalletInformationTool(rpcClient *rpc.Client) *WalletInformationTool {
	jar := tls_client.NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Chrome_131),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		panic(err)
	}

	return &WalletInformationTool{
		rpcClient:  rpcClient,
		gmgnClient: gmgn.New(client),
	}
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

	output := WalletInformationOutput{
		// Wallet Address
		Wallet: wallet.String(),

		// Social Information
		TwitterUsername: walletInfo.Data.TwitterUsername,
		TwitterName:     walletInfo.Data.TwitterName,
		TwitterFansNum:  walletInfo.Data.TwitterFansNum,
		ENS:             walletInfo.Data.ENS,
		Avatar:          walletInfo.Data.Avatar,
		Name:            walletInfo.Data.Name,

		// Balance Information
		SolBalance: walletInfo.Data.SolBalance,
		TotalValue: walletInfo.Data.TotalValue,
		TokenNum:   walletInfo.Data.TokenNum,

		// Trading Performance
		UnrealizedProfit: walletInfo.Data.UnrealizedProfit,
		RealizedProfit:   walletInfo.Data.RealizedProfit,
		PNL:              walletInfo.Data.PNL,
		PNL1d:            walletInfo.Data.PNL1d,
		PNL7d:            walletInfo.Data.PNL7d,
		PNL30d:           walletInfo.Data.PNL30d,
		Winrate:          walletInfo.Data.Winrate,

		// Trading Activity
		Buy1d:   walletInfo.Data.Buy1d,
		Sell1d:  walletInfo.Data.Sell1d,
		Buy7d:   walletInfo.Data.Buy7d,
		Sell7d:  walletInfo.Data.Sell7d,
		Buy30d:  walletInfo.Data.Buy30d,
		Sell30d: walletInfo.Data.Sell30d,
		Buy:     walletInfo.Data.Buy,
		Sell:    walletInfo.Data.Sell,

		// Risk Information
		Risk: Risk{
			TokenHoneypotRatio: walletInfo.Data.Risk.TokenHoneypotRatio,
			NoBuyHoldRatio:     walletInfo.Data.Risk.NoBuyHoldRatio,
			SellPassBuyRatio:   walletInfo.Data.Risk.SellPassBuyRatio,
			FastTxRatio:        walletInfo.Data.Risk.FastTxRatio,
		},

		// Additional Information
		Tags:           walletInfo.Data.Tags,
		TagRank:        walletInfo.Data.TagRank,
		FollowersCount: walletInfo.Data.FollowersCount,
		IsContract:     walletInfo.Data.IsContract,
		UpdatedAt:      walletInfo.Data.UpdatedAt,
	}

	return json.Marshal(output)
}
