// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package device

import (
	"context"

	"github/ikun2021/bestpractice/app/user/api/internal/svc"
	"github/ikun2021/bestpractice/app/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeviceInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDeviceInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeviceInfoLogic {
	return &GetDeviceInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDeviceInfoLogic) GetDeviceInfo(req *types.GetDeviceInfoReq) (resp *types.GetDeviceInfoResp, err error) {
	// todo: add your logic here and delete this line

	return
}
