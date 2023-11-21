package main

import (
	"dapdap-job/common/log"
	"dapdap-job/common/shutdown"
	"dapdap-job/conf"
	"dapdap-job/service"
	"flag"
	"os"
	"os/signal"
	"syscall"
)

// nolint:all
func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}

	log.Init(conf.Conf.Log, conf.Conf.Debug)
	log.Info("dapdap-job service start")

	service.Init(conf.Conf)
	log.Info("dapdap-job service init")

	go func() {
		service.DapdapService.StartActionTask()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	s := <-c
	log.Info("dapdap-job exit for signal %v", s)

	shutdown.StopAndWaitAll()
}
