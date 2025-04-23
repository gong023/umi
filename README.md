# Umi - ウミガメのスープ Discord Bot

A Discord bot for playing ウミガメのスープ (Umigame no Soup) quiz games.

## Features

- `/ping` - Bot responds with "はい!" (Currently implemented)
- `/create` - Create a new quiz
- `/q` - Ask a question about the current quiz
- `/answer` - Submit an answer to the current quiz
- `/info` - Get information about the current quiz
- `/clue` - Get a clue about the current quiz
- `/quit` - Quit the current quiz
- `/help` - Get help with bot commands

## Setup

### Prerequisites

- Go 1.21 or higher
- Discord Bot Token

### Environment Variables

- `DISCORD_TOKEN` - Your Discord bot token

### Running the Bot

1. Clone the repository
2. Set the environment variables
3. Run the bot:

```bash
go run cmd/bot/main.go
```

## Project Structure

This project follows clean architecture principles:

- `cmd/` - Entry points and configuration
- `domain/` - Core domain models and interfaces
- `usecase/` - Application business logic
- `infra/` - External dependencies implementation
- `memo/` - Bot memory storage (gitignored)
- `.github/` - GitHub Actions workflows

## Development

### Testing

Run tests with:

```bash
go test ./...
```

### Linting

Run linting with:

```bash
golangci-lint run ./...
```

## License

MIT
