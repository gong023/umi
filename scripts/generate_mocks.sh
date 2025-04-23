#!/bin/bash

# This script generates mock files for interfaces using mockgen

# Ensure the mock directory exists
mkdir -p infra/mock

# Generate mock for DiscordClient
mockgen -destination=infra/mock/discord.go -package=mock github.com/gong023/umi/domain DiscordClient

# Generate mock for CommandHandler
mockgen -destination=infra/mock/command.go -package=mock github.com/gong023/umi/domain CommandHandler

# Generate mock for Session
mockgen -destination=infra/mock/session.go -package=mock github.com/gong023/umi/domain Session

# Generate mock for Logger
mockgen -destination=infra/mock/logger.go -package=mock github.com/gong023/umi/domain Logger

echo "Mock files generated successfully"
