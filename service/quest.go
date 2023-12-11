package service

import (
	"dapdap-job/common/log"
	"dapdap-job/model"
	"strings"
	"time"
)

var (
	maxQuestActionRecordId uint64
	//campaignInfo           *model.QuestCampaignInfo
	categories   map[string]int //map[name]id
	networks     map[int]int    //map[chain_id]id
	accountIdMap map[string]int
)

func (s *Service) InitQuest() (err error) {
	maxQuestActionRecordId, err = s.dao.FindQuestActionMaxRecordId()
	if err != nil {
		log.Error("InitQuest s.dao.FindQuestActionMaxRecordId error: %v", err)
		return
	}
	accountIdMap, err = s.dao.FindAllAccountId()
	if err != nil {
		log.Error("QuestTask s.dao.FindAllAccountId error: %v", err)
		return
	}
	//campaignInfo, err = s.dao.FindQuestCampaignInfo()
	//if err != nil {
	//	log.Error("InitQuest s.dao.FindQuestCampaignInfo error: %v", err)
	//	return
	//}
	//err = s.InitQuestCampaignReward()
	//if err != nil {
	//	log.Error("InitQuest InitQuestCampaignReward error: %v", err)
	//	return
	//}
	return
}

func (s *Service) StartQuestTask() (err error) {
	var (
		questCampaigns     []*model.QuestCampaign
		allQuest           map[int]*model.Quest
		allQuestAction     map[int]*model.QuestAction
		allDappQuestAction = map[int]*model.QuestAction{}
		sourceQuestAction  = map[int]*model.QuestAction{}
	)
	questCampaigns, err = s.dao.FindAllQuestCampaign()
	if err != nil {
		log.Error("QuestTask s.dao.FindAllQuestCampaign error: %v", err)
		return
	}
	if len(questCampaigns) == 0 {
		return
	}
	allQuest, err = s.dao.FindAllQuest(0)
	if err != nil {
		log.Error("QuestTask s.dao.FindAllQuest error: %v", err)
		return
	}
	if len(allQuest) == 0 {
		return
	}
	allQuestAction, err = s.dao.FindAllQuestAction()
	if err != nil {
		log.Error("QuestTask s.dao.FindAllQuestAction error: %v", err)
		return
	}
	if len(allQuestAction) == 0 {
		return
	}
	for questActionId, questAction := range allQuestAction {
		if _, ok := allQuest[questAction.QuestId]; !ok {
			delete(allQuestAction, questActionId)
		} else if questAction.Category == model.QuestActionDapp {
			allDappQuestAction[questActionId] = questAction
		} else if len(questAction.Source) > 0 {
			sourceQuestAction[questActionId] = questAction
		}
	}

	data, err := s.dao.FindAllActionRecords(maxQuestActionRecordId + 1)
	if err != nil {
		log.Error("QuestTask s.dao.FindAllActionRecords error: %v", err)
		return
	}
	if len(data) > 0 {
		maxActionRecordId := data[len(data)-1].Id
		categories, err = s.dao.FindActionCategory()
		if err != nil {
			log.Error("QuestTask s.dao.FindActionCategory error: %v", err)
			return
		}
		networks, err = s.dao.FindNetworks()
		if err != nil {
			log.Error("QuestTask s.dao.FindNetworks error: %v", err)
			return
		}
		log.Info("QuestTask action record %d - %d", maxQuestActionRecordId+1, maxActionRecordId)
		err = s.QuestActionTask(allQuest, allDappQuestAction, data)
		if err != nil {
			return
		}
		for {
			err = s.dao.UpdateActionRecord(maxActionRecordId)
			if err != nil {
				log.Error("QuestTask s.dao.UpdateActionRecord error: %v", err)
				time.Sleep(time.Second * 5)
			}
			break
		}
		maxQuestActionRecordId = maxActionRecordId
	}

	//err = s.CheckQuestStatus(allQuest, sourceQuestAction)
	//if err != nil {
	//	return
	//}

	//
	//sourceRecordData, err := s.dao.FindAllSourceRecords(maxQuestActionRecordId + 1)
	//if err != nil {
	//	log.Error("QuestTask s.dao.FindAllSourceRecords error: %v", err)
	//	return
	//}
	//if len(sourceRecordData) > 0 {
	//	maxSourceRecordId := data[len(sourceRecordData)-1].Id
	//	log.Info("QuestTask source record %d - %d", maxQuestSourceRecordId+1, maxSourceRecordId)
	//	err = s.QuestSourceTask(sourceRecordData)
	//	if err != nil {
	//		return
	//	}
	//	for {
	//		err = s.dao.UpdateSourceRecord(maxSourceRecordId)
	//		if err != nil {
	//			log.Error("QuestTask s.dao.UpdateSourceRecord error: %v", err)
	//			time.Sleep(time.Second * 5)
	//		}
	//		break
	//	}
	//	maxQuestSourceRecordId = maxSourceRecordId
	//}
	return
}

