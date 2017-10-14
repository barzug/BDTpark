package models

import (
	"github.com/jackc/pgx"

	"../../utils"

)

type Users struct {
	UID      int64  `json:"uid"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Fullname string `json:"fullname"`
	About    string `json:"about"`
}

func (user *Users) CreateUserQuery(pool *pgx.ConnPool) error {
	var id int64
	err := pool.QueryRow(`INSERT INTO users(nickname, email, fullname, about)`+
		 `VALUES ($1, $2, $3, $4) RETURNING "uID";`,
		user.Nickname, user.Email, user.Fullname, user.About).Scan(&id)
	if err != nil {
		if pgerr, ok := err.(pgx.PgError); ok {
			if pgerr.ConstraintName == "users_nickname_key" || pgerr.ConstraintName == "users_email_key"  {
				return utils.UniqueError
			} else {
				return err
			}
		}
		return err
	}
	return nil
}
