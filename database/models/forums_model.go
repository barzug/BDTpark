package models

import (
	"github.com/jackc/pgx"

	"../../utils"
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
