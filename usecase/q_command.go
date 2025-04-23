package usecase

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gong023/umi/domain"
)

type QCommandHandler struct {
	openaiClient domain.OpenAIClient
	logger       domain.Logger
}

func NewQCommandHandler(openaiClient domain.OpenAIClient, logger domain.Logger) *QCommandHandler {
	return &QCommandHandler{
		openaiClient: openaiClient,
		logger:       logger,
	}
}

func (h *QCommandHandler) Handle(s domain.Session, i *domain.InteractionCreate) {
	h.logger.Info("Handling q command")

	// Extract the message from the command options
	var message string
	if i.Data != nil && len(i.Data.Options) > 0 {
		if val, ok := i.Data.Options[0].Value.(string); ok {
			message = val
		}
	}

	if message == "" {
		h.logger.Error("No message provided in q command")
		// Send a response indicating that a message is required
		response := &domain.InteractionResponse{
			Type: int(domain.InteractionResponseChannelMessageWithSource),
			Data: &domain.InteractionResponseData{
				Content: "質問を入力してください。例: `/q 男性は何を飲んでいましたか？`",
			},
		}
		if err := s.InteractionRespond(i, response); err != nil {
			h.logger.Error("Failed to respond to interaction: %v", err)
		}
		return
	}

	// Create a response to acknowledge the command
	response := &domain.InteractionResponse{
		Type: int(domain.InteractionResponseChannelMessageWithSource),
		Data: &domain.InteractionResponseData{
			Content: "質問を処理しています...",
		},
	}

	// Send the initial response
	if err := s.InteractionRespond(i, response); err != nil {
		h.logger.Error("Failed to respond to interaction: %v", err)
		return
	}

	// Check if a quiz exists
	contextDir := filepath.Join("memo", "context")
	files, err := os.ReadDir(contextDir)
	if err != nil {
		if os.IsNotExist(err) {
			h.logger.Info("Context directory does not exist, no quiz available")
			// Send a response indicating that no quiz is available
			followupMessage := "現在クイズが存在しません。`/create` コマンドで新しいクイズを作成してください。"
			h.logger.Info("No quiz available, suggesting /create command: %s", followupMessage)
			return
		}
		h.logger.Error("Failed to read context directory: %v", err)
		return
	}

	// Find the most recent quiz file
	var latestQuizFile string
	var latestModTime int64
	for _, file := range files {
		if file.IsDir() || !strings.HasPrefix(file.Name(), "quiz_") {
			continue
		}

		filePath := filepath.Join(contextDir, file.Name())
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			h.logger.Error("Failed to get file info for %s: %v", filePath, err)
			continue
		}

		modTime := fileInfo.ModTime().Unix()
		if modTime > latestModTime {
			latestModTime = modTime
			latestQuizFile = filePath
		}
	}

	if latestQuizFile == "" {
		h.logger.Info("No quiz file found")
		// Send a response indicating that no quiz is available
		followupMessage := "現在クイズが存在しません。`/create` コマンドで新しいクイズを作成してください。"
		h.logger.Info("No quiz available, suggesting /create command: %s", followupMessage)
		return
	}

	// Read the quiz file
	quizContent, err := os.ReadFile(latestQuizFile)
	if err != nil {
		h.logger.Error("Failed to read quiz file: %v", err)
		return
	}

	// Read the prompt file
	promptPath := filepath.Join("memo", "prompt", "umigame.txt")
	promptContent, err := os.ReadFile(promptPath)
	if err != nil {
		h.logger.Error("Failed to read prompt file: %v", err)
		return
	}

	// Create a request to the OpenAI API
	req := &domain.ChatCompletionRequest{
		Model: "gpt-4-turbo",
		Messages: []domain.ChatMessage{
			{
				Role:    "system",
				Content: string(promptContent) + "\n\nあなたは質問に対して「はい」「いいえ」「わからない/関係ない」のいずれかで答えてください。質問が現在のクイズの解決に関連する場合は、適切な回答を選んでください。質問が現在のクイズの解決に関連しない場合は「わからない/関係ない」と答えてください。",
			},
			{
				Role:    "user",
				Content: "クイズ: " + string(quizContent) + "\n\n質問: " + message,
			},
		},
		Temperature: 0.7,
	}

	// Send the request to the OpenAI API
	h.logger.Info("Sending request to OpenAI API")
	resp, err := h.openaiClient.CreateChatCompletion(req)
	if err != nil {
		h.logger.Error("Failed to create chat completion: %v", err)
		return
	}

	// Extract the answer from the response
	if len(resp.Choices) == 0 {
		h.logger.Error("No choices in response")
		return
	}

	answer := resp.Choices[0].Message.Content
	h.logger.Info("Received answer: %s", answer)

	// Format the answer
	formattedAnswer := fmt.Sprintf("**質問**: %s\n\n**回答**: %s", message, strings.TrimSpace(answer))

	// Create a follow-up message with the answer
	// Note: In a real implementation, you would need to use the Discord API to send a follow-up message
	// This is just a placeholder to demonstrate the concept
	h.logger.Info("Answer created: %s", formattedAnswer)
	
	// In a real implementation, we would send a follow-up message with the answer
	// For now, we'll just log it
	// Example:
	// discordSession.ChannelMessageSend(channelID, formattedAnswer)
}
