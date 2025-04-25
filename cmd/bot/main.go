package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gong023/umi/domain"
	"github.com/gong023/umi/infra"
	"github.com/gong023/umi/usecase"
)

func main() {
	logger := domain.NewSimpleLogger()
	logger.Info("Starting umi bot")

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		logger.Error("DISCORD_TOKEN environment variable is not set")
		os.Exit(1)
	}

	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	if openaiAPIKey == "" {
		logger.Error("OPENAI_API_KEY environment variable is not set")
		os.Exit(1)
	}

	discordClient, err := infra.NewDiscordClient(token, logger)
	if err != nil {
		logger.Error("Failed to create Discord client: %v", err)
		os.Exit(1)
	}

	openaiClient := infra.NewOpenAIClient(openaiAPIKey, logger)

	botService := usecase.NewBotService(discordClient, openaiClient, logger)

	// Create file system
	fileSystem := infra.NewFileSystem(logger)

	// Create command handlers
	pingCommandHandler := usecase.NewPingCommandHandler(logger)
	createCommandHandler := usecase.NewCreateCommandHandler(openaiClient, fileSystem, logger)
	qCommandHandler := usecase.NewQCommandHandler(openaiClient, fileSystem, logger)
	answerCommandHandler := usecase.NewAnswerCommandHandler(openaiClient, fileSystem, logger)
	infoCommandHandler := usecase.NewInfoCommandHandler(openaiClient, fileSystem, logger)
	clueCommandHandler := usecase.NewClueCommandHandler(openaiClient, fileSystem, logger)
	quitCommandHandler := usecase.NewQuitCommandHandler(fileSystem, logger)
	helpCommandHandler := usecase.NewHelpCommandHandler(logger)

	// Register commands
	botService.RegisterCommand("ping", pingCommandHandler)
	botService.RegisterCommand("create", createCommandHandler)
	botService.RegisterCommand("q", qCommandHandler)
	botService.RegisterCommand("answer", answerCommandHandler)
	botService.RegisterCommand("info", infoCommandHandler)
	botService.RegisterCommand("clue", clueCommandHandler)
	botService.RegisterCommand("quit", quitCommandHandler)
	botService.RegisterCommand("help", helpCommandHandler)

	if err := botService.Start(); err != nil {
		logger.Error("Failed to start bot: %v", err)
		os.Exit(1)
	}

	logger.Info("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	if err := botService.Stop(); err != nil {
		logger.Error("Failed to stop bot: %v", err)
		os.Exit(1)
	}

	logger.Info("Bot has been stopped")
}
