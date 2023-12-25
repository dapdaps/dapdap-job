package service

import (
	"dapdap-job/common/log"
	"dapdap-job/conf"
	"dapdap-job/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"sync"
	"time"
)

var (
	bot                   *tgbotapi.BotAPI
	telegramQuestCategory = "telegram_join"
	telegramQuestAction   *model.QuestAction
	joinUsers             = sync.Map{} //map[tg user id]
)

func (s *Service) InitTelegram() {
	var (
		err error
	)
	if conf.Conf.Telegram == nil || len(conf.Conf.Telegram.BotToken) == 0 || conf.Conf.Telegram.ChatId == 0 {
		return
	}

	for {
		telegramQuestAction, _, err = s.GetQuestActionByCategory(telegramQuestCategory)
		if err != nil {
			log.Error("Telegram GetQuestActionByCategory error: %v", err)
			time.Sleep(time.Second * 5)
			continue
		}
		break
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

	//for {
	//	err = s.RecoverTelegram()
	//	if err != nil {
	//		log.Error("Telegram RecoverTelegram error: %v", err)
	//		time.Sleep(time.Second * 5)
	//		continue
	//	}
	//	break
	//}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = int(conf.Conf.Timeout)

	updates := bot.GetUpdatesChan(u)

	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}
			if update.Message.Chat != nil && update.Message.Chat.ID == conf.Conf.Telegram.ChatId && len(update.Message.NewChatMembers) > 0 {
				for _, newChatMember := range update.Message.NewChatMembers {
					if newChatMember.IsBot {
						continue
					}
					joinUsers.Store(strconv.Itoa(int(newChatMember.ID)), true)
				}
			}
		}
	}()
}

//func (s *Service) RecoverTelegram() (err error) {
//	var (
//		questAction *model.QuestAction
//		quest       *model.Quest
//	)
//	questAction, quest, err = s.GetQuestActionByCategory(telegramQuestCategory)
//	if err != nil || questAction == nil || quest == nil {
//		return
//	}
//	for _, accountExt := range allAccountExt {
//		var (
//			userQuestAction *model.UserQuestAction
//			tgUserId        int
//		)
//		if len(accountExt.TelegramUserId) <= 0 {
//			continue
//		}
//		userQuestAction, err = s.dao.FindUserQuestAction(accountExt.AccountId, questAction.Id)
//		if err != nil {
//			log.Error("Telegram s.dao.FindUserQuestAction error: %v", err)
//			continue
//		}
//		if userQuestAction != nil {
//			continue
//		}
//		tgUserId, _ = strconv.Atoi(accountExt.TelegramUserId)
//		chatMember, e := bot.GetChatMember(tgbotapi.GetChatMemberConfig{
//			ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
//				ChatID: conf.Conf.Telegram.ChatId,
//				UserID: int64(tgUserId),
//			},
//		})
//		if e != nil {
//			log.Error("Telegram bot.GetChatMember account:%d tgUserId:%d error: %v", accountExt.AccountId, accountExt.TelegramUserId, e)
//			continue
//		}
//		if chatMember.IsMember {
//			e = s.UpdateTelegramQuest(accountExt, questAction, quest)
//			if e != nil {
//				continue
//			}
//		}
//	}
//	return
//}

func (s *Service) CheckTelegramQuest(accountExt *model.AccountExt) {
	_, ok := joinUsers.Load(accountExt.TelegramUserId)
	if !ok {
		if !isFirstStartQuest {
			return
		}
		if telegramQuestAction == nil {
			return
		}
		var (
			userQuestAction *model.UserQuestAction
			tgUserId        int
			err             error
		)
		userQuestAction, err = s.dao.FindUserQuestAction(accountExt.AccountId, telegramQuestAction.Id)
		if err != nil {
			log.Error("Telegram s.dao.FindUserQuestAction error: %v", err)
			return
		}
		if userQuestAction != nil {
			return
		}
		tgUserId, err = strconv.Atoi(accountExt.TelegramUserId)
		if err != nil {
			log.Error("Telegram strconv.Atoi error: %v", err)
			return
		}
		chatMember, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{
			ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
				ChatID: conf.Conf.Telegram.ChatId,
				UserID: int64(tgUserId),
			},
		})
		if err != nil {
			log.Error("Telegram bot.GetChatMember account:%d tgUserId:%d error: %v", accountExt.AccountId, accountExt.TelegramUserId, err)
			return
		}
		if len(chatMember.Status) == 0 || chatMember.Status == "left" || chatMember.Status == "kicked" {
			return
		}
	}
	err := s.OnChannelJoin(accountExt)
	if err == nil {
		roleUsers.Delete(accountExt.TelegramUserId)
	}
}

func (s *Service) OnChannelJoin(accountExt *model.AccountExt) (err error) {
	var (
		questAction *model.QuestAction
		quest       *model.Quest
	)
	questAction, quest, err = s.GetQuestActionByCategory(telegramQuestCategory)
	if err != nil || questAction == nil || quest == nil {
		return
	}
	err = s.UpdateTelegramQuest(accountExt, questAction, quest)
	return
}

func (s *Service) UpdateTelegramQuest(accountExt *model.AccountExt, questAction *model.QuestAction, quest *model.Quest) (err error) {
	var (
		userQuestAction *model.UserQuestAction
		userQuest       *model.UserQuest
		completed       = 1
	)
	userQuest, err = s.dao.FindUserQuest(accountExt.AccountId, questAction.Id)
	if err != nil {
		log.Error("Telegram s.dao.FindUserQuest error: %v", err)
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
	err = s.dao.UpdateUserQuest([]*model.UserQuest{userQuest}, []*model.UserQuestAction{userQuestAction})
	if err != nil {
		log.Error("Telegram s.dao.UpdateUserQuest error: %v", err)
		return
	}
	return
}
