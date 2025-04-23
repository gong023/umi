package domain

type DiscordClient interface {
	Start() error
	
	Stop() error
	
	RegisterHandler(handler interface{}) func()
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
