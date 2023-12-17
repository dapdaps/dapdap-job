package service

import (
	"dapdap-job/common/log"
	"dapdap-job/model"
	"time"
)

var (
	allAccountExt  map[int]*model.AccountExt
	maxUpdatedTime *time.Time
)

func (s *Service) StartSocialQuest() {
	var err error

	for {
		allAccountExt, maxUpdatedTime, err = s.dao.FindAllAccountExt(maxUpdatedTime)
		if err != nil {
			log.Error("Social s.dao.FindAllAccountExt error: %v", err)
			time.Sleep(time.Second * 5)
			continue
		}
		log.Info("Social FindAllAccountExt maxUpdateTime: %s", maxUpdatedTime.Format(model.TimeFormat))
		break
	}

	go func() {
		s.StartTelegram()
	}()

	go func() {
		s.StartDiscord()
	}()

	//go func() {
	//	s.StartTwitter()
	//}()
}
