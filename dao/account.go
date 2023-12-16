package dao

import (
	"dapdap-job/common/dal"
	"dapdap-job/model"
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
		userIdSql sql.NullInt64
	)
	err = d.db.QueryRow(dal.FindAccountIdByTgSql, tgUserId).Scan(&userIdSql)
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

func (d *Dao) FindAccountIdByDiscord(discordUserId string) (accountId int, err error) {
	var (
		userIdSql sql.NullInt64
	)
	err = d.db.QueryRow(dal.FindAccountIdByDiscordSql, discordUserId).Scan(&userIdSql)
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

func (d *Dao) FindAllAccountExt() (data map[int]*model.AccountExt, err error) {
	data = map[int]*model.AccountExt{}
	rows, err := d.db.Query(dal.FindAccountExtSql)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var (
			twitter_user_id              sql.NullString
			twitter_access_token_type    sql.NullString
			twitter_access_token_expires sql.NullString
			twitter_access_token         sql.NullString
			twitter_refresh_token        sql.NullString
			telegram_user_id             sql.NullString
			discord_user_id              sql.NullString
			accountExt                   = &model.AccountExt{}
		)
		if err = rows.Scan(&accountExt.AccountId, &twitter_user_id, &twitter_access_token_type, &twitter_access_token_expires, &twitter_access_token, &twitter_refresh_token, &telegram_user_id, &discord_user_id); err != nil {
			return
		}
		if twitter_user_id.Valid {
			accountExt.TwitterUserId = twitter_user_id.String
		}
		if twitter_access_token_type.Valid {
			accountExt.TwitterAccessTokenType = twitter_access_token_type.String
		}
		if twitter_access_token.Valid {
			accountExt.TwitterAccessToken = twitter_access_token.String
		}
		if twitter_refresh_token.Valid {
			accountExt.TwitterRefreshToken = twitter_refresh_token.String
		}
		if telegram_user_id.Valid {
			tgUserId, e := strconv.Atoi(telegram_user_id.String)
			if e == nil {
				accountExt.TelegramUserId = int64(tgUserId)
			}
		}
		if discord_user_id.Valid {
			accountExt.DiscordUserId = discord_user_id.String
		}
		data[accountExt.AccountId] = accountExt
	}
	return
}
