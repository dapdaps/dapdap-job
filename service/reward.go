package service

import (
	"dapdap-job/common/log"
	"dapdap-job/model"
	"sort"
	"sync"
)

var (
	rwLock         sync.RWMutex
	allUserRewards = map[int]*model.UserReward{}
)

func (s *Service) UpdateUserReward(accountId int, reward int) {
	if accountId <= 0 || reward <= 0 {
		return
	}
	rwLock.Lock()
	defer rwLock.Unlock()
	userReward := allUserRewards[accountId]
	if userReward == nil {
		userReward = &model.UserReward{}
		allUserRewards[accountId] = userReward
	}
	userReward.Reward += reward
}

func (s *Service) StartRankTask() {
	s.updateRank()
}

func (s *Service) updateRank() {
	rwLock.RLock()
	var data []*model.UserReward
	for _, userReward := range allUserRewards {
		data = append(data, userReward)
	}
	rwLock.RUnlock()
	sort.Slice(data, func(i, j int) bool {
		if data[i].Reward > data[j].Reward {
			return true
		} else if data[i].Reward < data[j].Reward {
			return false
		} else {
			return data[i].AccountId > data[j].AccountId
		}
	})
	err := s.dao.UpdateUserRewardRank(data)
	if err != nil {
		log.Error("RankTask updateRank s.dao.UpdateUserRewardRank error: %v", err)
		return
	}
	for index, userReward := range data {
		userReward.Rank = index + 1
	}
	return
}
