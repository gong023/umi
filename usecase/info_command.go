package usecase

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gong023/umi/domain"
)

type InfoCommandHandler struct {
	openaiClient domain.OpenAIClient
	logger       domain.Logger
}

func NewInfoCommandHandler(openaiClient domain.OpenAIClient, logger domain.Logger) *InfoCommandHandler {
	return &InfoCommandHandler{
		openaiClient: openaiClient,
		logger:       logger,
	}
}

func (h *InfoCommandHandler) Handle(s domain.Session, i *domain.InteractionCreate) {
	h.logger.Info("Handling info command")

	// Create a response to acknowledge the command
	response := &domain.InteractionResponse{
		Type: int(domain.InteractionResponseChannelMessageWithSource),
		Data: &domain.InteractionResponseData{
			Content: "情報を取得しています...",
		},
	}

	// Send the initial response
	if err := s.InteractionRespond(i, response); err != nil {
		h.logger.Error("Failed to respond to interaction: %v", err)
		return
	}

	// Check if a quiz exists
	contextPath := filepath.Join("memo", "context.txt")
	quizExists := false
	var contextContent []byte
	var conversationHistory []string

	// Check if the context file exists
	if _, err := os.Stat(contextPath); err == nil {
		// Read the existing context file
		contextContent, err = os.ReadFile(contextPath)
		if err != nil {
			h.logger.Error("Failed to read context file: %v", err)
			return
		}

		// If the file is not empty, a quiz exists
		if len(contextContent) > 0 {
			quizExists = true
			// Split the content by new lines to get the conversation history
			conversationHistory = strings.Split(string(contextContent), "\n")
		}
	}

	if !quizExists {
		h.logger.Info("No quiz found")
		// Send a response indicating that no quiz is available
		followupMessage := "現在クイズが存在しません。`/create` コマンドで新しいクイズを作成してください。"

		// Send the response with the message
		if err := s.FollowupMessage(i, followupMessage); err != nil {
			h.logger.Error("Failed to send follow-up message: %v", err)
		}

		h.logger.Info("No quiz available, suggesting /create command: %s", followupMessage)
		return
	}

	// Read the prompt file
	promptPath := filepath.Join("memo", "prompt", "onInfo.txt")
	promptContent, err := os.ReadFile(promptPath)
	if err != nil {
		h.logger.Error("Failed to read prompt file: %v", err)
		return
	}

	// Create a request to the OpenAI API with the conversation history
	messages := []domain.ChatMessage{
		{
			Role:    "system",
			Content: string(promptContent),
		},
	}

	// Add all conversation history
	for i, message := range conversationHistory {
		if message == "" {
			continue
		}

		role := "assistant"
		if i > 0 {
			// Alternate between user and assistant for the conversation history
			if i%2 == 1 {
				role = "user"
			}
		}

		messages = append(messages, domain.ChatMessage{
			Role:    role,
			Content: message,
		})
	}

	req := &domain.ChatCompletionRequest{
		Model:       "gpt-4-turbo",
		Messages:    messages,
		Temperature: 0.7,
	}

	// Send the request to the OpenAI API
	h.logger.Info("Sending request to OpenAI API")
	resp, err := h.openaiClient.CreateChatCompletion(req)
	if err != nil {
		h.logger.Error("Failed to create chat completion: %v", err)
		return
	}

	// Extract the info from the response
	if len(resp.Choices) == 0 {
		h.logger.Error("No choices in response")
		return
	}

	info := resp.Choices[0].Message.Content
	h.logger.Info("Received info: %s", info)

	// Format the info
	formattedInfo := fmt.Sprintf("**クイズ情報**\n\n%s", strings.TrimSpace(info))

	// Append the info to the context file
	// First, read the existing content
	existingContent := ""
	if len(contextContent) > 0 {
		existingContent = string(contextContent)
	}

	// Append the new info with a newline
	updatedContent := existingContent
	if len(updatedContent) > 0 && !strings.HasSuffix(updatedContent, "\n") {
		updatedContent += "\n"
	}
	updatedContent += info

	// Write the updated content back to the context file
	if err := os.WriteFile(contextPath, []byte(updatedContent), 0644); err != nil {
		h.logger.Error("Failed to update context file: %v", err)
		// Continue with the response even if we fail to update the context file
	} else {
		h.logger.Info("Updated context file with new info")
	}

	// Send the response with the info
	if err := s.FollowupMessage(i, formattedInfo); err != nil {
		h.logger.Error("Failed to send follow-up message: %v", err)
	}

	h.logger.Info("Info created: %s", formattedInfo)
}
