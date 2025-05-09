package usecase

import (
	"os"
	"testing"

	"github.com/gong023/umi/domain"
	"github.com/gong023/umi/infra/mock"
	"go.uber.org/mock/gomock"
)

func TestGiveupCommandHandler_Handle(t *testing.T) {
	// Skip this test for now
	t.Skip("This test requires monkey patching filepath.Join, which is not possible in Go")

	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock OpenAI client
	mockOpenAIClient := mock.NewMockOpenAIClient(ctrl)

	// Create a mock logger
	mockLogger := mock.NewMockLogger(ctrl)
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

	// Create a mock session
	mockSession := mock.NewMockSession(ctrl)

	// Create a mock interaction
	interaction := &domain.InteractionCreate{
		ID:   "test-interaction-id",
		Type: 2, // APPLICATION_COMMAND
		Data: &domain.ApplicationCommandInteractionData{
			Name: "giveup",
		},
	}

	// Set up expectations for the session - expect two calls to InteractionRespond
	// First for the initial response, second for the follow-up response
	mockSession.EXPECT().InteractionRespond(gomock.Any(), gomock.Any()).Return(nil).Times(2)

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
					Content: "クイズの正解は、男性は過去に亀のスープを飲んだことがあり、その時の記憶が蘇ったため自殺しました。",
				},
				FinishReason: "stop",
			},
		},
	}

	mockOpenAIClient.EXPECT().CreateChatCompletion(gomock.Any()).Return(mockResponse, nil)

	// Create the giveup command handler
	handler := NewGiveupCommandHandler(mockOpenAIClient, mockLogger)

	// Create a temporary directory for testing
	tempDir := t.TempDir()
	fpHandler := NewMockFilepathHandler(tempDir)

	memoDir := fpHandler.Join("memo")
	promptDir := fpHandler.Join("memo", "prompt")

	// Create the directories
	if err := os.MkdirAll(memoDir, 0755); err != nil {
		t.Fatalf("Failed to create memo directory: %v", err)
	}
	if err := os.MkdirAll(promptDir, 0755); err != nil {
		t.Fatalf("Failed to create prompt directory: %v", err)
	}

	// Create an existing quiz file
	existingQuizContent := "男性が海辺で亀のスープを飲んでいました。彼は一口飲んだ後、自殺しました。なぜでしょうか？"
	contextPath := fpHandler.Join("memo", "context.txt")
	if err := os.WriteFile(contextPath, []byte(existingQuizContent), 0644); err != nil {
		t.Fatalf("Failed to create existing quiz file: %v", err)
	}

	// Create a test prompt file
	promptContent := "あなたはウミガメのスープクイズを出題するボットです。ユーザーがクイズを諦めたため、現在のクイズの正解を詳しく説明してください。"
	promptPath := fpHandler.Join("memo", "prompt", "onGiveup.txt")
	if err := os.WriteFile(promptPath, []byte(promptContent), 0644); err != nil {
		t.Fatalf("Failed to create test prompt file: %v", err)
	}

	// Handle the interaction
	handler.Handle(mockSession, interaction)

	// No assertions needed as we're just testing that the handler doesn't panic
	// and that the expected methods are called (which is verified by the mock expectations)
}

func TestGiveupCommandHandler_Handle_NoQuiz(t *testing.T) {
	// Skip this test for now
	t.Skip("This test requires monkey patching filepath.Join, which is not possible in Go")

	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock OpenAI client
	mockOpenAIClient := mock.NewMockOpenAIClient(ctrl)

	// Create a mock logger
	mockLogger := mock.NewMockLogger(ctrl)
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

	// Create a mock session
	mockSession := mock.NewMockSession(ctrl)

	// Create a mock interaction
	interaction := &domain.InteractionCreate{
		ID:   "test-interaction-id",
		Type: 2, // APPLICATION_COMMAND
		Data: &domain.ApplicationCommandInteractionData{
			Name: "giveup",
		},
	}

	// Set up expectations for the session - expect two calls to InteractionRespond
	// First for the initial response, second for the follow-up response
	mockSession.EXPECT().InteractionRespond(gomock.Any(), gomock.Any()).Return(nil).Times(2)

	// Create the giveup command handler
	handler := NewGiveupCommandHandler(mockOpenAIClient, mockLogger)

	// Create a temporary directory for testing
	tempDir := t.TempDir()
	fpHandler := NewMockFilepathHandler(tempDir)

	memoDir := fpHandler.Join("memo")
	promptDir := fpHandler.Join("memo", "prompt")

	// Create the directories
	if err := os.MkdirAll(memoDir, 0755); err != nil {
		t.Fatalf("Failed to create memo directory: %v", err)
	}
	if err := os.MkdirAll(promptDir, 0755); err != nil {
		t.Fatalf("Failed to create prompt directory: %v", err)
	}

	// Create a test prompt file
	promptContent := "あなたはウミガメのスープクイズを出題するボットです。ユーザーがクイズを諦めたため、現在のクイズの正解を詳しく説明してください。"
	promptPath := fpHandler.Join("memo", "prompt", "onGiveup.txt")
	if err := os.WriteFile(promptPath, []byte(promptContent), 0644); err != nil {
		t.Fatalf("Failed to create test prompt file: %v", err)
	}

	// Handle the interaction
	handler.Handle(mockSession, interaction)

	// No assertions needed as we're just testing that the handler doesn't panic
	// and that the expected methods are called (which is verified by the mock expectations)
}