func (s *Service) QuestActionTask(allQuest map[int]*model.Quest, allQuestAction map[int]*model.QuestAction, data []*model.Action) (err error) {
	var (
		userAddress               []string
		accountIdArr              []int
		allUserQuestActionRecords = map[string][]*model.UserQuestAction{}
		allUserQuests             map[int][]*model.UserQuest
		allUserQuestActions       map[int][]*model.UserQuestAction
	)
	for _, action := range data {
		for _, questAction := range allQuestAction {
			var ok bool
			if questAction.CategoryId != categories[action.ActionType] {
				continue
			}
			if len(questAction.Source) > 0 && !strings.EqualFold(questAction.Source, action.Source) {
				continue
			}
			if len(questAction.DappsMap) > 0 {
				if _, ok = questAction.DappsMap[action.DappId]; !ok {
					continue
				}
			}
			if len(questAction.NetworksMap) > 0 {
				if _, ok = questAction.NetworksMap[networks[action.ChainId]]; !ok {
					continue
				}
			}
			if len(questAction.ToNetworksMap) > 0 {
				if _, ok = questAction.ToNetworksMap[networks[action.ToCahinId]]; !ok {
					continue
				}
			}
			if allUserQuestActionRecords[action.AccountId] == nil {
				allUserQuestActionRecords[action.AccountId] = []*model.UserQuestAction{}
				userAddress = append(userAddress, action.AccountId)
			}
			var sameQuestAction bool
			for _, userQuestActionRecord := range allUserQuestActionRecords[action.AccountId] {
				if userQuestActionRecord.QuestActionId == questAction.Id {
					sameQuestAction = true
					userQuestActionRecord.Times++
					break
				}
			}
			if !sameQuestAction {
				allUserQuestActionRecords[action.AccountId] = append(allUserQuestActionRecords[action.AccountId], &model.UserQuestAction{
					QuestActionId:   questAction.Id,
					QuestId:         questAction.QuestId,
					QuestCampaignId: questAction.QuestCampaignId,
					Times:           1,
				})
			}
		}
	}

	err = s.UpdateCacheAccountId(userAddress)
	if err != nil {
		log.Error("QuestTask s.UpdateCacheAccountId error: %v", err)
		return
	}

	for _, addr := range userAddress {
		if accountId, ok := accountIdMap[addr]; ok {
			accountIdArr = append(accountIdArr, accountId)
		}
	}
	allUserQuests, err = s.dao.FindUserQuest(0, accountIdArr)
	if err != nil {
		log.Error("QuestTask s.dao.FindUserQuest error: %v", err)
		return
	}
	allUserQuestActions, err = s.dao.FindUserQuestAction(0, accountIdArr)
	if err != nil {
		log.Error("QuestTask s.dao.FindUserQuestAction error: %v", err)
		return
	}

	for address, userQuestActionRecords := range allUserQuestActionRecords {
		var (
			accountId               = accountIdMap[address]
			userQuests              = allUserQuests[accountId]
			userQuestActions        = allUserQuestActions[accountId]
			saveOrUpdateUserQuests  []*model.UserQuest
			saveOrUpdateQuestAction []*model.UserQuestAction
			reward                  int
			completedQuest          int
			//newParticipant          = len(userQuests) == 0
			//updateCampaignInfo      *model.QuestCampaignInfo
		)
		if accountId == 0 {
			continue
		}
		for _, userQuestActionRecord := range userQuestActionRecords {
			var (
				quest                   *model.Quest
				questAction             *model.QuestAction
				userQuestAction         *model.UserQuestAction
				userQuest               *model.UserQuest
				hasQuestCompleted       bool
				hasQuestActionCompleted bool
				times                   int
				questActionStatus       string
				actionCompleted         int
				questStatus             string
			)

			if quest = allQuest[userQuestActionRecord.QuestId]; quest == nil {
				log.Info("QuestTask not find quest id: %d", userQuestActionRecord.QuestId)
				continue
			}
			if questAction = allQuestAction[userQuestActionRecord.QuestActionId]; questAction == nil {
				log.Info("QuestTask not find quest action id: %d", userQuestActionRecord.QuestActionId)
				continue
			}

			for _, mUserQuest := range userQuests {
				if mUserQuest.QuestId == userQuestActionRecord.QuestId {
					if mUserQuest.Status == model.UserQuestCompletedStatus || mUserQuest.Status == model.UserQuestExpiredStatus {
						hasQuestCompleted = true
					} else {
						userQuest = mUserQuest
					}
					break
				}
			}
			if hasQuestCompleted {
				continue
			}

			for _, mUserQuestAction := range userQuestActions {
				if mUserQuestAction.QuestActionId == userQuestActionRecord.QuestActionId {
					if strings.EqualFold(mUserQuestAction.Status, model.UserQuestActionCompletedStatus) || strings.EqualFold(mUserQuestAction.Status, model.UserQuestActionExpiredStatus) {
						hasQuestActionCompleted = true
					} else {
						userQuestAction = mUserQuestAction
					}
					break
				}
			}
			if hasQuestActionCompleted {
				continue
			}

			if userQuestAction != nil {
				times = userQuestAction.Times
			}
			times += userQuestActionRecord.Times
			if times >= questAction.Times {
				questActionStatus = model.UserQuestActionCompletedStatus
			} else {
				questActionStatus = model.UserQuestActionInProcessStatus
			}
			if userQuestAction != nil {
				userQuestAction.Times = times
				userQuestAction.Status = questActionStatus
			} else {
				userQuestAction = &model.UserQuestAction{
					QuestActionId:   userQuestActionRecord.QuestActionId,
					QuestId:         userQuestActionRecord.QuestId,
					QuestCampaignId: userQuestActionRecord.QuestCampaignId,
					AccountId:       accountId,
					Times:           times,
					Status:          questActionStatus,
				}
			}
			saveOrUpdateQuestAction = append(saveOrUpdateQuestAction, userQuestAction)

			if times >= questAction.Times {
				if userQuest != nil {
					actionCompleted = userQuest.ActionCompleted
				}
				actionCompleted++
				if actionCompleted >= quest.TotalAction {
					questStatus = model.UserQuestCompletedStatus
				} else {
					questStatus = model.UserQuestInProcessStatus
				}
				if userQuest == nil {
					userQuest = &model.UserQuest{
						QuestId:         userQuestActionRecord.QuestId,
						QuestCampaignId: userQuestActionRecord.QuestCampaignId,
						AccountId:       accountId,
					}
					if userQuests == nil {
						userQuests = []*model.UserQuest{}
					}
					userQuests = append(userQuests, userQuest)
					allUserQuests[accountId] = userQuests
				}
				userQuest.ActionCompleted = actionCompleted
				userQuest.Status = questStatus

				var saveOrUpdateQuestExist bool
				for _, saveOrUpdateQuest := range saveOrUpdateUserQuests {
					if saveOrUpdateQuest.QuestId == userQuest.QuestId {
						saveOrUpdateQuestExist = true
						break
					}
				}
				if !saveOrUpdateQuestExist {
					saveOrUpdateUserQuests = append(saveOrUpdateUserQuests, userQuest)
				}
			}
		}

		for _, saveOrUpdateQuest := range saveOrUpdateUserQuests {
			if saveOrUpdateQuest.Status == model.UserQuestCompletedStatus {
				quest := allQuest[saveOrUpdateQuest.QuestId]
				reward += quest.Reward
				completedQuest++
			}
		}
		//if newParticipant || completedQuest > 0 || reward > 0 {
		//	updateCampaignInfo = &model.QuestCampaignInfo{}
		//	updateCampaignInfo.TotalQuestExecution = campaignInfo.TotalQuestExecution + completedQuest
		//	updateCampaignInfo.TotalReward = campaignInfo.TotalReward + reward
		//	if newParticipant {
		//		updateCampaignInfo.TotalUsers = campaignInfo.TotalUsers + 1
		//	}
		//}

		err = s.dao.UpdateUserQuest(accountId, reward, saveOrUpdateUserQuests, saveOrUpdateQuestAction)
		if err != nil {
			log.Error("QuestTask s.dao.UpdateUserQuest error: %v", err)
			return
		}
		//if updateCampaignInfo != nil {
		//	campaignInfo.TotalUsers = updateCampaignInfo.TotalUsers
		//	campaignInfo.TotalReward = updateCampaignInfo.TotalReward
		//	campaignInfo.TotalQuestExecution = updateCampaignInfo.TotalQuestExecution
		//}
	}
	return
}

