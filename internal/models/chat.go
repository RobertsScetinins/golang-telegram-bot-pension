package models

import "time"

type Chat struct {
	ID        int64     `db:"id"`
	ChatID    int64     `db:"chat_id"`
	Type      string    `db:"type"`
	Title     *string   `db:"title"`
	CreatedAt time.Time `db:"created_at"`
}
