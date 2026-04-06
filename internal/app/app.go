package app

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/repository"
	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/service"
)

type App struct {
	DB *pgxpool.Pool

	MessageRepository *repository.MessageRepository
	GeminiService     service.GeminiService
}

func New(db *pgxpool.Pool) *App {
	return &App{
		DB: db,

		MessageRepository: repository.NewMessageRepository(db),
		GeminiService:     *service.NewGeminiService(),
	}
}
