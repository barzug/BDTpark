package models

import (
	"time"
	"github.com/jackc/pgx"
	"../../utils"
)

type Posts struct {
	PID      int64     `json:"id"`
	Author   string    `json:"author"`
	Created  time.Time `json:"created"`
	Forum    string    `json:"forum"`
	IsEdited bool      `json:"isEdited"`
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

		var nickname string
		err := pool.QueryRow(`SELECT nickname FROM users WHERE nickname = $1`, posts[i].Author).Scan(&nickname)
		if err != nil {
			tx.Rollback()
			return utils.NotFoundError
		}

		err = tx.QueryRow(`INSERT INTO posts ("pID", message, thread, parent, author, created, forum, path)
										VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING "pID"`,
			&posts[i].PID, &posts[i].Message, &posts[i].Thread, &posts[i].Parent, &posts[i].Author, &posts[i].Created, &posts[i].Forum, &path).Scan(&posts[i].PID);
		if err != nil {
			return err;
		}

		AddMember(tx, posts[i].Forum, posts[i].Author)
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



func (post *Posts) GetPostById(pool *pgx.ConnPool) (Posts, error) {
	resultPost := Posts{PID:post.PID}
	err := pool.QueryRow(`SELECT author, created, forum, message, thread, "isEdited" FROM posts WHERE "pID" = $1`,
		post.PID).Scan(&resultPost.Author, &resultPost.Created, &resultPost.Forum, &resultPost.Message, &resultPost.Thread, &resultPost.IsEdited)

	if err != nil {
		return resultPost, err
	}
	return resultPost, nil
}


func (post *Posts) UpdatePost(pool *pgx.ConnPool) error {
	var id int64
	err := pool.QueryRow(`UPDATE posts SET message = $1, "isEdited" = true`+
		` WHERE "pID" = $2 RETURNING "pID", author, created, forum, "isEdited", thread;`,
		post.Message, post.PID).Scan(&id, &post.Author, &post.Created, &post.Forum, &post.IsEdited, &post.Thread)
	if err != nil {
		if pgerr, ok := err.(pgx.PgError); ok {
			if pgerr.ConstraintName == "post_pk"  {
				return utils.UniqueError
			} else {
				return err
			}
		}
		return err
	}
	return nil
}

func PostsCount(pool *pgx.ConnPool) (int32, error) {
	var count int32
	err := pool.QueryRow("SELECT COUNT(*) FROM posts").Scan(&count)
	return count, err
}