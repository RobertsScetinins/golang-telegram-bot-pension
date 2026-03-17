package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/app"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/helpers"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/logger"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/service"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Summary(ctx context.Context, b *bot.Bot, update *models.Update, app *app.App) {

	logger.DebugJson(update)
	geminiService := service.NewGeminiService()

	messages, err := app.MessageRepository.GetLastMessages(ctx, update.Message.Chat.ID, 200)
	if err != nil {
		fmt.Println("Failed to get messages:", err)
		return
	}

	processedMessages, _ := helpers.ProcessSummary(messages)

	summaryJSON, err := json.Marshal(processedMessages)
	if err != nil {
		return
	}

	ans, err := geminiService.GenResponseWithPreset(ctx, string(summaryJSON), service.PromptTypeSummary)

	logger.DebugJson(ans)

	utils.Reply(ctx, b, update, ans)
}
