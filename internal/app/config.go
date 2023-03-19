package app

import (
	"github.com/core-go/core"
	"github.com/core-go/log"
	mid "github.com/core-go/log/middleware"
)

type Config struct {
	Server     core.ServerConf `mapstructure:"server"`
	Cql        Cassandra       `mapstructure:"cassandra"`
	Log        log.Config      `mapstructure:"log"`
	MiddleWare mid.LogConfig   `mapstructure:"middleware"`
}

type Cassandra struct {
	PublicIp	string	`mapstructure:"public_ip"`
	UserName	string	`mapstructure:"user_name"`
	Password	string	`mapstructure:"password"`
}
