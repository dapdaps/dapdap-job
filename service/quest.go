package service

import (
	"dapdap-job/common/log"
	"dapdap-job/model"
	"strings"
	"time"
)

var (
	maxQuestActionRecordId uint64
)

func (s *Service) InitQuest() (err error) {
	maxQuestActionRecordId, err = s.dao.FindQuestActionMaxRecordId()
	if err != nil {
		return
	}
	log.Info("InitQuest maxQuestActionRecordId: %d", maxQuestActionRecordId)
	return
}

func (s *Service) StartQuestTask() {
	questCampaigns, err := s.dao.FindAllQuestCampaign()
	if err != nil {
		log.Error("QuestTask s.dao.FindAllQuestCampaign error: %v", err)
		return
	}
	if len(questCampaigns) == 0 {
		return
	}
	data, err := s.dao.FindAllActions(maxQuestActionRecordId + 1)
	if err != nil {
		log.Error("QuestTask s.dao.FindAllActions error: %v", err)
		return
	}
	if len(data) == 0 {
		return
	}
	maxRecordId := data[len(data)-1].Id
	for _, questCampaign := range questCampaigns {
		err = s.QuestTask(questCampaign.Id, data)
		if err != nil {
			return
		}
	}
	for {
		err = s.dao.UpdateActionRecord(maxRecordId)
		if err != nil {
			log.Error("QuestTask s.dao.UpdateActionRecord error: %v", err)
			time.Sleep(time.Second * 5)
		}
		break
	}
	maxQuestActionRecordId = maxRecordId
}

func (s *Service) QuestTask(questCampaignId int, data []*model.Action) (err error) {
	var (
		allQuest                  map[int]*model.Quest
		allQuestAction            map[int]*model.QuestAction
		userAddress               []string
		allUserQuestActionRecords = map[string][]*model.UserQuestAction{}
		accountIdMap              map[string]int
		accountIdArr              []int
		allUserQuests             map[int][]*model.UserQuest
		allUserQuestActions       map[int][]*model.UserQuestAction
	)
	allQuest, err = s.dao.FindAllQuest(questCampaignId)
	if err != nil {
		log.Error("QuestTask s.dao.FindAllQuest error: %v", err)
		return
	}
	if len(allQuest) == 0 {
		return
	}
	allQuestAction, err = s.dao.FindAllQuestAction(questCampaignId)
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
		}
	}

	for _, action := range data {
		for _, questAction := range allQuestAction {
			var ok bool
			if questAction.CategoryId != action.DappCategoryId {
				continue
			}
			if len(questAction.Source) > 0 && !strings.EqualFold(questAction.Source, action.Source) {
				continue
			}
			if _, ok = questAction.DappsMap[action.DappId]; !ok {
				continue
			}
			if _, ok = questAction.NetworksMap[action.NetworkId]; !ok {
				continue
			}
			if len(questAction.ToNetworksMap) > 0 {
				if _, ok = questAction.ToNetworksMap[action.ToNetworkId]; !ok {
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

	accountIdMap, accountIdArr, err = s.dao.FindAccountId(userAddress)
	if err != nil {
		log.Error("QuestTask s.dao.FindAccountId error: %v", err)
		return
	}

	allUserQuests, err = s.dao.FindUserQuest(accountIdArr)
	if err != nil {
		log.Error("QuestTask s.dao.FindUserQuest error: %v", err)
		return
	}
	allUserQuestActions, err = s.dao.FindUserQuestAction(accountIdArr)
	if err != nil {
		log.Error("QuestTask s.dao.FindUserQuestAction error: %v", err)
		return
	}

	for address, userQuestActionRecords := range allUserQuestActionRecords {
		var (
			accountId               = accountIdMap[address]
			userQuests              = allUserQuests[accountId]
			userQuestActions        = allUserQuestActions[accountId]
			saveOrUpdateQuest       []*model.UserQuest
			saveOrUpdateQuestAction []*model.UserQuestAction
			reward                  int
			userReward              int
			userQuestCampaignReward int
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

			if questAction = allQuestAction[userQuestActionRecord.QuestActionId]; questAction == nil {
				log.Info("QuestTask not find quest action id: %d", userQuestActionRecord.QuestActionId)
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
					allUserQuests[accountId] = append(allUserQuests[accountId], userQuest)
				}
				userQuest.ActionCompleted = actionCompleted
				userQuest.Status = questStatus
				saveOrUpdateQuest = append(saveOrUpdateQuest, userQuest)
			}
		}

		for _, updateUserQuest := range saveOrUpdateQuest {
			if updateUserQuest.Status == model.UserQuestCompletedStatus {
				quest := allQuest[updateUserQuest.QuestId]
				reward += quest.Reward
			}
		}

		if reward > 0 {
			userReward, err = s.dao.FindUserReward(accountId)
			if err != nil {
				log.Error("QuestTask s.dao.FindUserReward error: %v", err)
				return
			}
			userReward += reward
			userQuestCampaignReward, err = s.dao.FindQuestCampaignReward(accountId, questCampaignId)
			if err != nil {
				log.Error("QuestTask s.dao.FindQuestCampaignReward error: %v", err)
				return
			}
			userQuestCampaignReward += reward
		}
		err = s.dao.UpdateUserQuest(accountId, questCampaignId, userReward, userQuestCampaignReward, saveOrUpdateQuest, saveOrUpdateQuestAction)
		if err != nil {
			log.Error("QuestTask s.dao.UpdateUserQuest error: %v", err)
			return
		}
	}
	return
}