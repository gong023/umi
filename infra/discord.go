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
	return c.session.Open()
}

func (c *DiscordClient) Stop() error {
	c.logger.Info("Stopping Discord client")
	return c.session.Close()
}

func (c *DiscordClient) RegisterHandler(handler interface{}) func() {
	return c.session.AddHandler(handler)
}

type Session struct {
	session *discordgo.Session
}

func NewSession(session *discordgo.Session) *Session {
	return &Session{
		session: session,
	}
}

func (s *Session) InteractionRespond(i *domain.InteractionCreate, r *domain.InteractionResponse) error {
	return s.session.InteractionRespond(
		&discordgo.Interaction{
			ID:   i.ID,
			Type: discordgo.InteractionType(i.Type),
			Data: &discordgo.ApplicationCommandInteractionData{
				Name: i.Data.Name,
			},
		},
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseType(r.Type),
			Data: &discordgo.InteractionResponseData{
				Content: r.Data.Content,
			},
		},
	)
}

func ConvertInteraction(i *discordgo.InteractionCreate) *domain.InteractionCreate {
	if i == nil || i.Interaction == nil {
		return nil
	}

	result := &domain.InteractionCreate{
		ID:   i.Interaction.ID,
		Type: int(i.Interaction.Type),
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
