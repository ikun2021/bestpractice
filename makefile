urpc:
	$(MAKE) -C ../xpb gen
	goctl rpc protoc -I ../xpb/userpb ../xpb/userpb/user.proto --go_out=. --go_opt=module=github/ikun2021/bestpractice --go-grpc_out=. --go-grpc_opt=module=github/ikun2021/bestpractice --zrpc_out=app/user/cmd/rpc --client=false --multiple
	protoc -I ../xpb/userpb --go_out=. --go_opt=module=github/ikun2021/bestpractice ../xpb/userpb/device.proto
	protoc -I ../xpb/userpb --go_out=. --go_opt=module=github/ikun2021/bestpractice ../xpb/userpb/social.proto
