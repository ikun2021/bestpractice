package main

import (
	"flag"
	"fmt"

	"github/ikun2021/bestpractice/app/user/rpc/internal/config"
	deviceserviceServer "github/ikun2021/bestpractice/app/user/rpc/internal/server/deviceservice"
	socialserviceServer "github/ikun2021/bestpractice/app/user/rpc/internal/server/socialservice"
	"github/ikun2021/bestpractice/app/user/rpc/internal/svc"
	"github/ikun2021/bestpractice/pb/userpb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		userpb.RegisterDeviceServiceServer(grpcServer, deviceserviceServer.NewDeviceServiceServer(ctx))
		userpb.RegisterSocialServiceServer(grpcServer, socialserviceServer.NewSocialServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
