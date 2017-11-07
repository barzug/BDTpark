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
	//pool.QueryRow(`INSERT INTO votes ("user", thread, voice) VALUES ($1, $2, $3)
	// ON CONFLICT ("user", thread) DO UPDATE SET voice = $3 RETURNING voice`,
	//	vote.User, vote.Thread, vote.Voice).Scan(&vote.Voice)
	//
	var prevVote int32

	//
	//pool.QueryRow(`UPDATE threads SET votes = (SELECT SUM(voice) AS voiceSum FROM votes WHERE thread = $1) WHERE "tID" = $1 RETURNING votes`, vote.Thread).Scan(&votesNumber)

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
