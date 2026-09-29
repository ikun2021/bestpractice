package socialservicelogic

import (
	"context"

	"github/ikun2021/bestpractice/app/user/rpc/internal/svc"
	"github/ikun2021/bestpractice/pb/userpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPostInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPostInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostInfoLogic {
	return &GetPostInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPostInfoLogic) GetPostInfo(in *userpb.GetPostInfoReq) (*userpb.GetPostInfoResp, error) {
	// todo: add your logic here and delete this line

	return &userpb.GetPostInfoResp{}, nil
}
