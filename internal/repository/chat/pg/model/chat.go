package model

import (
	"database/sql"
	"time"
)

// UserIdentity Информация по реципиентам
type UserIdentity struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ChatInfo Краткая информация по чату
type ChatInfo struct {
	Recipients []*UserIdentity `db:"recipients"`
	Body       string          `db:"body"`
}

// Chat Полная запись пользователя
type Chat struct {
	ID        int64        `db:"id"`
	Info      ChatInfo     ``
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
