package usecase

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gong023/umi/domain"
)

type QCommandHandler struct {
	openaiClient domain.OpenAIClient
	fileSystem   domain.FileSystem
	logger       domain.Logger
}

func NewQCommandHandler(openaiClient domain.OpenAIClient, fileSystem domain.FileSystem, logger domain.Logger) *QCommandHandler {
	return &QCommandHandler{
		openaiClient: openaiClient,
		fileSystem:   fileSystem,
		logger:       logger,
	}
}

func (h *QCommandHandler) Handle(s domain.Session, i *domain.InteractionCreate) {
	h.logger.Info("Handling q command")

	// Extract the message from the command options
	var message string

	// Log the entire interaction data for debugging
	h.logger.Info("Interaction data: %+v", i)

	// Check if the interaction data contains the message directly
	if i.Original != nil {
		originalInteraction, ok := i.Original.(*discordgo.InteractionCreate)
		if ok {
			h.logger.Info("Original interaction: %+v", originalInteraction)

			// Try to get the command data
			if originalInteraction.Type == discordgo.InteractionApplicationCommand {
				data := originalInteraction.ApplicationCommandData()
				h.logger.Info("Application command data: %+v", data)

				// Check if there are options
				if len(data.Options) > 0 {
					h.logger.Info("Command options found: %d options", len(data.Options))
					for idx, opt := range data.Options {
						h.logger.Info("Option %d: Name=%s, Value=%v", idx, opt.Name, opt.Value)
					}

					// Try to extract the message from the options
					for _, opt := range data.Options {
						if opt.Name == "message" {
							if val, ok := opt.Value.(string); ok {
								message = val
								h.logger.Info("Extracted question from user: %s", message)
							} else {
								h.logger.Error("Failed to extract question: Value is not a string: %T", opt.Value)
							}
							break
						}
					}
				} else {
					h.logger.Error("No options found in ApplicationCommandData")
				}
			} else {
				h.logger.Error("Interaction is not an ApplicationCommand: %d", originalInteraction.Type)
			}
		} else {
			h.logger.Error("Original interaction is not of type *discordgo.InteractionCreate: %T", i.Original)
		}
	} else if i.Data != nil && len(i.Data.Options) > 0 {
		h.logger.Info("Command options found in domain model: %d options", len(i.Data.Options))
		for idx, opt := range i.Data.Options {
			h.logger.Info("Option %d: Name=%s, Value=%v", idx, opt.Name, opt.Value)
		}

		if val, ok := i.Data.Options[0].Value.(string); ok {
			message = val
			h.logger.Info("Extracted question from user: %s", message)
		} else {
			h.logger.Error("Failed to extract question: Value is not a string: %T", i.Data.Options[0].Value)
		}
	} else {
		h.logger.Error("No command options found: Data=%v", i.Data)
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
	contextPath := h.fileSystem.JoinPath("memo", "context.txt")
	quizExists := false
	var contextContent []byte
	var conversationHistory []string

	// Check if the context file exists
	exists, err := h.fileSystem.FileExists(contextPath)
	if err != nil {
		h.logger.Error("Failed to check if context file exists: %v", err)
		return
	}

	if exists {
		// Read the existing context file
		contextContent, err = h.fileSystem.ReadFile(contextPath)
		if err != nil {
			h.logger.Error("Failed to read context file: %v", err)
			return
		}

		// If the file is not empty, a quiz exists
		if len(contextContent) > 0 {
			quizExists = true
			// Split the content by new lines to get the conversation history
			conversationHistory = strings.Split(string(contextContent), "\n")
		}
	}

	if !quizExists {
		h.logger.Info("No quiz found")
		// Send a response indicating that no quiz is available
		followupMessage := "現在クイズが存在しません。`/create` コマンドで新しいクイズを作成してください。"

		// Send the response with the message
		if err := s.FollowupMessage(i, followupMessage); err != nil {
			h.logger.Error("Failed to send follow-up message: %v", err)
		}

		h.logger.Info("No quiz available, suggesting /create command: %s", followupMessage)
		return
	}

	// Read the prompt file
	promptPath := h.fileSystem.JoinPath("memo", "prompt", "onQ.txt")
	promptContent, err := h.fileSystem.ReadFile(promptPath)
	if err != nil {
		h.logger.Error("Failed to read prompt file: %v", err)
		return
	}

	// Create a request to the OpenAI API with the conversation history
	messages := []domain.ChatMessage{
		{
			Role:    "system",
			Content: string(promptContent),
		},
	}

	// Process the conversation history and add it to the messages
	if len(conversationHistory) > 0 {
		// The first line is always the quiz (assistant's role)
		messages = append(messages, domain.ChatMessage{
			Role:    "assistant",
			Content: conversationHistory[0],
		})

		// Process the rest of the conversation history
		currentRole := "user" // Start with user after the quiz
		var currentMessage string

		for i := 1; i < len(conversationHistory); i++ {
			line := strings.TrimSpace(conversationHistory[i])
			if line == "" {
				continue // Skip empty lines
			}

			// Check if this line starts a new message
			if strings.HasPrefix(line, "質問: ") {
				// If we have a message in progress, add it
				if currentMessage != "" {
					messages = append(messages, domain.ChatMessage{
						Role:    currentRole,
						Content: currentMessage,
					})
				}
				// Start a new user message
				currentRole = "user"
				currentMessage = line
			} else {
				// If this is not a question line, it's part of the assistant's response
				if currentRole == "user" && currentMessage != "" {
					// Add the completed user message
					messages = append(messages, domain.ChatMessage{
						Role:    currentRole,
						Content: currentMessage,
					})
					// Start a new assistant message
					currentRole = "assistant"
					currentMessage = line
				} else if currentRole == "assistant" {
					// Continue the assistant's message
					currentMessage += "\n" + line
				} else {
					// Start a new assistant message
					currentRole = "assistant"
					currentMessage = line
				}
			}
		}

		// Add the last message if there is one
		if currentMessage != "" {
			messages = append(messages, domain.ChatMessage{
				Role:    currentRole,
				Content: currentMessage,
			})
		}
	}

	// Add the current question
	messages = append(messages, domain.ChatMessage{
		Role:    "user",
		Content: "質問: " + message,
	})

	// Log the messages being sent to OpenAI
	h.logger.Info("Sending the following messages to OpenAI:")
	for i, msg := range messages {
		h.logger.Info("Message %d - Role: %s, Content: %s", i, msg.Role, msg.Content)
	}

	req := &domain.ChatCompletionRequest{
		Model:       "chatgpt-4o-latest",
		Messages:    messages,
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

	// Append the question and answer to the context file
	// First, read the existing content
	existingContent := ""
	if len(contextContent) > 0 {
		existingContent = string(contextContent)
	}

	// Append the new question and answer with newlines
	updatedContent := existingContent
	if len(updatedContent) > 0 && !strings.HasSuffix(updatedContent, "\n") {
		updatedContent += "\n"
	}

	// Add the user's question
	userQuestion := "質問: " + message
	updatedContent += userQuestion + "\n"

	// Add the assistant's answer
	updatedContent += answer

	// Write the updated content back to the context file
	if err := h.fileSystem.WriteFile(contextPath, []byte(updatedContent), 0644); err != nil {
		h.logger.Error("Failed to update context file: %v", err)
		// Continue with the response even if we fail to update the context file
	} else {
		h.logger.Info("Updated context file with new question and answer")
	}

	// Send the response with the answer
	if err := s.FollowupMessage(i, formattedAnswer); err != nil {
		h.logger.Error("Failed to send follow-up message: %v", err)
	}

	h.logger.Info("Answer created: %s", formattedAnswer)
}
