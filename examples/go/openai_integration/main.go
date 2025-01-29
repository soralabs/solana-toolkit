package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	toolkit "github.com/soralabs/solana-toolkit/go/toolkit"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get environment variables
	rpcURL := os.Getenv("SOLANA_RPC_URL")
	openaiKey := os.Getenv("OPENAI_API_KEY")

	if rpcURL == "" || openaiKey == "" {
		log.Fatal("Please set SOLANA_RPC_URL and OPENAI_API_KEY environment variables")
	}

	// Initialize toolkit
	tk, err := toolkit.New(rpcURL)
	if err != nil {
		log.Fatal("Failed to initialize toolkit:", err)
	}

	// Get toolkit functions for OpenAI
	toolkitFunctions := tk.GetTools()

	openaiFunctions := []openai.FunctionDefinition{}
	for _, tool := range toolkitFunctions {
		openaiFunctions = append(openaiFunctions, toolToOpenAIFunction(tool))
	}

	// Initialize OpenAI client
	openaiClient := openai.NewClient(openaiKey)

	// Example queries to demonstrate different toolkit capabilities
	queries := []string{
		"Show me info about the transaction '52yGM7UrquuXAngCLPKuBUDquiwAHz6SQLggsNZfbKxGuegK2oYDCijYaeEno9sYxecVsxyh9nQrMhqqZq5oQgaa'",
	}

	for _, query := range queries {
		fmt.Printf("\n\nProcessing query: %s\n", query)
		fmt.Println("----------------------------------------")

		completion, err := openaiClient.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT4,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: query,
					},
				},
				Functions: openaiFunctions,
			},
		)

		if err != nil {
			log.Printf("Error getting completion: %v\n", err)
			continue
		}

		// Process function call response
		if completion.Choices[0].Message.FunctionCall != nil {
			functionCall := completion.Choices[0].Message.FunctionCall
			fmt.Printf("Function called: %s\n", functionCall.Name)

			// Pretty print the arguments
			var args map[string]interface{}
			if err := json.Unmarshal([]byte(functionCall.Arguments), &args); err == nil {
				prettyArgs, _ := json.MarshalIndent(args, "", "  ")
				fmt.Printf("Arguments:\n%s\n", string(prettyArgs))
			}

			// Execute the function
			tool, err := tk.GetTool(functionCall.Name)
			if err != nil {
				log.Printf("Error getting tool: %v\n", err)
				continue
			}

			fmt.Printf("Tool: %s\n", tool.GetName())

			result, err := tool.Execute(context.Background(), json.RawMessage(functionCall.Arguments))
			if err != nil {
				log.Printf("Error executing tool: %v\n", err)
				continue
			}

			// Pretty print the result
			prettyResult, _ := json.MarshalIndent(result, "", "  ")
			fmt.Printf("\nResult:\n%s\n", string(prettyResult))
		} else {
			fmt.Println("No function was called in the response")
		}
	}
}
