package service

import (
	"dapdap-job/common/log"
	"dapdap-job/model"
	"time"
)

var (
	tfQuestAction                *model.QuestAction
	tlQuestAction                *model.QuestAction
	trQuestAction                *model.QuestAction
	tqQuestAction                *model.QuestAction
	tcQuestAction                *model.QuestAction
	tQuest                       *model.Quest
	allTwitterUserQuestCompleted map[int]bool
)

func (s *Service) StartTwitter() {
	var (
		err error
	)
	for {
		var quest *model.Quest
		tfQuestAction, quest, err = s.GetQuestActionByCategory("twitter_follow")
		if err != nil {
			time.Sleep(time.Second * 5)
			continue
		}
		if tQuest == nil {
			tQuest = quest
		}
		break
	}
	for {
		var quest *model.Quest
		tlQuestAction, quest, err = s.GetQuestActionByCategory("twitter_like")
		if err != nil {
			time.Sleep(time.Second * 5)
			continue
		}
		if tQuest == nil {
			tQuest = quest
		}
		break
	}
	for {
		var quest *model.Quest
		trQuestAction, quest, err = s.GetQuestActionByCategory("twitter_retweet")
		if err != nil {
			time.Sleep(time.Second * 5)
			continue
		}
		if tQuest == nil {
			tQuest = quest
		}
		break
	}
	for {
		var quest *model.Quest
		tqQuestAction, quest, err = s.GetQuestActionByCategory("twitter_quote")
		if err != nil {
			time.Sleep(time.Second * 5)
			continue
		}
		if tQuest == nil {
			tQuest = quest
		}
		break
	}
	for {
		var quest *model.Quest
		tcQuestAction, quest, err = s.GetQuestActionByCategory("twitter_create")
		if err != nil {
			time.Sleep(time.Second * 5)
			continue
		}
		if tQuest == nil {
			tQuest = quest
		}
		break
	}

	for {
		s.CheckTwitter()
		time.Sleep(time.Second * 10)
	}
}

func (s *Service) CheckTwitter() {
	var (
		allAccountExt map[int]*model.AccountExt
		err           error
	)
	allAccountExt, err = s.dao.FindAllAccountExt()
	if err != nil {
		log.Error("Twitter s.dao.FindAllAccountExt error: %v", err)
		return
	}
	for _, accountExt := range allAccountExt {
		var (
			userQuest        *model.UserQuest
			userQuestActions []*model.UserQuestAction
		)
		if len(accountExt.TwitterUserId) == 0 {
			continue
		}
		if _, ok := allTwitterUserQuestCompleted[accountExt.AccountId]; ok {
			continue
		}
		userQuest, err = s.dao.FindUserQuest(accountExt.AccountId, tQuest.Id)
		if err != nil {
			log.Error("Twitter s.dao.FindUserQuest error: %v", err)
			continue
		}
		if userQuest.Status == model.UserQuestCompletedStatus {
			continue
		}
		userQuestActions, err = s.dao.FindUserQuestActionByQuestId(accountExt.AccountId, tQuest.Id)
		if err != nil {
			log.Error("Twitter s.dao.FindUserQuestActionByQuestId error: %v", err)
			continue
		}
		s.CheckTwitterFollow(accountExt, userQuestActions)
		s.CheckTwitterFollow(accountExt, userQuestActions)
		s.CheckTwitterFollow(accountExt, userQuestActions)
		s.CheckTwitterFollow(accountExt, userQuestActions)
		s.CheckTwitterFollow(accountExt, userQuestActions)
	}
}

func (s *Service) CheckTwitterFollow(accountExt *model.AccountExt, userQuestActions []*model.UserQuestAction) (updateQuestAction *model.UserQuestAction) {
	if tfQuestAction == nil {
		return
	}
	var hasCompleted bool
	for _, userQuestAction := range userQuestActions {
		if userQuestAction.QuestActionId == tfQuestAction.Id {
			hasCompleted = true
			break
		}
	}
	if hasCompleted {
		return
	}
	return
}

func (s *Service) CheckTwitterLike(accountExt *model.AccountExt, userQuestActions []*model.UserQuestAction) (updateQuestAction *model.UserQuestAction) {
	if tlQuestAction == nil {
		return
	}
	var hasCompleted bool
	for _, userQuestAction := range userQuestActions {
		if userQuestAction.QuestActionId == tlQuestAction.Id {
			hasCompleted = true
			break
		}
	}
	if hasCompleted {
		return
	}
	return
}

func (s *Service) CheckTwitterRetweet(accountExt *model.AccountExt, userQuestActions []*model.UserQuestAction) (updateQuestAction *model.UserQuestAction) {
	if trQuestAction == nil {
		return
	}
	var hasCompleted bool
	for _, userQuestAction := range userQuestActions {
		if userQuestAction.QuestActionId == trQuestAction.Id {
			hasCompleted = true
			break
		}
	}
	if hasCompleted {
		return
	}
	return
}

func (s *Service) CheckTwitterQuote(accountExt *model.AccountExt, userQuestActions []*model.UserQuestAction) (updateQuestAction *model.UserQuestAction) {
	if tqQuestAction == nil {
		return
	}
	var hasCompleted bool
	for _, userQuestAction := range userQuestActions {
		if userQuestAction.QuestActionId == tqQuestAction.Id {
			hasCompleted = true
			break
		}
	}
	if hasCompleted {
		return
	}
	return
}

func (s *Service) CheckTwitterCreate(accountExt *model.AccountExt, userQuestActions []*model.UserQuestAction) (updateQuestAction *model.UserQuestAction) {
	if tcQuestAction == nil {
		return
	}
	var hasCompleted bool
	for _, userQuestAction := range userQuestActions {
		if userQuestAction.QuestActionId == tcQuestAction.Id {
			hasCompleted = true
			break
		}
	}
	if hasCompleted {
		return
	}
	return
}
