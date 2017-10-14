package database

import (
	"github.com/jackc/pgx"

	"io/ioutil"
)

const schema = "./init.sql"

type DbFacade struct {
	Pool *pgx.ConnPool
}

func (db *DbFacade) InitDB(connConfig pgx.ConnConfig, poolSize int, function func(conn *pgx.Conn) error) error {
	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig:     connConfig,
		MaxConnections: poolSize,
		AfterConnect:   function,
	}
	var err error
	db.Pool, err = pgx.NewConnPool(connPoolConfig)
	return err
}

func (db *DbFacade) InitSchema() error {
	buf, err := ioutil.ReadFile(schema)

	if err != nil {
		return err
	}

	schema := string(buf)

	_, err = db.Pool.Exec(schema)

	if err != nil {
		return err
	}
	return nil
}
