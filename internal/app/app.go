package app

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/repository"
)

type App struct {
	DB *pgxpool.Pool

	MessageRepository *repository.MessageRepository
}

func New(db *pgxpool.Pool) *App {
	return &App{
		DB: db,

		MessageRepository: repository.NewMessageRepository(db),
	}
}
