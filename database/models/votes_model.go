package models

import (
	"github.com/jackc/pgx"
)

type Votes struct {
	Voice  int32 `json:"voice"`
	User   string `json:"nickname"`
	Thread int64 `json:"thread"`
}

func (vote *Votes) VoteForThreadAndReturningVotes(pool *pgx.ConnPool) (int32, error) {
	pool.QueryRow(`INSERT INTO votes ("user", thread, voice) VALUES ($1, $2, $3)
	 ON CONFLICT ("user", thread) DO UPDATE SET voice = $4`,
		vote.User, vote.Thread, vote.Voice, vote.Voice)


	var votesNumber int32
	//_, err = tx.Exec(`UPDATE forums SET posts = posts + $1 WHERE slug = $2`, len(posts), forum)
	//_, err = tx.Exec(`UPDATE thread SET votes = (SELECT SUM(voice) AS voiceSum FROM votes WHERE "thread" = $1)`, len(posts), forum)
	//`UPDATE threads SET votes = (SELECT SUM(voice) AS voiceSum FROM votes WHERE thread = 408) RETURNING votes`

	pool.QueryRow(`UPDATE threads SET votes = (SELECT SUM(voice) AS voiceSum FROM votes WHERE thread = $1) RETURNING votes`, vote.Thread).Scan(&votesNumber)
	return votesNumber, nil
}
