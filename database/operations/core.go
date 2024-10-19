package operations

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/vbulash/chat-server/config"
)

func InitDb() (*sqlx.DB, error) {
	return sqlx.Connect("postgres", fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		config.Config.Host,
		config.Config.Port,
		config.Config.Database,
		config.Config.Username,
		config.Config.Password,
	))
}
