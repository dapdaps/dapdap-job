package service

import (
	"dapdap-job/common/log"
	"dapdap-job/conf"
	"dapdap-job/model"
	"github.com/bwmarrin/discordgo"
	"sync"
	"time"
)

var (
	dg                   *discordgo.Session
	discordQuestCategory = "discord_role"
	roleUsers            = sync.Map{} //map[discord user id]
	discordQuestAction   *model.QuestAction
	discordQuest         *model.Quest
)

func (s *Service) InitDiscord() {
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
		discordQuestAction, discordQuest, err = s.GetQuestActionByCategory(discordQuestCategory)
		if err != nil {
			log.Error("Discord GetQuestActionByCategory error: %v", err)
			time.Sleep(time.Second * 5)
			continue
		}
		break
	}

	//for {
	//	err = s.RecoverDiscord()
	//	if err != nil {
	//		log.Error("Discord RecoverDiscord error: %v", err)
	//		time.Sleep(time.Second * 5)
	//		continue
	//	}
	//	break
	//}
}

//func (s *Service) RecoverDiscord() (err error) {
//	var (
//		questAction *model.QuestAction
//		quest       *model.Quest
//	)
//	questAction, quest, err = s.GetQuestActionByCategory(discordQuestCategory)
//	if err != nil || questAction == nil || quest == nil {
//		return
//	}
//	for _, accountExt := range allAccountExt {
//		var (
//			userQuestAction *model.UserQuestAction
//			member          *discordgo.Member
//			e               error
//		)
//		if len(accountExt.DiscordUserId) == 0 {
//			continue
//		}
//		userQuestAction, err = s.dao.FindUserQuestAction(accountExt.AccountId, questAction.Id)
//		if err != nil {
//			log.Error("Discord s.dao.FindUserQuestAction error: %v", err)
//			continue
//		}
//		if userQuestAction != nil {
//			continue
//		}
//		member, e = dg.GuildMember(conf.Conf.Discord.GuildId, accountExt.DiscordUserId)
//		if e != nil {
//			log.Error("Discord dg.GuildMember error: %v", e)
//			continue
//		}
//		for _, roleID := range member.Roles {
//			var role *discordgo.Role
//			role, e = dg.State.Role(conf.Conf.Discord.GuildId, roleID)
//			if err != nil {
//				log.Error("Discord dg.State.Role error: %v", e)
//				continue
//			}
//
//			if role.Name == conf.Conf.Discord.Role {
//				_ = s.UpdateDiscordQuest(accountExt, questAction, quest)
//				break
//			}
//		}
//	}
//	return
//}

func (s *Service) CheckDiscordQuest(accountExt *model.AccountExt) {
	_, ok := roleUsers.Load(accountExt.DiscordUserId)
	if !ok {
		if !isFristStartQuest {
			return
		}
		if discordQuestAction == nil {
			return
		}
		var (
			userQuestAction *model.UserQuestAction
			member          *discordgo.Member
			hasRole         bool
			err             error
		)
		userQuestAction, err = s.dao.FindUserQuestAction(accountExt.AccountId, discordQuestAction.Id)
		if err != nil {
			log.Error("Discord s.dao.FindUserQuestAction error: %v", err)
			return
		}
		if userQuestAction != nil {
			return
		}
		member, err = dg.GuildMember(conf.Conf.Discord.GuildId, accountExt.DiscordUserId)
		if err != nil {
			log.Error("Discord dg.GuildMember error: %v", err)
			return
		}
		for _, roleID := range member.Roles {
			var role *discordgo.Role
			role, err = dg.State.Role(conf.Conf.Discord.GuildId, roleID)
			if err != nil {
				log.Error("Discord dg.State.Role error: %v", err)
				continue
			}
			if role.Name == conf.Conf.Discord.Role {
				hasRole = true
				roleUsers.Store(accountExt.DiscordUserId, true)
				break
			}
		}
		if !hasRole {
			return
		}
	}
	err := s.OnRoleUpdate(accountExt)
	if err == nil {
		roleUsers.Delete(accountExt.DiscordUserId)
	}
}

func (s *Service) OnMemberUpdate(ds *discordgo.Session, m *discordgo.GuildMemberUpdate) {
	if m.GuildID != conf.Conf.Discord.GuildId {
		return
	}
	for _, roleID := range m.Roles {
		var (
			role *discordgo.Role
			err  error
		)
		role, err = ds.State.Role(m.GuildID, roleID)
		if err != nil {
			log.Error("Discord ds.State.Role error: %v", err)
			continue
		}
		if role.Name != conf.Conf.Discord.Role {
			continue
		}
		log.Info("Discord 用户 %s:%s 获得了角色 %s", m.User.ID, m.User.Username, conf.Conf.Discord.Role)
		roleUsers.Store(m.User.ID, true)
		return
	}
}

func (s *Service) OnRoleUpdate(accountExt *model.AccountExt) (err error) {
	var (
		userQuestAction *model.UserQuestAction
		questAction     *model.QuestAction
		quest           *model.Quest
	)
	questAction, quest, err = s.GetQuestActionByCategory(discordQuestCategory)
	if err != nil || questAction == nil || quest == nil {
		return
	}
	userQuestAction, err = s.dao.FindUserQuestAction(accountExt.AccountId, questAction.Id)
	if err != nil {
		log.Error("Discord s.dao.FindUserQuestAction error: %v", err)
		return
	}
	if userQuestAction != nil {
		return
	}
	err = s.UpdateDiscordQuest(accountExt, questAction, quest)
	return
}

func (s *Service) UpdateDiscordQuest(accountExt *model.AccountExt, questAction *model.QuestAction, quest *model.Quest) (err error) {
	var (
		userQuestAction *model.UserQuestAction
		userQuest       *model.UserQuest
		completed       = 1
		reward          int
	)
	userQuest, err = s.dao.FindUserQuest(accountExt.AccountId, questAction.Id)
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
			AccountId:       accountExt.AccountId,
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
		AccountId:       accountExt.AccountId,
		Times:           1,
		Status:          model.UserQuestActionCompletedStatus,
	}
	err = s.dao.UpdateUserQuest(accountExt.AccountId, reward, []*model.UserQuest{userQuest}, []*model.UserQuestAction{userQuestAction})
	if err != nil {
		log.Error("Discord s.dao.UpdateUserQuest error: %v", err)
		return
	}
	return
}
