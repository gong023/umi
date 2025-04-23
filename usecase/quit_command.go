package usecase

import (
	"os"
	"path/filepath"

	"github.com/gong023/umi/domain"
)

type QuitCommandHandler struct {
	logger domain.Logger
}

func NewQuitCommandHandler(logger domain.Logger) *QuitCommandHandler {
	return &QuitCommandHandler{
		logger: logger,
	}
}

func (h *QuitCommandHandler) Handle(s domain.Session, i *domain.InteractionCreate) {
	h.logger.Info("Handling quit command")

	// Create a response to acknowledge the command
	response := &domain.InteractionResponse{
		Type: int(domain.InteractionResponseChannelMessageWithSource),
		Data: &domain.InteractionResponseData{
			Content: "クイズを終了しています...",
		},
	}

	// Send the initial response
	if err := s.InteractionRespond(i, response); err != nil {
		h.logger.Error("Failed to respond to interaction: %v", err)
		return
	}

	// Check if a quiz exists
	contextPath := filepath.Join("memo", "context.txt")
	quizExists := false

	// Check if the context file exists
	if _, err := os.Stat(contextPath); err == nil {
		quizExists = true
	}

	if !quizExists {
		h.logger.Info("No quiz found, nothing to quit")

		// Send a response indicating that no quiz is available
		followupMessage := "現在クイズが存在しません。`/create` コマンドで新しいクイズを作成してください。"

		// Send the response with the message
		followupResponse := &domain.InteractionResponse{
			Type: int(domain.InteractionResponseChannelMessageWithSource),
			Data: &domain.InteractionResponseData{
				Content: followupMessage,
			},
		}

		if err := s.InteractionRespond(i, followupResponse); err != nil {
			h.logger.Error("Failed to send follow-up response: %v", err)
		}

		h.logger.Info("No quiz available, suggesting /create command: %s", followupMessage)
		return
	}

	// Delete the context file
	if err := os.Remove(contextPath); err != nil {
		h.logger.Error("Failed to delete context file: %v", err)

		// Send a response indicating that the quit command failed
		errorMessage := "クイズの終了に失敗しました。"

		// Send the response with the message
		errorResponse := &domain.InteractionResponse{
			Type: int(domain.InteractionResponseChannelMessageWithSource),
			Data: &domain.InteractionResponseData{
				Content: errorMessage,
			},
		}

		if err := s.InteractionRespond(i, errorResponse); err != nil {
			h.logger.Error("Failed to send error response: %v", err)
		}

		return
	}

	// Send a response indicating that the quit command succeeded
	successMessage := "クイズを終了しました。新しいクイズを始めるには `/create` コマンドを使用してください。"

	// Send the response with the message
	successResponse := &domain.InteractionResponse{
		Type: int(domain.InteractionResponseChannelMessageWithSource),
		Data: &domain.InteractionResponseData{
			Content: successMessage,
		},
	}

	if err := s.InteractionRespond(i, successResponse); err != nil {
		h.logger.Error("Failed to send success response: %v", err)
	}

	h.logger.Info("Quiz quit successfully")
}
