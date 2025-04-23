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
	
	// Start the Discord client first so we have a valid session
	if err := s.discordClient.Start(); err != nil {
		return err
	}
	
	// Register the interaction create handler
	s.discordClient.RegisterHandler(s.handleInteractionCreate)
	
	// Register commands with Discord API
	commands := make([]*domain.ApplicationCommand, 0, len(s.commands))
	for name := range s.commands {
		commands = append(commands, &domain.ApplicationCommand{
			Name:        name,
			Description: name + " command",
		})
	}
	
	if err := s.discordClient.RegisterCommands(commands); err != nil {
		return err
	}
	
	return nil
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
	s.logger.Info("Received interaction event")
	
	discordInteraction, ok := i.(*discordgo.InteractionCreate)
	if !ok {
		s.logger.Error("Failed to convert interaction to *discordgo.InteractionCreate: %T", i)
		return
	}
	
	s.logger.Info("Interaction type: %d", discordInteraction.Interaction.Type)
	
	// Check if this is an application command (slash command)
	// Type 2 is APPLICATION_COMMAND
	if discordInteraction.Interaction.Type != discordgo.InteractionType(2) {
		s.logger.Info("Ignoring non-application command interaction: %d", discordInteraction.Interaction.Type)
		return
	}
	
	s.logger.Info("Received application command interaction")
	
	interaction := infra.ConvertInteraction(discordInteraction)
	if interaction == nil {
		s.logger.Error("Failed to convert discordgo.InteractionCreate to domain.InteractionCreate")
		return
	}
	
	s.logger.Info("Converted interaction: ID=%s, Type=%d", interaction.ID, interaction.Type)
	
	discordSession := infra.NewSession(session.(*discordgo.Session))
	
	if interaction.Data == nil {
		s.logger.Debug("Interaction is not a command (Data is nil)")
		return
	}
	
	commandName := interaction.Data.Name
	s.logger.Info("Command name: %s", commandName)
	
	handler, ok := s.commands[commandName]
	if !ok {
		s.logger.Debug("No handler for command: %s", commandName)
		return
	}
	
	s.logger.Info("Found handler for command: %s, calling Handle", commandName)
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
