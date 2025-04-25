package usecase

import (
	"testing"

	"github.com/gong023/umi/domain"
	"github.com/gong023/umi/infra/mock"
	"go.uber.org/mock/gomock"
)

func TestQCommandHandler_Handle(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock OpenAI client
	mockOpenAIClient := mock.NewMockOpenAIClient(ctrl)

	// Create a mock logger
	mockLogger := mock.NewMockLogger(ctrl)
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

	// Create a mock file system
	mockFileSystem := mock.NewMockFileSystem(ctrl)

	// Create a mock session
	mockSession := mock.NewMockSession(ctrl)

	// Create a mock interaction with a question
	interaction := &domain.InteractionCreate{
		ID:   "test-interaction-id",
		Type: 2, // APPLICATION_COMMAND
		Data: &domain.ApplicationCommandInteractionData{
			Name: "q",
			Options: []*domain.ApplicationCommandInteractionDataOption{
				{
					Name:  "message",
					Value: "男性は何を飲んでいましたか？",
				},
			},
		},
	}

	// Set up expectations for the session
	mockSession.EXPECT().InteractionRespond(gomock.Any(), gomock.Any()).Return(nil)
	mockSession.EXPECT().FollowupMessage(gomock.Any(), gomock.Any()).Return(nil)

	// Set up expectations for the file system
	contextPath := "memo/context.txt"
	promptPath := "memo/prompt/onQ.txt"

	mockFileSystem.EXPECT().JoinPath("memo", "context.txt").Return(contextPath)
	mockFileSystem.EXPECT().FileExists(contextPath).Return(true, nil)

	quizContent := "男性が海辺で亀のスープを飲んでいました。彼は一口飲んだ後、自殺しました。なぜでしょうか？"
	mockFileSystem.EXPECT().ReadFile(contextPath).Return([]byte(quizContent), nil)

	mockFileSystem.EXPECT().JoinPath("memo", "prompt", "onQ.txt").Return(promptPath)

	promptContent := "あなたはウミガメのスープクイズを出題するボットです。日本語で短い問題を作成してください。問題は謎めいていて、「はい」「いいえ」で答えられる質問によって解決できるものにしてください。問題は論理的で解決可能なものにしてください。"
	mockFileSystem.EXPECT().ReadFile(promptPath).Return([]byte(promptContent), nil)

	// Expect the file to be written with updated content
	mockFileSystem.EXPECT().WriteFile(contextPath, gomock.Any(), 0644).Return(nil)

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
					Content: "はい",
				},
				FinishReason: "stop",
			},
		},
	}

	// Mock the OpenAI client to return our mock response
	mockOpenAIClient.EXPECT().CreateChatCompletion(gomock.Any()).Return(mockResponse, nil)

	// Create the q command handler
	handler := NewQCommandHandler(mockOpenAIClient, mockFileSystem, mockLogger)

	// Handle the interaction
	handler.Handle(mockSession, interaction)
}

func TestQCommandHandler_Handle_NoMessage(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock OpenAI client
	mockOpenAIClient := mock.NewMockOpenAIClient(ctrl)

	// Create a mock logger
	mockLogger := mock.NewMockLogger(ctrl)
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

	// Create a mock file system
	mockFileSystem := mock.NewMockFileSystem(ctrl)

	// Create a mock session
	mockSession := mock.NewMockSession(ctrl)

	// Create a mock interaction without a question
	interaction := &domain.InteractionCreate{
		ID:   "test-interaction-id",
		Type: 2, // APPLICATION_COMMAND
		Data: &domain.ApplicationCommandInteractionData{
			Name:    "q",
			Options: []*domain.ApplicationCommandInteractionDataOption{},
		},
	}

	// Set up expectations for the session
	mockSession.EXPECT().InteractionRespond(gomock.Any(), gomock.Any()).Return(nil)

	// Create the q command handler
	handler := NewQCommandHandler(mockOpenAIClient, mockFileSystem, mockLogger)

	// Handle the interaction
	handler.Handle(mockSession, interaction)

	// No assertions needed as we're just testing that the handler doesn't panic
	// and that the expected methods are called (which is verified by the mock expectations)
}

func TestQCommandHandler_Handle_NoQuiz(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock OpenAI client
	mockOpenAIClient := mock.NewMockOpenAIClient(ctrl)

	// Create a mock logger
	mockLogger := mock.NewMockLogger(ctrl)
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

	// Create a mock file system
	mockFileSystem := mock.NewMockFileSystem(ctrl)

	// Create a mock session
	mockSession := mock.NewMockSession(ctrl)

	// Create a mock interaction with a question
	interaction := &domain.InteractionCreate{
		ID:   "test-interaction-id",
		Type: 2, // APPLICATION_COMMAND
		Data: &domain.ApplicationCommandInteractionData{
			Name: "q",
			Options: []*domain.ApplicationCommandInteractionDataOption{
				{
					Name:  "message",
					Value: "男性は何を飲んでいましたか？",
				},
			},
		},
	}

	// Set up expectations for the session
	mockSession.EXPECT().InteractionRespond(gomock.Any(), gomock.Any()).Return(nil)
	mockSession.EXPECT().FollowupMessage(gomock.Any(), gomock.Any()).Return(nil)

	// Set up expectations for the file system - no quiz exists
	contextPath := "memo/context.txt"
	mockFileSystem.EXPECT().JoinPath("memo", "context.txt").Return(contextPath)
	mockFileSystem.EXPECT().FileExists(contextPath).Return(false, nil)

	// Create the q command handler
	handler := NewQCommandHandler(mockOpenAIClient, mockFileSystem, mockLogger)

	// Handle the interaction
	handler.Handle(mockSession, interaction)
}
