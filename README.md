# Solana Toolkit

<div align="center">
  <img src="./img/sora_readme_banner.png" alt="Sora Banner" width="100%" />
  <p align="center">
    <strong>A comprehensive toolkit for building on and interacting with the Solana blockchain</strong>
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
  - PumpFun token detection
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
  - GMGN integration for enhanced wallet data

### Transaction Tools
- **Transaction Analysis**: Parse and understand transaction data
  - Detailed swap information across multiple DEXes (Jupiter, Pump.fun, OKX, etc.)
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
  - Automatic Associated Token Account (ATA) handling

## Installation

### Prerequisites
- Go 1.23.3 or later
- Solana RPC endpoint

### Go
```bash
go get github.com/soralabs/solana-toolkit/go
```

## Usage

### Initialize Toolkit
```go
import (
    toolkit "github.com/soralabs/solana-toolkit/go/toolkit"
)

// Initialize toolkit
tk, err := toolkit.New("your-rpc-url")
if err != nil {
    log.Fatal(err)
}
```

For advanced usage and OpenAI Function Calling integration examples, please see the complete implementation in the [examples/go/openai_integration](examples/go/openai_integration) directory.

The toolkit provides built-in functions that can be directly used with OpenAI's function calling feature. These functions include:
- Token information retrieval
- Wallet analysis
- Transaction parsing
- Trading operations
- And more

The `GetTools()` method returns these functions in the format required by OpenAI's API, making it seamless to integrate blockchain functionality with AI.

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.