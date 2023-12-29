package dao

import (
	"dapdap-job/common/dal"
	"dapdap-job/model"
	"database/sql"
	"errors"
	"strconv"
	"time"
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

func (d *Dao) FindAccountIdByTg(tgUserId string) (accountId int, err error) {
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

func (d *Dao) FindAllAccountExt(minUpdateTime *time.Time) (data map[int]*model.AccountExt, maxUpdateTime *time.Time, err error) {
	var (
		rows *sql.Rows
	)
	data = map[int]*model.AccountExt{}
	if minUpdateTime == nil {
		rows, err = d.db.Query(dal.FindAccountExtSql)
	} else {
		rows, err = d.db.Query(dal.FindAccountExtSql+` where updated_at>$1`, *minUpdateTime)
	}
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
			telegram_user_id             sql.NullString
			discord_user_id              sql.NullString
			updateTime                   sql.NullTime
			accountExt                   = &model.AccountExt{}
		)
		if err = rows.Scan(&accountExt.AccountId, &twitter_user_id, &twitter_access_token_type, &twitter_access_token_expires, &twitter_access_token, &telegram_user_id, &discord_user_id, &updateTime); err != nil {
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
		if telegram_user_id.Valid {
			accountExt.TelegramUserId = telegram_user_id.String
		}
		if discord_user_id.Valid {
			accountExt.DiscordUserId = discord_user_id.String
		}
		if updateTime.Valid {
			if maxUpdateTime == nil || updateTime.Time.After(*maxUpdateTime) {
				maxUpdateTime = &updateTime.Time
			}
		}
		data[accountExt.AccountId] = accountExt
	}
	return
}
