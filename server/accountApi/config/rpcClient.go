package config

import "github/lunxun9527/bestpractice/pkg/xetcd"

type RpcClient struct {
	EtcdConf       xetcd.EtcdConf
	TargetConfList []*TargetConf
}

type TargetConf struct {
	Key     string
	TimeOut int64
}
