package dao

import (
	"dapdap-job/common/dal"
	"dapdap-job/model"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"
)

func (d *Dao) InitQuestCampaignInfo() (err error) {
	var count sql.NullInt64
	err = d.db.QueryRow(dal.FindCountQuestCampaignInfoSql).Scan(&count)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return
	}
	err = nil
	if count.Valid && count.Int64 > 0 {
		return
	}
	_, err = d.db.Exec(dal.InitQuestCampaignInfoSql)
	return
}

func (d *Dao) FindQuestActionMaxRecordId() (maxRecordIdActionQuest uint64, err error) {
	var (
		maxRecordIdActionQuestSql sql.NullInt64
	)
	err = d.db.QueryRow(dal.FindQuestActionRecordIdSql).Scan(&maxRecordIdActionQuestSql)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		} else {
			return
		}
	}
	if maxRecordIdActionQuestSql.Valid {
		maxRecordIdActionQuest = uint64(maxRecordIdActionQuestSql.Int64)
	}
	return
}

func (d *Dao) FindQuestCampaignInfo() (data *model.QuestCampaignInfo, err error) {
	var (
		totalUsers          sql.NullInt64
		totalReward         sql.NullInt64
		totalQuestExecution sql.NullInt64
	)
	data = &model.QuestCampaignInfo{}
	err = d.db.QueryRow(dal.FindQuestCampaignInfoSql).Scan(&totalUsers, &totalReward, &totalQuestExecution)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	if totalUsers.Valid {
		data.TotalUsers = int(totalUsers.Int64)
	}
	if totalReward.Valid {
		data.TotalReward = int(totalReward.Int64)
	}
	if totalQuestExecution.Valid {
		data.TotalQuestExecution = int(totalQuestExecution.Int64)
	}
	return
}

func (d *Dao) FindAllQuestCampaign() (data []*model.QuestCampaign, err error) {
	rows, err := d.db.Query(dal.FindQuestCampaignByStatusSql, model.QuestCampaignOnGoingStatus)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var (
			questCampaign = &model.QuestCampaign{}
		)
		if err = rows.Scan(&questCampaign.Id); err != nil {
			return
		}
		data = append(data, questCampaign)
	}
	return
}

func (d *Dao) FindAllQuest(questCampaignId int) (data map[int]*model.Quest, err error) {
	data = map[int]*model.Quest{}
	rows, err := d.db.Query(dal.FindAllQuestByStatusIdSql, questCampaignId, model.QuestOnGoingStatus)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var (
			quest  = &model.Quest{}
			status sql.NullString
		)
		if err = rows.Scan(&quest.Id, &quest.QuestCampaignId, &quest.QuestCategoryId, &quest.TotalAction, &status, &quest.Reward); err != nil {
			return
		}
		if status.Valid {
			quest.Status = status.String
		}
		data[quest.Id] = quest
	}
	return
}

func (d *Dao) FindAllQuestAction(questCampaignId int) (data map[int]*model.QuestAction, err error) {
	data = map[int]*model.QuestAction{}
	rows, err := d.db.Query(dal.FindAllQuestActionByCategoryIdSql, questCampaignId, "dapp")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var (
			questAction = &model.QuestAction{}
			source      sql.NullString
			dapps       sql.NullString
			networks    sql.NullString
			toNetworks  sql.NullString
		)
		if err = rows.Scan(&questAction.Id, &questAction.QuestCampaignId, &questAction.QuestId, &questAction.Times, &questAction.CategoryId, &source, &dapps, &networks, &toNetworks); err != nil {
			return
		}
		questAction.DappsMap = map[int]int{}
		if dapps.Valid {
			questAction.Dapps = dapps.String
			dappsArr := strings.Split(questAction.Dapps, ",")
			for _, dappId := range dappsArr {
				if dappIdInt, e := strconv.Atoi(dappId); e == nil {
					questAction.DappsMap[dappIdInt] = dappIdInt
				}
			}
		}

		questAction.NetworksMap = map[int]int{}
		if networks.Valid {
			questAction.Networks = networks.String
			networksArr := strings.Split(questAction.Networks, ",")
			for _, networkId := range networksArr {
				if networkIdInt, e := strconv.Atoi(networkId); e == nil {
					questAction.NetworksMap[networkIdInt] = networkIdInt
				}
			}
		}

		questAction.ToNetworksMap = map[int]int{}
		if toNetworks.Valid {
			questAction.ToNetworks = toNetworks.String
			toNetworksArr := strings.Split(questAction.ToNetworks, ",")
			for _, networkId := range toNetworksArr {
				if networkIdInt, e := strconv.Atoi(networkId); e == nil {
					questAction.ToNetworksMap[networkIdInt] = networkIdInt
				}
			}
		}
		data[questAction.Id] = questAction
	}
	return
}

