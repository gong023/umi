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

	pingCommandHandler := usecase.NewPingCommandHandler(logger)
	createCommandHandler := usecase.NewCreateCommandHandler(openaiClient, logger)

	botService.RegisterCommand("ping", pingCommandHandler)
	botService.RegisterCommand("create", createCommandHandler)

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
