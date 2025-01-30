package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	toolkit "github.com/soralabs/solana-toolkit/go/toolkit"
)

// ChatMessage represents a message in the chat history
type ChatMessage struct {
	Role    string
	Content string
}

func main() {
	log.Println("Starting interactive Solana toolkit demo with OpenAI integration")

	// Load environment variables
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

	log.Printf("Connecting to Solana RPC URL: %s\n", rpcURL)

	// Initialize toolkit
	tk, err := toolkit.New(rpcURL)
	if err != nil {
		log.Fatal("Failed to initialize toolkit:", err)
	}
	log.Println("Successfully initialized Solana toolkit")

	// Get toolkit functions for OpenAI
	toolkitFunctions := tk.GetTools()
	log.Printf("Loaded %d toolkit functions\n", len(toolkitFunctions))

	openaiFunctions := []openai.FunctionDefinition{}
	for _, tool := range toolkitFunctions {
		openaiFunctions = append(openaiFunctions, toolToOpenAIFunction(tool))
		log.Printf("Registered function: %s\n", tool.GetName())
	}

	// Initialize OpenAI client
	openaiClient := openai.NewClient(openaiKey)
	log.Println("Successfully initialized OpenAI client")

	// Initialize chat history
	chatHistory := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are a helpful assistant that interprets Solana blockchain data. You maintain context of the conversation and can reference previous queries and results.",
		},
	}

	log.Println("\n=== Solana Toolkit Interactive Demo ===")
	log.Println("Type your questions about Solana (type 'exit' to quit)")
	log.Println("Example queries:")
	log.Println("- Show me info about the transaction [SIGNATURE]")
	log.Println("- What's the current balance of [WALLET_ADDRESS]")
	log.Println("- Tell me about the token [TOKEN_ADDRESS]")
	log.Println("----------------------------------------")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		log.Print("\nEnter your query: ")
		if !scanner.Scan() {
			break
		}

		query := scanner.Text()
		if strings.ToLower(query) == "exit" {
			log.Println("User requested exit")
			break
		}

		log.Printf("Processing user query: %s\n", query)

		// Add user query to chat history
		chatHistory = append(chatHistory, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: query,
		})

		completion, err := openaiClient.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:     openai.GPT4oMini,
				Messages:  chatHistory,
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
			log.Printf("OpenAI called function: %s\n", functionCall.Name)

			// Add the assistant's message with function call to chat history
			chatHistory = append(chatHistory, completion.Choices[0].Message)

			// Pretty print the arguments
			var args map[string]interface{}
			if err := json.Unmarshal([]byte(functionCall.Arguments), &args); err == nil {
				prettyArgs, _ := json.MarshalIndent(args, "", "  ")
				log.Printf("Function arguments:\n%s\n", string(prettyArgs))
			}

			// Execute the function
			tool, err := tk.GetTool(functionCall.Name)
			if err != nil {
				log.Printf("Error getting tool: %v\n", err)
				continue
			}

			result, err := tool.Execute(context.Background(), json.RawMessage(functionCall.Arguments))
			if err != nil {
				log.Printf("Error executing tool: %v\n", err)
				continue
			}

			// Pretty print the result
			prettyResult, _ := json.MarshalIndent(result, "", "  ")
			log.Printf("Function result:\n%s\n", string(prettyResult))

			// Add the function result to chat history
			chatHistory = append(chatHistory, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleFunction,
				Name:    functionCall.Name,
				Content: string(prettyResult),
			})

			chatHistory = append(chatHistory, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf("Interpret the results in a conversational way"),
			})

			// Get AI interpretation of the results
			interpretation, err := openaiClient.CreateChatCompletion(
				context.Background(),
				openai.ChatCompletionRequest{
					Model:    openai.GPT4oMini,
					Messages: chatHistory,
				},
			)

			if err != nil {
				log.Printf("Error getting interpretation: %v\n", err)
				continue
			} else {
				interpretationResponse := interpretation.Choices[0].Message.Content
				log.Println(interpretationResponse)

				// Add the interpretation to chat history
				chatHistory = append(chatHistory, interpretation.Choices[0].Message)
			}
		} else {
			log.Println("No function was called in the OpenAI response")
		}
	}
}
