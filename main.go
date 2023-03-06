// Package main provides a command-line tool that uses the OpenAI API to suggest improvements on function names,
// variable names, comments, and log content of a given code while keeping its structure.
package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

const (
	stdin = "-"
)

var logger = log.New(os.Stderr, "", log.LstdFlags)

func main() {
	// Set up command line flags
	var (
		inputFilePath = flag.String("f", stdin, `source file path or stdin if "-"`)
		writeToFile   = flag.Bool("w", false, `write result to (source) file instead of stdout`)
	)
	flag.Parse()

	// Read file content
	content, err := getFileContent(*inputFilePath)
	if err != nil {
		logger.Fatalf("Failed to read input file: %v", err)
	}

	// Prepare input text for OpenAI API
	inputText := fmt.Sprintf("In order to improve code without changing the structure, please polish the following code with emphasis on function names, variable names, comments and log content:\n%s", string(content))

	// Initialize OpenAI API client
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		logger.Fatal("Please set OPENAI_API_KEY environment variable")
	}
	client := openai.NewClient(apiKey)

	// Generate chat completion using GPT-3.5 Turbo model
	chatRequest := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: inputText,
			},
		},
	}
	chatResponse, err := client.CreateChatCompletion(context.Background(), chatRequest)
	if err != nil {
		logger.Fatalf("Failed to generate chat completion: %v", err)
	}

	improvedCode := chatResponse.Choices[0].Message.Content
	if *writeToFile && *inputFilePath != stdin {
		ioutil.WriteFile(*inputFilePath, []byte(improvedCode), 0)
	} else {
		fmt.Println(improvedCode)
	}
}

// getFileContent reads file content from a given file path or os.Stdin if "-" is passed as the file path.
func getFileContent(filePath string) ([]byte, error) {
	if filePath == stdin {
		return ioutil.ReadAll(os.Stdin)
	}
	return ioutil.ReadFile(filePath)
}
