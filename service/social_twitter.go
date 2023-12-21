package service

import (
	"context"
	"dapdap-job/common/log"
	"dapdap-job/conf"
	"dapdap-job/model"
	"fmt"
	"github.com/g8rswimmer/go-twitter/v2"
	"net/http"
	"strings"
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
		userQuest                   *model.UserQuest
		userQuestActions            []*model.UserQuestAction
		completed                   int
		followQuestActionCompleted  *model.UserQuestAction
		likeQuestActionCompleted    *model.UserQuestAction
		retweetQuestActionCompleted *model.UserQuestAction
		quoteQuestActionCompleted   *model.UserQuestAction
		createQuestActionCompleted  *model.UserQuestAction
		completedUserQuestActions   []*model.UserQuestAction
		err                         error
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
	if userQuest != nil {
		completed = userQuest.ActionCompleted
	}
	followQuestActionCompleted = s.CheckTwitterFollow(accountExt, userQuestActions)
	if followQuestActionCompleted != nil {
		completed++
		completedUserQuestActions = append(completedUserQuestActions, followQuestActionCompleted)
	}
	likeQuestActionCompleted = s.CheckTwitterLike(accountExt, userQuestActions)
	if likeQuestActionCompleted != nil {
		completed++
		completedUserQuestActions = append(completedUserQuestActions, likeQuestActionCompleted)
	}
	retweetQuestActionCompleted = s.CheckTwitterRetweet(accountExt, userQuestActions)
	if retweetQuestActionCompleted != nil {
		completed++
		completedUserQuestActions = append(completedUserQuestActions, retweetQuestActionCompleted)
	}
	quoteQuestActionCompleted, createQuestActionCompleted = s.CheckTwitterQuoteAndCreate(accountExt, userQuestActions)
	if quoteQuestActionCompleted != nil {
		completed++
		completedUserQuestActions = append(completedUserQuestActions, quoteQuestActionCompleted)
	}
	if createQuestActionCompleted != nil {
		completed++
		completedUserQuestActions = append(completedUserQuestActions, createQuestActionCompleted)
	}

	if completed >= tQuest.TotalAction || len(completedUserQuestActions) > 0 {
		if userQuest == nil {
			userQuest = &model.UserQuest{
				QuestId:         tQuest.Id,
				QuestCampaignId: tQuest.QuestCampaignId,
				AccountId:       accountExt.AccountId,
			}
		}
		userQuest.ActionCompleted = completed
		if userQuest.ActionCompleted >= tQuest.TotalAction {
			userQuest.Status = model.UserQuestCompletedStatus
		} else {
			userQuest.Status = model.UserQuestInProcessStatus
		}
		err = s.dao.UpdateUserQuest([]*model.UserQuest{userQuest}, completedUserQuestActions)
		if err != nil {
			log.Error("CheckTwitterQuest s.dao.UpdateUserQuest error: %v", err)
			return
		}
	}
}

// CheckTwitterFollow
// App rate limit (Application-only): 15 requests per 15-minute window shared among all users of your app
// User rate limit (User context): 15 requests per 15-minute window per each authenticated user
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
	client := getTwitterClient(accountExt.TwitterAccessToken)
	opts := twitter.UserFollowedListsOpts{
		UserFields: []twitter.UserField{twitter.UserFieldID},
		MaxResults: 100,
	}
	data, err := client.UserFollowedLists(context.Background(), accountExt.TwitterUserId, opts)
	if err != nil {
		log.Error("Twitter CheckTwitterFollow client.UserFollowedLists error: %v", err)
		return
	}
	var hasFollow = false
	if data.Raw != nil && len(data.Raw.Lists) > 0 {
		for _, user := range data.Raw.Lists {
			if user.ID == conf.Conf.Twitter.UserId {
				hasFollow = true
				break
			}
		}
	}
	if hasFollow {
		updateQuestAction = &model.UserQuestAction{
			QuestActionId:   tfQuestAction.Id,
			QuestId:         tfQuestAction.QuestId,
			QuestCampaignId: tfQuestAction.QuestCampaignId,
			AccountId:       accountExt.AccountId,
			Times:           1,
			Status:          model.UserQuestActionCompletedStatus,
		}
	}
	return
}

