package service

import (
	"dapdap-job/common/log"
	"sort"
)

var (
// rwLock sync.RWMutex
// allUserRewards = map[int]*model.UserReward{}
)

//func (s *Service) UpdateUserReward(accountId int, reward int) {
//	if accountId <= 0 || reward <= 0 {
//		return
//	}
//	rwLock.Lock()
//	defer rwLock.Unlock()
//	userReward := allUserRewards[accountId]
//	if userReward == nil {
//		userReward = &model.UserReward{}
//		allUserRewards[accountId] = userReward
//	}
//	userReward.Reward += reward
//}

func (s *Service) StartRankTask() {
	s.updateRank()
}

func (s *Service) updateRank() {
	//rwLock.RLock()
	//var data []*model.UserReward
	//for _, userReward := range allUserRewards {
	//	data = append(data, userReward)
	//}
	//rwLock.RUnlock()
	//sort.Slice(data, func(i, j int) bool {
	//	if data[i].Reward > data[j].Reward {
	//		return true
	//	} else if data[i].Reward < data[j].Reward {
	//		return false
	//	} else {
	//		return data[i].AccountId > data[j].AccountId
	//	}
	//})
	userRewards, totalReward, err := s.dao.FindAllUserReward()
	if err != nil {
		log.Error("RankTask s.dao.FindAllUserReward error: %v", err)
		return
	}
	err = s.dao.UpdateTotalReward(totalReward)
	if err != nil {
		log.Error("RankTask s.dao.UpdateTotalReward error: %v", err)
		return
	}
	sort.Slice(userRewards, func(i, j int) bool {
		if userRewards[i].Reward > userRewards[j].Reward {
			return true
		} else if userRewards[i].Reward < userRewards[j].Reward {
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
	return
}
