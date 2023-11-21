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
	Debug   bool
	Timeout int64
	Log     *log.Config
	Pgsql   *Pgsql
}

type Pgsql struct {
	DB *conf.Pgsql
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
	return
}
