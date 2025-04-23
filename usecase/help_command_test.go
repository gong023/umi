package usecase

import (
	"testing"

	"github.com/gong023/umi/domain"
	"github.com/gong023/umi/infra/mock"
	"go.uber.org/mock/gomock"
)

func TestHelpCommandHandler_Handle(t *testing.T) {
	// Create a mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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
			Name: "help",
		},
	}

	// Set up expectations for the session
	mockSession.EXPECT().InteractionRespond(gomock.Any(), gomock.Any()).Return(nil)

	// Create the help command handler
	handler := NewHelpCommandHandler(mockLogger)

	// Handle the interaction
	handler.Handle(mockSession, interaction)

	// No assertions needed as we're just testing that the handler doesn't panic
	// and that the expected methods are called (which is verified by the mock expectations)
}
