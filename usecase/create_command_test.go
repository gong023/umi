package usecase

import (
	"testing"

	"github.com/gong023/umi/domain"
	"github.com/gong023/umi/infra/mock"
	"go.uber.org/mock/gomock"
)

func TestCreateCommandHandler_Handle_NewQuiz(t *testing.T) {
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

	// Create a mock interaction
	interaction := &domain.InteractionCreate{
		ID:   "test-interaction-id",
		Type: 2, // APPLICATION_COMMAND
		Data: &domain.ApplicationCommandInteractionData{
			Name: "create",
		},
	}

	// Set up expectations for the session
	mockSession.EXPECT().InteractionRespond(gomock.Any(), gomock.Any()).Return(nil)
	mockSession.EXPECT().FollowupMessage(gomock.Any(), gomock.Any()).Return(nil)

	// Set up the file system mock
	contextPath := "memo/context.txt"
	promptPath := "memo/prompt/oncreate.txt"

	// Mock file existence check
	mockFileSystem.EXPECT().FileExists(contextPath).Return(false, nil)

	// Mock prompt file read
	promptContent := "あなたはウミガメのスープクイズを出題するボットです。"
	mockFileSystem.EXPECT().ReadFile(promptPath).Return([]byte(promptContent), nil)

	// Mock path joining
	mockFileSystem.EXPECT().JoinPath("memo", "context.txt").Return(contextPath).AnyTimes()
	mockFileSystem.EXPECT().JoinPath("memo", "prompt", "oncreate.txt").Return(promptPath).AnyTimes()

	// Set up OpenAI mock
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

	mockOpenAIClient.EXPECT().CreateChatCompletion(gomock.Any()).Return(mockResponse, nil)

	// Create the create command handler
	handler := NewCreateCommandHandler(mockOpenAIClient, mockFileSystem, mockLogger)

	// Capture the data written to the file
	var savedQuizData []byte
	mockFileSystem.EXPECT().WriteFile(contextPath, gomock.Any(), gomock.Any()).DoAndReturn(
		func(path string, data []byte, perm int) error {
			savedQuizData = data
			return nil
		})

	// Handle the interaction
	handler.Handle(mockSession, interaction)

	// Verify that the quiz was saved correctly
	if string(savedQuizData) != mockResponse.Choices[0].Message.Content {
		t.Errorf("Expected quiz to be saved as '%s', but got '%s'",
			mockResponse.Choices[0].Message.Content, string(savedQuizData))
	}
}

func TestCreateCommandHandler_Handle_ExistingQuiz(t *testing.T) {
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

	// Create a mock interaction
	interaction := &domain.InteractionCreate{
		ID:   "test-interaction-id",
		Type: 2, // APPLICATION_COMMAND
		Data: &domain.ApplicationCommandInteractionData{
			Name: "create",
		},
	}

	// Set up expectations for the session
	mockSession.EXPECT().InteractionRespond(gomock.Any(), gomock.Any()).Return(nil)
	mockSession.EXPECT().FollowupMessage(gomock.Any(), gomock.Any()).Return(nil)

	// Set up the file system mock
	contextPath := "memo/context.txt"

	// Mock file existence check
	mockFileSystem.EXPECT().FileExists(contextPath).Return(true, nil)

	// Mock file read
	existingQuizContent := "これは既存のクイズです。"
	mockFileSystem.EXPECT().ReadFile(contextPath).Return([]byte(existingQuizContent), nil)

	// Mock path joining
	mockFileSystem.EXPECT().JoinPath("memo", "context.txt").Return(contextPath).AnyTimes()

	// Create the create command handler
	handler := NewCreateCommandHandler(mockOpenAIClient, mockFileSystem, mockLogger)

	// Handle the interaction
	handler.Handle(mockSession, interaction)

	// No need to verify OpenAI calls since it shouldn't be called when a quiz already exists
}
