package models

import "time"

type Message struct {
	ID        int64     `db:"id"`
	ChatID    int64     `db:"chat_id"`
	MessageId int64     `db:"message_id"`
	Username  *string   `db:"username"`
	Text      *string   `db:"text"`
	IsEdited  bool      `db:"is_edited"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
