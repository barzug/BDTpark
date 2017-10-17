package models

import (
	"time"
	"github.com/jackc/pgx"
	"log"
)

type Posts struct {
	PID      int64     `json:"id"`
	Author   string    `json:"author"`
	Created  time.Time `json:"created"`
	Forum    string    `json:"forum"`
	IsEdited bool      `json:"isedited"`
	Message  string    `json:"message"`
	Parent   int64     `json:"parent"`
	Thread   int64     `json:"thread"`
}

func CreatePostsBySlice(pool *pgx.ConnPool, posts []Posts, threadId int64, created time.Time, forum string) error {
	tx, err := pool.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()


	for i := 0; i < len(posts); i++ {
		path := []int64{}
		pool.QueryRow(`SELECT nextval('"posts_pID_seq"')`).Scan(&posts[i].PID);
		if posts[i].Parent != 0 {
			var parentPath []int64
			err := pool.QueryRow(`SELECT path FROM posts WHERE "pID"=$1 AND thread=$2`, posts[i].Parent, threadId).Scan(&parentPath);
			if err != nil {
				return err
			}
			path = append(path, parentPath...)
		}
		path = append(path, posts[i].PID)

		posts[i].Forum = forum
		posts[i].Thread = threadId
		posts[i].Created = created

		err = tx.QueryRow(`INSERT INTO posts ("pID", message, thread, parent, author, created, forum, path)
										VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING "pID"`,
			&posts[i].PID, &posts[i].Message, &posts[i].Thread, &posts[i].Parent, &posts[i].Author, &posts[i].Created, &posts[i].Forum, &path).Scan(&posts[i].PID);
		if err != nil {
			log.Print(err)
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