// CheckTwitterLike
// App rate limit (Application-only): 5 requests per 15-minute window shared among all users of your app
// User rate limit (User context): 5 requests per 15-minute window per each authenticated user
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
	client := getTwitterClient(accountExt.TwitterAccessToken)
	opts := twitter.UserLikesLookupOpts{
		TweetFields: []twitter.TweetField{twitter.TweetFieldID, twitter.TweetFieldAuthorID},
		MaxResults:  100,
	}
	data, err := client.UserLikesLookup(context.Background(), accountExt.TwitterUserId, opts)
	if err != nil {
		log.Error("Twitter CheckTwitterLike client.UserLikesLookup error: %v", err)
		return
	}
	var hasLike = false
	if data.Raw != nil && len(data.Raw.Tweets) > 0 {
		for _, tweet := range data.Raw.Tweets {
			if tweet.AuthorID == conf.Conf.Twitter.UserId {
				hasLike = true
				break
			}
		}
	}
	if hasLike {
		updateQuestAction = &model.UserQuestAction{
			QuestActionId:   tlQuestAction.Id,
			QuestId:         tlQuestAction.QuestId,
			QuestCampaignId: tlQuestAction.QuestCampaignId,
			AccountId:       accountExt.AccountId,
			Times:           1,
			Status:          model.UserQuestActionCompletedStatus,
		}
	}
	return
}

// CheckTwitterRetweet
// App rate limit (Application-only): 5 requests per 15-minute window shared among all users of your app
// User rate limit (User context): 5 requests per 15-minute window per each authenticated user
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
	client := getTwitterClient(accountExt.TwitterAccessToken)
	opts := twitter.UserRetweetLookupOpts{
		UserFields: []twitter.UserField{twitter.UserFieldID},
		MaxResults: 100,
	}
	data, err := client.UserRetweetLookup(context.Background(), conf.Conf.Twitter.ReTweetId, opts)
	if err != nil {
		log.Error("Twitter CheckTwitterRetweet client.UserRetweetLookup error: %v", err)
		return
	}
	var hasRetweet = false
	if data.Raw != nil && len(data.Raw.Users) > 0 {
		for _, user := range data.Raw.Users {
			if user.ID == accountExt.TwitterUserId {
				hasRetweet = true
				break
			}
		}
	}
	if hasRetweet {
		updateQuestAction = &model.UserQuestAction{
			QuestActionId:   trQuestAction.Id,
			QuestId:         trQuestAction.QuestId,
			QuestCampaignId: trQuestAction.QuestCampaignId,
			AccountId:       accountExt.AccountId,
			Times:           1,
			Status:          model.UserQuestActionCompletedStatus,
		}
	}
	return
}

func (s *Service) CheckTwitterQuoteAndCreate(accountExt *model.AccountExt, userQuestActions []*model.UserQuestAction) (updateQuoteQuestAction *model.UserQuestAction, updateCreateQuestAction *model.UserQuestAction) {
	var (
		hasCompletedQuote  bool
		hasCompletedCreate bool
	)
	if tqQuestAction == nil {
		hasCompletedQuote = true
	} else {
		for _, userQuestAction := range userQuestActions {
			if userQuestAction.QuestActionId == tqQuestAction.Id {
				hasCompletedQuote = true
				break
			}
		}
	}
	if tcQuestAction == nil {
		hasCompletedCreate = true
	} else {
		for _, userQuestAction := range userQuestActions {
			if userQuestAction.QuestActionId == tcQuestAction.Id {
				hasCompletedCreate = true
				break
			}
		}
	}
	if hasCompletedQuote && hasCompletedCreate {
		return
	}
	client := getTwitterClient(accountExt.TwitterAccessToken)
	opts := twitter.UserTweetTimelineOpts{
		TweetFields: []twitter.TweetField{twitter.TweetFieldEntities, twitter.TweetFieldReferencedTweets},
		MaxResults:  100,
	}
	data, err := client.UserTweetTimeline(context.Background(), accountExt.TwitterUserId, opts)
	if err != nil {
		log.Error("Twitter CheckTwitterQuoteAndCreate client.UserTweetTimeline error: %v", err)
		return
	}
	if !hasCompletedQuote {
		updateQuoteQuestAction = s.CheckTwitterQuote(accountExt, data)
	}
	if !hasCompletedCreate {
		updateCreateQuestAction = s.CheckTwitterCreate(accountExt, data)
	}
	return
}

