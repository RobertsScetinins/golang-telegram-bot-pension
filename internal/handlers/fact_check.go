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

	parts := strings.SplitN(userText, " ", 2)
	hasTextAfterCommand := len(parts) > 1 && strings.TrimSpace(parts[1]) != ""
	hasReply := update.Message.ReplyToMessage != nil

	var claim string

	if hasReply && hasTextAfterCommand {
		repliedText := update.Message.ReplyToMessage.Text
		userComment := strings.TrimSpace(parts[1])
		claim = fmt.Sprintf("User comment:[%s] Replied content:[%s]", userComment, repliedText)
	} else if hasReply {
		claim = update.Message.ReplyToMessage.Text
	} else if hasTextAfterCommand {
		claim = strings.TrimSpace(parts[1])

	}

	if strings.TrimSpace(claim) == "" {
		fmt.Println("[WARN] No claim provided by user")
		utils.Reply(ctx, b, update, "Please provide a claim after /factcheck or chose a reply")
		return
	}

	fmt.Println("[INFO] User claim:", claim)

	// Call Gemini API
	check := service.GenResponse(claim)
	fmt.Println("[INFO] Gemini response received")

	utils.Reply(ctx, b, update, check)
}
