# Solana Toolkit

<div align="center">
  <img src="./img/sora_readme_banner.png" alt="Sora Banner" width="100%" />
  <p align="center">
    <strong>A comprehensive toolkit for building and managing Solana applications</strong>
  </p>
</div>

<div align="center">
  <a href="#features">Features</a> •
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#contributing">Contributing</a>
</div>

## Features

### Token Operations
- **Token Information**: Retrieve comprehensive token data including:
  - Market cap and price metrics
  - Metadata (name, symbol, etc.)
  - Holder statistics and distribution
  - Social media links
  - Price change tracking (5m, 1h, 6h, 24h)
- **Token Creation**: Launch new tokens through the pump.fun platform with customizable parameters
  - Set token name, symbol, and image
  - Configure initial buy amount and slippage
  - Automatic bonding curve setup

### Wallet Management
- **Wallet Information**: Get detailed wallet analytics
  - Token holdings and balances
  - Transaction history
  - Portfolio value tracking
  - Associated accounts

### Transaction Tools
- **Transaction Analysis**: Parse and understand transaction data
  - Detailed swap information
  - Transaction type identification
  - Complete transaction breakdown
  - Historical transaction lookup

### DeFi Integration
- **Trading Operations**: 
  - Execute token swaps through Jupiter
  - Support for all major Solana DEXes
  - Best price routing
  - Slippage protection
- **Token Transfers**: 
  - Send tokens between wallets
  - Support for SOL and SPL tokens
  - Associated token account handling

## Installation

### Go
```bash
go get github.com/soralabs/solana-toolkit/go
```

## Usage

### Token Information
```go
import "github.com/soralabs/solana-toolkit/go/token_information"

// Initialize the tool with your RPC client
tool := token_information.NewTokenInformationTool(rpcClient)

// Get token information
info, err := tool.Execute(context.Background(), TokenInformationInput{
    TokenAddress: "your-token-address",
})
```

### Create Token
```go
import "github.com/soralabs/solana-toolkit/go/onchain_actions"

// Initialize the tool
tool := onchain_actions.NewOnchainActionsTool(rpcClient)

// Create token parameters
params := CreateTokenParams{
    TokenInfo: pumpfun.CreateTokenInformation{
        Name:     "My Token",
        Symbol:   "TKN",
        ImageURI: "https://example.com/image.png",
    },
    Mint:            mintWallet,
    UserPrivateKey:  wallet,
    BuyAmount:       0.1,
    SlippagePercent: 1.0,
}

// Create the token
signature, err := tool.CreateToken(context.Background(), params)
```

### Wallet Information
```go
import "github.com/soralabs/solana-toolkit/go/wallet_information"

// Initialize the tool
tool := wallet_information.NewWalletInformationTool(rpcClient)

// Get wallet information
info, err := tool.Execute(context.Background(), WalletInformationInput{
    Wallet: "wallet-address",
})
```

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.