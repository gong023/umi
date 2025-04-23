package domain

type ApplicationCommand struct {
	Name        string
	Description string
}

type DiscordClient interface {
	Start() error
	
	Stop() error
	
	RegisterHandler(handler interface{}) func()
	
	RegisterCommands(commands []*ApplicationCommand) error
}

type CommandHandler interface {
	Handle(s Session, i *InteractionCreate)
}

type Session interface {
	InteractionRespond(i *InteractionCreate, r *InteractionResponse) error
}

type InteractionCreate struct {
	ID string
	
	Type int
	
	Data *ApplicationCommandInteractionData
	
	// Original is the original interaction object from the Discord API
	Original interface{}
}

type ApplicationCommandInteractionData struct {
	Name string
	
	Options []*ApplicationCommandInteractionDataOption
}

type ApplicationCommandInteractionDataOption struct {
	Name string
	
	Value interface{}
}

type InteractionResponse struct {
	Type int
	
	Data *InteractionResponseData
}

type InteractionResponseData struct {
	Content string
}

type InteractionResponseType int

const (
	InteractionResponseChannelMessageWithSource InteractionResponseType = 4
)
