package operations

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/vbulash/chat-server/config"
)

func Get(db *sqlx.DB) (*[]config.NoteType, error) {
	notes := []config.NoteType{}
	err := db.Select(&notes, "SELECT * FROM notes")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &notes, nil
}
