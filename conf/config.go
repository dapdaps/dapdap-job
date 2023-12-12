package conf

import (
	"dapdap-job/common/conf"
	"dapdap-job/common/log"
	"flag"
	"github.com/BurntSushi/toml"
)

var (
	confPath string
	Conf     = &Config{}
)

type Config struct {
	Debug           bool
	Timeout         int64
	QuestInterval   int64
	MaxInviteReward int64
	InviteReward    int64
	RankInterval    int64
	Log             *log.Config
	Pgsql           *Pgsql
	Telegram        *Telegram
}

type Pgsql struct {
	DB *conf.Pgsql
}

type Telegram struct {
	ChatId   int64
	BotToken string
}

func init() {
	flag.StringVar(&confPath, "conf", "", "config path")
}

func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	if err != nil {
		log.Error("error decoding [%v]:%v", confPath, err)
		return
	}
	if Conf.MaxInviteReward <= 0 {
		Conf.MaxInviteReward = 10000
	}
	if Conf.InviteReward <= 0 {
		Conf.InviteReward = 5
	}
	if Conf.QuestInterval <= 0 {
		Conf.QuestInterval = 5
	}
	if Conf.RankInterval <= 0 {
		Conf.RankInterval = 10
	}
	return
}
