package dao

import (
	"dapdap-job/common/dal"
	"database/sql"
	"errors"
)

func (d *Dao) FindUserReward(accountId int) (reward int, inviteReward int, err error) {
	var (
		rewardSql       sql.NullInt64
		rewardInviteSql sql.NullInt64
	)
	err = d.db.QueryRow(dal.FindUserRewardByIdSql, accountId).Scan(&rewardSql, &rewardInviteSql)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	if rewardSql.Valid {
		reward = int(rewardSql.Int64)
	}
	if rewardInviteSql.Valid {
		inviteReward = int(rewardInviteSql.Int64)
	}
	return
}
