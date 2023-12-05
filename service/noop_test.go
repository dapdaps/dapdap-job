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

	DapdapService.StartQuestTask()
}
