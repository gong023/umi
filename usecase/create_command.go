package usecase

import (
	"fmt"
	"strings"

	"github.com/gong023/umi/domain"
)

type CreateCommandHandler struct {
	openaiClient domain.OpenAIClient
	fileSystem   domain.FileSystem
	logger       domain.Logger
}

func NewCreateCommandHandler(openaiClient domain.OpenAIClient, fileSystem domain.FileSystem, logger domain.Logger) *CreateCommandHandler {
	return &CreateCommandHandler{
		openaiClient: openaiClient,
		fileSystem:   fileSystem,
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
	contextPath := h.fileSystem.JoinPath("memo", "context.txt")
	quizExists := false
	var existingQuiz string

	// Check if the context file exists
	exists, err := h.fileSystem.FileExists(contextPath)
	if err != nil {
		h.logger.Error("Failed to check if context file exists: %v", err)
		return
	}

	if exists {
		// Read the existing context file
		contextContent, err := h.fileSystem.ReadFile(contextPath)
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
		if err := s.FollowupMessage(i, formattedResponse); err != nil {
			h.logger.Error("Failed to send follow-up message: %v", err)
		}

		h.logger.Info("Returning existing quiz: %s", formattedResponse)
		return
	}

	// No existing quiz, create a new one
	h.logger.Info("No existing quiz found, creating a new one")

	// Read the prompt file
	promptPath := h.fileSystem.JoinPath("memo", "prompt", "oncreate.txt")
	promptContent, err := h.fileSystem.ReadFile(promptPath)
	if err != nil {
		h.logger.Error("Failed to read prompt file: %v", err)
		return
	}

	// Create a request to the OpenAI API
	req := &domain.ChatCompletionRequest{
		Model: "chatgpt-4o-latest",
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

	// Save the quiz to the context file
	if err := h.fileSystem.WriteFile(contextPath, []byte(quiz), 0644); err != nil {
		h.logger.Error("Failed to write quiz to context file: %v", err)
		return
	}

	h.logger.Info("Saved quiz to context file: %s", contextPath)

	// Format the quiz
	formattedQuiz := fmt.Sprintf("**新しいウミガメのスープクイズ**\n\n%s", strings.TrimSpace(quiz))

	// Send the response with the new quiz
	if err := s.FollowupMessage(i, formattedQuiz); err != nil {
		h.logger.Error("Failed to send follow-up message: %v", err)
	}

	h.logger.Info("Quiz created: %s", formattedQuiz)
}
