package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/service"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/utils"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func FactCheck(ctx context.Context, b *bot.Bot, update *models.Update) {
	userText := update.Message.Text
	geminiService := service.NewGeminiService()

	parts := strings.SplitN(userText, " ", 2)
	hasTextAfterCommand := len(parts) > 1 && strings.TrimSpace(parts[1]) != ""
	hasReply := update.Message.ReplyToMessage != nil

	var claim string
	var userComment string

	if hasTextAfterCommand {
		userComment = strings.TrimSpace(parts[1])
	}

	if hasReply && hasTextAfterCommand {
		claim = fmt.Sprintf(
			"Комментарий пользователя: [%s] Ответ на сообщение: [%s]",
			userComment,
			update.Message.ReplyToMessage.Text)
	} else if hasReply {
		// TODO: we might need to support other fields as well
		if strings.TrimSpace(update.Message.ReplyToMessage.Text) != "" {
			claim = update.Message.ReplyToMessage.Text
		} else {
			claim = update.Message.ReplyToMessage.Caption
		}
	} else if hasTextAfterCommand {
		claim = userComment
	}

	if strings.TrimSpace(claim) == "" {
		fmt.Println("[WARN] No claim provided by user")
		utils.Reply(ctx, b, update, "Пожалуйста, укажите утверждение после /factcheck или выберите сообщение для ответа")
		return
	}

	fmt.Println("[INFO] User claim:", claim)
	check, err := geminiService.GenResponseWithPreset(ctx, claim, "fact_check")
	if err != nil {
		fmt.Println("[ERROR] Gemini API failed:", err)
		utils.Reply(ctx, b, update, "⚠️ Не удалось проверить факт. Попробуйте позже.")
		return
	}
	fmt.Println("[INFO] Gemini response received")

	utils.Reply(ctx, b, update, check)
}
