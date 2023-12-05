package service

import (
	"dapdap-job/common/log"
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
			_ = s.ActionTask()
			time.Sleep(time.Second * 5)
		}
	}()

	//go func() {
	//	err := s.InitQuest()
	//	if err != nil {
	//		panic(err)
	//	}
	//	for {
	//		log.Info("QuestTask maxQuestActionRecordId:%d time:%d", maxQuestActionRecordId, time.Now().Unix())
	//		_ = s.QuestTask()
	//		time.Sleep(time.Minute * 1)
	//	}
	//}()
}
