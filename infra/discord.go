package infra

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gong023/umi/domain"
)

type DiscordClient struct {
	session *discordgo.Session
	logger  domain.Logger
}

func NewDiscordClient(token string, logger domain.Logger) (*DiscordClient, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &DiscordClient{
		session: session,
		logger:  logger,
	}, nil
}

func (c *DiscordClient) Start() error {
	c.logger.Info("Starting Discord client")
	c.session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildIntegrations
	return c.session.Open()
}

func (c *DiscordClient) Stop() error {
	c.logger.Info("Stopping Discord client")
	return c.session.Close()
}

func (c *DiscordClient) RegisterHandler(handler interface{}) func() {
	return c.session.AddHandler(handler)
}

func (c *DiscordClient) RegisterCommands(commands []*domain.ApplicationCommand) error {
	c.logger.Info("Registering %d commands", len(commands))
	
	for _, cmd := range commands {
		_, err := c.session.ApplicationCommandCreate(c.session.State.User.ID, "", &discordgo.ApplicationCommand{
			Name:        cmd.Name,
			Description: cmd.Description,
		})
		
		if err != nil {
			c.logger.Error("Failed to register command %s: %v", cmd.Name, err)
			return err
		}
		
		c.logger.Info("Registered command: %s", cmd.Name)
	}
	
	return nil
}

func (c *DiscordClient) DeleteCommands() error {
	c.logger.Info("Deleting all commands")
	
	// Get all the commands
	commands, err := c.session.ApplicationCommands(c.session.State.User.ID, "")
	if err != nil {
		c.logger.Error("Failed to get commands: %v", err)
		return err
	}
	
	c.logger.Info("Found %d commands to delete", len(commands))
	
	// Delete each command
	for _, cmd := range commands {
		c.logger.Info("Deleting command: %s", cmd.Name)
		err := c.session.ApplicationCommandDelete(c.session.State.User.ID, "", cmd.ID)
		if err != nil {
			c.logger.Error("Failed to delete command %s: %v", cmd.Name, err)
			return err
		}
	}
	
	return nil
}

type Session struct {
	session *discordgo.Session
	logger  domain.Logger
}

func NewSession(session *discordgo.Session) *Session {
	return &Session{
		session: session,
		logger:  domain.NewSimpleLogger(), // Use a simple logger for now
	}
}

func (s *Session) InteractionRespond(i *domain.InteractionCreate, r *domain.InteractionResponse) error {
	// Use the original interaction if available
	if i.Original != nil {
		originalInteractionCreate, ok := i.Original.(*discordgo.InteractionCreate)
		if ok {
			// Log the original interaction details
			s.logger.Info("Using original InteractionCreate: ID=%s, Type=%d", originalInteractionCreate.Interaction.ID, originalInteractionCreate.Interaction.Type)
			
			// Create the response
			response := &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseType(r.Type),
				Data: &discordgo.InteractionResponseData{
					Content: r.Data.Content,
				},
			}
			
			// Log the response
			s.logger.Info("Sending response: Type=%d, Content=%s", response.Type, response.Data.Content)
			
			// Send the response
			err := s.session.InteractionRespond(originalInteractionCreate.Interaction, response)
			if err != nil {
				s.logger.Error("Failed to respond to interaction: %v", err)
			}
			return err
		} else {
			s.logger.Error("Original interaction is not of type *discordgo.InteractionCreate: %T", i.Original)
		}
	} else {
		s.logger.Info("No original interaction available, falling back to creating a new one")
	}
	
	// Fallback to creating a new interaction
	interaction := &discordgo.Interaction{
		ID:   i.ID,
		Type: discordgo.InteractionType(i.Type),
		Data: &discordgo.ApplicationCommandInteractionData{
			Name: i.Data.Name,
		},
	}
	
	// Log the fallback interaction details
	s.logger.Info("Using fallback interaction: ID=%s, Type=%d", interaction.ID, interaction.Type)
	
	// Create the response
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseType(r.Type),
		Data: &discordgo.InteractionResponseData{
			Content: r.Data.Content,
		},
	}
	
	// Log the response
	s.logger.Info("Sending fallback response: Type=%d, Content=%s", response.Type, response.Data.Content)
	
	// Send the response
	err := s.session.InteractionRespond(interaction, response)
	if err != nil {
		s.logger.Error("Failed to respond to fallback interaction: %v", err)
	}
	return err
}

func ConvertInteraction(i *discordgo.InteractionCreate) *domain.InteractionCreate {
	if i == nil || i.Interaction == nil {
		return nil
	}

	result := &domain.InteractionCreate{
		ID:       i.Interaction.ID,
		Type:     int(i.Interaction.Type),
		Original: i, // Store the entire InteractionCreate object
	}

	if i.Interaction.Data == nil {
		return result
	}

	data, ok := i.Interaction.Data.(*discordgo.ApplicationCommandInteractionData)
	if !ok {
		return result
	}

	result.Data = &domain.ApplicationCommandInteractionData{
		Name: data.Name,
	}

	return result
}
