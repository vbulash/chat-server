package model

import (
	"database/sql"
	"time"
)

// UserIdentity Информация по реципиентам
type UserIdentity struct {
	ID    int64
	Name  string
	Email string
}

// ChatInfo Краткая информация по чату
type ChatInfo struct {
	Recipients []*UserIdentity
	Body       string
}

// Chat Полная запись пользователя
type Chat struct {
	ID        int64
	Info      ChatInfo
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
