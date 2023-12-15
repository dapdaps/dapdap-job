package service

import (
	"dapdap-job/common/log"
	"dapdap-job/conf"
	"dapdap-job/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

var (
	bot *tgbotapi.BotAPI
)

func (s *Service) StartTelegram() {
	var (
		err error
	)
	if conf.Conf.Telegram == nil || len(conf.Conf.Telegram.BotToken) == 0 || conf.Conf.Telegram.ChatId == 0 {
		return
	}

	for {
		bot, err = tgbotapi.NewBotAPI(conf.Conf.Telegram.BotToken)
		if err != nil {
			log.Error("Telegram tgbotapi.NewBotAPI error: %v", err)
			time.Sleep(time.Second * 5)
			continue
		}
		break
	}
	bot.Debug = conf.Conf.Debug

	for {
		err = s.RecoverTelegram()
		if err != nil {
			log.Error("Telegram RecoverTelegram error: %v", err)
			time.Sleep(time.Second * 5)
			continue
		}
		break
	}

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
				s.OnChannelJoin(&newChatMember)
			}
		}
	}
	return
}

func (s *Service) RecoverTelegram() (err error) {
	var (
		allAccountExt map[int]*model.AccountExt
		questAction   *model.QuestAction
		quest         *model.Quest
	)
	questAction, quest, err = s.GetQuestActionByCategory("tg_join")
	if err != nil || questAction == nil || quest == nil {
		return
	}
	allAccountExt, err = s.dao.FindAllAccountExt()
	if err != nil {
		return
	}
	for _, accountExt := range allAccountExt {
		var (
			userQuestAction *model.UserQuestAction
		)
		if accountExt.TelegramUserId <= 0 {
			continue
		}
		userQuestAction, err = s.dao.FindUserQuestAction(accountExt.AccountId, questAction.Id)
		if err != nil {
			log.Error("Telegram s.dao.FindUserQuestAction error: %v", err)
			continue
		}
		if userQuestAction != nil {
			continue
		}
		chatMember, e := bot.GetChatMember(tgbotapi.GetChatMemberConfig{
			ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
				ChatID: conf.Conf.Telegram.ChatId,
				UserID: accountExt.TelegramUserId,
			},
		})
		if e != nil {
			log.Error("Telegram bot.GetChatMember account:%d tgUserId:%d error: %v", accountExt.AccountId, accountExt.TelegramUserId, e)
			continue
		}
		if chatMember.IsMember {
			e = s.ChannelJoin(accountExt.AccountId, questAction, quest)
			if e != nil {
				continue
			}
		}
	}
	return
}

func (s *Service) OnChannelJoin(tgUser *tgbotapi.User) {
	var (
		accountId   int
		questAction *model.QuestAction
		quest       *model.Quest
		err         error
	)
	accountId, err = s.dao.FindAccountIdByTg(tgUser.ID)
	if err != nil {
		log.Error("Telegram s.dao.FindAccountIdByTg error: %v", err)
		return
	}
	if accountId <= 0 {
		return
	}
	questAction, quest, err = s.GetQuestActionByCategory("tg_join")
	if err != nil || questAction == nil || quest == nil {
		return
	}
	err = s.ChannelJoin(accountId, questAction, quest)
	return
}

func (s *Service) ChannelJoin(accountId int, questAction *model.QuestAction, quest *model.Quest) (err error) {
	var (
		userQuestAction *model.UserQuestAction
		userQuest       *model.UserQuest
		completed       = 1
		reward          int
	)
	userQuest, err = s.dao.FindUserQuest(accountId, questAction.Id)
	if err != nil {
		log.Error("Telegram s.dao.FindUserQuest error: %v", err)
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
		log.Error("Telegram s.dao.UpdateUserQuest error: %v", err)
		return
	}
	return
}

func (s *Service) GetQuestActionByCategory(category string) (questAction *model.QuestAction, quest *model.Quest, err error) {
	questAction, err = s.dao.FindQuestActionByCategory(category)
	if err != nil {
		log.Error("Telegram s.dao.FindQuestActionByCategory error: %v", err)
		return
	}
	if questAction == nil {
		return
	}
	quest, err = s.dao.FindQuest(questAction.QuestId)
	if err != nil {
		log.Error("Telegram s.dao.FindQuest error: %v", err)
		return
	}
	if quest == nil || quest.Status != model.QuestOnGoingStatus {
		quest = nil
		return
	}
	return
}
