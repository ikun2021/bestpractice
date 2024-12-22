package config

import "github/lunxun9527/bestpractice/pkg/xetcd"

type ServerConf struct {
	Port         int32                  `mapstructure:"port"`
	RegisterConf xetcd.EtcdRegisterConf `mapstructure:"registerConf"`
}
