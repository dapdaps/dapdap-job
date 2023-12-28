package service

import (
	"dapdap-job/common/log"
	"dapdap-job/model"
	"sync"
	"time"
)

var (
	allAccountExt     = map[int]*model.AccountExt{}
	maxUpdatedTime    *time.Time
	isFirstStartQuest = true
)

func (s *Service) InitSocialQuest() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.InitTelegram()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.InitDiscord()
	}()

	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	s.InitTwitter()
	//}()
	wg.Wait()
}

func (s *Service) StartSocialQuest() {
	var (
		accountExts    map[int]*model.AccountExt
		updatedTime    *time.Time
		totalDiscords  int
		totalTelegrams int
		err            error
	)
	defer func() {
		isFirstStartQuest = false
	}()
	roleUsers.Range(func(key, value any) bool {
		totalDiscords++
		return true
	})
	joinUsers.Range(func(key, value any) bool {
		totalTelegrams++
		return true
	})
	//if !isFirstStartQuest && totalDiscords == 0 && totalTelegrams == 0 {
	//	return
	//}
	accountExts, updatedTime, err = s.dao.FindAllAccountExt(maxUpdatedTime)
	if err != nil {
		log.Error("Social s.dao.FindAllAccountExt error: %v", err)
		return
	}
	if !isFirstStartQuest && len(accountExts) == 0 {
		return
	}
	if len(accountExts) > 0 {
		for _, accountExt := range accountExts {
			allAccountExt[accountExt.AccountId] = accountExt
		}
		log.Info("Social FindAllAccountExt maxUpdateTime: %s", updatedTime.Format(model.TimeFormat))
		maxUpdatedTime = updatedTime
	}
	for _, accountExt := range allAccountExt {
		var forceCheck = false
		for _, accountExtTemp := range accountExts {
			if accountExtTemp.AccountId == accountExt.AccountId {
				forceCheck = true
				break
			}
		}
		if !accountExt.DiscordQuestCompleted && len(accountExt.DiscordUserId) > 0 {
			s.CheckDiscordQuest(accountExt, isFirstStartQuest || forceCheck)
		}
		if !accountExt.TelegramQuestCompleted && len(accountExt.TelegramUserId) > 0 {
			s.CheckTelegramQuest(accountExt, isFirstStartQuest || forceCheck)
		}
		//if !accountExt.TwitterQuestCompleted && len(accountExt.TwitterUserId) > 0 {
		//	s.CheckTwitterQuest(accountExt)
		//}
	}
}
