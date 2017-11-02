package models

import (
	"github.com/jackc/pgx"
	"log"
)

type Votes struct {
	Voice  int32  `json:"voice"`
	User   string `json:"nickname"`
	Thread int64  `json:"thread"`
}

func (vote *Votes) VoteForThreadAndReturningVotes(pool *pgx.ConnPool) (int32, error) {
	err := pool.QueryRow(`INSERT INTO votes ("user", thread, voice) VALUES ($1, $2, $3)
	 ON CONFLICT ("user", thread) DO UPDATE SET voice = $3 RETURNING voice`,
		vote.User, vote.Thread, vote.Voice).Scan(&vote.Voice)
	if err != nil {
		log.Print(err)
	}

	var votesNumber int32

	err = pool.QueryRow(`UPDATE threads SET votes = (SELECT SUM(voice) AS voiceSum FROM votes WHERE thread = $1) WHERE "tID" = $1 RETURNING votes`, vote.Thread).Scan(&votesNumber)
	if err != nil {
		log.Print(err)
	}
	return votesNumber, nil
}
