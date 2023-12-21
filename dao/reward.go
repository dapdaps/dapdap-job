package dao

import (
	"dapdap-job/common/dal"
	"dapdap-job/model"
	"database/sql"
	"errors"
)

func (d *Dao) FindUserRewardMaxId() (id int, err error) {
	var (
		maxRecordIdSql sql.NullInt64
	)
	err = d.db.QueryRow(dal.FindUserRewardMaxIdSql).Scan(&maxRecordIdSql)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	if maxRecordIdSql.Valid {
		id = int(maxRecordIdSql.Int64)
	}
	return
}

func (d *Dao) FindAllUserReward() (data []*model.UserReward, totalReward int, err error) {
	var (
		id    = 0
		maxId int
		limit = 5000
		rows  *sql.Rows
	)
	maxId, err = d.FindUserRewardMaxId()
	if err != nil {
		return
	}
	if maxId == 0 {
		return
	}
	for {
		rows, err = d.db.Query(dal.FindUserRewardByBetweenSql, id, id+limit)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return
		}
		err = nil
		for rows.Next() {
			var userReward = &model.UserReward{}
			if err = rows.Scan(&userReward.Id, &userReward.AccountId, &userReward.ClaimedReward); err != nil {
				return
			}
			totalReward += userReward.ClaimedReward
			data = append(data, userReward)
		}
		_ = rows.Close()

		if len(data) > 0 && data[len(data)-1].Id >= maxId {
			return
		}
		id = id + limit + 1
	}
}
