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

	//pool.Query(`INSERT INTO votes ("user", thread, voice)`+
	//	`VALUES ($1, $2, $3);`,
	//	vote.User, vote.Thread, vote.Voice)

	var votesNumber int32
	//thread.setVotes(template.queryForObject("SELECT t.votes FROM threads t " +
	//	"WHERE t.id = ?", Integer.class, thread.getId()));
	pool.QueryRow(`SELECT SUM(voice) AS voiceSum FROM votes WHERE "thread" = $1`, vote.Thread).Scan(&votesNumber)
	return votesNumber, nil
}
