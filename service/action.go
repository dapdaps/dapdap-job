package service

import (
	"dapdap-job/common/log"
	"dapdap-job/model"
	"fmt"
)

var (
	maxDappId        uint64
	maxChainId       uint64
	dappParticipants = map[int]map[string]string{}
	actionsDapp      map[int]*model.ActionDapp
	actionsChain     map[string]*model.ActionChain
)

func (s *Service) InitAction() (err error) {
	maxDappId, maxChainId, err = s.dao.FindMaxRecordId()
	if err != nil {
		return
	}
	log.Info("InitAction maxDappId:%d  maxChainId:%d", maxDappId, maxChainId)
	data, err := s.dao.FindRangeActions(1, maxDappId)
	if err != nil {
		log.Error("InitAction s.dao.FindActions error: %v", err)
		return
	}
	for _, action := range data {
		var (
			dappParticipant map[string]string
			ok              bool
		)
		dappParticipant, ok = dappParticipants[action.DappId]
		if !ok {
			dappParticipant = map[string]string{}
			dappParticipants[action.DappId] = dappParticipant
		}
		dappParticipant[action.AccountId] = ""
	}

	actionsDapp, err = s.dao.FindActionsDapp()
	if err != nil {
		log.Error("InitAction s.dao.FindActionsDapp error: %v", err)
		return
	}

	actionsChain, err = s.dao.FindActionsChain()
	if err != nil {
		log.Error("InitAction s.dao.FindActionsChain error: %v", err)
		return
	}
	return
}

func (s *Service) StartActionTask() (err error) {
	var (
		minRecordId        uint64
		maxRecordId        uint64
		data               []*model.Action
		invitedUserAddress = map[string]string{}
		updateActionsDapp  = map[int]*model.ActionDapp{}
		updateActionsChain = map[string]*model.ActionChain{}
	)
	if maxDappId < maxChainId {
		minRecordId = maxDappId
	} else {
		minRecordId = maxChainId
	}
	data, err = s.dao.FindActions(minRecordId+1, 500)
	if err != nil {
		log.Error("ActionTask s.dao.FindActions error: %v", err)
		return
	}
	if len(data) == 0 {
		return
	}
	maxRecordId = data[len(data)-1].Id

	for _, action := range data {
		var (
			ok bool
		)

		if len(action.AccountId) > 0 {
			invitedUserAddress[action.AccountId] = action.AccountId
		}

		if action.Id > maxDappId {
			var (
				dappParticipant map[string]string
				actionDapp      *model.ActionDapp
			)
			if action.DappId > 0 {
				if actionDapp, ok = updateActionsDapp[action.DappId]; !ok {
					actionDapp = &model.ActionDapp{
						DappId: action.DappId,
					}
					updateActionsDapp[action.DappId] = actionDapp
				}
				actionDapp.RecordId = action.Id
				actionDapp.Count++

				if len(action.AccountId) > 0 {
					dappParticipant, ok = dappParticipants[action.DappId]
					if !ok {
						dappParticipant = map[string]string{}
						dappParticipants[action.DappId] = dappParticipant
					}
					_, ok = dappParticipant[action.AccountId]
					if !ok {
						actionDapp.Participants++
						dappParticipant[action.AccountId] = ""
					}
				}
			}
		}

		if action.Id > maxChainId {
			var (
				actionChain *model.ActionChain
				networkId   = networks[action.ChainId]
				chainKey    = fmt.Sprintf("%d_%d_%s", networkId, action.DappId, action.ActionTitle)
			)
			if networkId > 0 && action.DappId > 0 && len(action.ActionTitle) > 0 {
				actionChain, ok = updateActionsChain[chainKey]
				if !ok {
					actionChain = &model.ActionChain{
						NetworkId:   networkId,
						DappId:      action.DappId,
						ActionTitle: action.ActionTitle,
					}
					updateActionsChain[chainKey] = actionChain
				}
				actionChain.RecordId = action.Id
				actionChain.Count++
			}
		}
	}

	err = s.UpdateInviteReward(invitedUserAddress)
	if err != nil {
		return
	}

	for _, updateActionDapp := range updateActionsDapp {
		if actionDapp, ok := actionsDapp[updateActionDapp.DappId]; ok {
			updateActionDapp.Count += actionDapp.Count
			updateActionDapp.Participants += actionDapp.Participants
		}
	}

	for _, updateActionChain := range updateActionsChain {
		var chainKey = fmt.Sprintf("%d_%d_%s", updateActionChain.NetworkId, updateActionChain.DappId, updateActionChain.ActionTitle)
		if actionChain, ok := actionsChain[chainKey]; ok {
			updateActionChain.Count += actionChain.Count
		}
	}

	err = s.dao.UpdateActionsDapp(updateActionsDapp)
	if err != nil {
		log.Error("ActionTask s.dao.UpdateActionsDapp error: %v", err)
		return
	}
	maxDappId = maxRecordId

	err = s.dao.UpdateActionsChain(updateActionsChain)
	if err != nil {
		log.Error("ActionTask s.dao.UpdateActionsChain error: %v", err)
		return
	}
	maxChainId = maxRecordId
	return
}
