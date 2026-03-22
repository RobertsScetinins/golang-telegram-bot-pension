package handlers

import (
	"context"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/app"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Status(ctx context.Context, b *bot.Bot, update *models.Update, app *app.App) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Bot is Active 🚀",
	})
}
