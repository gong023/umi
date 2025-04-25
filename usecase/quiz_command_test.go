package usecase

import (
	"testing"

	"github.com/gong023/umi/domain"
	"github.com/gong023/umi/infra/mock"
	"go.uber.org/mock/gomock"
)

func TestQuizCommandHandler_Handle(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock OpenAI client
	mockOpenAIClient := mock.NewMockOpenAIClient(ctrl)

	// Create a mock logger
	mockLogger := mock.NewMockLogger(ctrl)
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

	// Create a mock session
	mockSession := mock.NewMockSession(ctrl)

	// Create a mock interaction
	interaction := &domain.InteractionCreate{
		ID:   "test-interaction-id",
		Type: 2, // APPLICATION_COMMAND
		Data: &domain.ApplicationCommandInteractionData{
			Name: "quiz",
		},
	}

	// Set up expectations for the session
	mockSession.EXPECT().InteractionRespond(gomock.Any(), gomock.Any()).Return(nil)

	// Set up expectations for the OpenAI client
	expectedRequest := &domain.ChatCompletionRequest{
		Model: "chatgpt-4o-latest",
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

	// Create a mock response
	mockResponse := &domain.ChatCompletionResponse{
		ID:      "test-response-id",
		Object:  "chat.completion",
		Created: 1619644200,
		Choices: []struct {
			Index        int                `json:"index"`
			Message      domain.ChatMessage `json:"message"`
			FinishReason string             `json:"finish_reason"`
		}{
			{
				Index: 0,
				Message: domain.ChatMessage{
					Role:    "assistant",
					Content: "男性が海辺で亀のスープを飲んでいました。彼は一口飲んだ後、自殺しました。なぜでしょうか？",
				},
				FinishReason: "stop",
			},
		},
	}

	mockOpenAIClient.EXPECT().CreateChatCompletion(gomock.Eq(expectedRequest)).Return(mockResponse, nil)

	// Create the quiz command handler
	handler := NewQuizCommandHandler(mockOpenAIClient, mockLogger)

	// Handle the interaction
	handler.Handle(mockSession, interaction)

	// No assertions needed as we're just testing that the handler doesn't panic
	// and that the expected methods are called (which is verified by the mock expectations)
}
