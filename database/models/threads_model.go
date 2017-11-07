package models

import (
	"strconv"
	"time"

	"../../utils"
	"github.com/jackc/pgx"
)

type Threads struct {
	TID     int64     `json:"id"`
	Author  string    `json:"author"`
	Created time.Time `json:"created"`
	Forum   string    `json:"forum"`
	Message string    `json:"message"`
	Slug    string    `json:"slug"`
	Title   string    `json:"title"`
	Votes   int32     `json:"votes"`
}

func (thread *Threads) CreateThread(pool *pgx.ConnPool) error {
	var id int64

	tx, err := pool.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	//waitData := &sync.WaitGroup{}
	//waitData.Add(1)

	//go func(waitData *sync.WaitGroup) {
	//	defer waitData.Done()
	tx.Exec("UPDATE forums SET threads=threads+1 WHERE slug=$1", thread.Forum)
	//}(waitData)

	err = tx.QueryRow(`INSERT INTO threads (author, created, message, slug, title, forum)`+
		`VALUES ($1, $2, $3, $4, $5, $6) RETURNING "tID", created;`,
		thread.Author, thread.Created, thread.Message, thread.Slug, thread.Title, thread.Forum).Scan(&id, &thread.Created)
	if err != nil {
		if pgerr, ok := err.(pgx.PgError); ok {
			if pgerr.ConstraintName == "threads_slug_key" {
				return utils.UniqueError
			} else {
				return err
			}
		}
		return err
	}

	//waitData.Wait() //почему здесь нк работает?

	err = tx.Commit()
	if err != nil {
		return err
	}

	AddMember(pool, thread.Forum, thread.Author)

	thread.TID = id

	return nil
}

func (thread *Threads) GetThreadBySlug(pool *pgx.ConnPool) error {
	return pool.QueryRow(`SELECT "tID", author, created, forum, message, title, votes, slug FROM threads WHERE slug = $1`,
		thread.Slug).Scan(&thread.TID, &thread.Author, &thread.Created, &thread.Forum,
		&thread.Message, &thread.Title, &thread.Votes, &thread.Slug)

}

