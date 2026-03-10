package handlers

import (
	"context"
	"log"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/repository"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func HandleLeaveChat(ctx context.Context, b *bot.Bot, update *models.Update, chatRepo *repository.ChatRepository) {
	botID := b.ID()
	newChatMember := update.MyChatMember.NewChatMember.Left

	if newChatMember.User != nil && newChatMember.User.ID == botID {
		chat := update.MyChatMember.Chat

		err := chatRepo.DeleteByChatId(ctx, chat.ID)
		if err != nil {
			log.Printf("Failed to delete chat %d: %v", chat.ID, err)
		}
	}
}
