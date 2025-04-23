package usecase

import (
	"errors"
	"testing"

	"github.com/gong023/umi/domain"
	"github.com/gong023/umi/infra/mock"
	"go.uber.org/mock/gomock"
)

func TestBotService_Start(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock Discord client
	mockDiscordClient := mock.NewMockDiscordClient(ctrl)

	// Create a mock logger
	logger := &MockLogger{}

	// Create the bot service
	botService := NewBotService(mockDiscordClient, logger)

	// Set up expectations
	mockDiscordClient.EXPECT().RegisterHandler(gomock.Any()).Return(func() {})
	mockDiscordClient.EXPECT().Start().Return(nil)

	// Start the bot
	err := botService.Start()

	// Check that there was no error
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check that the logger was called
	if !logger.InfoCalled {
		t.Error("Expected logger.Info to be called")
	}
}

func TestBotService_Start_Error(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock Discord client
	mockDiscordClient := mock.NewMockDiscordClient(ctrl)

	// Create a mock logger
	logger := &MockLogger{}

	// Create the bot service
	botService := NewBotService(mockDiscordClient, logger)

	// Set up expectations
	mockDiscordClient.EXPECT().RegisterHandler(gomock.Any()).Return(func() {})
	mockDiscordClient.EXPECT().Start().Return(errors.New("test error"))

	// Start the bot
	err := botService.Start()

	// Check that there was an error
	if err == nil {
		t.Error("Expected an error, got nil")
	}

	// Check that the logger was called
	if !logger.InfoCalled {
		t.Error("Expected logger.Info to be called")
	}
}

func TestBotService_Stop(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock Discord client
	mockDiscordClient := mock.NewMockDiscordClient(ctrl)

	// Create a mock logger
	logger := &MockLogger{}

	// Create the bot service
	botService := NewBotService(mockDiscordClient, logger)

	// Set up expectations
	mockDiscordClient.EXPECT().Stop().Return(nil)

	// Stop the bot
	err := botService.Stop()

	// Check that there was no error
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check that the logger was called
	if !logger.InfoCalled {
		t.Error("Expected logger.Info to be called")
	}
}

func TestBotService_RegisterCommand(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock Discord client
	mockDiscordClient := mock.NewMockDiscordClient(ctrl)

	// Create a mock logger
	logger := &MockLogger{}

	// Create the bot service
	botService := NewBotService(mockDiscordClient, logger)

	// Create a mock command handler
	mockCommandHandler := &MockCommandHandler{}

	// Register the command
	botService.RegisterCommand("test", mockCommandHandler)

	// Check that the logger was called
	if !logger.InfoCalled {
		t.Error("Expected logger.Info to be called")
	}

	// Check that the command was registered
	handler, ok := botService.commands["test"]
	if !ok {
		t.Error("Expected command to be registered")
	}
	if handler != mockCommandHandler {
		t.Error("Expected registered handler to be the mock handler")
	}
}

// MockCommandHandler is a mock implementation of the domain.CommandHandler interface
type MockCommandHandler struct {
	HandleCalled bool
	Session      domain.Session
	Interaction  *domain.InteractionCreate
}

func (h *MockCommandHandler) Handle(s domain.Session, i *domain.InteractionCreate) {
	h.HandleCalled = true
	h.Session = s
	h.Interaction = i
}
