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

	if len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		fmt.Println("[WARN] No claim provided by user")
		utils.Reply(ctx, b, update, "Please provide a claim after /factcheck")
		return
	}

	claim := strings.TrimSpace(parts[1])
	fmt.Println("[INFO] User claim:", claim)

	// Call Gemini API
	check := service.GenResponse(claim)
	fmt.Println("[INFO] Gemini response received")

	utils.Reply(ctx, b, update, check)

}
