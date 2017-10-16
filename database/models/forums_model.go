package models

import (
	"github.com/jackc/pgx"

	"../../utils"
	"strconv"
)

type Forums struct {
	FID     int64  `json:"fid"`
	Posts   int64  `json:"posts"`
	Slug    string `json:"slug"`
	Threads int32  `json:"threads"`
	Title   string `json:"title"`
	Author  string `json:"user"`
}

func (forum *Forums) CreateForum(pool *pgx.ConnPool) error {
	var id int64
	err := pool.QueryRow(`INSERT INTO forums(slug, title, author)`+
		`VALUES ($1, $2, $3) RETURNING "fID";`,
		forum.Slug, forum.Title, forum.Author).Scan(&id)
	if err != nil {
		if pgerr, ok := err.(pgx.PgError); ok {
			if pgerr.ConstraintName == "forums_slug_key" {
				return utils.UniqueError
			} else {
				return err
			}
		}
		return err
	}
	return nil
}

func (forum *Forums) GetForumBySlug(pool *pgx.ConnPool) (Forums, error) {
	resultForum := Forums{}
	err := pool.QueryRow(`SELECT slug, title, author, posts, threads  FROM forums WHERE slug = $1`,
		forum.Slug).Scan(&resultForum.Slug, &resultForum.Title, &resultForum.Author, &resultForum.Posts, &resultForum.Threads)

	if err != nil {
		return resultForum, err
	}
	return resultForum, nil
}

func (forum *Forums) GetAll(pool *pgx.ConnPool, limit, since, desc string) ([]Threads, error) {
	queryRow := `SELECT "tID", author, created, forum, message, slug, title, votes FROM threads WHERE forum = $1`

	var params []interface{}
	params = append(params, forum.Slug)
	if since != "" {
		if desc == "true" {
			queryRow += ` AND created <= $` + strconv.Itoa(len(params) + 1)
		} else {
			queryRow += ` AND created >= $` + strconv.Itoa(len(params) + 1)
		}
		params = append(params, since)
	}
	if desc == "true" {
		queryRow += ` ORDER BY created DESC`
	} else {
		queryRow += ` ORDER BY created ASC`
	}
	if limit != "" {
		queryRow += ` LIMIT $` + strconv.Itoa(len(params) + 1)
		params = append(params, limit)
	}



	rows, err := pool.Query(queryRow, params...)
	if err != nil {
		return nil, err
	}


	resultThreads := []Threads{}


	currentThreadInRows := Threads{}
	for rows.Next() {
		rows.Scan(&currentThreadInRows.TID, &currentThreadInRows.Author, &currentThreadInRows.Created, &currentThreadInRows.Forum,
			&currentThreadInRows.Message, &currentThreadInRows.Slug, &currentThreadInRows.Title, &currentThreadInRows.Votes)
		resultThreads = append(resultThreads, currentThreadInRows)
	}
	return resultThreads, nil
}
