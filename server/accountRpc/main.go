package main

import (
	"flag"
	"fmt"
	"github.com/luxun9527/zlog"
	accountPb "github/lunxun9527/bestpractice/pb/account"
	"github/lunxun9527/bestpractice/server/accountRpc/global"
	"github/lunxun9527/bestpractice/server/accountRpc/initializer"
	accountService "github/lunxun9527/bestpractice/server/accountRpc/service/account"
	"google.golang.org/grpc"
	"net"
)

var (
	path = flag.String("f", "example/server/accountRpc/conf/config.yaml", "config file path")
)

func main() {
	flag.Parse()
	initializer.Init(*path)
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", global.Config.Server.Port))
	if err != nil {
		zlog.Panicf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	accountPb.RegisterAccountSrvServer(s, accountService.AccountRpc)
	zlog.Infof("start rpc server on %d", global.Config.Server.Port)
	if err := s.Serve(listener); err != nil {
		zlog.Panicf("failed to serve: %v", err)
	}
}
