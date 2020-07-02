package cfg

import (
	"fmt"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/ecode"
	"github.com/douyu/jupiter/pkg/util/xtime"
	"github.com/douyu/jupiter/pkg/xlog"
	"os"
)

var Cfg cfg

//
type cfg struct {
	Server    Server
	HeartBeat HeartBeat
}

// DefaultConfig ...
func defaultConfig() cfg {
	return cfg{
		Server: Server{
			Http: ServerSchema{
				Enable: true,
				Host:   "0.0.0.0",
				Port:   50010,
			},
			Grpc: ServerSchema{
				Enable: true,
				Host:   "0.0.0.0",
				Port:   50011,
			},
			Govern: ServerSchema{
				Enable: true,
				Host:   "0.0.0.0",
				Port:   50012,
			},
		},
		HeartBeat: HeartBeat{
			Enable:     true,
			Debug:      false,
			Addr:       "",
			Internal:   xtime.Duration("60s"),
			HostName:   "",
			RegionCode: "",
			RegionName: "",
			ZoneCode:   "",
			ZoneName:   "",
			Env:        "",
		},
	}
}

// StdConfig ...
func InitCfg() {
	var (
		err    error
		config = defaultConfig()
	)
	if err = conf.UnmarshalKey("", &config); err != nil {
		xlog.Panic("parse cfg error", xlog.FieldErrKind(ecode.ErrKindUnmarshalConfigErr), xlog.FieldErr(err), xlog.FieldKey("system cfg"), xlog.FieldValueAny(config))
	}

	config.parseHeartBeat()
	Cfg = config
}

func (c *cfg) parseHeartBeat() {
	c.HeartBeat.RegionCode = os.Getenv(c.HeartBeat.RegionCode)
	c.HeartBeat.RegionName = os.Getenv(c.HeartBeat.RegionName)
	c.HeartBeat.ZoneCode = os.Getenv(c.HeartBeat.ZoneCode)
	c.HeartBeat.ZoneName = os.Getenv(c.HeartBeat.ZoneName)
	c.HeartBeat.HostName = GetHostName(c.HeartBeat.HostName)
	env := os.Getenv(c.HeartBeat.Env)
	if env != "" {
		c.HeartBeat.Env = env
	}
}

// GetHostName ...
func GetHostName(hostName string) string {
	if host := os.Getenv(hostName); host != "" {
		return host
	}
	name, err := os.Hostname()
	if err != nil {
		return fmt.Sprintf("error:%s", err.Error())
	}
	return name
}
