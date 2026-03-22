package dispatcher

import (
	"context"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/app"
	ch "github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/handlers/command"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Commandhandler func(ctx context.Context, b *bot.Bot, update *models.Update, app *app.App)

var commandHandlers = map[string]Commandhandler{
	"factcheck": ch.FactCheck,
	"summary":   ch.Summary,
	"look":      ch.Look,
	"ask":       ch.Ask,
	"status":    ch.Status,
}
