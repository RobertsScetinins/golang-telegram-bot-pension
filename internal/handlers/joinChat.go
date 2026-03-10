package handlers

import (
	"context"
	"log"
	"time"

	chatModel "github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/models"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/repository"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func HandleJoinChat(ctx context.Context, b *bot.Bot, update *models.Update, chatRepo *repository.ChatRepository) {
	botID := b.ID()
	newChatMember := update.MyChatMember.NewChatMember.Member

	if newChatMember.User != nil && newChatMember.User.ID == botID {
		chat := update.MyChatMember.Chat

		currentChat := &chatModel.Chat{
			ChatID:    chat.ID,
			Type:      string(chat.Type),
			Title:     &chat.Title,
			CreatedAt: time.Now(),
		}

		err := chatRepo.Save(ctx, currentChat)
		if err != nil {
			log.Printf("Failed to save chat %d: %v", chat.ID, err)
		}
	}
}
