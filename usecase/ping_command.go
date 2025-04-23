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
	
	response := &domain.InteractionResponse{
		Type: int(domain.InteractionResponseChannelMessageWithSource),
		Data: &domain.InteractionResponseData{
			Content: "はい!",
		},
	}
	
	if err := s.InteractionRespond(i, response); err != nil {
		h.logger.Error("Failed to respond to ping command: %v", err)
	}
}
