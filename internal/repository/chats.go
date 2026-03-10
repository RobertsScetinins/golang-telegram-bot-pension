package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/Dmitrijs-Vasilevskis/go-telegram-bot/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatRepository struct {
	db *pgxpool.Pool
}

func NewChatRepository(db *pgxpool.Pool) *ChatRepository {
	return &ChatRepository{db}
}

func (r *ChatRepository) Save(ctx context.Context, m *models.Chat) error {
	query := `
	INSERT INTO chats
	(chat_id, type, title, created_at)
	VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		m.ChatID,
		m.Type,
		m.Title,
		m.CreatedAt,
	)

	return err
}

func (r *ChatRepository) GetChatByID(ctx context.Context, chatID int64) (*models.Chat, error) {
	query := `
	SELECT id, chat_id, type, title
	FROM chats 
	WHERE chat_id=$1
	`

	var chat models.Chat

	err := r.db.QueryRow(ctx, query, chatID).Scan(
		&chat.ID,
		&chat.ChatID,
		&chat.Type,
		&chat.Title,
	)

	if err != nil {
		return nil, err
	}

	return &chat, nil
}

func (r *ChatRepository) DeleteByChatId(ctx context.Context, chatId int64) error {
	query := `DELETE FROM chats WHERE chat_id=$1`

	log.Println(">> DeleteByChatId")

	result, err := r.db.Exec(ctx, query, chatId)
	if err != nil {
		return fmt.Errorf("Failed to save chat %w", err)
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
