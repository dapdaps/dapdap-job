package dao

import (
	"dapdap-job/common/dal"
	"database/sql"
	"errors"
	"strconv"
)

func (d *Dao) FindAccountId(address []string) (data map[string]int, dataArr []int, err error) {
	var (
		findSql = dal.FindAccountIdByAddressSql
		args    []interface{}
	)
	data = map[string]int{}
	if len(address) == 0 {
		return
	}
	for index, addr := range address {
		findSql += `$` + strconv.Itoa(index+1) + `,`
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
		dataArr = append(dataArr, accountId)
	}
	return
}
