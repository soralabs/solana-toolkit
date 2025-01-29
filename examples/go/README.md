# Solana Toolkit Examples

This directory contains example code demonstrating various use cases of the Solana Toolkit.

## Directory Structure

- `openai_integration/` - Examples of using the toolkit with OpenAI's function calling

## Prerequisites

Before running the examples, make sure you have:

1. Go 1.23.3 or later installed
2. Required environment variables set:
   ```bash
   export SOLANA_RPC_URL="your-rpc-endpoint"
   export OPENAI_API_KEY="your-openai-key"  # Only needed for OpenAI examples
   ```

## Running the Examples

Each example can be run from its directory:

```bash
# Run OpenAI integration example
cd openai_integration
go run main.go
```

## Example Descriptions

### OpenAI Integration
Demonstrates how to use the toolkit's functions with OpenAI's function calling feature. This example shows:
- Getting toolkit functions in OpenAI-compatible format
- Processing natural language queries about Solana data
- Executing toolkit functions based on AI responses