package service

import (
	"dapdap-job/common/log"
	"dapdap-job/conf"
	"dapdap-job/model"
)

func (s *Service) UpdateInviteReward(invitedUserAddress map[string]string) (err error) {
	var (
		invitedUserIds map[string]int
		invites        map[string]*model.Invite
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
	for _, invite := range invites {
		if invite.InvitedReward+conf.Conf.InviteReward > conf.Conf.MaxInviteReward {
			continue
		}
		err = s.dao.UpdateInvite(invite)
		if err != nil {
			log.Error("InviteReward s.dao.UpdateInvite error: %v", err)
			return
		}
	}
	return
}
