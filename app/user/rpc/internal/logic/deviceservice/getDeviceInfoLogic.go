package deviceservicelogic

import (
	"context"

	"github/ikun2021/bestpractice/app/user/rpc/internal/svc"
	"github/ikun2021/bestpractice/pb/userpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeviceInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDeviceInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeviceInfoLogic {
	return &GetDeviceInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDeviceInfoLogic) GetDeviceInfo(in *userpb.GetDeviceInfoReq) (*userpb.GetDeviceInfoResp, error) {
	// todo: add your logic here and delete this line

	return &userpb.GetDeviceInfoResp{}, nil
}
