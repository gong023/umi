package domain

type BotService interface {
	Start() error
	
	Stop() error
}

type CommandRegistry interface {
	RegisterCommand(name string, handler CommandHandler)
}
