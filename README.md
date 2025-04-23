note: this repository is mostly written by AI.

# ウミガメのスープ Discord Bot

This is a Discord bot for the ウミガメのスープ (Umigame no Soup) quiz game.

## Features

- `/ping` command: A simple ping command to check if the bot is running
- `/quiz` command: Generates a new ウミガメのスープ quiz using OpenAI

## Technical Stack

- Go 1.23.2
- [discordgo](https://github.com/bwmarrin/discordgo): Discord API client for Go
- [OpenAI API](https://platform.openai.com/docs/api-reference): Used for generating quizzes
- [go.uber.org/mock](https://github.com/uber-go/mock): For generating mocks for testing

## Architecture

This project follows the clean architecture pattern:

- `domain`: Contains interfaces and domain models
- `infra`: Contains implementations of external dependencies (Discord, OpenAI)
- `usecase`: Contains the core business logic
- `cmd`: Contains the entry points for the application

## Setup

### Prerequisites

- Go 1.23.2 or higher
- Discord Bot Token
- OpenAI API Key

### Environment Variables

The following environment variables are required:

- `DISCORD_TOKEN`: Your Discord bot token
- `OPENAI_API_KEY`: Your OpenAI API key

### Running the Bot

```bash
go run cmd/bot/main.go
```

## Development

### Generating Mocks

To generate mocks for testing, run:

```bash
bash scripts/generate_mocks.sh
```

### Running Tests

To run tests, use:

```bash
go test ./...
```

## OpenAI Integration

The bot uses the OpenAI API to generate ウミガメのスープ quizzes. The integration is implemented as follows:

1. `domain/openai.go`: Defines the OpenAI client interface and related types
2. `infra/openai.go`: Implements the OpenAI client
3. `usecase/quiz_command.go`: Uses the OpenAI client to generate quizzes

The OpenAI client is initialized in `cmd/bot/main.go` and passed to the `QuizCommandHandler`.

### Quiz Generation

The bot uses the GPT-4 Turbo model to generate quizzes. The prompt is designed to create ウミガメのスープ style quizzes in Japanese.

## License

This project is licensed under the MIT License.
