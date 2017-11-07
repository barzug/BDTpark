package models

import (
	"github.com/jackc/pgx"
	"strconv"
	"log"
)

type Votes struct {
	Voice  int32  `json:"voice"`
	User   string `json:"nickname"`
	Thread int64  `json:"thread"`
}


//UPDATE t
//SET
//t.votes = t.votes + v.voice
//FROM
//threads AS t
//INNER JOIN votes AS v ON v.thread = t."tID"
//WHERE "user" = $1 AND t."tID" = $2
func (vote *Votes) VoteForThreadAndReturningVotes(pool *pgx.ConnPool, slugOrId string) (Threads, error) {
	var prevVote int32

	thread := Threads{}

	err := pool.QueryRow(`SELECT voice FROM votes WHERE "user" = $1 AND thread = $2`,
		vote.User, vote.Thread).Scan(&prevVote)

	tx, err := pool.Begin()
	if err != nil {
		return thread, err
	}
	defer tx.Rollback()

	if err != nil {
		tx.Exec(`INSERT INTO votes ("user", thread, voice) VALUES ($1, $2, $3)`,
			vote.User, vote.Thread, vote.Voice)
	} else {
		tx.Exec(`UPDATE votes SET voice = $1 WHERE "user" = $2 AND thread = $3`,
			vote.Voice, vote.User, vote.Thread)
		vote.Voice -= prevVote
	}

	log.Print(vote.Voice)
	log.Print(slugOrId)

	if id, parseErr := strconv.ParseInt(slugOrId, 10, 64); parseErr == nil {
		thread.TID = id
		tx.QueryRow(`UPDATE threads SET votes = votes + $1 WHERE "tID" = $2 RETURNING "tID", author, created, forum, message, title, votes, slug`,
			vote.Voice, thread.TID).Scan(&thread.TID, &thread.Author, &thread.Created, &thread.Forum,
			&thread.Message, &thread.Title, &thread.Votes, &thread.Slug)
	} else {
		thread.Slug = slugOrId
		tx.QueryRow(`UPDATE threads SET votes = votes + $1 WHERE slug = $2 RETURNING "tID", author, created, forum, message, title, votes, slug`,
			vote.Voice, thread.Slug).Scan(&thread.TID, &thread.Author, &thread.Created, &thread.Forum,
			&thread.Message, &thread.Title, &thread.Votes, &thread.Slug)
	}

	err = tx.Commit()
	if err != nil {
		return thread, err
	}

	return thread, nil
}
