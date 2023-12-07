package dao

import (
	"dapdap-job/common/dal"
	"database/sql"
	"errors"
)

func (d *Dao) FindActionCategory() (data map[string]int, err error) {
	data = map[string]int{}
	rows, err := d.db.Query(dal.FindActionCategorySql)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var (
			id   int
			name string
		)
		if err = rows.Scan(&id, &name); err != nil {
			return
		}
		data[name] = id
	}
	return
}

func (d *Dao) FindNetworks() (data map[int]int, err error) {
	data = map[int]int{}
	rows, err := d.db.Query(dal.FindNetworksSql)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var (
			id      int
			chainId int
		)
		if err = rows.Scan(&id, &chainId); err != nil {
			return
		}
		data[chainId] = id
	}
	return
}
