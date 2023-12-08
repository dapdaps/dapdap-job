package service

import (
	"dapdap-job/conf"
	"dapdap-job/dao"
	"time"
)

type Service struct {
	Timeout time.Duration
	dao     *dao.Dao
}

var (
	DapdapService *Service
)

func Init(c *conf.Config) (err error) {
	DapdapService = &Service{
		Timeout: time.Duration(c.Timeout),
		dao:     dao.New(c),
	}

	err = DapdapService.dao.InitQuestCampaignInfo()
	if err != nil {
		return
	}
	err = DapdapService.dao.InitQuestActionRecord()
	if err != nil {
		return
	}
	return
}
