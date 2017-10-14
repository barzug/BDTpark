package services


import (
	"../models"
	"github.com/jackc/pgx"
	"encoding/json"
)

func GetForumBySlug(pool *pgx.ConnPool, slug string) ([]byte, error) {
	forum := new(models.Forums)
	err := pool.QueryRow("SELECT * FROM forums WHERE slug = $1", slug).
		Scan(&forum.Posts, &forum.Slug, &forum.Threads, &forum.Title, &forum.User)
	byteData, _ := json.Marshal(forum)
	return byteData, err
}

func CreateForum(pool *pgx.ConnPool, body []byte) error {
	forum := new(models.Forums)
	//проверка на то что пользователь не найден 404 или присутствует форум 409
	if err := json.Unmarshal(body, forum); err != nil {
		return err
	}

	var id int64 //надо ли?
	if err := pool.QueryRow("INSERT INTO forums (slug, title, user ) VALUES ($1, $2, $3);",
		forum.Slug, forum.Title, forum.User).Scan(&id); err != nil {
		return err
	}
	return nil
}

