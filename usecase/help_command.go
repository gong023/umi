package usecase

import (
	"github.com/gong023/umi/domain"
)

type HelpCommandHandler struct {
	logger domain.Logger
}

func NewHelpCommandHandler(logger domain.Logger) *HelpCommandHandler {
	return &HelpCommandHandler{
		logger: logger,
	}
}

func (h *HelpCommandHandler) Handle(s domain.Session, i *domain.InteractionCreate) {
	h.logger.Info("Handling help command")

	// Create a response with the help message
	helpMessage := `**ウミガメのスープクイズボットの使い方**

以下のコマンドが利用可能です：

- **/create** - 新しいクイズを作成します。クイズが既に存在する場合は、現在のクイズを表示します。
- **/q [質問]** - クイズに関する質問をします。回答は「はい」「いいえ」「わからない/関係ない」のいずれかになります。
- **/answer [回答]** - クイズの答えを提出します。正解の場合はクイズが終了し、不正解の場合はクイズが続行されます。
- **/info** - 現在のクイズとこれまでの質問と回答の履歴を要約します。
- **/clue** - 現在のクイズに関するヒントを提供します。
- **/quit** - 現在のクイズを終了します。
- **/ping** - ボットが応答可能かどうかを確認します。
- **/help** - このヘルプメッセージを表示します。

クイズを始めるには、まず **/create** コマンドを使用してください。`

	response := &domain.InteractionResponse{
		Type: int(domain.InteractionResponseChannelMessageWithSource),
		Data: &domain.InteractionResponseData{
			Content: helpMessage,
		},
	}

	// Send the response
	if err := s.InteractionRespond(i, response); err != nil {
		h.logger.Error("Failed to respond to interaction: %v", err)
		return
	}

	h.logger.Info("Help message sent")
}
