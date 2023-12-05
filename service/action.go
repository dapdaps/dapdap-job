package service

import (
	"dapdap-job/common/log"
	"dapdap-job/model"
	"fmt"
	"strings"
)

var (
	maxDappId        uint64
	maxChainId       uint64
	dappParticipants = map[string]map[string]string{}
	actionsDapp      map[string]*model.ActionDapp
	actionsChain     map[string]*model.ActionChain
)

func (s *Service) InitAction() (err error) {
	maxDappId, maxChainId, err = s.dao.FindMaxRecordId()
	if err != nil {
		return
	}
	log.Info("InitAction maxDappId:%d  maxChainId:%d", maxDappId, maxChainId)
	data, err := s.dao.FindAllActions(1)
	if err != nil {
		log.Error("InitAction s.dao.FindActions error: %v", err)
		return
	}
	for _, action := range data {
		var (
			dappParticipant map[string]string
			ok              bool
		)
		dappParticipant, ok = dappParticipants[action.Template]
		if !ok {
			dappParticipant = map[string]string{}
			dappParticipants[action.Template] = dappParticipant
		}
		dappParticipant[strings.ToLower(action.AccountId)] = ""
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

func (s *Service) ActionTask() (err error) {
	var (
		minRecordId        uint64
		maxRecordId        uint64
		data               []*model.Action
		updateActionsDapp  = map[string]*model.ActionDapp{}
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
		if action.Id > maxDappId {
			var (
				dappParticipant map[string]string
				actionDapp      *model.ActionDapp
			)
			actionDapp, ok = updateActionsDapp[action.Template]
			if !ok {
				actionDapp = &model.ActionDapp{
					Template: action.Template,
				}
				updateActionsDapp[action.Template] = actionDapp
			}
			actionDapp.RecordId = action.Id
			actionDapp.Count++

			dappParticipant, ok = dappParticipants[action.Template]
			if !ok {
				dappParticipant = map[string]string{}
				dappParticipants[action.Template] = dappParticipant
			}
			_, ok = dappParticipant[strings.ToLower(action.AccountId)]
			if !ok {
				actionDapp.Participants++
				dappParticipant[strings.ToLower(action.AccountId)] = ""
			}
		}

		if action.Id > maxChainId {
			var (
				actionChain *model.ActionChain
				chainKey    = fmt.Sprintf("%s_%s_%s", action.ActionNetworkId, action.Template, action.ActionTitle)
			)
			actionChain, ok = updateActionsChain[chainKey]
			if !ok {
				actionChain = &model.ActionChain{
					ActionNetworkId: action.ActionNetworkId,
					Template:        action.Template,
					ActionTitle:     action.ActionTitle,
				}
				updateActionsChain[chainKey] = actionChain
			}
			actionChain.RecordId = action.Id
			actionChain.Count++
		}
	}

	for _, updateActionDapp := range updateActionsDapp {
		if actionDapp, ok := actionsDapp[updateActionDapp.Template]; ok {
			updateActionDapp.Count += actionDapp.Count
			updateActionDapp.Participants += actionDapp.Participants
		}
	}

	for _, updateActionChain := range updateActionsChain {
		var chainKey = fmt.Sprintf("%s_%s_%s", updateActionChain.ActionNetworkId, updateActionChain.Template, updateActionChain.ActionTitle)
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
