package service

import (
	"dapdap-job/common/log"
	"dapdap-job/model"
	"encoding/json"
)

func (s *Service) UpdateInviteReward(invitedUserAddress map[string]string) (err error) {
	var (
		invitedUserIds map[string]int
		invites        map[int64][]*model.Invite
		quest          *model.QuestLong
		inviteRule     = &model.InviteQuestRule{}
	)
	invitedUserIds, _, err = s.dao.FindAccountIds(invitedUserAddress)
	if err != nil {
		log.Error("InviteReward s.dao.FindAccountIds error: %v", err)
		return
	}
	invites, err = s.dao.FindInvites(invitedUserIds)
	if err != nil {
		log.Error("InviteReward s.dao.FindInvites error: %v", err)
		return
	}
	quest, err = s.dao.FindLongQuest("invite")
	if err != nil {
		log.Error("InviteReward s.dao.FindLongQuest error: %v", err)
		return
	}
	err = json.Unmarshal([]byte(quest.Rule), inviteRule)
	if err != nil {
		log.Error("InviteReward json.Unmarshal error: %v", err)
		return
	}
	for accountId, accountInvites := range invites {
		var (
			totalInviteReward int64
		)
		totalInviteReward, err = s.dao.FindTotalInviteReward(accountId)
		if err != nil {
			log.Error("InviteReward s.dao.FindTotalInviteReward error: %v", err)
			return
		}
		for _, invite := range accountInvites {
			if totalInviteReward+inviteRule.Reward > inviteRule.MaxReward {
				invite.Reward = 0
			} else {
				invite.Reward = inviteRule.Reward
				totalInviteReward += inviteRule.Reward
			}
		}
		err = s.dao.UpdateInviteReward(accountInvites)
		if err != nil {
			log.Error("InviteReward s.dao.UpdateInvite error: %v", err)
			return
		}
	}
	return
}
