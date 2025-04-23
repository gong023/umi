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
			Content: "クイズを確認しています...",
		},
	}

	// Send the initial response
	if err := s.InteractionRespond(i, response); err != nil {
		h.logger.Error("Failed to respond to interaction: %v", err)
		return
	}

	// Check if a quiz already exists
	contextDir := filepath.Join("memo", "context")
	files, err := os.ReadDir(contextDir)
	if err != nil && !os.IsNotExist(err) {
		h.logger.Error("Failed to read context directory: %v", err)
		return
	}

	// Find the most recent quiz file
	var latestQuizFile string
	var latestModTime int64
	if err == nil { // Directory exists
		for _, file := range files {
			if file.IsDir() || !strings.HasPrefix(file.Name(), "quiz_") {
				continue
			}

			filePath := filepath.Join(contextDir, file.Name())
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				h.logger.Error("Failed to get file info for %s: %v", filePath, err)
				continue
			}

			modTime := fileInfo.ModTime().Unix()
			if modTime > latestModTime {
				latestModTime = modTime
				latestQuizFile = filePath
			}
		}
	}

	// If a quiz already exists, return it and introduce the /quit command
	if latestQuizFile != "" {
		h.logger.Info("Quiz already exists: %s", latestQuizFile)

		// Read the existing quiz
		quizContent, err := os.ReadFile(latestQuizFile)
		if err != nil {
			h.logger.Error("Failed to read quiz file: %v", err)
			return
		}

		// Format the response with the existing quiz and introduce the /quit command
		formattedResponse := fmt.Sprintf("**現在のウミガメのスープクイズ**\n\n%s\n\n現在のクイズを終了するには `/quit` コマンドを使用してください。", strings.TrimSpace(string(quizContent)))

		h.logger.Info("Returning existing quiz: %s", formattedResponse)
		return
	}

	// No existing quiz, create a new one
	h.logger.Info("No existing quiz found, creating a new one")

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

	// Create the context directory if it doesn't exist
	if err := os.MkdirAll(contextDir, 0755); err != nil {
		h.logger.Error("Failed to create context directory: %v", err)
		return
	}

	// Save the quiz to the context directory
	timestamp := time.Now().Format("20060102_150405")
	contextPath := filepath.Join(contextDir, fmt.Sprintf("quiz_%s.txt", timestamp))

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
