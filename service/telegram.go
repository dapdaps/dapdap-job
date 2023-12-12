package service

import (
	"dapdap-job/common/log"
	"dapdap-job/conf"
	"dapdap-job/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Service) InitTelegram() (err error) {
	if conf.Conf.Telegram == nil || len(conf.Conf.Telegram.BotToken) == 0 || conf.Conf.Telegram.ChatId == 0 {
		return
	}
	bot, err := tgbotapi.NewBotAPI(conf.Conf.Telegram.BotToken)
	if err != nil {
		return
	}
	bot.Debug = conf.Conf.Debug

	//log.Printf("Authorized on account %s \n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = int(conf.Conf.Timeout)

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		//log.Printf("[%s] %s \n", update.Message.From.UserName, update.Message.Text)
		if update.Message.Chat != nil && update.Message.Chat.ID == conf.Conf.Telegram.ChatId && len(update.Message.NewChatMembers) > 0 {
			for _, newChatMember := range update.Message.NewChatMembers {
				if newChatMember.IsBot {
					continue
				}
				s.ChannelJoin(&newChatMember)
			}
		}
		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID
		//
		//bot.Send(msg)
	}
	return
}

func (s *Service) ChannelJoin(tgUser *tgbotapi.User) {
	var (
		accountId       int
		completed       = 1
		reward          int
		questAction     *model.QuestAction
		quest           *model.Quest
		userQuestAction *model.UserQuestAction
		userQuest       *model.UserQuest
		err             error
	)
	accountId, err = s.dao.FindAccountIdByTg(tgUser.ID)
	if err != nil {
		log.Error("ChannelJoin s.dao.FindAccountIdByTg error: %v", err)
		return
	}
	if accountId <= 0 {
		return
	}
	questAction, err = s.dao.FindQuestActionByCategory("tg_join")
	if err != nil {
		log.Error("ChannelJoin s.dao.FindQuestActionByCategory error: %v", err)
		return
	}
	if questAction == nil {
		return
	}
	quest, err = s.dao.FindQuest(questAction.QuestId)
	if err != nil {
		log.Error("ChannelJoin s.dao.FindQuest error: %v", err)
		return
	}
	if quest == nil || quest.Status != model.QuestOnGoingStatus {
		return
	}
	userQuestAction, err = s.dao.FindUserQuestAction(accountId, questAction.Id)
	if err != nil {
		log.Error("ChannelJoin s.dao.FindUserQuestAction error: %v", err)
		return
	}
	if userQuestAction != nil {
		return
	}
	userQuest, err = s.dao.FindUserQuest(accountId, questAction.Id)
	if err != nil {
		log.Error("ChannelJoin s.dao.FindUserQuest error: %v", err)
		return
	}
	if userQuest.Status == model.UserQuestCompletedStatus {
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
		log.Error("ChannelJoin s.dao.UpdateUserQuest error: %v", err)
		return
	}
}
