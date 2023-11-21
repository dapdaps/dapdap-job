package dao

import (
	m "dapdap-job/common/pgsql"
	"dapdap-job/conf"
	"database/sql"
)

type Dao struct {
	db *sql.DB
}

func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db: m.NewPgsql(c.Pgsql.DB),
	}
	return
}

func (d *Dao) WithTrx(f func(db *sql.Tx) (err error)) (err error) {
	return m.WithTrx(d.db, f)
}
