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
	path = flag.String("f", "server/accountApi/conf/config.yaml", "config file path")
)

// 这部分 @Tag 设置用于排序, 需要排序的接口请按照下面的格式添加
// swag init 对 @Tag 只会从入口文件解析, 默认 main.go
// 也可通过 --generalInfo flag 指定其他文件
// @Tag.Name        Base
// @Tag.Name        SysUser
// @Tag.Description 用户
// @title                       AccountApi接口
// @version                     v2.7.9
// @description                 AccountApi接口
// @in                          header
// @name                        x-token
// @BasePath                    /
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
