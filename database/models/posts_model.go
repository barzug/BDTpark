package models

import (
	"sync"
	"time"

	"../../utils"
	"github.com/jackc/pgx"
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

	waitData := &sync.WaitGroup{}

	muteForInsert := &sync.Mutex{}

	var notFoundErr, conflictErr error
	for index := 0; index < len(posts); index++ {
		waitData.Add(1)

		var valFromSeq int64
		err := pool.QueryRow(`SELECT nextval('"posts_pID_seq"')`).Scan(&valFromSeq)

		go func(waitData *sync.WaitGroup, i int, pIDfFromSeq int64) {
			posts[i].PID = pIDfFromSeq
			defer waitData.Done()
			path := []int64{}

			if posts[i].Parent != 0 {
				var parentPath []int64
				err := pool.QueryRow(`SELECT path FROM posts WHERE "pID"=$1 AND thread=$2`, posts[i].Parent, threadId).Scan(&parentPath)
				if err != nil {
					conflictErr = err
				}
				path = append(path, parentPath...)
			}
			path = append(path, posts[i].PID)

			posts[i].Forum = forum
			posts[i].Thread = threadId

			var nickname string
			err = pool.QueryRow(`SELECT nickname FROM users WHERE nickname = $1`, posts[i].Author).Scan(&nickname)
			if err != nil {
				notFoundErr = utils.NotFoundError
			}

			muteForInsert.Lock()
			err = tx.QueryRow(`INSERT INTO posts ("pID", message, thread, parent, author, created, forum, path)
										VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING "pID", created`,
				&posts[i].PID, &posts[i].Message, &posts[i].Thread, &posts[i].Parent, &posts[i].Author, created.Format(time.RFC3339), &posts[i].Forum, &path).Scan(&posts[i].PID, &posts[i].Created)
			if err != nil {
				conflictErr = err
			}
			muteForInsert.Unlock()

		}(waitData, index, valFromSeq)

	}

	waitData.Wait()

	if notFoundErr != nil {
		return notFoundErr
	}

	if conflictErr != nil {
		return conflictErr
	}
	tx.Exec(`UPDATE forums SET posts = posts + $1 WHERE slug = $2`, len(posts), forum)
	err = tx.Commit()
	if err != nil {
		return err
	}

	for i := 0; i < len(posts); i++ {
		go AddMember(pool, posts[i].Forum, posts[i].Author)
	}
	return nil
}

func (post *Posts) GetPostById(pool *pgx.ConnPool) (Posts, error) {
	resultPost := Posts{PID: post.PID}
	err := pool.QueryRow(`SELECT author, created, forum, message, thread, "isEdited", parent FROM posts WHERE "pID" = $1`,
		post.PID).Scan(&resultPost.Author, &resultPost.Created, &resultPost.Forum, &resultPost.Message, &resultPost.Thread, &resultPost.IsEdited, &resultPost.Parent)

	if err != nil {
		return resultPost, err
	}
	return resultPost, nil
}

func (post *Posts) UpdatePost(pool *pgx.ConnPool) error {
	err := pool.QueryRow(`UPDATE posts SET message = $1, "isEdited" = true`+
		` WHERE "pID" = $2 RETURNING author, created, forum, "isEdited", thread;`,
		post.Message, post.PID).Scan(&post.Author, &post.Created, &post.Forum, &post.IsEdited, &post.Thread)
	if err != nil {
		if pgerr, ok := err.(pgx.PgError); ok {
			if pgerr.ConstraintName == "post_pk" {
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
