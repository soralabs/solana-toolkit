package main

import (
	"github.com/sashabaranov/go-openai"
	toolkit "github.com/soralabs/toolkit/go"
)

func toolToOpenAIFunction(tool toolkit.Tool) openai.FunctionDefinition {
	return openai.FunctionDefinition{
		Name:        tool.GetName(),
		Description: tool.GetDescription(),
		Parameters:  tool.GetSchema().Parameters,
	}
}