func (d *Dao) FindUserQuest(questCampaignId int, accountIds []int) (data map[int][]*model.UserQuest, err error) {
	var (
		findSql = dal.FindUserQuestSql
		args    []interface{}
	)
	data = map[int][]*model.UserQuest{}
	if len(accountIds) == 0 {
		return
	}
	findSql += ` where `
	findSql += ` account_id in(`
	for index, accountId := range accountIds {
		findSql += `$` + strconv.Itoa(index+1) + `,`
		args = append(args, accountId)
	}
	args = append(args, questCampaignId)
	findSql = findSql[0:len(findSql)-1] + `) and quest_campaign_id=$` + strconv.Itoa(len(args))
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
			userQuest = &model.UserQuest{}
		)
		if err = rows.Scan(&userQuest.Id, &userQuest.QuestCampaignId, &userQuest.QuestId, &userQuest.AccountId, &userQuest.ActionCompleted, &userQuest.Status); err != nil {
			return
		}
		if _, ok := data[userQuest.AccountId]; !ok {
			data[userQuest.AccountId] = []*model.UserQuest{}
		}
		data[userQuest.AccountId] = append(data[userQuest.AccountId], userQuest)
	}
	return
}

func (d *Dao) FindUserQuestAction(questCampaignId int, accountIds []int) (data map[int][]*model.UserQuestAction, err error) {
	var (
		findSql = dal.FindUserQuestActionSql
		args    []interface{}
	)
	data = map[int][]*model.UserQuestAction{}
	if len(accountIds) == 0 {
		return
	}
	findSql += ` where `
	findSql += ` account_id in(`
	for index, accountId := range accountIds {
		findSql += `$` + strconv.Itoa(index+1) + `,`
		args = append(args, accountId)
	}
	args = append(args, questCampaignId)
	findSql = findSql[0:len(findSql)-1] + `) and quest_campaign_id=$` + strconv.Itoa(len(args))
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
			userQuestAction = &model.UserQuestAction{}
		)
		if err = rows.Scan(&userQuestAction.Id, &userQuestAction.QuestCampaignId, &userQuestAction.QuestId, &userQuestAction.QuestActionId, &userQuestAction.AccountId, &userQuestAction.Times, &userQuestAction.Status); err != nil {
			return
		}
		if _, ok := data[userQuestAction.AccountId]; !ok {
			data[userQuestAction.AccountId] = []*model.UserQuestAction{}
		}
		data[userQuestAction.AccountId] = append(data[userQuestAction.AccountId], userQuestAction)
	}
	return
}

func (d *Dao) UpdateActionRecord(id uint64) (err error) {
	lastId, err := d.FindQuestActionMaxRecordId()
	if err != nil {
		return
	}
	if lastId == 0 {
		_, err = d.db.Exec(dal.SaveQuestActionRecordIdSql, id)
	} else {
		_, err = d.db.Exec(dal.UpdateQuestActionRecordIdSql, id)
	}
	return
}

