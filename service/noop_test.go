package service

import (
	"dapdap-job/common/log"
	"dapdap-job/conf"
	"flag"
	"fmt"
	"sort"
	"testing"
)

func TestActionInit(t *testing.T) {
	var err error
	flag.Set("conf", "../testdata/config.toml")
	if err = conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log, conf.Conf.Debug)
	Init(conf.Conf)

	err = DapdapService.InitAction()
	if err != nil {
		panic(err)
	}

	err = DapdapService.StartActionTask()
	if err != nil {
		panic(err)
	}
}

func TestQuestInit(t *testing.T) {
	var err error
	flag.Set("conf", "../testdata/config.toml")
	if err = conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log, conf.Conf.Debug)
	Init(conf.Conf)

	err = DapdapService.InitQuest()
	if err != nil {
		panic(err)
	}

	err = DapdapService.StartQuestTask()
	if err != nil {
		panic(err)
	}
}

func TestRankInit(t *testing.T) {
	var err error
	flag.Set("conf", "../testdata/config.toml")
	if err = conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log, conf.Conf.Debug)
	Init(conf.Conf)

	questCampaigns, err := DapdapService.dao.FindAllQuestCampaign()
	if err != nil {
		log.Error("InitQuest s.dao.FindAllQuestCampaign error: %v", err)
		return
	}
	err = DapdapService.InitQuestCampaignReward(questCampaigns)
	if err != nil {
		log.Error("InitQuest InitQuestCampaignReward error: %v", err)
		return
	}
	questCampaignReward := questCampaignRewards[1]
	sort.Slice(questCampaignReward, func(i, j int) bool {
		if questCampaignReward[i].Reward > questCampaignReward[j].Reward {
			return true
		} else if questCampaignReward[i].Reward < questCampaignReward[j].Reward {
			return false
		} else {
			return questCampaignReward[i].AccountId > questCampaignReward[j].AccountId
		}
	})
	err = DapdapService.dao.UpdateRewardRank(1, questCampaignReward)
	if err != nil {
		panic(err)
	}
	fmt.Println(questCampaignReward)
}
