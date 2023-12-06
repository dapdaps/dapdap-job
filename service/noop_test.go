package service

import (
	"dapdap-job/common/log"
	"dapdap-job/conf"
	"flag"
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

	DapdapService.StartRankTask()
}

func TestSelectForUpdate(t *testing.T) {
	//var err error
	//flag.Set("conf", "../testdata/config.toml")
	//if err = conf.Init(); err != nil {
	//	panic(err)
	//}
	//log.Init(conf.Conf.Log, conf.Conf.Debug)
	//Init(conf.Conf)
	//timestamp := time.Now()
	//err = DapdapService.dao.WithTrx(func(db *sql.Tx) (err error) {
	//	var userId sql.NullInt64
	//	err = db.QueryRow(dal.FindAccountForUpdateSql, 1).Scan(&userId)
	//	if err != nil {
	//		return
	//	}
	//	reward, inviteReward, err := DapdapService.dao.FindUserReward(2)
	//	if err != nil {
	//		return
	//	}
	//	_, err = db.Exec(dal.UpdateUserInviteRewardByIdSql, 2, reward, inviteReward, timestamp)
	//	if err != nil {
	//		return
	//	}
	//	return
	//})
	//if err != nil {
	//	panic(err)
	//}
}