func (s *Service) UpdateCacheAccountId(userAddress []string) (err error) {
	var accountAddress []string
	for _, address := range userAddress {
		if _, ok := accountIdMap[address]; !ok {
			accountAddress = append(accountAddress, address)
		}
	}
	if len(accountAddress) > 0 {
		var data map[string]int
		data, _, err = s.dao.FindAccountId(accountAddress)
		if err != nil {
			log.Error("QuestTask UpdateCacheAccountId s.dao.FindAccountId error: %v", err)
			return
		}
		for addr, accountId := range data {
			accountIdMap[addr] = accountId
		}
	}
	return
}

//func (s *Service) CheckQuestStatus(allQuest map[int]*model.Quest, allQuestAction map[int]*model.QuestAction) (err error) {
//	return
//}

//func (s *Service) QuestSourceTask(data []*model.QuestSourceRecord) (err error) {
//	var (
//		allQuest            map[int]*model.Quest
//		accountIdMap        map[int]int
//		accountIdArr        []int
//		allUserQuests       map[int][]*model.UserQuest
//		allUserQuestActions map[int][]*model.UserQuestAction
//	)
//	allQuest, err = s.dao.FindAllQuest(0)
//	if err != nil {
//		log.Error("QuestTask s.dao.FindAllQuest error: %v", err)
//		return
//	}
//	if len(allQuest) == 0 {
//		return
//	}
//	for _, record := range data {
//		if _, ok := accountIdMap[record.AccountId]; ok {
//			continue
//		}
//		accountIdArr = append(accountIdArr, record.AccountId)
//	}
//	allUserQuests, err = s.dao.FindUserQuest(0, accountIdArr)
//	if err != nil {
//		log.Error("QuestTask s.dao.FindUserQuest error: %v", err)
//		return
//	}
//	allUserQuestActions, err = s.dao.FindUserQuestAction(0, accountIdArr)
//	if err != nil {
//		log.Error("QuestTask s.dao.FindUserQuestAction error: %v", err)
//		return
//	}
//	for _, record := range data {
//		var (
//			accountId               = record.AccountId
//			userQuests              = allUserQuests[accountId]
//			userQuestActions        = allUserQuestActions[accountId]
//			quest                   *model.Quest
//			userQuestAction         *model.UserQuestAction
//			userQuest               *model.UserQuest
//			reward                  int
//			completedQuest          int
//			newParticipant          bool
//			updateCampaignInfo      *model.QuestCampaignInfo
//			hasQuestCompleted       bool
//			hasQuestActionCompleted bool
//			questActionStatus       string
//			questStatus             string
//		)
//
//		for _, mUserQuest := range userQuests {
//			if mUserQuest.QuestId == record.QuestId {
//				if mUserQuest.Status == model.UserQuestCompletedStatus || mUserQuest.Status == model.UserQuestExpiredStatus {
//					hasQuestCompleted = true
//				} else {
//					userQuest = mUserQuest
//				}
//				break
//			}
//		}
//		if hasQuestCompleted {
//			continue
//		}
//
//		for _, mUserQuestAction := range userQuestActions {
//			if mUserQuestAction.QuestActionId == record.QuestActionId {
//				if strings.EqualFold(mUserQuestAction.Status, model.UserQuestActionCompletedStatus) || strings.EqualFold(mUserQuestAction.Status, model.UserQuestActionExpiredStatus) {
//					hasQuestActionCompleted = true
//				} else {
//					userQuestAction = mUserQuestAction
//				}
//				break
//			}
//		}
//		if hasQuestActionCompleted {
//			continue
//		}
//
//		newParticipant = len(userQuests) == 0
//		quest = allQuest[record.QuestId]
//		questActionStatus = model.UserQuestActionCompletedStatus
//		userQuestAction = &model.UserQuestAction{
//			QuestActionId:   record.QuestActionId,
//			QuestId:         record.QuestId,
//			QuestCampaignId: record.QuestCampaignId,
//			AccountId:       accountId,
//			Times:           1,
//			Status:          questActionStatus,
//		}
//		if userQuest == nil {
//			userQuest = &model.UserQuest{
//				QuestId:         record.QuestId,
//				QuestCampaignId: record.QuestCampaignId,
//				AccountId:       accountId,
//			}
//			if userQuests == nil {
//				userQuests = []*model.UserQuest{}
//			}
//			userQuests = append(userQuests, userQuest)
//			allUserQuests[accountId] = userQuests
//		}
//		userQuest.ActionCompleted++
//		if userQuest.ActionCompleted >= quest.TotalAction {
//			questStatus = model.UserQuestCompletedStatus
//		} else {
//			questStatus = model.UserQuestInProcessStatus
//		}
//		userQuest.Status = questStatus
//		if userQuest.Status == model.UserQuestCompletedStatus {
//			reward += quest.Reward
//			completedQuest++
//		}
//		if newParticipant || completedQuest > 0 || reward > 0 {
//			updateCampaignInfo = &model.QuestCampaignInfo{}
//			updateCampaignInfo.TotalQuestExecution = campaignInfo.TotalQuestExecution + completedQuest
//			updateCampaignInfo.TotalReward = campaignInfo.TotalReward + reward
//			if newParticipant {
//				updateCampaignInfo.TotalUsers = campaignInfo.TotalUsers + 1
//			}
//		}
//		err = s.dao.UpdateUserQuest(accountId, record.QuestCampaignId, reward, updateCampaignInfo, []*model.UserQuest{userQuest}, []*model.UserQuestAction{userQuestAction})
//		if err != nil {
//			log.Error("QuestTask s.dao.UpdateUserQuest error: %v", err)
//			return
//		}
//		if updateCampaignInfo != nil {
//			campaignInfo.TotalUsers = updateCampaignInfo.TotalUsers
//			campaignInfo.TotalReward = updateCampaignInfo.TotalReward
//			campaignInfo.TotalQuestExecution = updateCampaignInfo.TotalQuestExecution
//		}
//	}
//	return
//}
