package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/app"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/database"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/handlers"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/helpers"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/router"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	dbUrl := os.Getenv("DATABASE_URL")

	if token == "" {
		log.Fatal("Telegram bot token is missing in configuration file")
	}

	if dbUrl == "" {
		log.Fatal("Database url is missing in configuration file")
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	db, err := database.PostgresPool(dbUrl)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}
	defer db.Close()

	app := app.New(db)

	log.Println("Database connected")

	botClient, err := bot.New(token)
	if err != nil {
		log.Fatal(err)
	}

	r := router.NewRouter()

	r.Register("instagram", handlers.Instagram)
	r.Register("tiktok", handlers.TikTok)
	r.Register("", func(ctx context.Context, b *bot.Bot, update *models.Update) {
		handlers.RecordMessage(ctx, b, update, db)
	})

	botClient.RegisterHandler(bot.HandlerTypeMessageText, "/status", bot.MatchTypeExact,
		func(ctx context.Context, botClient *bot.Bot, update *models.Update) {
			botClient.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Bot is Active 🚀",
			})
		})

	botClient.RegisterHandler(bot.HandlerTypeMessageText, "/factcheck", bot.MatchTypePrefix,
		func(ctx context.Context, botClient *bot.Bot, update *models.Update) {
			handlers.FactCheck(ctx, botClient, update)
		})

	botClient.RegisterHandlerMatchFunc(func(update *models.Update) bool {
		return update.Message != nil && (strings.HasPrefix(update.Message.Text, "/look") || strings.HasPrefix(update.Message.Caption, "/look"))
	},
		func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			handlers.Look(ctx, bot, update)
		})

	botClient.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypeContains,
		func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			text := update.Message.Text

			if helpers.IsToxic(text) {
				handlers.Clown(ctx, bot, update)
			}

			// depricated for now, maybe one day ...
			// if update.Message.Chat.Type == models.ChatTypePrivate {
			// 	handlers.Duplicator(ctx, bot, update)
			// }

			r.Handle(ctx, bot, update)
		})

	botClient.RegisterHandlerMatchFunc(func(update *models.Update) bool {
		return update != nil && update.MyChatMember != nil &&
			(update.MyChatMember.NewChatMember.Member != nil || update.MyChatMember.OldChatMember.Member != nil)
	},
		func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			botID := bot.ID()
			myChatMember := update.MyChatMember

			newMember := myChatMember.NewChatMember

			if newMember.Type == models.ChatMemberTypeLeft && newMember.Left.User.IsBot &&
				(newMember.Left.User.ID == botID) {
				handlers.HandleLeaveChat(ctx, bot, update, app.ChatRepository)
			}

			if newMember.Type == models.ChatMemberTypeMember && newMember.Member.User.IsBot &&
				(newMember.Member.User.ID == botID) {
				handlers.HandleJoinChat(ctx, bot, update, app.ChatRepository)
			}
		})

	botClient.RegisterHandlerMatchFunc(func(update *models.Update) bool {
		return update != nil && update.EditedMessage != nil
	}, func(ctx context.Context, bot *bot.Bot, update *models.Update) {
		handlers.UpdateMessage(ctx, bot, update, db)
	})

	log.Println("Bot started")

	botClient.Start(ctx)
}
