package models

import (
	"time"
	"github.com/jackc/pgx"
	"../../utils"
)

type Threads struct {
	TID     int64     `json:"id"`
	Author  string    `json:"author"`
	Created time.Time `json:"created"`
	Forum   string    `json:"forum"`
	Message string    `json:"message"`
	Slug    string    `json:"slug"`
	Title   string    `json:"title"`
	Votes   int32     `json:"votes"`
}

func (thread *Threads) CreateThread(pool *pgx.ConnPool) error {
	var id int64

	slug := thread.Slug
	if slug == "" {
		slug = thread.Forum
	}

	err := pool.QueryRow(`INSERT INTO threads (author, created, message, slug, title, forum)`+
		`VALUES ($1, $2, $3, $4, $5, $6) RETURNING "tID";`,
		thread.Author, thread.Created, thread.Message, slug, thread.Title, thread.Forum).Scan(&id)
	if err != nil {
		if pgerr, ok := err.(pgx.PgError); ok {
			if pgerr.ConstraintName == "threads_slug_key" {
				return utils.UniqueError
			} else {
				return err
			}
		}
		return err
	}
	thread.TID = id
	return nil
}

func (thread *Threads) GetThreadBySlug(pool *pgx.ConnPool) (Threads, error) {
	resultThread := Threads{}
	//resultThread.Slug = thread.Slug
	err := pool.QueryRow(`SELECT "tID", author, created, forum, message, title, votes, slug FROM threads WHERE slug = $1`,
		thread.Slug).Scan(&resultThread.TID, &resultThread.Author, &resultThread.Created, &resultThread.Forum,
		&resultThread.Message, &resultThread.Title, &resultThread.Votes, &resultThread.Slug)

	if err != nil {
		return resultThread, err
	}
	return resultThread, nil
}


func (thread *Threads) GetThreadById(pool *pgx.ConnPool) (Threads, error) {
	resultThread := Threads{}
	resultThread.TID = thread.TID
	err := pool.QueryRow(`SELECT author, created, forum, message, slug, title, votes FROM threads WHERE "tID" = $1`,
		thread.TID).Scan(&resultThread.Author, &resultThread.Created, &resultThread.Forum,
		&resultThread.Message, &resultThread.Slug, &resultThread.Title, &resultThread.Votes)

	if err != nil {
		return resultThread, err
	}
	return resultThread, nil
}
