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

func (d *Dao) InitQuestActionRecord() (err error) {
	var countSql sql.NullInt64
	err = d.db.QueryRow(dal.FindQuestRecordCountSql).Scan(&countSql)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return
	}
	err = nil
	if countSql.Valid && countSql.Int64 > 0 {
		return
	}
	_, err = d.db.Exec(dal.InitQuestRecordSql, 0, 0)
	return
}

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

func (d *Dao) FindQuestActionMaxRecordId() (maxActionRecordId uint64, err error) {
	var (
		maxActionRecordIdSql sql.NullInt64
	)
	err = d.db.QueryRow(dal.FindMaxQuestRecordIdSql).Scan(&maxActionRecordIdSql)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	if maxActionRecordIdSql.Valid {
		maxActionRecordId = uint64(maxActionRecordIdSql.Int64)
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
	var (
		findSql = dal.FindAllQuestByStatusIdSql
		rows    *sql.Rows
	)
	data = map[int]*model.Quest{}
	if questCampaignId > 0 {
		findSql += ` and quest_campaign_id=$2 order by id asc`
		rows, err = d.db.Query(findSql, model.QuestOnGoingStatus, questCampaignId)
	} else {
		findSql += ` order by id asc`
		rows, err = d.db.Query(findSql, model.QuestOnGoingStatus)
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

func (d *Dao) FindQuest(id int) (quest *model.Quest, err error) {
	var status sql.NullString
	quest = &model.Quest{}
	err = d.db.QueryRow(dal.FindQuestByIdSql, id).Scan(&quest.Id, &quest.QuestCampaignId, &quest.QuestCategoryId, &quest.TotalAction, &status, &quest.Reward)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	if status.Valid {
		quest.Status = status.String
	}
	return
}

func (d *Dao) FindAllQuestAction() (data map[int]*model.QuestAction, err error) {
	data = map[int]*model.QuestAction{}
	rows, err := d.db.Query(dal.FindAllQuestActionSql)
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
			category    sql.NullString
			categoryId  sql.NullInt64
			source      sql.NullString
			dapps       sql.NullString
			networks    sql.NullString
			toNetworks  sql.NullString
		)
		if err = rows.Scan(&questAction.Id, &questAction.QuestCampaignId, &questAction.QuestId, &questAction.Times, &category, &categoryId, &source, &dapps, &networks, &toNetworks); err != nil {
			return
		}
		if category.Valid {
			questAction.Category = category.String
		}
		if categoryId.Valid {
			questAction.CategoryId = int(categoryId.Int64)
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

func (d *Dao) FindUserQuests(questCampaignId int, accountIds []int) (data map[int][]*model.UserQuest, err error) {
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
	findSql = findSql[0:len(findSql)-1] + `)`

	if questCampaignId > 0 {
		args = append(args, questCampaignId)
		findSql += ` and quest_campaign_id=$` + strconv.Itoa(len(args))
	}

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

func (d *Dao) FindUserQuest(accountId int, questId int) (userQuest *model.UserQuest, err error) {
	userQuest = &model.UserQuest{}
	err = d.db.QueryRow(dal.FindUserQuestSql+` where account_id=$1 and quest_id=$2`, accountId, questId).Scan(&userQuest.Id, &userQuest.QuestCampaignId, &userQuest.QuestId, &userQuest.AccountId, &userQuest.ActionCompleted, &userQuest.Status)
	if err != nil {
		userQuest = nil
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	return
}

func (d *Dao) FindUserQuestActions(questCampaignId int, accountIds []int) (data map[int][]*model.UserQuestAction, err error) {
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
	findSql = findSql[0:len(findSql)-1] + `)`

	if questCampaignId > 0 {
		args = append(args, questCampaignId)
		findSql += ` and quest_campaign_id=$` + strconv.Itoa(len(args))
	}
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

func (d *Dao) FindUserQuestAction(accountId int, questActionId int) (userQuestAction *model.UserQuestAction, err error) {
	userQuestAction = &model.UserQuestAction{}
	err = d.db.QueryRow(dal.FindUserQuestActionSql+` where account_id=$1 and quest_action_id=$2`, accountId, questActionId).Scan(&userQuestAction.Id, &userQuestAction.QuestCampaignId, &userQuestAction.QuestId, &userQuestAction.QuestActionId, &userQuestAction.AccountId, &userQuestAction.Times, &userQuestAction.Status)
	if err != nil {
		userQuestAction = nil
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	return
}

func (d *Dao) FindUserQuestActionByQuestId(accountId int, questId int) (data []*model.UserQuestAction, err error) {
	rows, err := d.db.Query(dal.FindUserQuestActionSql+` where account_id=$1 and quest_id=$2`, accountId, questId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var userQuestAction = &model.UserQuestAction{}
		if err = rows.Scan(&userQuestAction.Id, &userQuestAction.QuestCampaignId, &userQuestAction.QuestId, &userQuestAction.QuestActionId, &userQuestAction.AccountId, &userQuestAction.Times, &userQuestAction.Status); err != nil {
			return
		}
		data = append(data, userQuestAction)
	}
	return
}

func (d *Dao) FindQuestActionByCategory(category string) (questAction *model.QuestAction, err error) {
	var (
		categoryId sql.NullInt64
		source     sql.NullString
		dapps      sql.NullString
		networks   sql.NullString
		toNetworks sql.NullString
	)
	questAction = &model.QuestAction{}
	err = d.db.QueryRow(dal.FindQuestActionByCategorySql, category).Scan(&questAction.Id, &questAction.QuestCampaignId, &questAction.QuestId, &questAction.Times, &questAction.Category, &categoryId, &source, &dapps, &networks, &toNetworks)
	if err != nil {
		questAction = nil
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	if categoryId.Valid {
		questAction.CategoryId = int(categoryId.Int64)
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
	return
}

func (d *Dao) UpdateActionRecord(id uint64) (err error) {
	_, err = d.db.Exec(dal.UpdateQuestActionRecordIdSql, id)
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
				if data[index].Rank != index+1 {
					_, err = db.Exec(dal.UpdateUserRewardRankSql, data[index].AccountId, data[index].ClaimedReward, index+1, timestamp)
					if err != nil {
						return
					}
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

func (d *Dao) FindQuestSourceRecordMaxId() (id uint64, err error) {
	var (
		maxRecordIdSql sql.NullInt64
	)
	err = d.db.QueryRow(dal.FindQuestSourceRecordMaxIdSql).Scan(&maxRecordIdSql)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	if maxRecordIdSql.Valid {
		id = uint64(maxRecordIdSql.Int64)
	}
	return
}

func (d *Dao) FindAllSourceRecords(id uint64) (data []*model.QuestSourceRecord, err error) {
	var (
		maxId uint64
		rows  *sql.Rows
		limit = uint64(5000)
	)
	maxId, err = d.FindQuestSourceRecordMaxId()
	if err != nil {
		return
	}
	if maxId < id {
		return
	}
	for {
		rows, err = d.db.Query(dal.FindQuestSourceRecordByBetweenSql, id, id+limit)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return
		}
		err = nil
		for rows.Next() {
			var record = &model.QuestSourceRecord{}
			if err = rows.Scan(&record.Id, &record.Source, &record.AccountId, &record.QuestActionId, &record.QuestId, &record.QuestCampaignId); err != nil {
				return
			}
			data = append(data, record)
		}
		_ = rows.Close()
		if len(data) > 0 && data[len(data)-1].Id >= maxId {
			return
		}
		id = id + limit + 1
	}
}

func (d *Dao) FindQuestTotalUsers() (totalUsers int64, err error) {
	var (
		totalUsersSql sql.NullInt64
	)
	err = d.db.QueryRow(dal.FindQuestTotalUsersSql).Scan(&totalUsersSql)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	if totalUsersSql.Valid {
		totalUsers = totalUsersSql.Int64
	}
	return
}

func (d *Dao) FindQuestCampaignTotalUsers() (campaignTotalUsers map[int]map[int]bool, totalUsers int64, err error) {
	campaignTotalUsers = map[int]map[int]bool{}
	rows, err := d.db.Query(dal.FindQuestCampaignTotalUsersSql)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var (
			accountId         int
			questCampaignId   int
			campaignTotalUser map[int]bool
			hasExist          bool
			ok                bool
		)
		if err = rows.Scan(&accountId, &questCampaignId); err != nil {
			return
		}
		if campaignTotalUser, ok = campaignTotalUsers[questCampaignId]; !ok {
			campaignTotalUser = map[int]bool{}
			campaignTotalUsers[questCampaignId] = campaignTotalUser
		}
		hasExist = false
		for _, campaignTotalUser = range campaignTotalUsers {
			if _, ok = campaignTotalUser[accountId]; ok {
				hasExist = true
				break
			}
		}
		campaignTotalUser[accountId] = true
		if !hasExist {
			totalUsers++
		}
	}
	return
}

func (d *Dao) FindQuestTotalExecutions() (totalExecutions int64, err error) {
	var (
		totalExecutionsSql sql.NullInt64
	)
	err = d.db.QueryRow(dal.FindQuestTotalExecutionSql, model.UserQuestCompletedStatus).Scan(&totalExecutionsSql)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	if totalExecutionsSql.Valid {
		totalExecutions = totalExecutionsSql.Int64
	}
	return
}

func (d *Dao) UpdateCampaignInfo(reward int, totalUsers int64, totalReward int64) (err error) {
	timestamp := time.Now()
	_, err = d.db.Exec(dal.UpdateCampaignInfoSql, reward, totalUsers, totalReward, timestamp)
	return
}

func (d *Dao) UpdateCampaignTotalUsers(id int, totalUsers int) (err error) {
	timestamp := time.Now()
	_, err = d.db.Exec(dal.UpdateQuestCampaignTotalUsersSql, totalUsers, timestamp, id)
	return
}

func (d *Dao) UpdateUserQuest(userQuests []*model.UserQuest, userQuestActions []*model.UserQuestAction) (err error) {
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
		return
	})
	return
}

func (d *Dao) FindLongQuest(category string) (quest *model.QuestLong, err error) {
	quest = &model.QuestLong{}
	err = d.db.QueryRow(dal.FindQuestLongSql, category, model.QuestOnGoingStatus).Scan(&quest.Id, &quest.Rule)
	if err != nil {
		quest = nil
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
		return
	}
	return
}
