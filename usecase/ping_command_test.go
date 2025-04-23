package usecase

import (
	"errors"
	"testing"

	"github.com/gong023/umi/domain"
)

// MockLogger is a mock implementation of the domain.Logger interface
type MockLogger struct {
	InfoCalled  bool
	ErrorCalled bool
	InfoMsg     string
	ErrorMsg    string
}

func (l *MockLogger) Info(format string, args ...interface{}) {
	l.InfoCalled = true
	l.InfoMsg = format
}

func (l *MockLogger) Error(format string, args ...interface{}) {
	l.ErrorCalled = true
	l.ErrorMsg = format
}

func (l *MockLogger) Debug(format string, args ...interface{}) {
	// Not used in this test
}

// MockSession is a mock implementation of the domain.Session interface
type MockSession struct {
	RespondCalled bool
	RespondError  error
	Interaction   *domain.InteractionCreate
	Response      *domain.InteractionResponse
}

func (s *MockSession) InteractionRespond(i *domain.InteractionCreate, r *domain.InteractionResponse) error {
	s.RespondCalled = true
	s.Interaction = i
	s.Response = r
	return s.RespondError
}

func TestPingCommandHandler_Handle(t *testing.T) {
	// Create a mock logger
	logger := &MockLogger{}

	// Create a mock session
	session := &MockSession{}

	// Create a mock interaction
	interaction := &domain.InteractionCreate{
		ID:   "123",
		Type: 2,
		Data: &domain.ApplicationCommandInteractionData{
			Name: "ping",
		},
	}

	// Create the ping command handler
	handler := NewPingCommandHandler(logger)

	// Test successful response
	handler.Handle(session, interaction)

	// Check that the logger was called
	if !logger.InfoCalled {
		t.Error("Expected logger.Info to be called")
	}

	// Check that the session was called
	if !session.RespondCalled {
		t.Error("Expected session.InteractionRespond to be called")
	}

	// Check the response content
	if session.Response.Data.Content != "はい!" {
		t.Errorf("Expected response content to be 'はい!', got '%s'", session.Response.Data.Content)
	}

	// Test error response
	session = &MockSession{
		RespondError: errors.New("test error"),
	}
	logger = &MockLogger{}
	handler = NewPingCommandHandler(logger)

	handler.Handle(session, interaction)

	// Check that the error logger was called
	if !logger.ErrorCalled {
		t.Error("Expected logger.Error to be called")
	}
}
