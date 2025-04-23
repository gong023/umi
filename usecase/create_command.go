package usecase

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
			Content: "クイズを確認しています...",
		},
	}

	// Send the initial response
	if err := s.InteractionRespond(i, response); err != nil {
		h.logger.Error("Failed to respond to interaction: %v", err)
		return
	}

	// Check if a quiz already exists
	contextPath := filepath.Join("memo", "context.txt")
	quizExists := false
	var existingQuiz string

	// Check if the context file exists
	if _, err := os.Stat(contextPath); err == nil {
		// Read the existing context file
		contextContent, err := os.ReadFile(contextPath)
		if err != nil {
			h.logger.Error("Failed to read context file: %v", err)
			return
		}

		// If the file is not empty, a quiz exists
		if len(contextContent) > 0 {
			quizExists = true
			existingQuiz = string(contextContent)
		}
	}

	// If a quiz already exists, return it and introduce the /quit command
	if quizExists {
		h.logger.Info("Quiz already exists")

		// Format the response with the existing quiz and introduce the /quit command
		formattedResponse := fmt.Sprintf("**現在のウミガメのスープクイズ**\n\n%s\n\n現在のクイズを終了するには `/quit` コマンドを使用してください。", strings.TrimSpace(existingQuiz))

		// Send the response with the existing quiz
		followupResponse := &domain.InteractionResponse{
			Type: int(domain.InteractionResponseChannelMessageWithSource),
			Data: &domain.InteractionResponseData{
				Content: formattedResponse,
			},
		}

		if err := s.InteractionRespond(i, followupResponse); err != nil {
			h.logger.Error("Failed to send follow-up response: %v", err)
		}

		h.logger.Info("Returning existing quiz: %s", formattedResponse)
		return
	}

	// No existing quiz, create a new one
	h.logger.Info("No existing quiz found, creating a new one")

	// Read the prompt file
	promptPath := filepath.Join("memo", "prompt", "oncreate.txt")
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

	// Create the memo directory if it doesn't exist
	memoDir := filepath.Dir(contextPath)
	if err := os.MkdirAll(memoDir, 0755); err != nil {
		h.logger.Error("Failed to create memo directory: %v", err)
		return
	}

	// Save the quiz to the context file
	if err := os.WriteFile(contextPath, []byte(quiz), 0644); err != nil {
		h.logger.Error("Failed to write quiz to context file: %v", err)
		return
	}

	h.logger.Info("Saved quiz to context file: %s", contextPath)

	// Format the quiz
	formattedQuiz := fmt.Sprintf("**新しいウミガメのスープクイズ**\n\n%s", strings.TrimSpace(quiz))

	// Send the response with the new quiz
	followupResponse := &domain.InteractionResponse{
		Type: int(domain.InteractionResponseChannelMessageWithSource),
		Data: &domain.InteractionResponseData{
			Content: formattedQuiz,
		},
	}

	if err := s.InteractionRespond(i, followupResponse); err != nil {
		h.logger.Error("Failed to send follow-up response: %v", err)
	}

	h.logger.Info("Quiz created: %s", formattedQuiz)
}
