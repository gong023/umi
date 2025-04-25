package infra

import (
	"fmt"

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
		var options []*discordgo.ApplicationCommandOption

		// Add appropriate options based on the command name
		switch cmd.Name {
		case "q", "answer":
			// Add a string option for the question/answer
			options = append(options, &discordgo.ApplicationCommandOption{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "message",
				Description: "The message to send",
				Required:    true,
			})
		}

		_, err := c.session.ApplicationCommandCreate(c.session.State.User.ID, "", &discordgo.ApplicationCommand{
			Name:        cmd.Name,
			Description: cmd.Description,
			Options:     options,
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

	commands, err := c.session.ApplicationCommands(c.session.State.User.ID, "")
	if err != nil {
		c.logger.Error("Failed to get commands: %v", err)
		return err
	}

	c.logger.Info("Found %d commands to delete", len(commands))

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
	if i.Original != nil {
		originalInteractionCreate, ok := i.Original.(*discordgo.InteractionCreate)
		if ok {
			s.logger.Info("Using original InteractionCreate: ID=%s, Type=%d", originalInteractionCreate.ID, originalInteractionCreate.Type)

			response := &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseType(r.Type),
				Data: &discordgo.InteractionResponseData{
					Content: r.Data.Content,
				},
			}

			s.logger.Info("Sending response: Type=%d, Content=%s", response.Type, response.Data.Content)

			// Send the response - this is the key part
			err := s.session.InteractionRespond(originalInteractionCreate.Interaction, response)
			if err != nil {
				s.logger.Error("Failed to respond to interaction: %v", err)
			}
			return err
		} else {
			s.logger.Error("Original interaction is not of type *discordgo.InteractionCreate: %T", i.Original)
		}
	}

	s.logger.Error("No original interaction available, cannot respond")
	return fmt.Errorf("no original interaction available")
}

func (s *Session) FollowupMessage(i *domain.InteractionCreate, content string) error {
	if i.Original != nil {
		originalInteractionCreate, ok := i.Original.(*discordgo.InteractionCreate)
		if ok {
			s.logger.Info("Sending followup message for interaction: ID=%s", originalInteractionCreate.ID)

			// Create a webhook message
			_, err := s.session.FollowupMessageCreate(originalInteractionCreate.Interaction, true, &discordgo.WebhookParams{
				Content: content,
			})

			if err != nil {
				s.logger.Error("Failed to send followup message: %v", err)
			}
			return err
		} else {
			s.logger.Error("Original interaction is not of type *discordgo.InteractionCreate: %T", i.Original)
		}
	}

	s.logger.Error("No original interaction available, cannot send followup message")
	return fmt.Errorf("no original interaction available")
}

func ConvertInteraction(i *discordgo.InteractionCreate) *domain.InteractionCreate {
	if i == nil || i.Interaction == nil {
		return nil
	}

	result := &domain.InteractionCreate{
		ID:       i.ID,
		Type:     int(i.Type),
		Original: i, // Store the entire InteractionCreate object
	}

	// Check if this is an application command interaction
	if i.Type == discordgo.InteractionApplicationCommand {
		// Get the application command data directly from the interaction
		data := i.ApplicationCommandData()

		// Create the domain ApplicationCommandInteractionData
		result.Data = &domain.ApplicationCommandInteractionData{
			Name: data.Name,
		}

		// Convert options if they exist
		if len(data.Options) > 0 {
			options := make([]*domain.ApplicationCommandInteractionDataOption, len(data.Options))
			for j, opt := range data.Options {
				options[j] = &domain.ApplicationCommandInteractionDataOption{
					Name:  opt.Name,
					Value: opt.Value,
				}
			}
			result.Data.Options = options
		}
	}

	return result
}
