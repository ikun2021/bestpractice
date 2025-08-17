package account

import (
	"context"
	"github.com/luxun9527/zlog"
	accountPb "github/lunxun9527/bestpractice/pb/account"
	accountApiModel "github/lunxun9527/bestpractice/server/accountApi/model/account"
	"github/lunxun9527/bestpractice/server/accountApi/rpcClient"
)

var UserService = &userService{}

type userService struct{}

func (*userService) GetAccountInfo(ctx context.Context, req *accountApiModel.GetAccountInfoReq) (*accountApiModel.GetAccountInfoResp, error) {
	accountInfo, err := rpcClient.AccountClient.GetAccountInfo(ctx, &accountPb.GetAccountInfoReq{AccountId: req.AccountId})
	if err != nil {
		zlog.Errorf("GetAccountInfoReq err: %v", err)
		return nil, err
	}
	return &accountApiModel.GetAccountInfoResp{
		AccountId:   accountInfo.AccountId,
		AccountName: accountInfo.AccountName,
	}, nil
}
func (*userService) AddAccount(ctx context.Context, req *accountApiModel.AddAccountReq) error {
	_, err := rpcClient.AccountClient.RegisterUser(ctx, &accountPb.RegisterUserReq{
		AccountId:   req.AccountID,
		AccountName: req.AccountName,
		Password:    req.Password,
	})
	if err != nil {
		zlog.Errorf("GetAccountInfoReq err: %v", err)
		return err
	}
	return nil
}
