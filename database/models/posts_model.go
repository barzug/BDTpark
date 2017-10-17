package models

import (
	"time"
	"github.com/jackc/pgx"

)

type Posts struct {
	PID      int64     `json:"id"`
	Author   string    `json:"author"`
	Created  time.Time `json:"created"`
	Forum    string    `json:"forum"`
	IsEdited bool      `json:"isedited"`
	Message  string    `json:"message"`
	Parent   int64     `json:"parents"`
	Thread   int64     `json:"thread"`
}

func CreatePostsBySlice(pool *pgx.ConnPool, posts []Posts, threadId int64, created time.Time, forum string) error {
	tx, err := pool.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()


	for i := 0; i < len(posts); i++ {
		var id int64
		if posts[i].Parent != 0 {
			err := pool.QueryRow("SELECT parent FROM post WHERE id=$1 AND thread=$2", posts[i].Parent, threadId).Scan(&id);
			if err != nil {
				return err
			}
		}

		posts[i].Forum = forum
		posts[i].Thread = threadId
		posts[i].Created = created

		err = tx.QueryRow(`INSERT INTO posts (message, thread, parent, author, created, forum)
										VALUES ($1, $2, $3, $4, $5, $6) RETURNING "pID"`,
			&posts[i].Message, &posts[i].Thread, &posts[i].Parent, &posts[i].Author, &posts[i].Created, &posts[i].Forum).Scan(&posts[i].PID);
		if err != nil {
			return err;
		}
	}

	_, err = tx.Exec(`UPDATE forums SET posts = posts + $1 WHERE slug = $2`, len(posts), forum)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
