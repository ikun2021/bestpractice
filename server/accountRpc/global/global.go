package global

import (
	"github.com/redis/go-redis/v9"
	"github/lunxun9527/bestpractice/server/accountRpc/config"
	accountQuery "github/lunxun9527/bestpractice/server/accountRpc/dao/account/query"
)

var (
	AccountDB *accountQuery.Query
	Config    *config.Config
	RedisCli  *redis.Client
)
