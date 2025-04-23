package usecase

import (
	"fmt"
	"strings"

	"github.com/gong023/umi/domain"
)

type QuizCommandHandler struct {
	openaiClient domain.OpenAIClient
	logger       domain.Logger
}

func NewQuizCommandHandler(openaiClient domain.OpenAIClient, logger domain.Logger) *QuizCommandHandler {
	return &QuizCommandHandler{
		openaiClient: openaiClient,
		logger:       logger,
	}
}

func (h *QuizCommandHandler) Handle(s domain.Session, i *domain.InteractionCreate) {
	h.logger.Info("Handling quiz command")

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

	// Create a request to the OpenAI API
	req := &domain.ChatCompletionRequest{
		Model: "gpt-4-turbo",
		Messages: []domain.ChatMessage{
			{
				Role:    "system",
				Content: "あなたはウミガメのスープクイズを出題するボットです。日本語で短い問題を作成してください。",
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

	// Format the quiz
	formattedQuiz := fmt.Sprintf("**新しいウミガメのスープクイズ**\n\n%s", strings.TrimSpace(quiz))

	// Create a follow-up message with the quiz
	// Note: In a real implementation, you would need to use the Discord API to send a follow-up message
	// This is just a placeholder to demonstrate the concept
	h.logger.Info("Quiz created: %s", formattedQuiz)
}
