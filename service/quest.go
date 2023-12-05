package service

import (
	"dapdap-job/common/log"
	"dapdap-job/model"
	"strings"
)

var (
	maxQuestActionRecordId uint64
)

func (s *Service) InitQuest() (err error) {
	maxQuestActionRecordId, err = s.dao.FindQuestActionMaxRecordId()
	if err != nil {
		return
	}
	log.Info("InitQuest maxQuestActionRecordId:%d", maxQuestActionRecordId, maxChainId)
	return
}

func (s *Service) QuestTask() (err error) {
	var (
		maxRecordId    uint64
		data           []*model.Action
		allQuest       []*model.Quest
		allQuestAction []*model.QuestAction
		userAddress    []string
		userActions    = map[string][]*model.QuestAction{}
		//accountIds     map[string]int
	)
	data, err = s.dao.FindAllActions(maxQuestActionRecordId + 1)
	if err != nil {
		log.Error("QuestTask s.dao.FindAllActions error: %v", err)
		return
	}
	if len(data) == 0 {
		return
	}
	maxRecordId = data[len(data)-1].Id

	allQuest, err = s.dao.FindAllQuest()
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
			if userActions[action.AccountId] == nil {
				userActions[action.AccountId] = []*model.QuestAction{}
				userAddress = append(userAddress, action.AccountId)
			}
			userActions[action.AccountId] = append(userActions[action.AccountId], questAction)
		}
	}

	//accountIds, err = s.dao.FindAccountId(userAddress)
	if err != nil {
		log.Error("QuestTask s.dao.FindAccountId error: %v", err)
		return
	}

	//for address, actions := range userActions {
	//
	//}

	maxQuestActionRecordId = maxRecordId
	return
}
