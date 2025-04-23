package usecase

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gong023/umi/domain"
)

type CreateCommandHandler struct {
	openaiClient domain.OpenAIClient
	logger       domain.Logger
}

func NewCreateCommandHandler(openaiClient domain.OpenAIClient, logger domain.Logger) *CreateCommandHandler {
	return &CreateCommandHandler{
		openaiClient: openaiClient,
		logger:       logger,
	}
}

func (h *CreateCommandHandler) Handle(s domain.Session, i *domain.InteractionCreate) {
	h.logger.Info("Handling create command")

	// Create a response to acknowledge the command
	response := &domain.InteractionResponse{
		Type: int(domain.InteractionResponseChannelMessageWithSource),
		Data: &domain.InteractionResponseData{
			Content: "クイズを考えています...",
		},
	}

	// Send the initial response
	if err := s.InteractionRespond(i, response); err != nil {
		h.logger.Error("Failed to respond to interaction: %v", err)
		return
	}

	// Read the prompt file
	promptPath := filepath.Join("memo", "prompt", "umigame.txt")
	promptContent, err := os.ReadFile(promptPath)
	if err != nil {
		h.logger.Error("Failed to read prompt file: %v", err)
		return
	}

	// Create a request to the OpenAI API
	req := &domain.ChatCompletionRequest{
		Model: "gpt-4-turbo",
		Messages: []domain.ChatMessage{
			{
				Role:    "system",
				Content: string(promptContent),
			},
			{
				Role:    "user",
				Content: "新しいウミガメのスープクイズを考えてください。",
			},
		},
		Temperature: 0.7,
	}

	// Send the request to the OpenAI API
	h.logger.Info("Sending request to OpenAI API")
	resp, err := h.openaiClient.CreateChatCompletion(req)
	if err != nil {
		h.logger.Error("Failed to create chat completion: %v", err)
		return
	}

	// Extract the quiz from the response
	if len(resp.Choices) == 0 {
		h.logger.Error("No choices in response")
		return
	}

	quiz := resp.Choices[0].Message.Content
	h.logger.Info("Received quiz: %s", quiz)

	// Save the quiz to the context directory
	timestamp := time.Now().Format("20060102_150405")
	contextPath := filepath.Join("memo", "context", fmt.Sprintf("quiz_%s.txt", timestamp))

	if err := os.MkdirAll(filepath.Dir(contextPath), 0755); err != nil {
		h.logger.Error("Failed to create context directory: %v", err)
		return
	}

	if err := os.WriteFile(contextPath, []byte(quiz), 0644); err != nil {
		h.logger.Error("Failed to write quiz to context file: %v", err)
		return
	}

	h.logger.Info("Saved quiz to context file: %s", contextPath)

	// Format the quiz
	formattedQuiz := fmt.Sprintf("**新しいウミガメのスープクイズ**\n\n%s", strings.TrimSpace(quiz))

	// Create a follow-up message with the quiz
	// Note: In a real implementation, you would need to use the Discord API to send a follow-up message
	// This is just a placeholder to demonstrate the concept
	h.logger.Info("Quiz created: %s", formattedQuiz)
}
