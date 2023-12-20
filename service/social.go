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
	isFristStartQuest = true
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
	go func() {
		s.StartSocialQuest()
		isFristStartQuest = false
		time.Sleep(time.Minute * 3)
	}()
}

func (s *Service) StartSocialQuest() {
	var (
		accountExts map[int]*model.AccountExt
		updatedTime *time.Time
		err         error
	)
	accountExts, updatedTime, err = s.dao.FindAllAccountExt(maxUpdatedTime)
	if err != nil {
		log.Error("Social s.dao.FindAllAccountExt error: %v", err)
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
		if !accountExt.DiscordQuestCompleted && len(accountExt.DiscordUserId) > 0 {
			s.CheckDiscordQuest(accountExt)
		}
		if !accountExt.TelegramQuestCompleted && len(accountExt.TelegramUserId) > 0 {
			s.CheckTelegramQuest(accountExt)
		}
		//if !accountExt.TwitterQuestCompleted && len(accountExt.TwitterUserId) > 0 {
		//	s.CheckTwitterQuest(accountExt)
		//}
	}
}
