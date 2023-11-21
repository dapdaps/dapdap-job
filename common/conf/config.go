package conf

import "dapdap-job/common/xtime"

type Pgsql struct {
	DSN         string
	Opts        string
	Active      int
	Idle        int
	IdleTimeout xtime.Duration
}
