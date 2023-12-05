package service

import (
	"dapdap-job/common/log"
	"dapdap-job/model"
	"sort"
	"strings"
	"time"
)

var (
	maxQuestActionRecordId uint64
	questCampaignRewards   = map[int][]*model.QuestCampaignReward{}
)

func (s *Service) InitQuest() (err error) {
	maxQuestActionRecordId, err = s.dao.FindQuestActionMaxRecordId()
	if err != nil {
		return
	}
	log.Info("InitQuest maxQuestActionRecordId: %d", maxQuestActionRecordId)
	questCampaigns, err := s.dao.FindAllQuestCampaign()
	if err != nil {
		log.Error("InitQuest s.dao.FindAllQuestCampaign error: %v", err)
		return
	}
	err = s.InitQuestCampaignReward(questCampaigns)
	if err != nil {
		log.Error("InitQuest InitQuestCampaignReward error: %v", err)
		return
	}
	return
}

func (s *Service) InitQuestCampaignReward(questCampaigns []*model.QuestCampaign) (err error) {
	if len(questCampaigns) == 0 {
		return
	}
	for _, questCampaign := range questCampaigns {
		if _, ok := questCampaignRewards[questCampaign.Id]; ok {
			continue
		}
		var questCampaignReward []*model.QuestCampaignReward
		questCampaignReward, err = s.dao.FindQuestCampaignRewardById(questCampaign.Id)
		if err != nil {
			return
		}
		questCampaignRewards[questCampaign.Id] = questCampaignReward
	}
	return
}

func (s *Service) StartQuestTask() (err error) {
	questCampaigns, err := s.dao.FindAllQuestCampaign()
	if err != nil {
		log.Error("QuestTask s.dao.FindAllQuestCampaign error: %v", err)
		return
	}
	if len(questCampaigns) == 0 {
		return
	}
	err = s.InitQuestCampaignReward(questCampaigns)
	if err != nil {
		log.Error("QuestTask InitQuestCampaignReward error: %v", err)
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
		log.Info("QuestTask record %d - %d", maxQuestActionRecordId+1, maxRecordId)
		err = s.QuestTask(questCampaign, data)
		if err != nil {
			return
		}
		err = s.UpdateRank(questCampaign.Id)
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
	return
}

func (s *Service) GetUserQuestCampaignReward(questCampaignId int, accountId int) (data *model.QuestCampaignReward) {
	questCampaignReward := questCampaignRewards[questCampaignId]
	if questCampaignReward == nil {
		return
	}
	for _, userQuestReward := range questCampaignReward {
		if userQuestReward.AccountId == accountId {
			data = userQuestReward
			break
		}
	}
	return
}

func (s *Service) SetUserQuestCampaignReward(data *model.QuestCampaignReward) {
	questCampaignReward := questCampaignRewards[data.QuestCampaignId]
	if questCampaignReward == nil {
		questCampaignReward = []*model.QuestCampaignReward{}
		questCampaignRewards[data.QuestCampaignId] = questCampaignReward
	}
	questCampaignRewards[data.QuestCampaignId] = append(questCampaignRewards[data.QuestCampaignId], data)
}

func (s *Service) QuestTask(questCampaign *model.QuestCampaign, data []*model.Action) (err error) {
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
	allQuest, err = s.dao.FindAllQuest(questCampaign.Id)
	if err != nil {
		log.Error("QuestTask s.dao.FindAllQuest error: %v", err)
		return
	}
	if len(allQuest) == 0 {
		return
	}
	allQuestAction, err = s.dao.FindAllQuestAction(questCampaign.Id)
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

	allUserQuests, err = s.dao.FindUserQuest(questCampaign.Id, accountIdArr)
	if err != nil {
		log.Error("QuestTask s.dao.FindUserQuest error: %v", err)
		return
	}
	allUserQuestActions, err = s.dao.FindUserQuestAction(questCampaign.Id, accountIdArr)
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
			userReward              int
			userQuestCampaignReward int
			completedQuest          int
			newParticipant          bool
			updateQuestCampaign     *model.QuestCampaign
			questCampaignReward     *model.QuestCampaignReward
		)
		if accountId == 0 {
			continue
		}
		newParticipant = len(userQuests) == 0
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
		if newParticipant || completedQuest > 0 || reward > 0 {
			updateQuestCampaign = questCampaign
			updateQuestCampaign.TotalQuestExecution += completedQuest
			updateQuestCampaign.TotalReward += reward
			if newParticipant {
				updateQuestCampaign.TotalUsers += 1
			}
		}
		if reward > 0 {
			userReward, err = s.dao.FindUserReward(accountId)
			if err != nil {
				log.Error("QuestTask s.dao.FindUserReward error: %v", err)
				return
			}
			userReward += reward

			questCampaignReward = s.GetUserQuestCampaignReward(questCampaign.Id, accountId)
			if questCampaignReward != nil {
				userQuestCampaignReward = questCampaignReward.Reward
			}
			userQuestCampaignReward += reward
		}
		err = s.dao.UpdateUserQuest(accountId, questCampaign.Id, userReward, userQuestCampaignReward, updateQuestCampaign, saveOrUpdateUserQuests, saveOrUpdateQuestAction)
		if err != nil {
			log.Error("QuestTask s.dao.UpdateUserQuest error: %v", err)
			return
		}
		if questCampaignReward == nil {
			questCampaignReward = &model.QuestCampaignReward{
				AccountId:       accountId,
				QuestCampaignId: questCampaign.Id,
				Reward:          userQuestCampaignReward,
			}
			s.SetUserQuestCampaignReward(questCampaignReward)
		}
	}
	return
}

func (s *Service) UpdateRank(questCampaignId int) (err error) {
	questCampaignReward := questCampaignRewards[questCampaignId]
	if len(questCampaignReward) == 0 {
		return
	}
	sort.Slice(questCampaignReward, func(i, j int) bool {
		if questCampaignReward[i].Reward > questCampaignReward[j].Reward {
			return true
		} else if questCampaignReward[i].Reward < questCampaignReward[j].Reward {
			return false
		} else {
			return questCampaignReward[i].AccountId > questCampaignReward[j].AccountId
		}
	})
	err = s.dao.UpdateRewardRank(questCampaignId, questCampaignReward)
	if err != nil {
		log.Error("QuestTask UpdateRank s.dao.UpdateRewardRank error: %v", err)
		return
	}
	for index, userQuestCampaignReward := range questCampaignReward {
		userQuestCampaignReward.Rank = index + 1
	}
	return
}
