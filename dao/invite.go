package dao

import (
	"dapdap-job/common/dal"
	"dapdap-job/conf"
	"dapdap-job/model"
	"database/sql"
	"errors"
	"strconv"
	"time"
)

func (d *Dao) FindInvites(inviteAddress map[string]int) (data map[string]*model.Invite, err error) {
	var (
		index   int
		findSql = dal.FindInvitesSql + ` where i.used_user_id in(`
		args    []interface{}
	)
	data = map[string]*model.Invite{}
	for _, accountId := range inviteAddress {
		index++
		findSql += `$` + strconv.Itoa(index) + `,`
		args = append(args, accountId)
	}
	findSql = findSql[0:len(findSql)-1] + `)`
	index++
	findSql += ` and i.is_used=$` + strconv.Itoa(index)
	args = append(args, true)

	index++
	findSql += ` and i.status!=$` + strconv.Itoa(index)
	args = append(args, "Active")

	index++
	findSql += ` and u.invite_reward<$` + strconv.Itoa(index)
	args = append(args, conf.Conf.MaxInviteReward)

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
			invitedReward sql.NullInt64
		)
		if err = rows.Scan(&invite.Id, &creatorUserId, &usedUserId, &invitedReward); err != nil {
			return
		}
		if !creatorUserId.Valid || !usedUserId.Valid {
			continue
		}
		invite.CreatorUserId = creatorUserId.Int64
		invite.CreatorUserId = creatorUserId.Int64
		if invitedReward.Valid {
			invite.InvitedReward = invitedReward.Int64
		}
		for address, accountId := range inviteAddress {
			if accountId == int(usedUserId.Int64) {
				data[address] = invite
			}
		}
	}
	return
}

func (d *Dao) UpdateInvite(invite *model.Invite) (err error) {
	timestamp := time.Now()
	err = d.WithTrx(func(db *sql.Tx) (err error) {
		err = d.SelectForUpdate(db, int(invite.CreatorUserId))
		if err != nil {
			return
		}
		_, err = db.Exec(dal.UpdateInviteSql, "Active", invite.Id)
		if err != nil {
			return
		}
		reward, inviteReward, err := d.FindUserReward(int(invite.CreatorUserId))
		if err != nil {
			return
		}
		_, err = db.Exec(dal.UpdateUserInviteRewardByIdSql, invite.CreatorUserId, reward+int(conf.Conf.InviteReward), inviteReward+int(conf.Conf.InviteReward), timestamp)
		if err != nil {
			return
		}
		return
	})
	return
}
