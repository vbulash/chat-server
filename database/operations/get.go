package operations

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	// pq Используется sqlx для работы с postgres
	_ "github.com/lib/pq"
	"github.com/vbulash/chat-server/config"
)

// Get Получение данных из таблицы chats
func Get(db *sqlx.DB) (*[]config.ChatType, error) {
	chats := []config.ChatType{}
	err := db.Select(&chats, "SELECT id, title, body, created_at, updated_at FROM chats")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &chats, nil
}