// CheckTwitterQuote
// App rate limit (Application-only): 5 requests per 15-minute window shared among all users of your app
// User rate limit (User context): 5 requests per 15-minute window per each authenticated user
func (s *Service) CheckTwitterQuote(accountExt *model.AccountExt, data *twitter.UserTweetTimelineResponse) (updateQuestAction *model.UserQuestAction) {
	//if tqQuestAction == nil {
	//	return
	//}
	//var hasCompleted bool
	//for _, userQuestAction := range userQuestActions {
	//	if userQuestAction.QuestActionId == tqQuestAction.Id {
	//		hasCompleted = true
	//		break
	//	}
	//}
	//if hasCompleted {
	//	return
	//}
	//client := getTwitterClient(accountExt.TwitterAccessToken)
	//opts := twitter.QuoteTweetsLookupOpts{
	//	TweetFields: []twitter.TweetField{twitter.TweetFieldAuthorID},
	//	MaxResults:  100,
	//}
	//data, err := client.QuoteTweetsLookup(context.Background(), conf.Conf.Twitter.QuoteTweetId, opts)
	//if err != nil {
	//	log.Error("Twitter CheckTwitterQuote client.QuoteTweetsLookup error: %v", err)
	//	return
	//}
	//var hasQuote = false
	//if data.Raw != nil && len(data.Raw.Tweets) > 0 {
	//	for _, tweet := range data.Raw.Tweets {
	//		if tweet.AuthorID == accountExt.TwitterUserId {
	//			hasQuote = true
	//			break
	//		}
	//	}
	//}
	//opts := twitter.UserTweetTimelineOpts{
	//	TweetFields: []twitter.TweetField{twitter.TweetFieldText, twitter.TweetFieldEntities, twitter.TweetFieldReferencedTweets},
	//	MaxResults:  100,
	//}
	//data, err := client.UserTweetTimeline(context.Background(), accountExt.TwitterUserId, opts)
	//if err != nil {
	//	log.Error("Twitter CheckTwitterCreate client.UserTweetTimeline error: %v", err)
	//	return
	//}
	var hasQuote = false
	if data.Raw != nil && len(data.Raw.Tweets) > 0 {
		for _, tweet := range data.Raw.Tweets {
			var hasReferenced = false
			if len(tweet.ReferencedTweets) > 0 {
				for _, referencedTweet := range tweet.ReferencedTweets {
					if strings.EqualFold(referencedTweet.Type, "quoted") && referencedTweet.ID == conf.Conf.Twitter.QuoteTweetId {
						hasReferenced = true
						break
					}
				}
				if hasReferenced && tweet.Entities != nil && len(tweet.Entities.Mentions) >= 3 {
					hasQuote = true
					break
				}
			}
		}
	}
	if hasQuote {
		updateQuestAction = &model.UserQuestAction{
			QuestActionId:   tqQuestAction.Id,
			QuestId:         tqQuestAction.QuestId,
			QuestCampaignId: tqQuestAction.QuestCampaignId,
			AccountId:       accountExt.AccountId,
			Times:           1,
			Status:          model.UserQuestActionCompletedStatus,
		}
	}
	return
}

// CheckTwitterCreate
// App rate limit (Application-only): 5 requests per 15-minute window shared among all users of your app
// User rate limit (User context): 10 requests per 15-minute window per each authenticated user
func (s *Service) CheckTwitterCreate(accountExt *model.AccountExt, data *twitter.UserTweetTimelineResponse) (updateQuestAction *model.UserQuestAction) {
	//if tcQuestAction == nil {
	//	return
	//}
	//var hasCompleted bool
	//for _, userQuestAction := range userQuestActions {
	//	if userQuestAction.QuestActionId == tcQuestAction.Id {
	//		hasCompleted = true
	//		break
	//	}
	//}
	//if hasCompleted {
	//	return
	//}
	//client := getTwitterClient(accountExt.TwitterAccessToken)
	//opts := twitter.UserTweetTimelineOpts{
	//	TweetFields: []twitter.TweetField{twitter.TweetFieldText},
	//	MaxResults:  100,
	//}
	//data, err := client.UserTweetTimeline(context.Background(), accountExt.TwitterUserId, opts)
	//if err != nil {
	//	log.Error("Twitter CheckTwitterCreate client.UserTweetTimeline error: %v", err)
	//	return
	//}
	var hasCreate = false
	if data.Raw != nil && len(data.Raw.Tweets) > 0 {
		for _, tweet := range data.Raw.Tweets {
			if tweet.Entities != nil && len(tweet.Entities.Mentions) >= 0 {
				for _, entity := range tweet.Entities.Mentions {
					if entity.UserName == conf.Conf.Twitter.Username {
						hasCreate = true
						break
					}
				}
			}
		}
	}
	if hasCreate {
		updateQuestAction = &model.UserQuestAction{
			QuestActionId:   tcQuestAction.Id,
			QuestId:         tcQuestAction.QuestId,
			QuestCampaignId: tcQuestAction.QuestCampaignId,
			AccountId:       accountExt.AccountId,
			Times:           1,
			Status:          model.UserQuestActionCompletedStatus,
		}
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
