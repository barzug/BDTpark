package daemon

import (
	"../database"
	"github.com/jackc/pgx"
	"github.com/fasthttp-contrib/render"
)

func afterConnect(conn *pgx.Conn) error {
	return nil
}

var (
	DB database.DbFacade
	Render = render.New()
)

func Init(host, database, user, password string, connections int) error {
	err := DB.InitDB(pgx.ConnConfig{
		Host:     host,
		User:     user,
		Password: password,
		Database: database,
	}, connections, afterConnect)
	if err != nil {
		return err
	}

	err = DB.InitSchema()
	if err != nil {
		return err
	}

	return nil
}

