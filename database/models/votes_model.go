package models

import (
	"github.com/jackc/pgx"
)

type Votes struct {
	Voice  int32  `json:"voice"`
	User   string `json:"nickname"`
	Thread int64  `json:"thread"`
}

func (vote *Votes) VoteForThreadAndReturningVotes(pool *pgx.ConnPool, votesThread int32) (int32, error) {

	var prevVote int32

	err := pool.QueryRow(`SELECT voice FROM votes WHERE "user" = $1 AND thread = $2`,
		vote.User, vote.Thread).Scan(&prevVote)

	if err != nil {
		pool.Exec(`INSERT INTO votes ("user", thread, voice) VALUES ($1, $2, $3)`,
			vote.User, vote.Thread, vote.Voice)
		votesThread += vote.Voice
	} else {
		pool.Exec(`UPDATE votes SET voice = $1 WHERE "user" = $2 AND thread = $3`,
			vote.Voice, vote.User, vote.Thread)
		votesThread += vote.Voice - prevVote
	}

	pool.Exec(`UPDATE threads SET votes = $1 WHERE "tID" = $2`, votesThread, vote.Thread)

	return votesThread, nil
}