func (d *Dao) UpdateUserQuest(accountId int, questCampaignId int, reward int, campaignInfo *model.QuestCampaignInfo, userQuests []*model.UserQuest, userQuestActions []*model.UserQuestAction) (err error) {
	timestamp := time.Now()
	err = d.WithTrx(func(db *sql.Tx) (err error) {
		for _, userQuest := range userQuests {
			_, err = db.Exec(dal.UpdateUserQuestSql, userQuest.AccountId, userQuest.QuestId, userQuest.QuestCampaignId, userQuest.ActionCompleted, userQuest.Status, timestamp)
			if err != nil {
				return
			}
		}
		for _, userQuestAction := range userQuestActions {
			_, err = db.Exec(dal.UpdateUserQuestActionSql, userQuestAction.AccountId, userQuestAction.QuestActionId, userQuestAction.QuestId, userQuestAction.QuestCampaignId, userQuestAction.Times, userQuestAction.Status, timestamp)
			if err != nil {
				return
			}
		}
		if campaignInfo != nil {
			_, err = db.Exec(dal.UpdateQuestCampaignInfoSql, campaignInfo.TotalUsers, campaignInfo.TotalQuestExecution, timestamp)
			if err != nil {
				return
			}
		}
		//if questCampaign != nil {
		//	_, err = db.Exec(dal.UpdateQuestCampaignSql, questCampaign.TotalUsers, 0, questCampaign.TotalQuestExecution, timestamp, questCampaignId)
		//	if err != nil {
		//		return
		//	}
		//}
		if reward > 0 {
			err = d.SelectForUpdate(db, accountId)
			if err != nil {
				return
			}
			var userReward int
			userReward, _, err = d.FindUserReward(accountId)
			if err != nil {
				return
			}
			_, err = db.Exec(dal.UpdateUserRewardByIdSql, accountId, reward+userReward, timestamp)
			if err != nil {
				return
			}
		}
		return
	})
	return
}

func (d *Dao) UpdateTotalReward(reward int) (err error) {
	timestamp := time.Now()
	_, err = d.db.Exec(dal.UpdateTotalRewardSql, reward, timestamp)
	return
}

func (d *Dao) FindUpdateStatusCampaign() (data []*model.QuestCampaign, err error) {
	rows, err := d.db.Query(dal.FindQuestCampaignByNotStatusSql, model.QuestCampaignEndedStatus)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	for rows.Next() {
		var (
			questCampaign = &model.QuestCampaign{}
			status        sql.NullString
		)
		if err = rows.Scan(&questCampaign.Id, &questCampaign.StartTime, &questCampaign.EndTime, &status); err != nil {
			return
		}
		if status.Valid {
			questCampaign.Status = status.String
		}
		data = append(data, questCampaign)
	}
	return
}

func (d *Dao) UpdateQuestCampaignStatus(id int, status string) (err error) {
	_, err = d.db.Exec(dal.UpdateQuestCampaignStatusSql, status, id)
	return
}

func (d *Dao) FindUpdateStatusQuest() (data []*model.Quest, err error) {
	rows, err := d.db.Query(dal.FindQuestByNotStatusSql, model.QuestEndedStatus)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	for rows.Next() {
		var (
			quest  = &model.Quest{}
			status sql.NullString
		)
		if err = rows.Scan(&quest.Id, &quest.StartTime, &quest.EndTime, &status); err != nil {
			return
		}
		if status.Valid {
			quest.Status = status.String
		}
		data = append(data, quest)
	}
	return
}

func (d *Dao) UpdateQuestStatus(id int, status string) (err error) {
	_, err = d.db.Exec(dal.UpdateQuestStatusSql, status, id)
	return
}

func (d *Dao) UpdateUserQuestStatus(questId int, status string) (err error) {
	_, err = d.db.Exec(dal.UpdateUserQuestStatusSql, status, questId)
	return
}

func (d *Dao) UpdateUserRewardRank(data []*model.UserReward) (err error) {
	var (
		total     = len(data)
		timestamp = time.Now()
		limit     = 50
		from      = 0
		end       = from + limit
	)
	if total == 0 {
		return
	}

	for {
		if from >= total {
			return
		}
		if end >= total {
			end = total
		}
		err = d.WithTrx(func(db *sql.Tx) (err error) {
			for index := from; index < end; index++ {
				_, err = db.Exec(dal.UpdateUserRewardRankSql, data[index].AccountId, data[index].Reward, index+1, timestamp)
				if err != nil {
					return
				}
			}
			return
		})
		if err != nil {
			return
		}
		from = end
		end = from + limit
	}
}