func (thread *Threads) GetThreadById(pool *pgx.ConnPool) error {
	return pool.QueryRow(`SELECT author, created, forum, message, slug, title, votes FROM threads WHERE "tID" = $1`,
		thread.TID).Scan(&thread.Author, &thread.Created, &thread.Forum,
		&thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
}

func (thread *Threads) GetPostsWithFlatSort(pool *pgx.ConnPool, limit, since, desc string) ([]Posts, error) {
	queryRow := `SELECT "pID", author, created, forum, message, thread, parent FROM posts WHERE thread = $1`

	var params []interface{}
	params = append(params, thread.TID)
	if since != "" {
		if desc == "true" {
			queryRow += ` AND "pID" < $` + strconv.Itoa(len(params)+1)
		} else {
			queryRow += ` AND "pID" > $` + strconv.Itoa(len(params)+1)
		}
		params = append(params, since)
	}
	if desc == "true" {
		queryRow += ` ORDER BY "pID" DESC`
	} else {
		queryRow += ` ORDER BY "pID" ASC`
	}
	if limit != "" {
		queryRow += ` LIMIT $` + strconv.Itoa(len(params)+1)
		params = append(params, limit)
	}

	rows, err := pool.Query(queryRow, params...)
	if err != nil {
		return nil, err
	}

	resultPosts := []Posts{}

	currentPostInRows := Posts{}
	for rows.Next() {
		rows.Scan(&currentPostInRows.PID, &currentPostInRows.Author, &currentPostInRows.Created, &currentPostInRows.Forum,
			&currentPostInRows.Message, &currentPostInRows.Thread, &currentPostInRows.Parent)
		resultPosts = append(resultPosts, currentPostInRows)
	}
	return resultPosts, nil
}

func (thread *Threads) GetPostsWithTreeSort(pool *pgx.ConnPool, limit, since, desc string) ([]Posts, error) {
	queryRow := `SELECT "pID", author, created, forum, message, thread, parent FROM posts WHERE thread = $1`

	var params []interface{}
	params = append(params, thread.TID)
	if since != "" {
		if desc == "true" {
			queryRow += ` AND path < (SELECT path FROM posts WHERE "pID" = $` + strconv.Itoa(len(params)+1) + `)`
		} else {
			queryRow += ` AND path > (SELECT path FROM posts WHERE "pID" = $` + strconv.Itoa(len(params)+1) + `)`
		}
		params = append(params, since)
	}
	if desc == "true" {
		queryRow += ` ORDER BY path DESC`
	} else {
		queryRow += ` ORDER BY path ASC`
	}
	if limit != "" {
		queryRow += ` LIMIT $` + strconv.Itoa(len(params)+1)
		params = append(params, limit)
	}

	rows, err := pool.Query(queryRow, params...)
	if err != nil {
		return nil, err
	}

	resultPosts := []Posts{}

	currentPostInRows := Posts{}
	for rows.Next() {
		rows.Scan(&currentPostInRows.PID, &currentPostInRows.Author, &currentPostInRows.Created, &currentPostInRows.Forum,
			&currentPostInRows.Message, &currentPostInRows.Thread, &currentPostInRows.Parent)
		resultPosts = append(resultPosts, currentPostInRows)
	}
	return resultPosts, nil
}

//SELECT "pID", author, created, forum, message, thread, parent FROM posts WHERE thread = 119 AND path[1] in (SELECT "pID" FROM posts
//WHERE thread = 119 AND parent = 0 AND path > (SELECT path FROM posts WHERE "pID" = 2038) ORDER BY path ASC  LIMIT 3) ORDER BY path ASC
func (thread *Threads) GetPostsWithParentTreeSort(pool *pgx.ConnPool, limit, since, desc string) ([]Posts, error) {
	queryRow := `SELECT "pID", author, created, forum, message, thread, parent FROM posts WHERE thread = $1 AND path[1] in (SELECT "pID" FROM posts
	WHERE thread = $1 AND parent = 0 `

	var params []interface{}
	params = append(params, thread.TID)

	//if since == "" { //почему так?
	//} else {
	//	if desc != "true" {
	//		queryRow += ` ORDER BY path DESC`
	//	} else {
	//		queryRow += ` ORDER BY path ASC`
	//	}
	//}

	if since != "" {
		if desc == "true" {
			queryRow += ` AND path < (SELECT path FROM posts WHERE "pID" = $` + strconv.Itoa(len(params)+1) + `)`
		} else {
			queryRow += ` AND path > (SELECT path FROM posts WHERE "pID" = $` + strconv.Itoa(len(params)+1) + `)`
		}
		params = append(params, since)
	}
	if desc == "true" {
		queryRow += ` ORDER BY path DESC`
	} else {
		queryRow += ` ORDER BY path ASC`
	}
	if limit != "" {
		queryRow += ` LIMIT $` + strconv.Itoa(len(params)+1)
		params = append(params, limit)
	}
	queryRow += `)`

	if desc == "true" {
		queryRow += ` ORDER BY path DESC`
	} else {
		queryRow += ` ORDER BY path ASC`
	}
	//if limit != "" {
	//	queryRow += ` LIMIT $` + strconv.Itoa(len(params)+1)
	//	params = append(params, limit)
	//}
	rows, err := pool.Query(queryRow, params...)
	if err != nil {

		return nil, err
	}

	resultPosts := []Posts{}

	currentPostInRows := Posts{}
	for rows.Next() {
		rows.Scan(&currentPostInRows.PID, &currentPostInRows.Author, &currentPostInRows.Created, &currentPostInRows.Forum,
			&currentPostInRows.Message, &currentPostInRows.Thread, &currentPostInRows.Parent)
		resultPosts = append(resultPosts, currentPostInRows)
	}
	return resultPosts, nil
}

func (thread *Threads) UpdateThread(pool *pgx.ConnPool) error {
	err := pool.QueryRow(`UPDATE threads SET author = $1, message = $2, title = $3, forum = $4 `+
		`WHERE slug = $5 RETURNING "tID", slug, created;`,
		thread.Author, thread.Message, thread.Title, thread.Forum, thread.Slug).Scan(&thread.TID, &thread.Slug, &thread.Created)
	if err != nil {
		return err
	}
	return nil
}

func ThreadsCount(pool *pgx.ConnPool) (int32, error) {
	var count int32
	err := pool.QueryRow("SELECT COUNT(*) FROM threads").Scan(&count)
	return count, err
}
