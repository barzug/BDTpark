package models

import (
	"github.com/jackc/pgx"
)

func AddMember(pool *pgx.ConnPool, forum, author string) {
	pool.Exec(`INSERT INTO members(forum, author)`+
		`VALUES ($1, $2) ON CONFLICT (forum, author) DO NOTHING`, forum, author)

}
