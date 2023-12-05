package dao

import (
	"dapdap-job/common/dal"
	"database/sql"
	"errors"
)

func (d *Dao) FindAccountId(address []string) (data map[string]int, err error) {
	var (
		findSql = dal.FindAccountIdByAddressSql
		args    []interface{}
	)
	data = map[string]int{}
	for _, addr := range address {
		findSql += `?,`
		args = append(args, addr)
	}
	findSql = findSql[0:len(findSql)-1] + `)`
	rows, err := d.db.Query(findSql, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var (
			accountId int
			addr      string
		)
		if err = rows.Scan(&accountId, &addr); err != nil {
			return
		}
		data[addr] = accountId
	}
	return
}
