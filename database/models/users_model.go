package models

import (
	"github.com/jackc/pgx"

	"../../utils"

)

type Users struct {
	//UID      int64
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Fullname string `json:"fullname"`
	About    string `json:"about"`
}

func (user *Users) CreateUser(pool *pgx.ConnPool) error {
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


func (user *Users) GetUserByLogin(pool *pgx.ConnPool) (Users, error) {
	resultUser := Users{}
	err := pool.QueryRow(`SELECT nickname, email, fullname, about FROM users WHERE nickname = $1`,
		user.Nickname).Scan(&resultUser.Nickname, &resultUser.Email, &resultUser.Fullname, &resultUser.About)

	if err != nil {
		return resultUser, err
	}
	return resultUser, nil
}

func (user *Users) GetUserByEmail(pool *pgx.ConnPool) (Users, error) {
	resultUser := Users{}
	err := pool.QueryRow(`SELECT nickname, email, fullname, about FROM users WHERE email = $1`,
		user.Email).Scan(&resultUser.Nickname, &resultUser.Email, &resultUser.Fullname, &resultUser.About)

	if err != nil {
		return resultUser, err
	}
	return resultUser, nil
}

func (user *Users) GetUserByLoginAndEmail(pool *pgx.ConnPool) ([]Users, error) {
	rows, err := pool.Query(`SELECT nickname, email, fullname, about FROM users WHERE nickname = $1 OR email = $2`,
		user.Nickname, user.Email)

	resultUsers := []Users{}
	if err != nil {
		return resultUsers, err
	}

	currentUserInRows := Users{}
	for rows.Next() {
		rows.Scan(&currentUserInRows.Nickname, &currentUserInRows.Email, &currentUserInRows.Fullname, &currentUserInRows.About)
		resultUsers = append(resultUsers,currentUserInRows)
	}
	return resultUsers, nil
}


func (user *Users) UpdateUser(pool *pgx.ConnPool) error {
	var id int64
	err := pool.QueryRow(`UPDATE users SET email = $1, fullname = $2, about = $3`+
		`WHERE nickname = $4 RETURNING "uID";`,
		user.Email, user.Fullname, user.About, user.Nickname).Scan(&id)
	if err != nil {
		if pgerr, ok := err.(pgx.PgError); ok {
			if pgerr.ConstraintName == "users_email_key"  {
				return utils.UniqueError
			} else {
				return err
			}
		}
		return err
	}
	return nil
}