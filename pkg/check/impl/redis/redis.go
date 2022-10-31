package redis

import (
	"encoding/json"
	"time"

	"github.com/douyu/juno-agent/pkg/check/view"
	"github.com/garyburd/redigo/redis"
	"github.com/spf13/cast"
)

// RedisHealthCheck redis check config
type RedisHealthCheck struct {
	// Default 1s
	DialTimeout time.Duration `json:"dialTimeout" toml:"dialTimeout"`
	// Default 1s
	ReadTimeout time.Duration `json:"readTimeout" toml:"readTimeout"`
	//Default 1s
	WriteTimeout time.Duration `json:"writeTimeout" toml:"writeTimeout"`
	//Default 0
	DB int `json:"db" toml:"db"`
	//redis://:password@ip:port
	Addr     string `json:"addr"`
	PassWord string `json:"passWord" toml:"passWord"`
}

// NewRedisHealthCheck new a instance
func NewRedisHealthCheck() *RedisHealthCheck {
	return DefaultRedisNodeConfig()
}

// LoadExtConfig parse redis config
func (h *RedisHealthCheck) LoadExtConfig(extConfig string) (err error) {
	if err = json.Unmarshal([]byte(extConfig), &h); err != nil {
		return
	}
	return
}

// DoHealthCheck check is invoked periodically to perform the mysql check
func (h *RedisHealthCheck) DoHealthCheck() (resHealthCheck *view.ResHealthCheck, err error) {
	dialOptions := []redis.DialOption{
		redis.DialConnectTimeout(h.DialTimeout),
		redis.DialReadTimeout(h.ReadTimeout),
		redis.DialWriteTimeout(h.WriteTimeout),
		redis.DialDatabase(h.DB),
	}
	if h.PassWord != "" {
		dialOptions = append(dialOptions, redis.DialPassword(h.PassWord))
	}
	// now := time.Now()
	_, err = redis.Dial("tcp", h.Addr, dialOptions...)
	// statDial(now, redis.Addr, err)
	if err != nil {
		return
	}
	resHealthCheck = view.HealthCheckResult("redis", true, "success")
	return
}

// DefaultRedisNodeConfig return default config
func DefaultRedisNodeConfig() *RedisHealthCheck {
	return &RedisHealthCheck{
		DialTimeout:  cast.ToDuration("2s"),
		ReadTimeout:  cast.ToDuration("2s"),
		WriteTimeout: cast.ToDuration("2s"),
		DB:           0,
	}
}
