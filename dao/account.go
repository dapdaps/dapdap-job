package dao

import (
	"dapdap-job/common/dal"
	"database/sql"
	"errors"
	"strconv"
)

func (d *Dao) FindAllAccountId() (data map[string]int, err error) {
	data = map[string]int{}
	rows, err := d.db.Query(dal.FindAccountIdSql)
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

func (d *Dao) FindAccountIds(addresses map[string]string) (data map[string]int, dataArr []int, err error) {
	var params []string
	for addr := range addresses {
		params = append(params, addr)
	}
	return d.FindAccountIdByAddress(params)
}

func (d *Dao) FindAccountIdByAddress(address []string) (data map[string]int, dataArr []int, err error) {
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

func (d *Dao) FindAccountIdByTg(tgUserId int64) (accountId int, err error) {
	var (
		userIdSql  sql.NullInt64
		addressSql sql.NullString
	)
	err = d.db.QueryRow(dal.FindAccountIdByTgSql, tgUserId).Scan(&userIdSql, &addressSql)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	if userIdSql.Valid {
		accountId = int(userIdSql.Int64)
	}
	return
}

func (d *Dao) SelectForUpdate(db *sql.Tx, accountId int) (err error) {
	var userId sql.NullInt64
	err = db.QueryRow(dal.FindAccountForUpdateSql, accountId).Scan(&userId)
	if err != nil {
		return
	}
	return
}
