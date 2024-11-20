package model

// UserIdentity Информация по отдельному реципиенту
type UserIdentity struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Chat Полная запись пользователя
type Chat struct {
	ID         int64  `redis:"id"`
	Recipients string `redis:"recipients"`
	Body       string `redis:"body"`
	CreatedAt  int64  `redis:"created_at"`
	UpdatedAt  *int64 `redis:"updated_at"`
}
