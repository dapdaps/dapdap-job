package dao

import (
	"dapdap-job/common/dal"
	"dapdap-job/model"
	"database/sql"
	"errors"
	"strconv"
	"time"
)

func (d *Dao) FindInvites(inviteAddress map[string]int) (data map[int64][]*model.Invite, err error) {
	var (
		index   int
		findSql = dal.FindInvitesSql + ` where used_user_id in(`
		args    []interface{}
	)
	data = map[int64][]*model.Invite{}
	if len(inviteAddress) == 0 {
		return
	}
	for _, accountId := range inviteAddress {
		index++
		findSql += `$` + strconv.Itoa(index) + `,`
		args = append(args, accountId)
	}
	findSql = findSql[0:len(findSql)-1] + `)`
	index++
	findSql += ` and is_used=$` + strconv.Itoa(index)
	args = append(args, true)

	index++
	findSql += ` and status!=$` + strconv.Itoa(index)
	args = append(args, "Active")

	findSql += ` order by id asc`

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
			invite        = &model.Invite{}
			creatorUserId sql.NullInt64
			usedUserId    sql.NullInt64
		)
		if err = rows.Scan(&invite.Id, &creatorUserId, &usedUserId); err != nil {
			return
		}
		if !creatorUserId.Valid || !usedUserId.Valid {
			continue
		}
		invite.CreatorUserId = creatorUserId.Int64
		invite.UsedUserId = usedUserId.Int64
		if data[invite.CreatorUserId] == nil {
			data[invite.CreatorUserId] = []*model.Invite{}
		}
		data[invite.CreatorUserId] = append(data[invite.CreatorUserId], invite)
	}
	return
}

func (d *Dao) FindTotalInviteReward(accountId int64) (reward int64, err error) {
	var rewardSql sql.NullInt64
	err = d.db.QueryRow(dal.FindTotalInviteRewardSql, accountId).Scan(&rewardSql)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	if rewardSql.Valid {
		reward = rewardSql.Int64
	}
	return
}

func (d *Dao) UpdateInviteReward(invites []*model.Invite) (err error) {
	timestamp := time.Now()
	err = d.WithTrx(func(db *sql.Tx) (err error) {
		for _, invite := range invites {
			_, err = db.Exec(dal.UpdateInviteSql, "Active", invite.Reward, timestamp, invite.Id)
			if err != nil {
				return
			}
		}
		return
	})
	return
}
