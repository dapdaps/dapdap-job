package service

import (
	"dapdap-job/common/log"
	"dapdap-job/model"
	"sort"
	"sync"
)

func (s *Service) StartRankTask() {
	s.updateRank()
}

func (s *Service) updateRank() {
	var (
		campaignTotalUsers map[int]map[int]bool
		totalUsers         int64
		totalExecutions    int64
		totalReward        int
		userRewards        []*model.UserReward
		wg                 sync.WaitGroup
		err                error
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		var e error
		campaignTotalUsers, totalUsers, e = s.dao.FindQuestCampaignTotalUsers()
		if e != nil {
			log.Error("RankTask s.dao.FindQuestTotalUsers error: %v", e)
			err = e
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var e error
		totalExecutions, e = s.dao.FindQuestTotalExecutions()
		if e != nil {
			log.Error("RankTask s.dao.FindQuestTotalExecutions error: %v", e)
			err = e
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var e error
		userRewards, totalReward, e = s.dao.FindAllUserReward()
		if err != nil {
			log.Error("RankTask s.dao.FindAllUserReward error: %v", e)
			err = e
			return
		}
	}()

	wg.Wait()
	if err != nil {
		return
	}

	err = s.dao.UpdateCampaignInfo(totalReward, totalUsers, totalExecutions)
	if err != nil {
		log.Error("RankTask s.dao.UpdateCampaignInfo error: %v", err)
		return
	}

	for campaignId, campaignTotalUser := range campaignTotalUsers {
		err = s.dao.UpdateCampaignTotalUsers(campaignId, len(campaignTotalUser))
		if err != nil {
			log.Error("RankTask s.dao.UpdateCampaignInfo error: %v", err)
			return
		}
	}

	sort.Slice(userRewards, func(i, j int) bool {
		if userRewards[i].ClaimedReward > userRewards[j].ClaimedReward {
			return true
		} else if userRewards[i].ClaimedReward < userRewards[j].ClaimedReward {
			return false
		} else {
			return userRewards[i].AccountId > userRewards[j].AccountId
		}
	})
	err = s.dao.UpdateUserRewardRank(userRewards)
	if err != nil {
		log.Error("RankTask s.dao.UpdateUserRewardRank error: %v", err)
		return
	}
	//for index, reward := range userRewards {
	//	reward.Rank = index + 1
	//}
	return
}
