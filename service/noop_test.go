package service

import (
	"dapdap-job/common/log"
	"dapdap-job/conf"
	"flag"
	"testing"
)

func TestActionInit(t *testing.T) {
	//var err error
	//flag.Set("conf", "../testdata/config.toml")
	//if err = conf.Init(); err != nil {
	//	panic(err)
	//}
	//log.Init(conf.Conf.Log, conf.Conf.Debug)
	//Init(conf.Conf)
	//
	//err = DapdapService.InitAction()
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = DapdapService.StartActionTask()
	//if err != nil {
	//	panic(err)
	//}
}

func TestQuestInit(t *testing.T) {
	//var err error
	//flag.Set("conf", "../testdata/config.toml")
	//if err = conf.Init(); err != nil {
	//	panic(err)
	//}
	//log.Init(conf.Conf.Log, conf.Conf.Debug)
	//Init(conf.Conf)
	//
	//err = DapdapService.InitQuest()
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = DapdapService.StartQuestTask()
	//if err != nil {
	//	panic(err)
	//}
}

func TestRank(t *testing.T) {
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

func TestTelegram(t *testing.T) {
	//var err error
	//flag.Set("conf", "../testdata/config.toml")
	//if err = conf.Init(); err != nil {
	//	panic(err)
	//}
	//log.Init(conf.Conf.Log, conf.Conf.Debug)
	//Init(conf.Conf)
	//
	//DapdapService.StartTelegram()
}

func TestDiscord(t *testing.T) {
	//var err error
	//flag.Set("conf", "../testdata/config.toml")
	//if err = conf.Init(); err != nil {
	//	panic(err)
	//}
	//log.Init(conf.Conf.Log, conf.Conf.Debug)
	//err = Init(conf.Conf)
	//if err != nil {
	//	panic(err)
	//}
	//DapdapService.StartDiscord()
}

func TestTwitter(t *testing.T) {
	//var err error
	//flag.Set("conf", "../testdata/config.toml")
	//if err = conf.Init(); err != nil {
	//	panic(err)
	//}
	//log.Init(conf.Conf.Log, conf.Conf.Debug)
	//err = Init(conf.Conf)
	//if err != nil {
	//	panic(err)
	//}
	//DapdapService.InitTwitter()
	//DapdapService.CheckTwitterQuest(&model.AccountExt{
	//	AccountId:             51,
	//	TwitterUserId:         "816926408",
	//	TwitterAccessToken:    "R29fUTZRRFhZYllHeGhXNURaVlIwRWF5a3Q0WkdBeUs5N3otdENtaWs4Rk4zOjE3MDMwMzg2MDEyMTc6MToxOmF0OjE",
	//	TwitterQuestCompleted: false,
	//})

	//client := getTwitterClient("ankzWXpxRmJYMkViZUo4VmFaZjcyZWx3QnJ5ZXU0U240WjE4MFdBWGNWZU42OjE3MDMwNDU4OTIzMjA6MTowOmF0OjE")
	//opts := twitter.UserLikesLookupOpts{
	//	TweetFields: []twitter.TweetField{twitter.TweetFieldID, twitter.TweetFieldAuthorID},
	//	MaxResults:  100,
	//}
	//data, err := client.UserLikesLookup(context.Background(), "816926408", opts)
	//if err != nil {
	//	t.Fatalf("Twitter client.UserLikesLookup error: %v", err)
	//	return
	//}
	//fmt.Println(data.RateLimit.Limit)
	//fmt.Println(data.RateLimit.Remaining)
	//fmt.Println(data.RateLimit.Reset.Time().Unix())
}
