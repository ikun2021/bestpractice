package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/luxun9527/zlog"
	"github.com/spf13/cast"
	"github/lunxun9527/bestpractice/server/accountApi/global"
	"github/lunxun9527/bestpractice/server/accountApi/initializer"
	"github/lunxun9527/bestpractice/server/accountApi/router"
)

var (
	path = flag.String("f", "example/server/accountApi/conf/config.yaml", "config file path")
)

func main() {
	flag.Parse()
	e := gin.New()
	gin.SetMode(gin.ReleaseMode)
	initializer.Init(*path)
	router.InitRouter(e)
	zlog.Infof("account api server success on %v", global.Config.Server.Port)
	if err := e.Run(fmt.Sprintf("0.0.0.0:" + cast.ToString(global.Config.Server.Port))); err != nil {
		zlog.Panicf("startup service failed, err:%v", err)
	}
}
