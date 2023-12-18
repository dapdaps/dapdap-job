package service

import (
	"dapdap-job/common/log"
	"dapdap-job/model"
	"fmt"
	"github.com/g8rswimmer/go-twitter/v2"
	"net/http"
	"time"
)

var (
	tfQuestAction *model.QuestAction
	tlQuestAction *model.QuestAction
	trQuestAction *model.QuestAction
	tqQuestAction *model.QuestAction
	tcQuestAction *model.QuestAction
	tQuest        *model.Quest
)

func (s *Service) InitTwitter() {
	var (
		quest *model.Quest
		err   error
	)
	for {
		tfQuestAction, quest, err = s.GetQuestActionByCategory("twitter_follow")
		if err != nil {
			time.Sleep(time.Second * 5)
			continue
		}
		if tQuest == nil {
			tQuest = quest
		}
		break
	}
	for {
		tlQuestAction, quest, err = s.GetQuestActionByCategory("twitter_like")
		if err != nil {
			time.Sleep(time.Second * 5)
			continue
		}
		if tQuest == nil {
			tQuest = quest
		}
		break
	}
	for {
		trQuestAction, quest, err = s.GetQuestActionByCategory("twitter_retweet")
		if err != nil {
			time.Sleep(time.Second * 5)
			continue
		}
		if tQuest == nil {
			tQuest = quest
		}
		break
	}
	for {
		tqQuestAction, quest, err = s.GetQuestActionByCategory("twitter_quote")
		if err != nil {
			time.Sleep(time.Second * 5)
			continue
		}
		if tQuest == nil {
			tQuest = quest
		}
		break
	}
	for {
		tcQuestAction, quest, err = s.GetQuestActionByCategory("twitter_create")
		if err != nil {
			time.Sleep(time.Second * 5)
			continue
		}
		if tQuest == nil {
			tQuest = quest
		}
		break
	}
}

func (s *Service) CheckTwitterQuest(accountExt *model.AccountExt) {
	var (
		userQuest        *model.UserQuest
		userQuestActions []*model.UserQuestAction
		err              error
	)
	userQuest, err = s.dao.FindUserQuest(accountExt.AccountId, tQuest.Id)
	if err != nil {
		log.Error("Twitter s.dao.FindUserQuest error: %v", err)
		return
	}
	if userQuest != nil && userQuest.Status == model.UserQuestCompletedStatus {
		accountExt.TwitterQuestCompleted = true
		return
	}
	userQuestActions, err = s.dao.FindUserQuestActionByQuestId(accountExt.AccountId, tQuest.Id)
	if err != nil {
		log.Error("Twitter s.dao.FindUserQuestActionByQuestId error: %v", err)
		return
	}
	s.CheckTwitterFollow(accountExt, userQuestActions)
	s.CheckTwitterFollow(accountExt, userQuestActions)
	s.CheckTwitterFollow(accountExt, userQuestActions)
	s.CheckTwitterFollow(accountExt, userQuestActions)
	s.CheckTwitterFollow(accountExt, userQuestActions)
}

func (s *Service) CheckTwitterFollow(accountExt *model.AccountExt, userQuestActions []*model.UserQuestAction) (updateQuestAction *model.UserQuestAction) {
	if tfQuestAction == nil {
		return
	}
	var hasCompleted bool
	for _, userQuestAction := range userQuestActions {
		if userQuestAction.QuestActionId == tfQuestAction.Id {
			hasCompleted = true
			break
		}
	}
	if hasCompleted {
		return
	}
	//client := getTwitterClient(accountExt.TwitterAccessToken)
	//opts := twitter.UserFollowedListsOpts{
	//	//Expansions:  []twitter.Expansion{twitter.ExpansionEntitiesMentionsUserName, twitter.ExpansionAuthorID},
	//	//TweetFields: []twitter.TweetField{twitter.TweetFieldCreatedAt, twitter.TweetFieldConversationID, twitter.TweetFieldAttachments},
	//}
	//data, err := client.UserFollowedLists(context.Background(), accountExt.TwitterUserId, opts)
	//if err != nil {
	//	log.Error("Twitter CheckTwitterFollow client.UserFollowedLists error: %v", err)
	//	return
	//}
	//log.Info("data: %v", data)
	return
}

func (s *Service) CheckTwitterLike(accountExt *model.AccountExt, userQuestActions []*model.UserQuestAction) (updateQuestAction *model.UserQuestAction) {
	if tlQuestAction == nil {
		return
	}
	var hasCompleted bool
	for _, userQuestAction := range userQuestActions {
		if userQuestAction.QuestActionId == tlQuestAction.Id {
			hasCompleted = true
			break
		}
	}
	if hasCompleted {
		return
	}
	return
}

func (s *Service) CheckTwitterRetweet(accountExt *model.AccountExt, userQuestActions []*model.UserQuestAction) (updateQuestAction *model.UserQuestAction) {
	if trQuestAction == nil {
		return
	}
	var hasCompleted bool
	for _, userQuestAction := range userQuestActions {
		if userQuestAction.QuestActionId == trQuestAction.Id {
			hasCompleted = true
			break
		}
	}
	if hasCompleted {
		return
	}
	return
}

func (s *Service) CheckTwitterQuote(accountExt *model.AccountExt, userQuestActions []*model.UserQuestAction) (updateQuestAction *model.UserQuestAction) {
	if tqQuestAction == nil {
		return
	}
	var hasCompleted bool
	for _, userQuestAction := range userQuestActions {
		if userQuestAction.QuestActionId == tqQuestAction.Id {
			hasCompleted = true
			break
		}
	}
	if hasCompleted {
		return
	}
	return
}

func (s *Service) CheckTwitterCreate(accountExt *model.AccountExt, userQuestActions []*model.UserQuestAction) (updateQuestAction *model.UserQuestAction) {
	if tcQuestAction == nil {
		return
	}
	var hasCompleted bool
	for _, userQuestAction := range userQuestActions {
		if userQuestAction.QuestActionId == tcQuestAction.Id {
			hasCompleted = true
			break
		}
	}
	if hasCompleted {
		return
	}
	return
}

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

func getTwitterClient(token string) (client *twitter.Client) {
	client = &twitter.Client{
		Authorizer: authorize{
			Token: token,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}
	return
}
