package service

import (
	"dapdap-job/common/log"
	"dapdap-job/conf"
	"time"
)

func (s *Service) StartTask() {
	go func() {
		err := s.InitAction()
		if err != nil {
			panic(err)
		}
		for {
			log.Info("ActionTask maxDappRecordId:%d maxChainRecordId:%d time:%d", maxDappId, maxChainId, time.Now().Unix())
			_ = s.StartActionTask()
			time.Sleep(time.Second * 5)
		}
	}()

	go func() {
		questInterval := int64(60)
		if conf.Conf.QuestInterval > 0 {
			questInterval = conf.Conf.QuestInterval
		}
		err := s.InitQuest()
		if err != nil {
			panic(err)
		}
		for {
			log.Info("QuestTask maxQuestActionRecordId:%d time:%d", maxQuestActionRecordId, time.Now().Unix())
			s.StartQuestTask()
			time.Sleep(time.Second * time.Duration(questInterval))
		}
	}()
}
