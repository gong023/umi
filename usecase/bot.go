package usecase

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gong023/umi/domain"
	"github.com/gong023/umi/infra"
)

type BotService struct {
	discordClient domain.DiscordClient
	logger        domain.Logger
	commands      map[string]domain.CommandHandler
}

func NewBotService(discordClient domain.DiscordClient, logger domain.Logger) *BotService {
	return &BotService{
		discordClient: discordClient,
		logger:        logger,
		commands:      make(map[string]domain.CommandHandler),
	}
}

func (s *BotService) Start() error {
	s.logger.Info("Starting bot service")
	
	s.discordClient.RegisterHandler(s.handleInteractionCreate)
	
	return s.discordClient.Start()
}

func (s *BotService) Stop() error {
	s.logger.Info("Stopping bot service")
	return s.discordClient.Stop()
}

func (s *BotService) RegisterCommand(name string, handler domain.CommandHandler) {
	s.logger.Info("Registering command: %s", name)
	s.commands[name] = handler
}

func (s *BotService) handleInteractionCreate(session interface{}, i interface{}) {
	discordInteraction, ok := i.(*discordgo.InteractionCreate)
	if !ok {
		s.logger.Error("Failed to convert interaction to *discordgo.InteractionCreate")
		return
	}
	
	interaction := infra.ConvertInteraction(discordInteraction)
	if interaction == nil {
		s.logger.Error("Failed to convert discordgo.InteractionCreate to domain.InteractionCreate")
		return
	}
	
	discordSession := infra.NewSession(session.(*discordgo.Session))
	
	if interaction.Data == nil {
		s.logger.Debug("Interaction is not a command")
		return
	}
	
	commandName := interaction.Data.Name
	
	handler, ok := s.commands[commandName]
	if !ok {
		s.logger.Debug("No handler for command: %s", commandName)
		return
	}
	
	handler.Handle(discordSession, interaction)
}

type Session struct {
	session domain.Session
}

func NewSession(session domain.Session) *Session {
	return &Session{
		session: session,
	}
}

func (s *Session) InteractionRespond(i *domain.InteractionCreate, r *domain.InteractionResponse) error {
	return s.session.InteractionRespond(i, r)
}
