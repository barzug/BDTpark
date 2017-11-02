package models

import (
	"github.com/jackc/pgx"
	"log"
)

func AddMember(pool *pgx.ConnPool, forum, author string) {
	_, err := pool.Exec(`INSERT INTO members(forum, author)`+
		`VALUES ($1, $2) ON CONFLICT (forum, author) DO NOTHING`, forum, author)
	if err != nil {
		log.Print(err)
	}
}
