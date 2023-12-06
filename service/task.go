package service

import (
	"dapdap-job/common/log"
	"dapdap-job/conf"
	"time"
)

func (s *Service) StartTask() {
	go func() {
		for {
			log.Info("StatusTask time:%d", time.Now().Unix())
			s.StartStatusTask()
			time.Sleep(time.Second * 5)
		}
	}()

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
		err := s.InitQuest()
		if err != nil {
			panic(err)
		}
		for {
			log.Info("QuestTask maxQuestActionRecordId:%d time:%d", maxQuestActionRecordId, time.Now().Unix())
			err = s.StartQuestTask()
			if err != nil {
				time.Sleep(time.Second * 5)
			} else {
				time.Sleep(time.Second * time.Duration(conf.Conf.QuestInterval))
			}
		}
	}()

	//go func() {
	//	for {
	//		log.Info("RankTask time:%d", time.Now().Unix())
	//		s.StartRankTask()
	//		time.Sleep(time.Second * time.Duration(conf.Conf.RankInterval))
	//	}
	//}()
}
