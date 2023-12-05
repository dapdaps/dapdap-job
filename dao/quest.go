package dao

import (
	"dapdap-job/common/dal"
	"dapdap-job/model"
	"database/sql"
	"errors"
	"strconv"
	"strings"
)

func (d *Dao) FindQuestActionMaxRecordId() (maxRecordIdActionQuest uint64, err error) {
	var (
		maxRecordIdActionQuestSql sql.NullInt64
	)
	err = d.db.QueryRow(dal.FindMaxRecordIdFromActionQuestSql).Scan(&maxRecordIdActionQuestSql)
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

func (d *Dao) FindAllQuest() (data []*model.Quest, err error) {
	rows, err := d.db.Query(dal.FindAllQuestSql)
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
		if err = rows.Scan(&quest.Id, &quest.QuestCampaignId, &quest.QuestCategoryId, &quest.StartTime, &quest.EndTime, &quest.TotalAction, &status); err != nil {
			return
		}
		if status.Valid {
			quest.Status = status.String
		}
		data = append(data, quest)
	}
	return
}

func (d *Dao) FindAllQuestAction() (data []*model.QuestAction, err error) {
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
		data = append(data, questAction)
	}
	return
}

func (d *Dao) FindUserQuest(accountIds []int) (data []*model.UserQuest, err error) {
	var (
		findSql = dal.FindUserQuestSql
		args    []interface{}
	)
	findSql += ` where `
	findSql += ` account_id in(`
	for _, accountId := range accountIds {
		findSql += `?,`
		args = append(args, accountId)
	}
	findSql = findSql[0:len(findSql)-1] + `)`
	rows, err := d.db.Query(findSql)
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
		data = append(data, userQuest)
	}
	return
}

func (d *Dao) FindUserQuestAction(accountIds []int) (data []*model.UserQuestAction, err error) {
	var (
		findSql = dal.FindUserQuestActionSql
		args    []interface{}
	)
	findSql += ` where `
	findSql += ` account_id in(`
	for _, accountId := range accountIds {
		findSql += `?,`
		args = append(args, accountId)
	}
	findSql = findSql[0:len(findSql)-1] + `)`
	rows, err := d.db.Query(findSql)
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
		data = append(data, userQuestAction)
	}
	return
}
