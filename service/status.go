package service

import (
	"dapdap-job/common/log"
	"dapdap-job/model"
	"time"
)

func (s *Service) StartStatusTask() {
	go func() {
		questCampaigns, err := s.dao.FindUpdateStatusCampaign()
		if err != nil {
			log.Error("StatusTask s.dao.FindUpdateStatusCampaign error: %v", err)
			return
		}
		if len(questCampaigns) == 0 {
			return
		}

		for _, questCampaign := range questCampaigns {
			var (
				now    = time.Now().UnixNano() / 1e6
				status string
			)
			if questCampaign.EndTime <= now {
				status = model.QuestCampaignEndedStatus
			} else if questCampaign.StartTime <= now {
				status = model.QuestCampaignOnGoingStatus
			} else {
				status = model.QuestCampaignUnStartStatus
			}
			if status != questCampaign.Status {
				err = s.dao.UpdateQuestCampaignStatus(questCampaign.Id, status)
				if err != nil {
					log.Error("StatusTask s.dao.UpdateQuestCampaignStatus error: %v", err)
					return
				}
			}
		}
	}()

	go func() {
		quests, err := s.dao.FindUpdateStatusQuest()
		if err != nil {
			log.Error("StatusTask s.dao.FindUpdateStatusQuest error: %v", err)
			return
		}
		if len(quests) == 0 {
			return
		}

		for _, quest := range quests {
			var (
				now    = time.Now().UnixNano() / 1e6
				status string
			)
			if quest.EndTime <= now {
				status = model.QuestEndedStatus
			} else if quest.StartTime <= now {
				status = model.QuestOnGoingStatus
			} else {
				status = model.QuestUnStartStatus
			}
			if status != quest.Status {
				if status == model.QuestEndedStatus {
					err = s.dao.UpdateUserQuestStatus(quest.Id, model.UserQuestExpiredStatus)
					if err != nil {
						log.Error("StatusTask s.dao.UpdateUserQuestStatus error: %v", err)
						return
					}
				}
				err = s.dao.UpdateQuestStatus(quest.Id, status)
				if err != nil {
					log.Error("StatusTask s.dao.UpdateQuestStatus error: %v", err)
					return
				}
			}
		}
	}()
}
