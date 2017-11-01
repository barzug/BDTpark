package models

import (
	"github.com/jackc/pgx"
)

func AddMember(tx *pgx.Tx, forum, author string) {
	tx.Exec(`INSERT INTO members(forum, author)`+
		`VALUES ($1, $2) ON CONFLICT (forum, author) DO NOTHING`, forum, author)
}