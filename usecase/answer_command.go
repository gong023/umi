package usecase

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gong023/umi/domain"
)

type AnswerCommandHandler struct {
	openaiClient domain.OpenAIClient
	logger       domain.Logger
}

func NewAnswerCommandHandler(openaiClient domain.OpenAIClient, logger domain.Logger) *AnswerCommandHandler {
	return &AnswerCommandHandler{
		openaiClient: openaiClient,
		logger:       logger,
	}
}

func (h *AnswerCommandHandler) Handle(s domain.Session, i *domain.InteractionCreate) {
	h.logger.Info("Handling answer command")

	// Extract the message from the command options
	var message string
	if i.Data != nil && len(i.Data.Options) > 0 {
		h.logger.Info("Command options found: %d options", len(i.Data.Options))
		for idx, opt := range i.Data.Options {
			h.logger.Info("Option %d: Name=%s, Value=%v", idx, opt.Name, opt.Value)
		}

		if val, ok := i.Data.Options[0].Value.(string); ok {
			message = val
			h.logger.Info("Extracted answer from user: %s", message)
		} else {
			h.logger.Error("Failed to extract answer: Value is not a string: %T", i.Data.Options[0].Value)
		}
	} else {
		h.logger.Error("No command options found: Data=%v", i.Data)
	}

	if message == "" {
		h.logger.Error("No message provided in answer command")
		// Send a response indicating that a message is required
		response := &domain.InteractionResponse{
			Type: int(domain.InteractionResponseChannelMessageWithSource),
			Data: &domain.InteractionResponseData{
				Content: "回答を入力してください。例: `/answer 男性は亀のスープを飲んだことがあり、妻が亀のスープを作ったことを思い出して自殺した`",
			},
		}
		if err := s.InteractionRespond(i, response); err != nil {
			h.logger.Error("Failed to respond to interaction: %v", err)
		}
		return
	}

	// Create a response to acknowledge the command
	response := &domain.InteractionResponse{
		Type: int(domain.InteractionResponseChannelMessageWithSource),
		Data: &domain.InteractionResponseData{
			Content: "回答を判定しています...",
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
	promptPath := filepath.Join("memo", "prompt", "onAnswer.txt")
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

	// Add the quiz (first line of the conversation history)
	if len(conversationHistory) > 0 {
		messages = append(messages, domain.ChatMessage{
			Role:    "assistant",
			Content: conversationHistory[0],
		})
	}

	// Add the current answer
	messages = append(messages, domain.ChatMessage{
		Role:    "user",
		Content: "回答: " + message,
	})

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

	// Extract the judgment from the response
	if len(resp.Choices) == 0 {
		h.logger.Error("No choices in response")
		return
	}

	judgment := resp.Choices[0].Message.Content
	h.logger.Info("Received judgment: %s", judgment)

	// Check if the answer is correct
	isCorrect := strings.Contains(judgment, "正解")

	// Format the judgment
	formattedJudgment := fmt.Sprintf("**回答**: %s\n\n**判定**: %s", message, strings.TrimSpace(judgment))

	if isCorrect {
		// If the answer is correct, delete the context file
		if err := os.Remove(contextPath); err != nil {
			h.logger.Error("Failed to delete context file: %v", err)
			// Continue with the response even if we fail to delete the context file
		} else {
			h.logger.Info("Deleted context file because the answer was correct")
		}
	} else {
		// If the answer is incorrect, append the judgment to the context file
		// First, read the existing content
		existingContent := ""
		if len(contextContent) > 0 {
			existingContent = string(contextContent)
		}

		// Append the new judgment with a newline
		updatedContent := existingContent
		if len(updatedContent) > 0 && !strings.HasSuffix(updatedContent, "\n") {
			updatedContent += "\n"
		}
		updatedContent += judgment

		// Write the updated content back to the context file
		if err := os.WriteFile(contextPath, []byte(updatedContent), 0644); err != nil {
			h.logger.Error("Failed to update context file: %v", err)
			// Continue with the response even if we fail to update the context file
		} else {
			h.logger.Info("Updated context file with new judgment")
		}
	}

	// Send the response with the judgment
	if err := s.FollowupMessage(i, formattedJudgment); err != nil {
		h.logger.Error("Failed to send follow-up message: %v", err)
	}

	h.logger.Info("Judgment created: %s", formattedJudgment)
}
