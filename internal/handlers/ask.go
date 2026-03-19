package handlers

import (
	"context"
	"fmt"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/helpers"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/service"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Ask(ctx context.Context, b *bot.Bot, update *models.Update) {
	message := update.Message
	if message == nil {
		return
	}

	hasMedia := message.Caption != ""

	token := b.Token()
	geminiService := service.NewGeminiService()

	var userInput string
	var mediaData helpers.MediaData

	if hasMedia {
		media, err := helpers.ProcessMedia(update)
		if err != nil {
			fmt.Println("[WARN] Failed to process media")
			utils.Reply(ctx, b, update, "Не удалось обработать медиафайл.")
			return
		}

		userInput = message.Caption
		mediaData = *media
	} else {
		userInput = message.Text
	}

	userComment, _ := helpers.GetCommandArgs(userInput)

	if hasMedia {
		fileLink, err := service.GetDownloadLink(b, helpers.GetFileUrl(mediaData.FileId, token))
		if err != nil {
			fmt.Println("[WARN]", fmt.Sprintf("Failed to get download link for file %s: %v", mediaData.FileId, err))
			utils.Reply(ctx, b, update, "Не удалось получить ссылку для скачивания. Попробуйте позже.")
			return
		}

		response, err := geminiService.GenResponseWithMediaPreset(ctx, userComment, fileLink, service.PromptTypeCustom)
		if err != nil {
			fmt.Println("[WARN] Gemini API failed:", err)
			utils.Reply(ctx, b, update, "Не удалось обработать изображение. Возможно, формат не поддерживается или сервис временно недоступен.")
			return
		}

		utils.Reply(ctx, b, update, response)

	} else {
		response, err := geminiService.GenResponseWithPreset(ctx, userComment, service.PromptTypeCustom)
		if err != nil {
			fmt.Println("[ERROR] Gemini API failed:", err)
			utils.Reply(ctx, b, update, "⚠️ Не удалось обработать запрос. Попробуйте позже.")
			return
		}

		utils.Reply(ctx, b, update, response)
	}
}
