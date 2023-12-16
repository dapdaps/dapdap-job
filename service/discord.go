package service

import (
	"dapdap-job/common/log"
	"dapdap-job/conf"
	"dapdap-job/model"
	"github.com/bwmarrin/discordgo"
	"time"
)

var (
	dg                   *discordgo.Session
	discordQuestCategory = "discord_role"
)

func (s *Service) StartDiscord() {
	var (
		err error
	)
	if conf.Conf.Discord == nil || len(conf.Conf.Discord.BotToken) == 0 || len(conf.Conf.Discord.GuildId) == 0 || len(conf.Conf.Discord.Role) == 0 {
		return
	}

	for {
		dg, err = discordgo.New("Bot " + conf.Conf.Discord.BotToken)
		if err != nil {
			log.Error("Discord discordgo.New error: %v", err)
			time.Sleep(time.Second * 5)
			continue
		}
		break
	}

	dg.AddHandler(s.OnMemberUpdate)
	err = dg.Open()
	if err != nil {
		log.Error("Discord dg.Open error: %v", err)
		return
	}

	for {
		err = s.RecoverDiscord()
		if err != nil {
			log.Error("Discord RecoverDiscord error: %v", err)
			time.Sleep(time.Second * 5)
			continue
		}
		break
	}
	return
}

func (s *Service) RecoverDiscord() (err error) {
	var (
		allAccountExt map[int]*model.AccountExt
		questAction   *model.QuestAction
		quest         *model.Quest
		e             error
	)
	questAction, quest, err = s.GetQuestActionByCategory(discordQuestCategory)
	if err != nil || questAction == nil || quest == nil {
		return
	}
	allAccountExt, err = s.dao.FindAllAccountExt()
	if err != nil {
		log.Error("Discord s.dao.FindAllAccountExt error: %v", err)
		return
	}
	for _, accountExt := range allAccountExt {
		var (
			userQuestAction *model.UserQuestAction
			member          *discordgo.Member
		)
		if len(accountExt.DiscordUserId) == 0 {
			continue
		}
		userQuestAction, err = s.dao.FindUserQuestAction(accountExt.AccountId, questAction.Id)
		if err != nil {
			log.Error("Discord s.dao.FindUserQuestAction error: %v", err)
			continue
		}
		if userQuestAction != nil {
			continue
		}

		member, e = dg.GuildMember(conf.Conf.Discord.GuildId, accountExt.DiscordUserId)
		if e != nil {
			log.Error("Discord dg.GuildMember error: %v", e)
			continue
		}

		for _, roleID := range member.Roles {
			var role *discordgo.Role
			role, e = dg.State.Role(conf.Conf.Discord.GuildId, roleID)
			if err != nil {
				log.Error("Discord dg.State.Role error: %v", e)
				continue
			}

			if role.Name == conf.Conf.Discord.Role {
				log.Info("Discord 用户 %s:%s 获得了角色 %s", member.User.ID, member.User.Username, conf.Conf.Discord.Role)
				s.OnRoleUpdate(member.User.ID)
				break
			}
		}
	}
	return
}

func (s *Service) OnMemberUpdate(ds *discordgo.Session, m *discordgo.GuildMemberUpdate) {
	if m.GuildID != conf.Conf.Discord.GuildId {
		return
	}
	for _, roleID := range m.Roles {
		role, err := ds.State.Role(m.GuildID, roleID)
		if err != nil {
			log.Error("Discord ds.State.Role error: %v", err)
			continue
		}

		if role.Name == conf.Conf.Discord.Role {
			log.Info("Discord 用户 %s:%s 获得了角色 %s", m.User.ID, m.User.Username, conf.Conf.Discord.Role)
			s.OnRoleUpdate(m.User.ID)
			return
		}
	}
}

func (s *Service) OnRoleUpdate(discordUserId string) {
	var (
		accountId   int
		questAction *model.QuestAction
		quest       *model.Quest
		err         error
	)
	accountId, err = s.dao.FindAccountIdByDiscord(discordUserId)
	if err != nil {
		log.Error("Discord s.dao.FindAccountIdByDiscord error: %v", err)
		return
	}
	if accountId <= 0 {
		return
	}
	questAction, quest, err = s.GetQuestActionByCategory(discordQuestCategory)
	if err != nil || questAction == nil || quest == nil {
		return
	}
	err = s.UpdateDiscordQuest(accountId, questAction, quest)
	return
}

func (s *Service) UpdateDiscordQuest(accountId int, questAction *model.QuestAction, quest *model.Quest) (err error) {
	var (
		userQuestAction *model.UserQuestAction
		userQuest       *model.UserQuest
		completed       = 1
		reward          int
	)
	userQuest, err = s.dao.FindUserQuest(accountId, questAction.Id)
	if err != nil {
		log.Error("Discord s.dao.FindUserQuest error: %v", err)
		return
	}
	if userQuest != nil && userQuest.Status == model.UserQuestCompletedStatus {
		return
	}
	if userQuest != nil {
		completed += userQuest.ActionCompleted
	} else {
		userQuest = &model.UserQuest{
			QuestId:         quest.Id,
			QuestCampaignId: quest.QuestCampaignId,
			AccountId:       accountId,
		}
	}
	userQuest.ActionCompleted = completed
	if completed >= quest.TotalAction {
		reward = quest.Reward
		userQuest.Status = model.UserQuestCompletedStatus
	} else {
		userQuest.Status = model.UserQuestInProcessStatus
	}
	userQuestAction = &model.UserQuestAction{
		QuestActionId:   questAction.Id,
		QuestId:         quest.Id,
		QuestCampaignId: quest.QuestCampaignId,
		AccountId:       accountId,
		Times:           1,
		Status:          model.UserQuestActionCompletedStatus,
	}
	err = s.dao.UpdateUserQuest(accountId, reward, []*model.UserQuest{userQuest}, []*model.UserQuestAction{userQuestAction})
	if err != nil {
		log.Error("Discord s.dao.UpdateUserQuest error: %v", err)
		return
	}
	return
}
