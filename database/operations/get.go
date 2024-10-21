package operations

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	// pq Используется sqlx для работы с postgres
	_ "github.com/lib/pq"
	"github.com/vbulash/chat-server/config"
)

// Get Получение данных из таблицы notes
func Get(db *sqlx.DB) (*[]config.NoteType, error) {
	notes := []config.NoteType{}
	err := db.Select(&notes, "SELECT * FROM notes")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &notes, nil
}
