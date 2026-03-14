package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/helpers"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/service"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Look(ctx context.Context, b *bot.Bot, update *models.Update) {
	token := b.Token()
	geminiService := service.NewGeminiService()

	userText := update.Message.Caption
	parts := strings.SplitN(userText, " ", 2)
	hasTextAfterCommand := len(parts) > 1 && strings.TrimSpace(parts[1]) != ""

	var userComment string

	if hasTextAfterCommand {
		userComment = strings.TrimSpace(parts[1])
	}

	mediaData, err := helpers.ProcessMedia(update)
	if err != nil {
		fmt.Println("[WARN] Failed to process media")
		utils.Reply(ctx, b, update, "Не удалось обработать медиафайл.")
		return
	}

	fileLink, err := service.GetDownloadLink(b, helpers.GetFileUrl(mediaData.FileId, token))
	if err != nil {
		fmt.Println("[WARN]", fmt.Sprintf("[WARN] Failed to get download link for file %s: %v", mediaData.FileId, err))
		utils.Reply(ctx, b, update, "Не удалось получить ссылку для скачивания. Попробуйте позже.")
		return
	}

	ans, err := geminiService.GenResponseWithMediaPreset(ctx, userComment, fileLink, service.PromptTypeAnalyze)
	if err != nil {
		fmt.Println("[WARN] Gemini API failed:", err)
		utils.Reply(ctx, b, update, "Не удалось обработать изображение. Возможно, формат не поддерживается или сервис временно недоступен.")
		return
	}

	utils.Reply(ctx, b, update, ans)
}
