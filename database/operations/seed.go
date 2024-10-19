package operations

import (
	"crypto/rand"
	"github.com/brianvoe/gofakeit"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"math/big"
)

func Seed(db *sqlx.DB) bool {
	tx := db.MustBegin()
	// Рандомное количество добавляемых записей
	nBig, err := rand.Int(rand.Reader, big.NewInt(29))
	if err != nil {
		panic(err)
	}
	records := make([]string, nBig.Int64()+1) // 1 .. 30
	//
	for range records {
		tx.MustExec(
			"INSERT INTO notes (title, body) VALUES ($1, $2)", gofakeit.Question(), gofakeit.Sentence(20))
	}
	err = tx.Commit()
	if err != nil {
		return false
	}
	return true
}
