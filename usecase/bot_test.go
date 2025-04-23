package usecase

import (
	"errors"
	"testing"

	"github.com/gong023/umi/infra/mock"
	"go.uber.org/mock/gomock"
)

func TestBotService_Start(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock Discord client
	mockDiscordClient := mock.NewMockDiscordClient(ctrl)

	// Create a mock OpenAI client
	mockOpenAIClient := mock.NewMockOpenAIClient(ctrl)

	// Create a mock logger
	mockLogger := mock.NewMockLogger(ctrl)
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

	// Create the bot service
	botService := NewBotService(mockDiscordClient, mockOpenAIClient, mockLogger)

	// Set up expectations
	mockDiscordClient.EXPECT().Start().Return(nil)
	mockDiscordClient.EXPECT().RegisterHandler(gomock.Any()).Return(func() {})
	mockDiscordClient.EXPECT().RegisterCommands(gomock.Any()).Return(nil)

	// Start the bot
	err := botService.Start()

	// Check that there was no error
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestBotService_Start_Error(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock Discord client
	mockDiscordClient := mock.NewMockDiscordClient(ctrl)

	// Create a mock OpenAI client
	mockOpenAIClient := mock.NewMockOpenAIClient(ctrl)

	// Create a mock logger
	mockLogger := mock.NewMockLogger(ctrl)
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

	// Create the bot service
	botService := NewBotService(mockDiscordClient, mockOpenAIClient, mockLogger)

	// Set up expectations
	mockDiscordClient.EXPECT().Start().Return(errors.New("test error"))

	// Start the bot
	err := botService.Start()

	// Check that there was an error
	if err == nil {
		t.Error("Expected an error, got nil")
	}
}

func TestBotService_Stop(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock Discord client
	mockDiscordClient := mock.NewMockDiscordClient(ctrl)

	// Create a mock OpenAI client
	mockOpenAIClient := mock.NewMockOpenAIClient(ctrl)

	// Create a mock logger
	mockLogger := mock.NewMockLogger(ctrl)
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

	// Create the bot service
	botService := NewBotService(mockDiscordClient, mockOpenAIClient, mockLogger)

	// Set up expectations
	mockDiscordClient.EXPECT().DeleteCommands().Return(nil)
	mockDiscordClient.EXPECT().Stop().Return(nil)

	// Stop the bot
	err := botService.Stop()

	// Check that there was no error
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestBotService_Stop_DeleteCommandsError(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock Discord client
	mockDiscordClient := mock.NewMockDiscordClient(ctrl)

	// Create a mock OpenAI client
	mockOpenAIClient := mock.NewMockOpenAIClient(ctrl)

	// Create a mock logger
	mockLogger := mock.NewMockLogger(ctrl)
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

	// Create the bot service
	botService := NewBotService(mockDiscordClient, mockOpenAIClient, mockLogger)

	// Set up expectations
	mockDiscordClient.EXPECT().DeleteCommands().Return(errors.New("delete commands error"))
	mockDiscordClient.EXPECT().Stop().Return(nil)

	// Stop the bot
	err := botService.Stop()

	// Check that there was no error (we continue even if DeleteCommands fails)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestBotService_RegisterCommand(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock Discord client
	mockDiscordClient := mock.NewMockDiscordClient(ctrl)

	// Create a mock OpenAI client
	mockOpenAIClient := mock.NewMockOpenAIClient(ctrl)

	// Create a mock logger
	mockLogger := mock.NewMockLogger(ctrl)
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

	// Create the bot service
	botService := NewBotService(mockDiscordClient, mockOpenAIClient, mockLogger)

	// Create a mock command handler
	mockCommandHandler := mock.NewMockCommandHandler(ctrl)

	// Register the command
	botService.RegisterCommand("test", mockCommandHandler)

	// Check that the command was registered
	handler, ok := botService.commands["test"]
	if !ok {
		t.Error("Expected command to be registered")
	}
	if handler != mockCommandHandler {
		t.Error("Expected registered handler to be the mock handler")
	}
}
