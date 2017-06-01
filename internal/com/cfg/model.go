package cfg

import (
	"github.com/larspensjo/config"
	"log"
	"strings"
)

type Conf struct {
	Srv            *Srv
	Rc             *Rc
	Log            *Log
	RpcHousekeeper *RpcHousekeeper
}

func NewConf() *Conf {
	return &Conf{
		Srv: &Srv{},
		Rc: &Rc{},
		Log:&Log{},
		RpcHousekeeper:&RpcHousekeeper{},
	}
}

// TCP server configs
type Srv struct {
	Version           string
	Port              int

	ValidTimeout      int // unit: s, TCP connected max valid seconds
	HbInterval        int // unit: s, Heartbeat interval
	ScavengeInterval  int // unit: s, Cron scavenge invalid sessions

	MsgWaitAckSeconds int // unit: s
}

// Redis consistence configs
type Rc struct {
	Servers     []string
	Password    string
	Db          int

	PoolSize    int
	PoolTimeout int // connections in poll's max waiting seconds
	VnodeNum    int // consistence vnode num
}

type Log struct {
	ChannelLen int
	Path       string
	Level      int
}

type RpcHousekeeper struct {
	Host string
	Port int
}

func parseConfigFile(path string) {
	cfg, err := config.ReadDefault(path)
	if err != nil {
		log.Fatalln("[E] ", err.Error(), " cfg file: ", path)
	}

	parseSrv(cfg)
	parseRc(cfg)
	parseLog(cfg)
	parseRpcHousekeeper(cfg)
}

func validSection(cfg *config.Config, section string)  {
	if !cfg.HasSection(section) {
		log.Fatalf("[E] config options of %s not found!\n", section)
	}
}

func parseSrv(cfg *config.Config) *Srv {
	section := "server"
	validSection(cfg, section)

	c := C.Srv
	c.Version, _ = cfg.String(section, "version")
	c.Port, _ = cfg.Int(section, "port")
	c.ValidTimeout, _ = cfg.Int(section, "valid_timeout")
	c.HbInterval, _ = cfg.Int(section, "heartbeat_interval")
	c.ScavengeInterval, _ = cfg.Int(section, "scavenge_interval")
	c.MsgWaitAckSeconds, _ = cfg.Int(section, "msg_wait_ack_seconds")

	return c
}

func parseRc(cfg *config.Config) *Rc {
	section := "rc"
	validSection(cfg, section)

	c := C.Rc
	servers, _ := cfg.String(section, "servers")
	c.Servers = strings.Split(strings.Trim(servers, ","), ",")
	c.Password, _ = cfg.String(section, "password")
	c.Db, _ = cfg.Int(section, "db")
	c.PoolSize, _ = cfg.Int(section, "pool_size")
	c.PoolTimeout, _ = cfg.Int(section, "pool_timeout")
	c.VnodeNum, _ = cfg.Int(section, "vnod_num")

	return c
}

func parseLog(cfg *config.Config) *Log {
	section := "log"
	validSection(cfg, section)

	c := C.Log
	c.ChannelLen, _ = cfg.Int(section, "channel_len")
	c.Path, _ = cfg.String(section, "path")
	c.Path = strings.TrimRight(c.Path, "/") + "/"
	c.Level, _ = cfg.Int(section, "level")

	return c
}

func parseRpcHousekeeper(cfg *config.Config) *RpcHousekeeper {
	section := "housekeeper"
	validSection(cfg, section)

	c := C.RpcHousekeeper
	c.Host, _ = cfg.String(section, "host")
	c.Port, _ = cfg.Int(section, "port")

	return c
}
