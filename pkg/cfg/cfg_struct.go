package cfg

import "time"

type ServerSchema struct {
	Enable bool
	Host   string
	Port   int
}

// Agent Server
type Server struct {
	Http   ServerSchema
	Grpc   ServerSchema
	Govern ServerSchema
}

type HeartBeat struct {
	Enable     bool          `json:"enable"`
	Debug      bool          `json:"debug"`
	Addr       string        `json:"addr"`
	Internal   time.Duration `json:"internal"`
	HostName   string        `json:"host_name"`
	RegionCode string        `json:"region_code"`
	RegionName string        `json:"region_name"`
	ZoneCode   string        `json:"zone_code"`
	ZoneName   string        `json:"zone_name"`
	Env        string        `json:"env"`
}
