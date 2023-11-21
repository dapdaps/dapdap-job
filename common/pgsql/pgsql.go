package pgsql

import (
	"dapdap-job/common/conf"
	"dapdap-job/common/log"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

func NewPgsql(c *conf.Pgsql) (db *sql.DB) {
	db, err := open(c)
	if err != nil {
		panic(err)
	}
	return
}

func open(c *conf.Pgsql) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", c.DSN)
	if err != nil {
		log.Error("error opening pgsql %+v: %v", c, err)
		return nil, err
	}
	db.SetMaxOpenConns(c.Active)
	db.SetMaxIdleConns(c.Idle)
	db.SetConnMaxLifetime(time.Duration(c.IdleTimeout))
	err = db.Ping()
	return
}

func WithTrx(db *sql.DB, f func(tx *sql.Tx) (err error)) (err error) {
	tx, err := db.Begin()
	if err != nil {
		log.Error("error getting transaction: %v", err)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("error exec trx : %v", r)
			_ = tx.Rollback()
			err = fmt.Errorf("%v", r)
		}
	}()

	err = f(tx)
	if err != nil {
		log.Error("error trx: %v, roll back", err)
		if e := tx.Rollback(); e != nil {
			log.Error("error rolling back: %v", e)
		}
		return
	}
	if err := tx.Commit(); err != nil {
		log.Error("error committing: %v", err)
	}
	return
}
