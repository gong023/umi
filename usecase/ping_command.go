package usecase

import (
	"github.com/gong023/umi/domain"
)

type PingCommandHandler struct {
	logger domain.Logger
}

func NewPingCommandHandler(logger domain.Logger) *PingCommandHandler {
	return &PingCommandHandler{
		logger: logger,
	}
}

func (h *PingCommandHandler) Handle(s domain.Session, i *domain.InteractionCreate) {
	h.logger.Info("Handling ping command")
	
	// Log interaction details
	h.logger.Info("Interaction ID: %s, Type: %d", i.ID, i.Type)
	if i.Data != nil {
		h.logger.Info("Interaction Data Name: %s", i.Data.Name)
	}
	if i.Original != nil {
		h.logger.Info("Original interaction is present")
	} else {
		h.logger.Info("Original interaction is nil")
	}
	
	response := &domain.InteractionResponse{
		Type: int(domain.InteractionResponseChannelMessageWithSource),
		Data: &domain.InteractionResponseData{
			Content: "はい!",
		},
	}
	
	h.logger.Info("Sending response with type: %d and content: %s", response.Type, response.Data.Content)
	
	if err := s.InteractionRespond(i, response); err != nil {
		h.logger.Error("Failed to respond to ping command: %v", err)
	} else {
		h.logger.Info("Successfully responded to ping command")
	}
}
