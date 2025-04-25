package usecase

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gong023/umi/domain"
)

type GiveupCommandHandler struct {
	openaiClient domain.OpenAIClient
	logger       domain.Logger
}

func NewGiveupCommandHandler(openaiClient domain.OpenAIClient, logger domain.Logger) *GiveupCommandHandler {
	return &GiveupCommandHandler{
		openaiClient: openaiClient,
		logger:       logger,
	}
}

func (h *GiveupCommandHandler) Handle(s domain.Session, i *domain.InteractionCreate) {
	h.logger.Info("Handling giveup command")

	// Create a response to acknowledge the command
	response := &domain.InteractionResponse{
		Type: int(domain.InteractionResponseChannelMessageWithSource),
		Data: &domain.InteractionResponseData{
			Content: "クイズの答えを取得しています...",
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
	promptPath := filepath.Join("memo", "prompt", "onGiveup.txt")
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

	// Process the conversation history and add it to the messages
	if len(conversationHistory) > 0 {
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
	}

	// Add the giveup message
	messages = append(messages, domain.ChatMessage{
		Role:    "user",
		Content: "クイズを諦めます。正解を教えてください。",
	})

	// Log the messages being sent to OpenAI
	h.logger.Info("Sending the following messages to OpenAI:")
	for i, msg := range messages {
		h.logger.Info("Message %d - Role: %s, Content: %s", i, msg.Role, msg.Content)
	}

	req := &domain.ChatCompletionRequest{
		Model:       "chatgpt-4o-latest",
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

	// Extract the answer from the response
	if len(resp.Choices) == 0 {
		h.logger.Error("No choices in response")
		return
	}

	answer := resp.Choices[0].Message.Content
	h.logger.Info("Received answer: %s", answer)

	// Format the answer
	formattedAnswer := fmt.Sprintf("**クイズの正解**\n\n%s", strings.TrimSpace(answer))

	// Delete the context file
	if err := os.Remove(contextPath); err != nil {
		h.logger.Error("Failed to delete context file: %v", err)
		// Continue with the response even if we fail to delete the context file
	} else {
		h.logger.Info("Deleted context file after giveup")
	}

	// Send the response with the answer
	if err := s.FollowupMessage(i, formattedAnswer); err != nil {
		h.logger.Error("Failed to send follow-up message: %v", err)
	}

	h.logger.Info("Answer provided: %s", formattedAnswer)
}
