package usecase

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gong023/umi/domain"
	"github.com/gong023/umi/infra"
)

type BotService struct {
	discordClient domain.DiscordClient
	openaiClient  domain.OpenAIClient
	logger        domain.Logger
	commands      map[string]domain.CommandHandler
}

func NewBotService(discordClient domain.DiscordClient, openaiClient domain.OpenAIClient, logger domain.Logger) *BotService {
	return &BotService{
		discordClient: discordClient,
		openaiClient:  openaiClient,
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
	// Use a function that matches the expected signature for discordgo handlers
	s.discordClient.RegisterHandler(func(session *discordgo.Session, i *discordgo.InteractionCreate) {
		s.handleInteractionCreate(session, i)
	})
	
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
	
	// Delete all commands
	s.logger.Info("Deleting all commands")
	if err := s.discordClient.DeleteCommands(); err != nil {
		s.logger.Error("Failed to delete commands: %v", err)
		// Continue with shutdown even if command deletion fails
	}
	
	// Stop the Discord client
	return s.discordClient.Stop()
}

func (s *BotService) RegisterCommand(name string, handler domain.CommandHandler) {
	s.logger.Info("Registering command: %s", name)
	s.commands[name] = handler
}

func (s *BotService) handleInteractionCreate(session *discordgo.Session, i *discordgo.InteractionCreate) {
	s.logger.Info("Received interaction event")
	
	// Check if this is an application command (slash command)
	// Type 2 is APPLICATION_COMMAND
	if i.Type != discordgo.InteractionApplicationCommand {
		s.logger.Info("Ignoring non-application command interaction: %d", i.Type)
		return
	}
	
	s.logger.Info("Received application command interaction")
	
	// Get the command name from the interaction data
	commandName := i.ApplicationCommandData().Name
	s.logger.Info("Command name: %s", commandName)
	
	// Find the handler for this command
	handler, ok := s.commands[commandName]
	if !ok {
		s.logger.Debug("No handler for command: %s", commandName)
		return
	}
	
	// Convert the interaction to our domain model
	interaction := infra.ConvertInteraction(i)
	if interaction == nil {
		s.logger.Error("Failed to convert discordgo.InteractionCreate to domain.InteractionCreate")
		return
	}
	
	// Create a session wrapper
	discordSession := infra.NewSession(session)
	
	// Call the handler
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
